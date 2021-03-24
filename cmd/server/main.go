package main

import (
	"encoding/json"
	"strconv"

	"github.com/hramov/battleship_server/pkg/battlefield"
	connection "github.com/hramov/battleship_server/pkg/connection"
	"github.com/hramov/battleship_server/pkg/ship"
	"github.com/hramov/battleship_server/pkg/utils"
)

var ships []ship.Ship

var clients = make(map[int]connection.Client)
var battlefields = make(map[int]battlefield.BattleField)

func main() {

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
		sh := ship.Ship{}
		json.Unmarshal([]byte(data), &sh)
		b := FindFieldByClientID(client.ID)
		err := b.CheckShip(sh)
		if err != nil {
			s.Emit(client, "wrongShip", err.Error())
			s.Emit(client, "placeShip", err.Error())
		} else {
			s.Emit(client, "rightShip", "Ваш корабль успешно добавлен!")
			b.CreateShip(sh)
			battlefields[client.ID] = b
			ships = append(ships, sh)

			shipData, _ := json.Marshal(b)

			s.Emit(client, "drawField", string(shipData))

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
		delete(battlefields, client.ID)
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

		b := battlefield.BattleField{}
		b.CreateField()
		battlefields[client.ID] = b

		go s.Speak(client)
		go s.Listen(client)
		go s.On(client, handlers)

		client.From <- "connect|" + strconv.Itoa(client.ID)
	}
}

func FindFieldByClientID(ID int) battlefield.BattleField {
	var b battlefield.BattleField
	for id, battlefield := range battlefields {
		if id == ID {
			b = battlefield
		}
	}
	return b
}
