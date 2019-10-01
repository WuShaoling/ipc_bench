package main

import (
	"encoding/json"
	"github.com/m3l/message"
	"log"
	"net"
	"os"
	"time"
)

func echoServer(c net.Conn) {

	buf := make([]byte, 4096)

	for {
		// read msg
		n, err := c.Read(buf)
		rTime := time.Now().UnixNano()

		// handler error
		if err != nil {
			log.Println("read error: ", err)
			return
		}

		// 转换消息
		msg := &message.Message{}
		if err = json.Unmarshal(buf[0:n], msg); err != nil {
			log.Println("json.Unmarshal error: ", n, err)
			continue
		}
		msg.ServerTime = rTime

		// 发送回去
		data, err := json.Marshal(msg)
		if err != nil {
			log.Println("json.Marshal error: ", err)
			continue
		}
		if n, err = c.Write(data); err != nil {
			log.Println("write error: ", err)
			return
		}
	}
}

func main() {
	log.Println("starting server")

	// 前期清理工作
	if _, e := os.Open("./go.socket"); e != nil {
		os.Remove("./go.socket")
	}

	// 开始监听
	ln, err := net.Listen("unix", "./go.socket")
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
