package test

import (
	"log"
	"strconv"
	"time"
)

type SliceMock struct {
	addr uintptr
	len  int
	cap  int
}

type Message struct {
	CSendTime    int64 //客户端发送时间戳
	CReceiveTime int64 //客户端接收时间戳
	ServerTime   int64 //服务端时间戳
	//Content      string //消息内容
}

func main() {
	t := time.Now().UnixNano()
	str := strconv.FormatInt(t, 10)
	log.Println(len(str), str)

	//s := time.Now().UnixNano()
	//fmt.Println(s)
	//header := make([]byte, 64)
	////var header [8]byte
	//binary.BigEndian.PutUint64(header[0:8], uint64(s))
	//fmt.Println(header, len(header))
	//
	//y := binary.BigEndian.Uint64(header)
	//fmt.Println(y)

	//var msg = &message.Message{
	//	CSendTime:    1,
	//	CReceiveTime: 1,
	//	ServerTime:   1,
	//	//Content:      "1234567890",
	//}
	//
	//msgLen := unsafe.Sizeof(*msg)
	//fmt.Println(msgLen)
	//
	//msgBytes := &SliceMock{
	//	addr: uintptr(unsafe.Pointer(msg)),
	//	cap:  int(msgLen),
	//	len:  int(msgLen),
	//}
	//
	//data := *(*[]byte)(unsafe.Pointer(msgBytes))
	//fmt.Println("[]byte is : ", data)
}
