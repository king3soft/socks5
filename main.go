package main

import (
	"io"
	"log"
	"net"
	"net/http"

	socks5 "github.com/armon/go-socks5"
)

func http_proxy() {
	http.HandleFunc("/", proxyHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	client := &http.Client{}
	req, err := http.NewRequest(r.Method, r.URL.String(), r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.Header = r.Header

	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	for name, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(name, value)
		}
	}

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func socks5_proxy() {
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


func main() {
	http_proxy()
}