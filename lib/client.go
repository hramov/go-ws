package lib

import "net"

type Client struct {
	ID      int
	EnemyID int
	Socket  net.Conn
	From    chan string
	To      chan string
	Turn    bool
}

func clear() {
	clientsMap = make(clients)
}

func findClientByID(ID int) *Client {
	cl := clientsMap[ID]
	return &cl
}

func getClientsNumber() int {
	return len(clientsMap)
}
