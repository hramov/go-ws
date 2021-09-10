package ws

import (
	"strconv"
)

var clients = make(Clients)

type Handlers map[string]func(client *Client, data string)
type Clients map[int]Client

func Execute(protocol, ip, port string) *Server {
	server := Server{protocol, ip, port, nil}
	server.createServer()
	return &server
}

func Proceed(s *Server, handlers *Handlers) {
	for {

		conn, _ := s.Ln.Accept()
		ID := len(clients) + 1
		client := Client{ID, 0, conn, make(chan string), make(chan string), false}
		clients[ID] = client

		Log("User " + strconv.Itoa(client.ID) + " connected!")

		go s.Speak(&client)
		go s.Listen(&client)
		go s.On(&client, handlers)

		client.From <- "connect|" + strconv.Itoa(client.ID)
	}
}
