package main

import (
	"log"
	"net"
)

func HandleConn(listener net.Listener, handleProtocol func(conn *net.Conn)) {
	defer listener.Close()
	for {
		client, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleProtocol(&client)
	}
}
