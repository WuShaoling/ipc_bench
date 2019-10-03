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
	Protocol   = "unix"
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
	receive := int64(0)
	total := int64(MessageCount * MessageSize * ConnectionCount)
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
		receive += int64(n)
		if receive >= total {
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
				show2(routineDurations)
				return
			}
		}
	}
}

func show2(routineDurations [][]int64) {
	// 计算总时间
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

	// 计算总数据量。总数据量(KB)=2*线程数*连接数*消息数*消息大小/2014
	totalDataSizeKB := int64(2 * RoutineCount * ConnectionCount * MessageCount * MessageSize / 1024)
	fmt.Printf("%.2f,", float64(totalDataSizeKB*1e6)/float64(1024*totalDuration))
}

func show(routineDurations [][]int64) {
	// 计算每个routine的时间及平均时间
	fmt.Printf("-----------------------------------------------------\n")
	dataSizeKBPerRoutine := int64(2 * ConnectionCount * MessageCount * MessageSize / 1024)
	for k, v := range routineDurations {
		duration := v[1] - v[0]
		fmt.Printf("routine %+v : %8d us, %2d MB/s\n", k, duration, dataSizeKBPerRoutine*1e6/(1024*duration))
	}

	// 计算总时间
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

	// 计算总数据量。总数据量(KB)=2*线程数*连接数*消息数*消息大小/2014
	totalDataSizeKB := int64(2 * RoutineCount * ConnectionCount * MessageCount * MessageSize / 1024)

	fmt.Printf("total: %13d us, %2d MB/s\n", totalDuration, totalDataSizeKB*1e6/(1024*totalDuration))
}

func parseFlag() {
	routineCount := flag.Int("r", 1, "routine count")
	connectionCount := flag.Int("conn", 1, "connection count")
	messageCount := flag.Int("c", 10000, "messageCount count")
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

	//fmt.Printf("-----------------------------------------------------\n")
	//fmt.Printf("RoutineCount ConnectionCount MessageCount MessageSize\n")
	//fmt.Printf("%d %12d %19d %12d\n", RoutineCount, ConnectionCount, MessageCount, MessageSize)

	// fill Message
	for i := 0; i < MessageSize; i++ {
		Message += "0"
	}
}

func main() {
	//fmt.Printf("start %+v test\n", Protocol)

	parseFlag()

	// 以RoutineCount为变量
	for i := 0; i < 20; i++ {
		RoutineCount += 5
		doTest()
	}

	////以MessageCount为变量，以1000为步长，测试 50 组数据
	//MessageCount = 10000
	//for i := 0; i < 100; i++ {
	//	MessageCount += 1000
	//	doTest()
	//}

	fmt.Println()
}
