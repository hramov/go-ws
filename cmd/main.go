package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
	"github.com/hramov/battleship_server/pkg/battlefield"
	"github.com/hramov/battleship_server/pkg/server"
)

const port = ":3000"

type Client struct {
	ID    string                  `json:"ID"`
	Field battlefield.BattleField `json:"field"`
}

func main() {

	var clients []Client

	server := server.CreateServer()
	defer server.Close()

	b := battlefield.BattleField{}

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		client := Client{s.ID(), b.CreateField()}
		data, _ := json.Marshal(client)
		fmt.Println(data)
		s.Emit("field", data)
		clients = append(clients, client)
		return nil
	})

	server.OnEvent("/", "shot", func(shot string) {
		s.Emit("reply", "have "+msg)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})

	go server.Serve()

	http.Handle("/", server)

	fmt.Printf("Сервер запущен на http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
