package lib

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
)

type server struct {
	protocol string
	ip       string
	port     string
	ln       net.Listener
}

func (s *server) createServer() { // +
	s.ln, _ = net.Listen(s.protocol, s.ip+":"+s.port)
	fmt.Println("Server is listening for connections on " + s.ip + ":" + s.port)
}

func (s *server) listen(client *Client) {
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

func (s *server) on(client *Client, handlers *handlers) {
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

func (s *server) speak(client *Client) {
	for {
		rawData := <-client.To
		event, data := Split(rawData, "|")
		client.Socket.Write([]byte(string(event) + "|" + string(data) + "\n"))
	}
}

func (s *server) emit(client *Client, event, data string) {
	time.Sleep(time.Second / 100)
	client.To <- string(event + "|" + data)
}

func (s *server) broadCast(clients *clients, event, data string) {
	for _, client := range *clients {
		time.Sleep(time.Second / 100)
		client.To <- string(event + "|" + data)
	}
}
