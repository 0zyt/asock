package main

import "log"

func main() {
	server, err := NewLocalServer("127.0.0.1:7474", "127.0.0.1:12346")
	if err != nil {
		log.Fatalln(err)
	}
	go server.Listen()
	newServer, _ := NewServer("127.0.0.1:12346")
	newServer.Listen()
}
