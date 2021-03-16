package connection

import (
	"bufio"
	"fmt"
	"net"

	"github.com/hramov/battleship_server/pkg/utils"
)

type Server struct {
	protocol string
	ip       string
	port     string
	ln       net.Listener
}

type Client struct {
	id          int
	name        string
	socket      net.Conn
	transmitter chan string
	receiver    chan string
}

func Execute(protocol, ip, port string) {
	server := Server{protocol, ip, port, nil}
	server.createServer()
	clients := make(map[int]Client)
	for {
		server.maintainConnections(clients)
	}
}

func (s *Server) createServer() {
	s.ln, _ = net.Listen(s.protocol, s.ip+":"+s.port)
	fmt.Println("Server is listening for connections on " + s.ip + ":" + s.port)
}

func (s *Server) listen(client Client) {
	rawData, _ := bufio.NewReader(client.socket).ReadString('\n')
	client.receiver <- rawData
}

func (s *Server) on(client Client, rawEvent string, callback func(data string)) {
	rawData := <-client.receiver
	event, data := utils.Split(rawData, ":")
	if event == rawEvent {
		callback(string(data))
	}
}

func (s *Server) speak(client Client) {
	rawData := <-client.transmitter
	event, data := utils.Split(rawData, ":")
	client.socket.Write([]byte(string(event) + ":" + string(data) + "\n"))
}

func (s *Server) emit(client Client, event string, data string) {
	rawData := event + ":" + data
	client.transmitter <- rawData
}

func (s *Server) maintainConnections(clients map[int]Client) {
	conn, _ := s.ln.Accept()

	id := len(clients) + 1
	client := Client{id, "", conn, make(chan string), make(chan string)}
	clients[id] = client

	go func() {
		go s.listen(client)
		go s.speak(client)
		rawData := "connect:" + string(client.id)
		client.receiver <- rawData
	}()

	go func() {
		s.emit(client, "whoami", "")

		s.on(client, "connect", func(data string) {
			fmt.Println("Connect!")
		})

		s.on(client, "sendName", func(data string) {
			client.name = data
			fmt.Println(client.name)
		})

		s.on(client, "disconnect", func(data string) {
			fmt.Println(client.name + " disconnected!")
			close(client.receiver)
			close(client.transmitter)
			delete(clients, client.id)
		})
	}()

}
