package test

import (
	"encoding/binary"
	"io"
	"log"
	"net"
	"os"
	"time"
)

func echoServer(c net.Conn) {

	buf := make([]byte, 4096)

	for {
		// read msg
		_, err := c.Read(buf)
		binary.BigEndian.PutUint64(buf[8:16], uint64(time.Now().UnixNano()))

		// handler error
		if err != nil {
			if err != io.EOF {
				log.Println("read error: ", err)
			}
			return
		}

		// 发送回去
		if _, err = c.Write(buf); err != nil {
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
