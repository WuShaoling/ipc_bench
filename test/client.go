package test

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type mmConn map[int][]byte
type mmRoutine map[int]mmConn

func handleResult(res []mmRoutine) {
	// 对于每一个routine

	sum1 := uint64(0)
	sum2 := uint64(0)
	sum3 := uint64(0)

	for routineId, routine := range res {

		routineLen := 0
		routineSum1 := uint64(0)
		routineSum2 := uint64(0)
		routineSum3 := uint64(0)

		// 对于每一个连接
		for _, conn := range routine {
			// 对于每一个消息
			for _, msg := range conn {
				t1 := binary.BigEndian.Uint64(msg[0:8])
				t2 := binary.BigEndian.Uint64(msg[8:16])
				t3 := binary.BigEndian.Uint64(msg[16:24])
				routineSum1 += (t2 - t1) / 1000
				routineSum2 += (t3 - t2) / 1000
				routineSum3 += (t3 - t1) / 1000
			}
			routineLen += len(conn)
		}

		routineAvg1 := routineSum1 / uint64(routineLen)
		routineAvg2 := routineSum2 / uint64(routineLen)
		routineAvg3 := routineSum3 / uint64(routineLen)
		fmt.Printf("%+v,%+v,%+v,%+v\n", routineId, routineAvg1, routineAvg2, routineAvg3)

		sum1 += routineAvg1
		sum2 += routineAvg2
		sum3 += routineAvg3
	}

	ll := uint64(len(res))
	if ll > 1 {
		fmt.Printf("--> %+v,%+v,%+v\n", sum1/ll, sum2/ll, sum3/ll)
	}
}

// 每次连接连续接收 K 个包，ClientSendTime, ServerReceiveTime, ClientReceiveTime
func reader(conn io.Reader, signal chan mmConn, K int) {
	count := 0
	res := make(mmConn)

	buffer := make([]byte, 4096)
	for {

		// read
		_, err := conn.Read(buffer)
		binary.BigEndian.PutUint64(buffer[16:24], uint64(time.Now().UnixNano()))

		// handle error
		if err != nil {
			signal <- nil
			return
		}

		// 赋值
		res[count] = buffer[0:24]

		// 判断是否达到次数
		count++
		if count >= K {
			signal <- res
			return
		}
	}
}

func doAConnection(K int) mmConn {
	// 建立连接
	conn, err := net.Dial("unix", "./go.socket")
	if err != nil {
		log.Fatal("dial error", err)
	}
	defer conn.Close()

	// 开始读进程
	signal := make(chan mmConn)
	go reader(conn, signal, K)

	// 开始写进程，发送K个包
	for i := 0; i < K; i++ {

		buffer := make([]byte, 4096)
		binary.BigEndian.PutUint64(buffer[0:8], uint64(time.Now().UnixNano()))

		if _, err = conn.Write(buffer); err != nil {
			fmt.Println("send error: ", err)
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

	count := 0
	for {
		select {
		case res := <-signal:
			result = append(result, res)

			count++
			if count == M {
				handleResult(result)
				return
			}
		}
	}
}

//var testSetM = []int{1, 		1, 		1, 		10,		10,		10,		100,	100, 	100,	200}
//var testSetN = []int{1, 		10, 	100,	1,		10,		100,	1,		10,		100,	100}
//var testSetK = []int{10000, 	10000, 	10000, 	10000,	10000,	1000,	10000,	10000,	10000,	1000}

func main() {

	M := 1
	N := 1
	K := 10000
	fmt.Printf("%+v,%+v,%+v\n", M, N, K)
	doTest(M, N, K)

	//for i := 0; i < len(testSetM); i++ {
	//
	//	M := testSetM[i]
	//	N := testSetN[i]
	//	K := testSetK[i]
	//
	//	fmt.Printf("--->%+v,%+v,%+v\n", M, N, K)
	//	doTest(M, N, K)
	//}
}
