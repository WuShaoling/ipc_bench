package message

type Message struct {
	CSendTime    int64 //客户端发送时间戳
	CReceiveTime int64 //客户端接收时间戳
	ServerTime   int64 //服务端时间戳
	Content      string
}
