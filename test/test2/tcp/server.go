package main

import (
	"io"
	"log"
	"net"
)

var ReceiveBufferSize = 1024

func echoServer(c net.Conn) {

	buffer := make([]byte, ReceiveBufferSize)
	//w := bufio.NewWriter(c)
	for {
		// read msg
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
		//if err := w.Flush(); err != nil {
		//	log.Println("flush error: ", err)
		//	return
		//}
	}
}

func main() {
	log.Println("starting tcp server")

	// 开始监听
	ln, err := net.Listen("tcp", ":8888")
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
