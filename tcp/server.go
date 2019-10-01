package main

import (
	"encoding/json"
	"github.com/m3l/message"
	"log"
	"net"
	"time"
)

func echoServer(c net.Conn) {

	buf := make([]byte, 4096)

	for {
		n, err := c.Read(buf)

		// handler error
		if err != nil {
			log.Println("read error: ", err)
			return
		}

		rTime := time.Now().UnixNano()

		// 转换消息
		msg := &message.Message{}
		if err = json.Unmarshal(buf[0:n], msg); err != nil {
			log.Println("json.Unmarshal error: ", n, err)
			continue
		}

		// 修改消息
		msg.ServerTime = rTime

		// 发送回去
		data, err := json.Marshal(msg)
		if err != nil {
			log.Println("json.Marshal error: ", err)
		}
		if n, err = c.Write(data); err != nil {
			log.Println("write error: ", err)
			return
		}
	}
}

func main() {
	log.Println("starting server")

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
