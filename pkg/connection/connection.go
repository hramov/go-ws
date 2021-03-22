package connection

import (
	"bufio"
	"fmt"
	"net"
	"strconv"

	"github.com/hramov/battleship_server/pkg/utils"
)

type Server struct {
	protocol string
	ip       string
	port     string
	ln       net.Listener
}

type Client struct {
	ID      int
	EnemyID int
	Name    string
	Socket  net.Conn
	From    chan string
	To      chan string
}

func Execute(protocol, ip, port string, handler func(s *Server, client Client, clients map[int]Client)) {
	server := Server{protocol, ip, port, nil}
	server.createServer()
	clients := make(map[int]Client)
	server.maintainConnections(clients, handler)
}

func (s *Server) createServer() {
	s.ln, _ = net.Listen(s.protocol, s.ip+":"+s.port)
	fmt.Println("Server is listening for connections on " + s.ip + ":" + s.port)
}

func (s *Server) listen(client Client) {
	for {
		rawData, _ := bufio.NewReader(client.Socket).ReadString('\n')
		client.From <- rawData
	}
}

func (s *Server) On(client Client, rawEvent string, callback func(data string)) {
	rawData := <-client.From
	event, data := utils.Split(rawData, "|")
	// fmt.Print()
	if event == rawEvent {
		callback(string(data))
	}
}

func (s *Server) speak(client Client) {
	for {
		rawData := <-client.To
		event, data := utils.Split(rawData, "|")
		client.Socket.Write([]byte(string(event) + "|" + string(data) + "\n"))
	}
}

func (s *Server) Emit(client Client, event string, data string) {
	client.To <- string(event + "|" + data)
}

func (s *Server) maintainConnections(clients map[int]Client, handler func(s *Server, client Client, clients map[int]Client)) {
	conn, _ := s.ln.Accept()

	ID := len(clients) + 1
	client := Client{ID, 0, "", conn, make(chan string, 10), make(chan string, 10)}
	clients[ID] = client

	go s.speak(client)
	go handler(s, client, clients)

	client.From <- "connect|" + strconv.Itoa(client.ID)
	s.listen(client)
}
