package main

import (
	"encoding/json"
	"strconv"

	"github.com/hramov/battleship_server/pkg/battlefield"
	connection "github.com/hramov/battleship_server/pkg/connection"
	"github.com/hramov/battleship_server/pkg/ship"
	"github.com/hramov/battleship_server/pkg/utils"
)

func main() {

	b := battlefield.BattleField{}
	sh := ship.Ship{}
	var ships []ship.Ship

	clients := make(map[int]connection.Client)

	s := connection.Execute("tcp", "127.0.0.1", "5000")

	handlers := make(map[string]func(client connection.Client, data string))

	handlers["connect"] = func(client connection.Client, data string) {
		s.Emit(client, "connect", "")
		s.Emit(client, "whoami", "")
	}

	handlers["sendName"] = func(client connection.Client, data string) {
		s.Emit(client, "enemy", "2")
		s.Emit(client, "drawField", "")
		s.Emit(client, "placeShip", "")
	}

	handlers["sendShip"] = func(client connection.Client, data string) {
		json.Unmarshal([]byte(data), &sh)
		err := b.CheckShip(sh)
		if err != nil {
			s.Emit(client, "placeShip", err.Error())
		} else {
			ships = append(ships, sh)
			if len(ships) < 10 {
				s.Emit(client, "placeShip", "")
			} else {
				s.Emit(client, "hit", "")
			}
		}
	}

	handlers["disconnect"] = func(client connection.Client, data string) {
		utils.Log("User: " + strconv.Itoa(client.ID) + " disconnected!")
		client.Socket.Close()
		delete(clients, client.ID)
		utils.Log("Users remain: " + strconv.Itoa(len(clients)))
	}

	maintainConnections(s, &clients, &handlers)
}

func maintainConnections(s *connection.Server, clients *map[int]connection.Client, handlers *map[string]func(client connection.Client, data string)) {
	for {
		conn, _ := s.Ln.Accept()
		ID := len(*clients) + 1
		client := connection.Client{ID, 0, conn, make(chan string), make(chan string)}
		(*clients)[ID] = client

		utils.Log("User " + strconv.Itoa(client.ID) + " connected!")

		go s.Speak(client)
		go s.Listen(client)
		go s.On(client, handlers)

		client.From <- "connect|" + strconv.Itoa(client.ID)
	}
}
