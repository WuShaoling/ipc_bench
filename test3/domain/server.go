package main

import (
	"io"
	"log"
	"net"
	"os"
)

var (
	ReceiveBufferSize = 1024
	ServeAt           = "./go.socket"
	ServeProtocol     = "unix"
)

func echoServer(c net.Conn) {
	buffer := make([]byte, ReceiveBufferSize)
	for {
		_, err := c.Read(buffer)

		// handler error
		if err != nil {
			if err != io.EOF {
				log.Println("read error: ", err)
			}
			return
		}

		// 发送回去
		if _, err = c.Write(buffer); err != nil {
			log.Println("write error: ", err)
			return
		}
	}
}

func main() {
	log.Printf("start %+v server at %+v", ServeProtocol, ServeAt)

	// 前期清理工作
	if _, e := os.Open(ServeAt); e != nil {
		os.Remove(ServeAt)
	}

	// 开始监听
	ln, err := net.Listen(ServeProtocol, ServeAt)
	if err != nil {
		log.Fatal("start server error: ", err)
	}

	// 多线程循环监听请求
	for {
		fd, err := ln.Accept()
		if err != nil {
			log.Println("accept error: ", err)
			continue
		}
		go echoServer(fd)
	}

}
