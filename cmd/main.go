package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
	"github.com/hramov/battleship_server/pkg/battlefield"
	"github.com/hramov/battleship_server/pkg/server"
	"github.com/hramov/battleship_server/pkg/ship"
)

const port = ":3000"

var clients = make(map[socketio.Conn]battlefield.Client)

func main() {

	server := server.CreateServer()
	defer server.Close()

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")

		fmt.Printf("Игрок %s подключился!\n", s.ID())

		c := battlefield.Client{}
		c.CreateField(s.ID())
		clients[s] = c

		data, _ := json.Marshal(c)

		s.Emit("whoami", s.ID())    // Игрок получил свой ID
		s.Emit("updateField", data) // У игрока отрисовалось поле
		s.Emit("placeShip")         // Команда на расстановку кораблей

		return nil
	})

	server.OnEvent("/", "sendShip", func(s socketio.Conn, data []byte) {
		client, _ := FindClientByID(s.ID())
		ship := ship.Ship{}
		json.Unmarshal(data, &ship)
		err := client.CheckShip(ship)

		if err != nil {
			fmt.Println(err)
			s.Emit("Ошибка!")
		} else {
			fmt.Printf("Корабль игрока %s успешно добавлен!", client.ID)
			client.CreateShip(ship)
			data, _ = json.Marshal(client)
			s.Emit("updateField", data)
			// s.Emit("placeShip")
		}
	})

	// server.OnEvent("/", "sendShot", func(Player bool, shot ship.Shot) {
	// 	b.CheckShot(Player, shot)
	// 	b.CheckHit(Player, shot)

	// })

	// server.OnDisconnect("/", func(s socketio.Conn, reason string) {
	// 	fmt.Println("closed", reason)
	// })

	go server.Serve()

	http.Handle("/", server)

	fmt.Printf("Сервер запущен на http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

func FindClientByID(ID string) (battlefield.Client, error) {
	for _, client := range clients {
		if client.ID == ID {
			return client, nil
		}
	}
	return battlefield.Client{}, fmt.Errorf("%s\n", "Игрок не найден")
}
