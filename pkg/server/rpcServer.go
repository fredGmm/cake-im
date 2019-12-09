package server

import (
	"context"
	"fmt"
	//"go.uber.org/zap"
	pb "cake-im/rpc/chat/grpc"
	"cake-im/pkg/db"
	"cake-im/pkg"
)

type RpcServer struct{}

//func (rs *RpcServer) Test(ctx context.Context, re *grpc.HelloRequest) (*grpc.HelloReply, error) {
//
//	//log.Printf("ping le")
//
//	return  &grpc.HelloReply{Message:"AS"}, nil
//}

func (s *RpcServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	//ZLogger.Info("rpc SayHello() Received", zap.String("in", in.String()))
	return &pb.HelloReply{Message: "Hello0000AAAA000000AAA0A0 " + in.Name}, nil
}

func (s *RpcServer) GetTotal(ctx context.Context, in *pb.RoomRequest) (*pb.RoomReply, error) {
	//ZLogger.Info("rpc GetTotal() Received", zap.String("in", in.String()))
	return &pb.RoomReply{Count: "12354343"}, nil
}

func (s *RpcServer) Ping(ctx context.Context, req *pb.PingReq) (*pb.PingReply, error) {
	//ZLogger.Info("ping() req")
	return &pb.PingReply{}, nil
}

func (s *RpcServer) PushMsg(ctx context.Context, req *pb.PushMsgReq) (*pb.PushMsgReply, error) {

	//ZLogger.Info("rpc PushMsg() req", zap.String("req", req.String()))
	//存入队列并做相当的逻辑
	msg := &pkg.Message{}
	msg.Text = req.Content
	msg.Type = req.Type
	msg.UserId = req.UserId
	msg.ToUserId = req.ToUserId
	msg.GroupId = req.GroupId
	redisParam := req.Redis
	redisContain := db.NewRedisContain(redisParam.Addr, redisParam.Password)

	err := redisContain.SetMessage(msg)
	if err != nil {
		fmt.Print("redis error", err)
	}
	return &pb.PushMsgReply{MsgId: "", ServerId: ""}, nil
}

//
//func (rs *RpcServer)Close(ctx context.Context, req *pb.CloseReq) (*pb.CloseReply, error){
//	return &pb.CloseReply{}, nil
//
//}
//
//func (rs *RpcServer)Connect(ctx context.Context, req *pb.ConnectReply) (*pb.ConnectReply, error) {
//	return &pb.ConnectReply{}, nil
//}
//
//func (rs *RpcServer)DisConnect(ctx context.Context, req *pb.DisconnectReq) (*pb.DisconnectReply, error) {
//	return &pb.DisconnectReply{}, nil
//}
//
//func (rs *RpcServer)Heartbeat(ctx context.Context, req pb.HeartbeatReq) (*pb.HeartbeatReply, error) {
//	return  &pb.HeartbeatReply{}, nil
//}
//func (rs *RpcServer)RenewOnline(ctx context.Context, req pb.OnlineReq) (*pb.OnlineReply, error) {
//	return &pb.OnlineReply{}, nil
//}
//func (rs *RpcServer)Receive(ctx context.Context, req pb.ReceiveReq) (*pb.ReceiveReply, error) {
//	return &pb.ReceiveReply{}, nil
//}
//func (rs *RpcServer)Nodes(ctx context.Context, req pb.NodesReq) (*pb.NodesReply, error) {
//	return &pb.NodesReply{}, nil
//}
