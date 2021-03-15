package main

import (
	connection "github.com/hramov/battleship_server/pkg/connection"
)

func main() {
	connection.Execute("tcp", "127.0.0.1", "5000")
}
