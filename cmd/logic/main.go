package main

import (
	"fmt"
	"google.golang.org/grpc"
	pb "imChong/api/logic/grpc"
	"github.com/fredGmm/imChong/pkg/connect"
	"net"
)

var addr = "127.0.0.1:9527"

func main() {

	fmt.Println("开始监听 grpc 服务端口,", addr)
	listen, err3 := net.Listen("tcp", addr)
	if err3 != nil {
		fmt.Println("发生错误")
	}
	s := grpc.NewServer()
	pb.RegisterLogicServer(s, &connect.RpcServer{})
	if err := s.Serve(listen); err != nil {
		fmt.Println("发生错误", err)
	}
}
