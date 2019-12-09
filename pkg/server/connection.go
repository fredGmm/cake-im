package server

import (
	"cake-im/pkg"
	"cake-im/pkg/client"
	"cake-im/pkg/db"
	chat "cake-im/rpc/chat/grpc"
	"context"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"log"
	"time"
)

type ConnectionPool struct {
	Connections map[string]*Connection
	Join        chan *Connection
	Left        chan string
}

type Connection struct {
	UserId string
	Conn *websocket.Conn
	Send chan pkg.Message
	Redis *db.RedisContain
}

func NewConnection(conn *websocket.Conn, userId string) *Connection {
	return &Connection{
		UserId:userId,
		Conn:conn,
		Send:make(chan pkg.Message, 10),
	}
}

func (c *Connection) Read(client *client.GrpcClient, g *Group)  {
	for {
		msg := pkg.Message{}
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			log.Print("读取连接消息时，发生错误")
			log.Print(err)
			break
		}
		msg.UserId = c.UserId
		zLog.Info("客户端新消息。",
			zap.String("user_id", msg.UserId),
			zap.String("name", msg.Name),
			zap.String("message", msg.Text),
			zap.String("GroupId", msg.GroupId),
			zap.String("time", msg.Time),

		)
		g.Send <- msg
		ctx, _ := context.WithTimeout(context.Background(), time.Second)
		_, pushErr := client.RpcClient.PushMsg(ctx, &chat.PushMsgReq{
			MsgId:    "0",
			Content:  msg.Text,
			UserId:   msg.UserId,
			GroupId:  msg.GroupId,
			ToUserId: msg.ToUserId,
			Type:     msg.Type,
			Redis:    &chat.RedisParam{Addr: client.Config.Redis.GetAddr(), Password: client.Config.Redis.Password},
		})

		if pushErr != nil {
			zLog.Fatal("rpc error :", zap.Error(pushErr))
			redisErr := c.Redis.SetMessage(&msg)
			if redisErr != nil {
				zLog.Error("存入发送消息到redis失败",
					zap.String("groupId:", msg.GroupId),
					zap.String("userId:", c.UserId),
					zap.Error(redisErr),
				)
			}
			log.Print("读取失败")
		}

		//log.Print("读取成功")
	}
}

func (c *Connection) Write(p *ConnectionPool)  {
	for {
		msg := <-c.Send
		if len(msg.Text) == 0 {
			continue
		}
		msg.OnlineCount = c.Redis.GetOnlineCount()
		log.Print("写入消息")
		// 发送消息给 链接
		if err := c.Conn.WriteJSON(msg); err != nil {
			zLog.Warn("写入消息失败",
				zap.String("userId", c.UserId),
				zap.Error(err),
			)
			err := c.Conn.Close()
			if err != nil {
				zLog.Warn("关闭链接失败",
					zap.String("userId", c.UserId),
					zap.Error(err),
				)
			}
			p.Left <- c.UserId
			redisErr := c.Redis.SetMessage(&msg)
			if redisErr != nil {
				zLog.Error("存入发送消息到redis失败",
					zap.String("userid:", c.UserId),
					zap.Error(redisErr),
				)
			}
		}

	}
}

func (p *ConnectionPool) Handle() {
	for {
		select {
		case c := <-p.Join:
			zLog.Warn("有人加入", zap.String("userId", c.UserId))
			p.Connections[c.UserId] = c
		case uid := <-p.Left:
			zLog.Warn("有人转身离开", zap.String("userId", uid))
			if _, ok := p.Connections[uid]; ok {
				delete(p.Connections, uid)
			}
		}
	}

}