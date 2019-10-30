package connect

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/keepalive"
	logic "github.com/fredGmm/imChong/api/logic/grpc"
	"github.com/fredGmm/imChong/config"
	"github.com/fredGmm/imChong/log"
	"time"
)

var ZLogger = Zlog.Logger()

const (
	minServerHeartbeat = time.Minute * 10
	maxServerHeartbeat = time.Minute * 30
	// grpc options
	grpcInitialWindowSize     = 1 << 24
	grpcInitialConnWindowSize = 1 << 24
	grpcMaxSendMsgSize        = 1 << 24
	grpcMaxCallMsgSize        = 1 << 24
	grpcKeepAliveTime         = time.Second * 10
	grpcKeepAliveTimeout      = time.Second * 3
	grpcBackoffMaxDelay       = time.Second * 3
)

type Server struct {
	c         *config.Config
	serverID  string
	RpcClient logic.LogicClient
	Exit      chan int
}

func newLogicClient() logic.LogicClient {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, "127.0.0.1:9527",
		[]grpc.DialOption{
			grpc.WithInsecure(),
			grpc.WithInitialWindowSize(grpcInitialWindowSize),
			grpc.WithInitialConnWindowSize(grpcInitialConnWindowSize),
			grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(grpcMaxCallMsgSize)),
			grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(grpcMaxSendMsgSize)),
			grpc.WithBackoffMaxDelay(grpcBackoffMaxDelay),
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                grpcKeepAliveTime,
				Timeout:             grpcKeepAliveTimeout,
				PermitWithoutStream: true,
			}),
			grpc.WithBalancerName(roundrobin.Name),
		}...)

	if err != nil {
		fmt.Printf("就是个pg")
		panic(err)
	}

	return logic.NewLogicClient(conn)
}

func NewServer(c *config.Config) *Server {
	s := &Server{
		c:         c,
		RpcClient: newLogicClient(),
		Exit:      make(chan int),
	}
	s.serverID = "127.0.0.1"
	return s

}
