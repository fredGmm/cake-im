package main

import (
	"cake-im/config"
	"cake-im/pkg/client"
	"cake-im/pkg/server"
	"cake-im/pkg/zapLog"
	"go.uber.org/zap"
	"log"
)

var zLog = zapLog.Logger()
var exit chan int

// cake-im 入口
func main() {
	conf, err := config.Init()
	if err != nil {
		log.Fatal("config fail")
		log.Fatalln(err)
		log.Fatalln(conf)
	}
	grpcClient := client.NewGrpcClient(conf)

	zLog.Info("启动tcp服务，监听地址端口:", zap.String("addr", conf.TCP.Bind[0]))
	if err := server.InitTcp(grpcClient); err != nil {
		log.Print("tcp 启动失败", err)
	}

	zLog.Info("启动websocket服务，监听地址端口:", zap.String("addr", conf.Websocket.Bind[0]))
	//启动websocket
	if err := server.InitWebsocket(grpcClient); err != nil {
		log.Print("ws启动失败")
		log.Print(err)
	}

	<-exit
}
