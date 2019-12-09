package server

import (
	"cake-im/pkg"
	"cake-im/pkg/client"
	"encoding/json"
	"go.uber.org/zap"
	"log"
	"net"
)

var (
	//zLog = zapLog.Logger()
	//err      error
	listener *net.TCPListener
	addr     *net.TCPAddr
)

func InitTcp(gClient *client.GrpcClient) (err error) {

	tcpTind := gClient.Config.TCP.Bind[0]
	if addr, err = net.ResolveTCPAddr("tcp", tcpTind); err != nil {
		zLog.Error("net.ResolveTCPAddr(tcp, bind), error ", zap.String("bind", tcpTind), zap.Error(err))
	}

	if listener, err = net.ListenTCP("tcp", addr); err != nil {
		zLog.Error("net.ListenTCP(tcp, %s) error(%v)", zap.String("addr", addr.String()), zap.Error(err))
	}

	go func() {
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				zLog.Error("TCP accept fail %s,maybe connnect has quit!", zap.Error(err))
			}
			go handleConn(gClient, conn)
		}
	}()
	return
}

func handleConn(grpcClient *client.GrpcClient, c *net.TCPConn) {

	for {
		var buf = make([]byte, 4096)
		n, err := c.Read(buf)
		if err != nil {
			zLog.Error("tcp conn read error , => %v", zap.Error(err))
			return
		}

		msg := pkg.Message{}
		log.Printf("read %d bytes, content is %s \n", n, string(buf))
		zLog.Info("new msg", zap.String("message", string(buf)))
		err = json.Unmarshal(buf[:n], msg)
		if err != nil {
			continue
		}

		//ctx, _ := context.WithTimeout(context.Background(), time.Second)
		//r, err := grpcClient.RpcClient.PushMsg(ctx, &chat.PushMsgReq{
		//	MsgId: "0",
		//	Content:msg.Text,
		//	UserId:msg.UserId,
		//	GroupId:msg.GroupId,
		//	ToUserId:msg.ToUserId,
		//	Type:msg.Type,
		//})
		//
		//if err != nil {
		//	zLog.Error("tcp push msg", zap.Error(err))
		//}
		//
		//zLog.Info("rpc return ", zap.String("r", r.String()))
	}
}
