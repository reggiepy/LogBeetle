package main

import (
	"fmt"
	"io"
	"net"
)

func handleSocks5Client(client net.Conn) {
	defer client.Close()

	// 读取客户端的请求
	buf := make([]byte, 256)
	n, err := client.Read(buf)
	if err != nil {
		fmt.Println("Error reading request:", err)
		return
	}

	// 解析客户端请求
	if buf[0] != 0x05 {
		fmt.Println("Invalid SOCKS5 version")
		return
	}
	nMethods := int(buf[1])
	if len(buf) < nMethods+2 {
		fmt.Println("Invalid SOCKS5 request")
		return
	}
	methods := buf[2 : nMethods+2]
	hasNoAuth := false
	for _, m := range methods {
		if m == 0x00 {
			hasNoAuth = true
			break
		}
	}
	if !hasNoAuth {
		fmt.Println("No supported authentication methods")
		return
	}

	// 响应客户端的请求
	client.Write([]byte{0x05, 0x00})

	// 读取客户端的连接请求
	n, err = client.Read(buf)
	if err != nil {
		fmt.Println("Error reading request:", err)
		return
	}
	if n < 7 {
		fmt.Println("Invalid SOCKS5 request")
		return
	}
	if buf[0] != 0x05 || buf[1] != 0x01 || buf[2] != 0x00 {
		fmt.Println("Invalid SOCKS5 request")
		return
	}
	var dstAddr string
	switch buf[3] {
	case 0x01:
		dstAddr = fmt.Sprintf("%d.%d.%d.%d:%d", buf[4], buf[5], buf[6], buf[7], (uint16(buf[8])<<8)|uint16(buf[9]))
	case 0x03:
		dstAddr = fmt.Sprintf("%s:%d", string(buf[5:n-2]), (uint16(buf[n-2])<<8)|uint16(buf[n-1]))
	default:
		fmt.Println("Unsupported address type")
		return
	}

	// 连接目标服务器
	server, err := net.Dial("tcp", dstAddr)
	if err != nil {
		fmt.Println("Error connecting to destination:", err)
		return
	}
	defer server.Close()

	// 响应客户端连接成功
	client.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})

	// 转发数据
	go io.Copy(server, client)
	io.Copy(client, server)
}

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:1080")
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()

	fmt.Println("SOCKS5 proxy server is running on 127.0.0.1:1080")

	for {
		client, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go handleSocks5Client(client)
	}
}
