# ipc bench

一件简单的ipc吞吐量性能测试，目前支持以下方式：

- tcp socket
- unix domain socket

## features

- 支持自定义测试线程数
- 支持自定义每个线程连接数
- 支持自定义包数量
- 支持自定义包大小

## getting start

### client 参数解析
- -r, 线程数, default 1
- -conn, 每个线程连接数, default 1
- -c, 发送包的个数, default 1000
- -s, 包大小, default 2048
- -host, server地址, tcp default 127.0.0.1:8888, domain default: ./go.socket

### run unix domain socket

```
cd test3 domain
go run server.go
go run client.go
```

### run tcp/ip socket
```
cd test3 socket
go run server.go
go run client.go
```

### run in docker
docker run -it -v $PWD:/go/src golang:stretch

## 结果说明

命名方式：domain-result-{进程数(起始大小-步长-次数)}-(连接数)-{消息数(起始大小-步长-次数)}-(包大小)-(测试组数)