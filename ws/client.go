package ws

import "net"

type Client struct {
	ID      int
	EnemyID int
	Socket  net.Conn
	From    chan string
	To      chan string
	Turn    bool
}

func Clear() {
	clients = make(Clients)
}

func FindClientByID(ID int) *Client {
	cl := clients[ID]
	return &cl
}

func GetClientsNumber() int {
	return len(clients)
}
