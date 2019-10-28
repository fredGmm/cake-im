package connect

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	logic "imChong/api/logic/grpc"
	"imChong/module/transfer"
	"net"
	"time"
)

func InitTCP(s *Server, bind string) (err error) {
	var (
		listener *net.TCPListener
		addr     *net.TCPAddr
	)

	if addr, err = net.ResolveTCPAddr("tcp", bind); err != nil {
		ZLogger.Error("net.ResolveTCPAddr(tcp, bind), error ", zap.String("bind", bind), zap.Error(err))
		return
	}
	if listener, err = net.ListenTCP("tcp", addr); err != nil {
		ZLogger.Error("net.ListenTCP(tcp, %s) error(%v)", zap.String("addr", addr.String()), zap.Error(err))
		return
	}
	ZLogger.Info("tcp listen success, ", zap.String("addr", addr.String()))
	go func() {
		for {
			conn, err2 := listener.AcceptTCP()
			if err2 != nil {
				fmt.Printf("TCP accept fail %s,maybe connnect has quit!", err2)
				break
			}
			go handleConn(s, conn)
		}
	}()
	return
}

func handleConn(s *Server, c *net.TCPConn) {
	ZLogger.Info("start to read from conn")
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second) // ?
	//defer cancel()
	for {
		var buf = make([]byte, 4096)
		n, err := c.Read(buf)
		if err != nil {
			ZLogger.Error("tcp conn read error:", zap.Error(err))
			return
		}
		m := &transfer.Message{}
		ZLogger.Info("new msg", zap.String("message", string(buf)))
		err2 := json.Unmarshal(buf[:n], m)
		if err2 != nil {
			ZLogger.Error("format error, m:%v", zap.Error(err2))
			continue
		}
		ctx, _ := context.WithTimeout(context.Background(), time.Second)
		r, err3 := s.RpcClient.PushMsg(ctx, &logic.PushMsgReq{
			MsgId:    "0",
			Content:  m.Message,
			UserId:   m.UserId,
			GroupId:  m.GroupId,
			ToUserId: m.ToUserId,
			Type:     m.Type,
		})
		if err3 != nil {
			ZLogger.Error("format error", zap.Error(err3))
		}
		ZLogger.Info("rpc return ", zap.String("r", r.String()))
	}
}

//func acceptTCP(server *Server, lis *net.TCPListener) {
//	var (
//		conn *net.TCPConn
//		err  error
//		//r    int
//	)
//
//	for {
//
//	}
//	if conn, err = lis.AcceptTCP(); err != nil {
//		log.Printf("listener.Accept(\"%s\") error(%v)", lis.Addr().String(), err)
//		return
//	}
//
//	if err = conn.SetKeepAlive(true); err != nil {
//		log.Printf("conn.setKeepAlive() error(%v)")
//		return
//	}
//
//	if err = conn.SetReadBuffer(4096); err != nil {
//		log.Printf("conn.SetReadBuffer() error(%v)", err)
//		return
//	}
//
//	if err = conn.SetWriteBuffer(4096); err != nil {
//		log.Printf("conn.SetWriteBuffer() error(%v)", err)
//		return
//	}
//
//	go serveTCP(server, conn)
//
//	//go serveTCP(server, conn, r)
//	//if r++; r == maxInt {
//	//	r = 0
//	//}
//
//}
//
//func serveTCP(s *Server, conn *net.TCPConn) {
//	var buf = make([]byte, 4096)
//
//	for {
//
//		n, err := conn.Read(buf)
//		if err != nil {
//			log.Println(n)
//			return
//		}
//		//k, err2 := conn.Write([]byte(`收到了`))
//		log.Println(string(buf))
//	}
//
//	//s.ServeTCP(conn, rp, wp, tr)
//}
