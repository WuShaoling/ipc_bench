package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

const (
	BufferSize = 1024
	Protocol   = "tcp"
)

var (
	RoutineCount    int    // 多少个线程
	ConnectionCount int    // 每个线程建立多少次连接
	MessageCount    int    // 每次连接发送多少次消息
	MessageSize     int    // 每个消息的大小
	ServerAddress   string // 服务器地址
	Message         = ""   // 消息
)

func reader(conn io.Reader, signal chan bool) {
	buffer := make([]byte, BufferSize)
	sum := 0
	total := MessageCount * MessageSize
	for {
		// read
		n, err := conn.Read(buffer)

		// handle error
		if err != nil {
			log.Println("receive error: ", err)
			signal <- false
			return
		}

		// return
		sum += n
		if sum >= total {
			signal <- true
			return
		}
	}
}

func doAConnection() {
	// 建立连接
	conn, err := net.Dial(Protocol, ServerAddress)
	if err != nil {
		log.Fatal("dial error", err)
	}
	defer conn.Close()

	// 开始读进程
	signal := make(chan bool, 1)
	go reader(conn, signal)

	// 开始写进程，发送MessageCount个包
	for i := 0; i < MessageCount; i++ {
		if _, err := conn.Write([]byte(Message)); err != nil {
			fmt.Println("send error: ", err)
		}
	}

	// 接收结果并返回
	select {
	case <-signal:
		return
	}
}

// 每个 routine 执行 ConnectionCount 次连接，每次连接发送 MessageCount 个包，串行执行
func doConnections(signal chan []int64) {
	var res []int64
	res = append(res, time.Now().UnixNano()/1000)
	for connId := 0; connId < ConnectionCount; connId++ {
		doAConnection()
	}
	res = append(res, time.Now().UnixNano()/1000)
	signal <- res
}

func doTest() {

	signal := make(chan []int64, RoutineCount)
	for routineId := 0; routineId < RoutineCount; routineId++ {
		go func(i int) {
			go doConnections(signal)
		}(routineId)
	}

	// 结果
	var routineDurations [][]int64
	count := 0
	for {
		select {
		case d := <-signal:

			routineDurations = append(routineDurations, d)

			// 边界条件
			count++
			if count == RoutineCount {
				show(routineDurations)
				return
			}
		}
	}
}

func show(routineDurations [][]int64) {
	// 计算总时间
	fmt.Printf("-----------------------------------------------------\n")
	totalStartTime := time.Now().UnixNano() / 1000
	totalEndTime := int64(0)
	for _, routine := range routineDurations {
		for _, t := range routine {
			if t < totalStartTime {
				totalStartTime = t
			} else if t > totalEndTime {
				totalEndTime = t
			}
		}
	}
	totalDuration := totalEndTime - totalStartTime
	fmt.Printf("total duration: %15dus\n", totalDuration)

	// 计算每个routine的时间及平均时间
	fmt.Printf("-----------------------------------\n")
	sum := int64(0)
	for k, v := range routineDurations {
		duration := v[1] - v[0]
		fmt.Printf("routine %+v duration: %11dus\n", k, duration)
		sum += duration
	}
	fmt.Printf("routine avg duration: %9dus\n", sum/int64(RoutineCount))
	fmt.Printf("-----------------------------------\n")

	// 计算吞吐率
	// 总数据量=线程数*连接数*消息数*消息大小/时间
	totalDataSize := int64(RoutineCount * ConnectionCount * MessageCount * MessageSize)
	fmt.Printf("throughput: %16d MB/s\n", totalDataSize/totalDuration)
}

func parseFlag() {
	routineCount := flag.Int("r", 10, "routine count")
	connectionCount := flag.Int("conn", 10, "connection count")
	messageCount := flag.Int("c", 1000, "messageCount count")
	messageSize := flag.Int("s", 2048, "message size")
	var serverAddress *string
	if Protocol == "tcp" {
		serverAddress = flag.String("host", "127.0.0.1:8888", "server Address")
	} else {
		serverAddress = flag.String("host", "./go.socket", "server Address")
	}

	flag.Parse()

	RoutineCount = *routineCount
	ConnectionCount = *connectionCount
	MessageCount = *messageCount
	MessageSize = *messageSize
	ServerAddress = *serverAddress

	fmt.Printf("-----------------------------------------------------\n")
	fmt.Printf("RoutineCount ConnectionCount MessageCount MessageSize\n")
	fmt.Printf("%d %12d %19d %12d\n", RoutineCount, ConnectionCount, MessageCount, MessageSize)

	// fill Message
	for i := 0; i < MessageSize; i++ {
		Message += "0"
	}
}

func main() {
	fmt.Printf("start %+v test\n", Protocol)
	parseFlag()
	doTest()
}
