package main

import (
	"fmt"

	connection "github.com/hramov/battleship_server/pkg/connection"
)

func main() {
	connection.Execute("tcp", "127.0.0.1", "5000",
		func(s *connection.Server, client connection.Client, clients map[int]connection.Client) {

			/** Обрабатываем событие подключения */
			s.On(client, "connect", func(data string) {
				fmt.Println("Connect!")
			})

			/** Узнаем имя игрока */
			s.Emit(client, "whoami", "")

			s.On(client, "sendName", func(data string) {
				client.Name = data
			})

			/** Обрабатываем событие отключения */
			s.On(client, "disconnect", func(data string) {
				fmt.Println(client.Name + " disconnected!")
				close(client.Receiver)
				close(client.Transmitter)
				delete(clients, client.ID)
			})
		}) // Запуск сокет-сервера
	// gameloop.Init()
}
