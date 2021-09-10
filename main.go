package main

import (
	gows "github.com/hramov/go-ws/lib"
)

func main() {
	s := gows.Execute("tcp", "127.0.0.1", "5000")

	gows.HandlersMap["connect"] = func(client *gows.Client, data string) {}
	gows.HandlersMap["disconnect"] = func(client *gows.Client, data string) {}

	gows.Proceed(s)
}
