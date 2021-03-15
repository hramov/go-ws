package connection

import (
	"bufio"
	"fmt"
	"net"

	utils "github.com/hramov/battleship_server/pkg/utils"
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
	transmitter chan Channel
	receiver    chan Channel
}

type Event string
type Data string

type Channel map[Event]Data

var clients []Client

func Execute(protocol, ip, port string) {
	server := Server{protocol, ip, port, nil}
	server.createServer()
	for {
		server.maintainConnections()
	}
}

func (s *Server) createServer() {
	s.ln, _ = net.Listen(s.protocol, s.ip+":"+s.port)
	fmt.Println("Server is listening for connections on " + s.ip + ":" + s.port)
}

func (s *Server) listen(client Client) {
	eventMap := make(Channel)
	rawData, _ := bufio.NewReader(client.socket).ReadString('\n')
	event, data := utils.Split(rawData, ":")
	eventMap[Event(event)] = Data(data)
	client.receiver <- eventMap
}

func (s *Server) on(client Client, rawEvent string, callback func(data string)) {
	eventMap := <-client.receiver
	for event, data := range eventMap {
		if event == Event(rawEvent) {
			callback(string(data))
		}
	}
}

func (s *Server) speak(client Client) {
	eventMap := <-client.transmitter
	for event, data := range eventMap {
		client.socket.Write([]byte(string(event) + ":" + string(data) + "\n"))
	}
}

func (s *Server) emit(client Client, event string, data string) {
	eventMap := make(Channel)
	eventMap[Event(event)] = Data(data)
	client.transmitter <- eventMap
}

func (s *Server) maintainConnections() {
	conn, _ := s.ln.Accept()

	client := Client{len(clients) + 1, "", conn, make(chan Channel), make(chan Channel)}
	fmt.Println("New connection with ID =", client.id)
	clients = append(clients, client)

	go s.listen(client)
	go s.speak(client)

	s.emit(client, "whoami", "")

	s.on(client, "sendName", func(data string) {
		client.name = data
		fmt.Println(client.name)
	})

}
