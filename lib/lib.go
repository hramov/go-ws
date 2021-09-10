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

func Execute(protocol, ip, port string) (*server, error) {
	server := server{protocol, ip, port, nil}
	err := server.createServer()
	if err != nil {
		return nil, err
	}
	return &server, nil
}

func Proceed(s *server) error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			return err
		}
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
