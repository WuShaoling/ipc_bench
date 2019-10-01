package main

import (
	"encoding/json"
	"github.com/m3l/message"
	"io"
	"log"
	"net"
	"time"
)

const Total = 30 * 10000
const ProcessCount = 150
const SendTimePerProcess = Total / ProcessCount

var FillStr string

type mmConn map[int]*message.Message
type mmRoutine map[int]mmConn

// 每次连接连续接收 K 个包，
func reader(conn io.Reader, signal chan mmConn, K int) {
	count := 0
	res := make(mmConn)
	buf := make([]byte, 4096)

	for {
		// read
		n, err := conn.Read(buf)
		receiveTime := time.Now().UnixNano()

		// handle error
		if err != nil {
			signal <- nil
			return
		}

		// 保存下来，等待分析
		msg := &message.Message{}
		err = json.Unmarshal(buf[0:n], msg)
		if err != nil {
			log.Println("json.Unmarshal error: ", n, err)
			continue
		}
		msg.CReceiveTime = receiveTime
		res[count] = msg

		// 判断是否达到次数
		if count >= K {
			signal <- res
			return
		}
		count++
	}
}

func doAConnection(K int) mmConn {
	// 建立连接
	conn, err := net.Dial("unix", "go.socket")
	if err != nil {
		log.Fatal("dial error", err)
	}
	defer conn.Close()

	// 开始读进程
	signal := make(chan mmConn)
	go reader(conn, signal, K)

	// 开始写进程，发送K个包
	msg := message.Message{Content: FillStr}
	for i := 0; i < K; i++ {

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

	// 接收结果并返回
	select {
	case res := <-signal:
		return res
	}
}

// 每个 routine 执行 N 次连接，每次连接发送 K 个包，串行执行
func doConnections(K, N int, signal chan mmRoutine) {

	mRoutine := make(mmRoutine)

	for connId := 0; connId < N; connId++ {
		mRoutine[connId] = doAConnection(K)
		log.Printf("connection %+v ok", connId)
	}

	signal <- mRoutine
}

// 每次测试执行 M 个 routine，并行执行
func doTest(M, N, K int) {

	result := []mmRoutine{}
	signal := make(chan mmRoutine, M)

	for routineId := 0; routineId < M; routineId++ {
		go func(i int) {
			go doConnections(K, N, signal)
		}(routineId)
	}

	select {
	case res := <-signal:
		log.Println("routine ok")
		result = append(result, res)
	}
}

func main() {
	FillStr = GenStr(4000)
	M := 1
	N := 1
	K := 100
	doTest(M, N, K)
}

func GenStr(n int) string {
	str := ""
	for i := 0; i < n; i++ {
		str += "a"
	}
	return str
}

// 输出微秒时间
//fmt.Printf("%+v,%+v,%+v\n",
//	(msg.ServerTime-msg.CSendTime)/1000,
//	(receiveTime-msg.ServerTime)/1000,
//	(receiveTime-msg.CSendTime)/1000,
//)
