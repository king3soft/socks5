package main

import (
	"log"
	"net"

	socks5 "github.com/armon/go-socks5"
)

func main() {
	// 创建一个新的SOCKS5服务器实例
	conf := &socks5.Config{}
	server, err := socks5.New(conf)
	if err != nil {
		log.Fatal(err)
	}

	// 在本地监听特定端口
	listener, err := net.Listen("tcp", "0.0.0.0:1080")
	if err != nil {
		log.Fatal(err)
	}

	// 使用循环接受和处理传入的连接
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		// 将连接交给SOCKS5服务器处理
		go server.ServeConn(conn)
	}
}
