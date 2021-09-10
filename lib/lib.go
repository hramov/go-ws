package lib

import (
	"strconv"
)

type handlers map[string]func(client *Client, data string)
type clients map[int]Client

var (
	clientsMap  clients
	HandlersMap handlers
)

func init() {
	clientsMap = make(clients)
	HandlersMap = make(handlers)
}

func Execute(protocol, ip, port string) *server {
	server := server{protocol, ip, port, nil}
	server.createServer()
	return &server
}

func Proceed(s *server) {
	for {

		conn, _ := s.ln.Accept()
		ID := len(clientsMap) + 1
		client := Client{ID, 0, conn, make(chan string), make(chan string), false}
		clientsMap[ID] = client

		Log("User " + strconv.Itoa(client.ID) + " connected!")

		go s.speak(&client)
		go s.listen(&client)
		go s.on(&client, &HandlersMap)

		client.From <- "connect|" + strconv.Itoa(client.ID)

	}
}
