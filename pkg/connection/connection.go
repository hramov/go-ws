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
	ID          int
	EnemyID     int
	Name        string
	Socket      net.Conn
	Transmitter chan string
	Receiver    chan string
}

func Execute(protocol, ip, port string, handler func(s *Server, client Client, clients map[int]Client)) {
	server := Server{protocol, ip, port, nil}
	server.createServer()
	clients := make(map[int]Client)
	for {
		server.maintainConnections(clients, handler)
	}
}

func (s *Server) createServer() {
	s.ln, _ = net.Listen(s.protocol, s.ip+":"+s.port)
	fmt.Println("Server is listening for connections on " + s.ip + ":" + s.port)
}

func (s *Server) listen(client Client) {
	rawData, _ := bufio.NewReader(client.Socket).ReadString('\n')
	client.Receiver <- rawData
}

func (s *Server) On(client Client, rawEvent string, callback func(data string)) {
	rawData := <-client.Receiver
	event, data := utils.Split(rawData, ":")
	if event == rawEvent {
		callback(string(data))
	}
}

func (s *Server) speak(client Client) {
	rawData := <-client.Transmitter
	event, data := utils.Split(rawData, ":")
	client.Socket.Write([]byte(string(event) + ":" + string(data) + "\n"))
}

func (s *Server) Emit(client Client, event string, data string) {
	rawData := event + ":" + data
	client.Transmitter <- rawData
}

func (s *Server) maintainConnections(clients map[int]Client, handler func(s *Server, client Client, clients map[int]Client)) {
	conn, _ := s.ln.Accept()

	ID := len(clients) + 1
	client := Client{ID, 0, "", conn, make(chan string), make(chan string)}
	clients[ID] = client

	go func() {
		go s.listen(client)
		go s.speak(client)
		rawData := "connect:" + string(client.ID)
		client.Receiver <- rawData
	}()

	go handler(s, client, clients)

}
