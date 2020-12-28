package main

import (
	"github.com/xtaci/kcp-go"
	"io"
	"log"
	"net"
)

type LocalConfig struct {
	local  *net.TCPAddr
	remote *net.UDPAddr
}

func NewLocalServer(listenAddr, remoteAddr string) (*LocalConfig, error) {
	local, err := net.ResolveTCPAddr("tcp", listenAddr)
	if err != nil {
		return nil, err
	}
	remote, err := net.ResolveUDPAddr("udp4", remoteAddr)
	if err != nil {
		return nil, err
	}
	return &LocalConfig{local, remote}, nil
}

func (local *LocalConfig) Listen() error {

	listener, err := net.ListenTCP("tcp", local.local)
	if err != nil {
		return err
	}
	protocol := func(conn *net.Conn) {
		Conn := (*conn).(*net.TCPConn)
		Conn.SetLinger(0)
		defer Conn.Close()
		//key := pbkdf2.Key([]byte("demo pass"), []byte("demo salt"), 1024, 32, sha1.New)
		//block, _ := kcp.NewAESBlockCrypt(key)
		//proxyServer, err := kcp.DialWithOptions(local.remote.String(), block, 10, 3)
		proxyServer, err := kcp.Dial(local.remote.String())
		if err != nil {
			log.Println(err)
			return
		}
		defer proxyServer.Close()
		go func() {
			for {
				written, _ := io.Copy(proxyServer, Conn)
				if written <= 0 {
					Conn.Close()
					proxyServer.Close()
				}
			}
		}()
		for {
			io.Copy(Conn, proxyServer)
		}
	}
	HandleConn(listener, protocol)
	return nil
}
