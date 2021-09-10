package main

import (
	"github.com/hramov/go-ws/ws"
)

func main() {

	s := ws.Execute("tcp", "127.0.0.1", "5000")

	handlers := make(ws.Handlers)

	handlers["connect"] = func(client *ws.Client, data string) {
	}

	handlers["disconnect"] = func(client *ws.Client, data string) {
	}

	ws.Proceed(s, &handlers)
}
