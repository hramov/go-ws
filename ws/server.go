package ws

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
)

type Server struct {
	Protocol string
	Ip       string
	Port     string
	Ln       net.Listener
}

func (s *Server) createServer() { // +
	s.Ln, _ = net.Listen(s.Protocol, s.Ip+":"+s.Port)
	fmt.Println("Server is listening for connections on " + s.Ip + ":" + s.Port)
}

func (s *Server) Listen(client *Client) {
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

func (s *Server) On(client *Client, handlers *Handlers) {
	for {
		time.Sleep(time.Second / 100)
		rawData := <-client.From
		rawEvent, data := Split(rawData, "|")
		for event, handler := range *handlers {
			if event == rawEvent {
				handler(client, data)
			}
		}
	}
}

func (s *Server) Speak(client *Client) {
	for {
		rawData := <-client.To
		event, data := Split(rawData, "|")
		client.Socket.Write([]byte(string(event) + "|" + string(data) + "\n"))
	}
}

func (s *Server) Emit(client *Client, event, data string) {
	time.Sleep(time.Second / 100)
	client.To <- string(event + "|" + data)
}

func (s *Server) BroadCast(clients *Clients, event, data string) {
	for _, client := range *clients {
		time.Sleep(time.Second / 100)
		client.To <- string(event + "|" + data)
	}
}
