package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hramov/battleship_server/pkg/battlefield"
	connection "github.com/hramov/battleship_server/pkg/connection"
	"github.com/hramov/battleship_server/pkg/ship"
)

func main() {
	connection.Execute("tcp", "127.0.0.1", "5000",
		func(s *connection.Server, client connection.Client, clients map[int]connection.Client) {

			b := battlefield.BattleField{}
			sh := ship.Ship{}
			var ships []ship.Ship

			for {
				/** Обрабатываем событие подключения */
				s.On(client, "connect", func(data string) {
					s.Emit(client, "connect", "")
					s.Emit(client, "whoami", "")
				})

				s.On(client, "sendName", func(data string) {
					s.Emit(client, "enemy", "2")
					client.Name = data
					/** Создаем и отправляем поле */
					b.CreateField(client.ID, 2)
					jsonData, err := json.Marshal(b)
					if err != nil {
						log.Fatal(err)
					}
					s.Emit(client, "drawField", string(jsonData))
					s.Emit(client, "placeShip", "")
				})

				s.On(client, "sendShip", func(data string) {
					json.Unmarshal([]byte(data), &sh)
					err := b.CheckShip(sh)
					if err != nil {
						s.Emit(client, "wrongShip", fmt.Sprint(err))
						s.Emit(client, "placeShip", "")
					} else {
						ships = append(ships, sh)
						if len(ships) < 10 {
							s.Emit(client, "placeShip", "")
						} else {
							s.Emit(client, "hit", "")
						}
					}
				})

				/** Обрабатываем событие отключения */
				s.On(client, "disconnect", func(data string) {
					fmt.Println(client.Name + " disconnected!")
					close(client.From)
					close(client.To)
					delete(clients, client.ID)
				})
			}
		}) // Запуск сокет-сервера

	// gameloop.Init()
}
