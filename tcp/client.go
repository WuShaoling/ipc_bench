package main

import (
	"encoding/json"
	"fmt"
	"github.com/m3l/message"
	"io"
	"log"
	"net"
	"time"
)

const Total = 30 * 10000
const ProcessCount = 50
const SendTimePerProcess = Total / ProcessCount

func reader(conn io.Reader) {
	buf := make([]byte, 4096)
	for {
		// read
		n, err := conn.Read(buf)

		// handle error
		if err != nil {
			return
		}

		// 保存下来，等待分析
		msg := &message.Message{}
		err = json.Unmarshal(buf[0:n], msg)
		if err != nil {
			log.Println("json.Unmarshal error: ", n, err)
			continue
		}

		// 输出微秒时间
		receiveTime := time.Now().UnixNano()
		fmt.Printf("%+v,%+v,%+v\n",
			(msg.ServerTime-msg.CSendTime)/1000,
			(receiveTime-msg.ServerTime)/1000,
			(receiveTime-msg.CSendTime)/1000,
		)
	}
}

func doAConnect(index int, str string) {
	// 建立连接
	conn, err := net.Dial("tcp", "localhost:8888")
	if err != nil {
		log.Fatal("dial error", err)
	}
	defer conn.Close()

	// 开始读进程
	go reader(conn)

	// 开始写进程
	msg := message.Message{
		Content: str,
	}
	for i := 0; i < SendTimePerProcess; i++ {
		msg.CSendTime = time.Now().UnixNano()
		msg.ServerTime = time.Now().UnixNano()

		sendData, err := json.Marshal(msg)
		if err != nil {
			log.Println("json.Marshal error: ", err)
			continue
		}

		if _, err = conn.Write(sendData); err != nil {
			log.Println("send error: ", err)
		}
	}

	select {}
}

func main() {

	str := GenStr(4000)
	for count := 0; count < ProcessCount; count++ {
		go func(c int) {
			go doAConnect(c, str)
		}(count)
	}
	select {}

}

func GenStr(n int) string {
	str := ""
	for i := 0; i < n; i++ {
		str += "a"
	}
	return str
}
