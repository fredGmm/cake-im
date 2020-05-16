package main

import (
	"cake-im/config"
	connServer "cake-im/pkg/server"
	"cake-im/pkg/zapLog"
	rpcServer "github.com/smallnest/rpcx/server"
	"go.uber.org/zap"
	"log"
	_ "net/http/pprof"
)

var zLog = zapLog.Logger()
var exit chan struct{}

// cake-im 入口
func main() {
	conf, err := config.Init()
	if err != nil {
		zLog.Fatal("config parse fail",zap.Error(err))
	}
	//grpcClient := client.NewGrpcClient(conf)
	initRpc(conf)

	// tcp通信
	//if err := server.InitTcp(grpcClient); err != nil {
	//	log.Print("tcp 启动失败", err)
	//}

	//websocket 通信
	zLog.Info("启动websocket服务，监听地址端口:", zap.String("addr", conf.Websocket.Bind[0]))
	if err := connServer.StartWebsocket(conf); err != nil {
		zLog.Fatal("websocket init error",zap.Error(err))
	}

	// 退出
	<-exit
}

func initRpc(conf *config.Config)  {
	s := rpcServer.NewServer()
	err := s.Register(new(connServer.Msg), "")
	if err != nil {
		log.Printf("rpc 注册出错")
	}
	go s.Serve("tcp", conf.Xrpc.Addr)
}
