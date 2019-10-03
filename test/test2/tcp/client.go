package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

var mm = [][]int64{}

const BufferSize = 1024

var MsgSize *int
var TcpHost *string
var Msg = ""

func handleResult() {

	// 对于每一个routine
	sum := int64(0)
	for routineId, routine := range mm {

		// 累加每一个线程的所有连接
		routineSum := int64(0)
		for _, conn := range routine {
			routineSum += conn
		}

		// 求平均
		routineAvg := routineSum / int64(len(routine))
		fmt.Printf("routine: %+v,%+v\n", routineId, routineAvg)

		sum += routineAvg
	}

	fmt.Printf("total:  %+v\n", sum/int64(len(mm)))
}

func reader(conn io.Reader, signal chan int64) {
	buffer := make([]byte, BufferSize)
	sum := 0
	for {
		// read
		n, err := conn.Read(buffer)

		// handle error
		if err != nil {
			log.Println("receive error: ", err)
			signal <- -1
			return
		}

		// return
		sum += n
		if sum >= *MsgSize {
			//log.Println("a connection ok")
			signal <- time.Now().UnixNano()
			return
		}
	}
}

func doAConnection() int64 {
	// 建立连接
	conn, err := net.Dial("tcp", *TcpHost)
	if err != nil {
		log.Fatal("dial error", err)
	}
	defer conn.Close()

	// 开始读进程
	signal := make(chan int64, 10)
	go reader(conn, signal)

	// 开始写进程，发送1个包
	beginTime := time.Now().UnixNano()
	if _, err := conn.Write([]byte(Msg)); err != nil {
		fmt.Println("send error: ", err)
	}

	// 接收结果并返回
	for {
		select {
		case endTime := <-signal:
			//log.Println("receive: ", endTime)
			return (endTime - beginTime) / 1000
		}
	}
}

// 每个 routine 执行 N 次连接，每次连接发送 K 个包，串行执行
func doConnections(N int, signal chan []int64) {

	res := []int64{}

	for connId := 0; connId < N; connId++ {
		res = append(res, doAConnection())
	}

	signal <- res
}

// 每次测试执行 M 个 routine，并行执行
func doTest(M, N int) {

	signal := make(chan []int64, M)

	for routineId := 0; routineId < M; routineId++ {
		go func(i int) {
			go doConnections(N, signal)
		}(routineId)
	}

	count := 0
	for {
		select {
		case res := <-signal:
			mm = append(mm, res)

			count++
			if count == M {
				handleResult()
				return
			}
		}
	}
}

//var testSetM = []int{1, 		1, 		1, 		10,		10,		10,		100,	100, 	100,	200}
//var testSetN = []int{1, 		10, 	100,	1,		10,		100,	1,		10,		100,	100}
//var testSetK = []int{10000, 	10000, 	10000, 	10000,	10000,	1000,	10000,	10000,	10000,	1000}

func main() {
	routineCount := flag.Int("r", 1, "routine counts")
	messageCount := flag.Int("c", 1000, "connection counts")
	MsgSize = flag.Int("s", 2048, "message size")
	TcpHost = flag.String("host", "127.0.0.1:8888", "N")
	flag.Parse()
	Msg = GenStr(*MsgSize)
	fmt.Printf("%+v,%+v,%+v,%+v\n", *routineCount, *messageCount, *MsgSize, *TcpHost)
	doTest(*routineCount, *messageCount)
}

func GenStr(n int) string {
	str := ""
	for i := 0; i < n; i++ {
		str += "0"
	}
	return str
}
