package example

import (
	"log"

	gows "github.com/hramov/go-ws/lib"
)

func defaultServer() {
	s, err := gows.Execute("tcp", "127.0.0.1", "5000")

	if err != nil {
		log.Fatal(err)
	}

	gows.HandlersMap["connect" /** Event title */] = func(client *gows.Client, data string) { /** Handler function code */ }
	gows.Proceed(s)
}
