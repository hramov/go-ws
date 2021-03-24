package connection

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/hramov/battleship_server/pkg/utils"
)

type Server struct {
	Protocol string
	Ip       string
	Port     string
	Ln       net.Listener
}

type Client struct {
	ID      int
	EnemyID int
	Socket  net.Conn
	From    chan string
	To      chan string
}

func Execute(protocol, ip, port string) *Server {
	server := Server{protocol, ip, port, nil}
	server.createServer()
	return &server
}

func (s *Server) createServer() { // +
	s.Ln, _ = net.Listen(s.Protocol, s.Ip+":"+s.Port)
	fmt.Println("Server is listening for connections on " + s.Ip + ":" + s.Port)
}

func (s *Server) Listen(client Client) {
	for {
		rawData, err := bufio.NewReader(client.Socket).ReadString('\n')
		if err != nil && err.Error() == "EOF" {
			client.From <- string("disconnect|" + strconv.Itoa(client.ID))
			return
		} else if err != nil {
			log.Fatal(err)
		}
		client.From <- rawData
	}
}

func (s *Server) On(client Client, handlers *map[string]func(client Client, data string)) {
	for {
		time.Sleep(time.Second / 1000)
		rawData := <-client.From
		rawEvent, data := utils.Split(rawData, "|")
		for event, handler := range *handlers {
			if event == rawEvent {
				handler(client, data)
			}
		}
	}
}

func (s *Server) Speak(client Client) {
	for {
		rawData := <-client.To
		event, data := utils.Split(rawData, "|")
		client.Socket.Write([]byte(string(event) + "|" + string(data) + "\n"))
	}
}

func (s *Server) Emit(client Client, event, data string) {
	time.Sleep(time.Second / 1000)
	client.To <- string(event + "|" + data)
}
