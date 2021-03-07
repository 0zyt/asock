package main

import (
	"log"
)

func main() {
	server, err := NewLocalServer(":2000", ":8888")
	if err != nil {
		log.Fatalln(err)
	}
	go server.Listen()
	newServer, _ := NewServer(":8888")
	newServer.Listen()
}
