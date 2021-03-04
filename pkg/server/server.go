package server

import (
	"log"

	socketio "github.com/googollee/go-socket.io"
)

func CreateServer() *socketio.Server {
	server, serveError := socketio.NewServer(nil)
	if serveError != nil {
		log.Fatalln(serveError)
	}
	return server
}
