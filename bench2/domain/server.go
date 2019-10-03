package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

const (
	ReceiveBufferSize = 1024
	ServeAt           = "./go.socket"
	ServeProtocol     = "unix"
)

var (
	ConnectionCount int // 每个线程建立多少次连接
	MessageCount    int // 每次连接发送多少次消息
	MessageSize     int // 每个消息的大小
)

func echoServer(c net.Conn) {
	total := int64(ConnectionCount * MessageCount * MessageSize)
	receive := int64(0)

	buffer := make([]byte, ReceiveBufferSize)
	for {
		n, err := c.Read(buffer)
		receive += int64(n)

		// handler error
		if err != nil {
			if err != io.EOF {
				log.Println("read error: ", err)
			}
			return
		}

		if receive >= total {
			// 发送回去截止时间
			if _, err = c.Write([]byte(strconv.FormatInt(time.Now().UnixNano()/1000, 10))); err != nil {
				log.Println("write error: ", err)
			}
			return
		}
	}
}

func doServe() {
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

func parseServerFlag() {
	connectionCount := flag.Int("conn", 1, "connection count")
	messageCount := flag.Int("c", 100000, "messageCount count")
	messageSize := flag.Int("s", 2048, "message size")

	flag.Parse()

	ConnectionCount = *connectionCount
	MessageCount = *messageCount
	MessageSize = *messageSize

	fmt.Printf("-----------------------------------------------------\n")
	fmt.Printf("ConnectionCount MessageCount MessageSize\n")
	fmt.Printf("%10d %10d %10d\n", ConnectionCount, MessageCount, MessageSize)
}

func main() {
	parseServerFlag()
	doServe()
}
