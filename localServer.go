package main

import (
	"crypto/sha1"
	"github.com/xtaci/kcp-go"
	"golang.org/x/crypto/pbkdf2"
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
	defer listener.Close()
	for {
		userConn, err := listener.AcceptTCP()
		if err != nil {
			log.Println(err)
			continue
		}
		userConn.SetLinger(0)
		go func(Conn *net.TCPConn) {
			defer Conn.Close()
			key := pbkdf2.Key([]byte("demo pass"), []byte("demo salt"), 1024, 32, sha1.New)
			block, _ := kcp.NewAESBlockCrypt(key)
			proxyServer, err := kcp.DialWithOptions(local.remote.String(), block, 10, 3)
			if err != nil {
				log.Println(err)
				return
			}
			defer proxyServer.Close()
			go func() {
				for {
					io.Copy(proxyServer, userConn)
				}
				userConn.Close()
				proxyServer.Close()
			}()
			for {
				io.Copy(userConn, proxyServer)
			}
		}(userConn)
	}
}
