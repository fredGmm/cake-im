package server

import (
	"cake-im/config"
	"cake-im/pkg"
	"cake-im/pkg/db"
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	RpcClient "github.com/smallnest/rpcx/client"
	"go.uber.org/zap"
	"log"
	"sync"
)
var MsgChannel   = make(chan *pkg.Message)

type ConnectionPool struct {
	Connections map[string]*Connection
	Join        chan *Connection
	Left        chan string
}

type Connection struct {
	UserId string
	WsAddr string
	Conn *websocket.Conn
	Send chan pkg.Message
	Redis *db.RedisContain
}

type ConnPool struct {
	Connections sync.Map  // userId : *conn
}



func NewConnection(conn *websocket.Conn, userId string, wsAddr string) *Connection {
	return &Connection{
		UserId:userId,
		WsAddr:wsAddr,
		Conn:conn,
		Send:make(chan pkg.Message, 10),
	}
}

func (c *Connection) Read(conf *config.Config, g *Group, pool *ConnectionPool)  {

	//d := RpcClient.NewPeer2PeerDiscovery("tcp@127.0.0.1:9527", "")
	//xclient := RpcClient.NewXClient("Msg", RpcClient.Failtry, RpcClient.RandomSelect,d,RpcClient.DefaultOption)
	for {
		msg := pkg.Message{}
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			log.Print("读取连接消息时，发生错误")
			log.Print(err)
			pool.Left <- c.UserId
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
		//ctx, _ := context.WithTimeout(context.Background(), time.Second)
		//_, pushErr := client.RpcClient.PushMsg(ctx, &chat.PushMsgReq{
		//	MsgId:    "0",
		//	Content:  msg.Text,
		//	UserId:   msg.UserId,
		//	GroupId:  msg.GroupId,
		//	ToUserId: msg.ToUserId,
		//	Type:     msg.Type,
		//	Redis:    &chat.RedisParam{Addr: client.Config.Redis.GetAddr(), Password: client.Config.Redis.Password},
		//})
		//
		//if pushErr != nil {
		//	zLog.Fatal("rpc error :", zap.Error(pushErr))
		//	redisErr := c.Redis.SetMessage(&msg)
		//	if redisErr != nil {
		//		zLog.Error("存入发送消息到redis失败",
		//			zap.String("groupId:", msg.GroupId),
		//			zap.String("userId:", c.UserId),
		//			zap.Error(redisErr),
		//		)
		//	}
		//	log.Print("读取失败")
		//}

		//args := &Req{}
		//reply := &Resp{}
		//err = xclient.Call(context.Background(), "Test",args,reply)
		//if err != nil {
		//	fmt.Print(err)
		//}
		//fmt.Print(reply.Id)
		//xclient.Close()
		//log.Print("读取成功")
	}
}

func (c *Connection) ReadSocket(conf *config.Config, g *Group, pool *ConnectionPool)  {
	redisContain := db.NewRedisContain(conf.Redis.GetAddr(), conf.Redis.Password)
	var xclient RpcClient.XClient
	var ok bool
	for {
		msg := pkg.Message{}
		err := c.Conn.ReadJSON(&msg)
		if err != nil {
			log.Print("读取连接消息时，发生错误")
			log.Print(err)
			pool.Left <- c.UserId
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
		//g.Send <- msg
		//MessageChannel <- &msg
		userIds := redisContain.GetUserMap(msg.GroupId)
		xclients := make(map[string]RpcClient.XClient)

		for uid,wsAddr := range userIds {

			if uid == msg.UserId {
				continue
			}

			if xclient,ok = xclients[wsAddr]; !ok{
				log.Println(addr)
				d := RpcClient.NewPeer2PeerDiscovery("tcp@" + wsAddr, "")
				xclient = RpcClient.NewXClient("Msg", RpcClient.Failtry, RpcClient.RandomSelect,d,RpcClient.DefaultOption)
				//defer xclient.Close()
				xclients[wsAddr] = xclient
			}
			args := &Req{
				UserId:msg.UserId,
				Name:msg.Name,
				GroupId:msg.GroupId,
				MessageId:"no",
				Text:msg.Text,
				Type:1,
				Time:msg.Time,
			}

			reply := &Resp{}
			err := xclient.Call(context.Background(), "Send",args,reply)
			if err != nil {
				fmt.Print(err)
			}
			if reply.Err != nil {
				// 发送失败
			}
			log.Println(reply.Id)
		}
	}
}

func (c *Connection) Write(p *ConnectionPool)  {
	for {
		msg := <-c.Send
		if len(msg.Text) == 0 {
			continue
		}
		msg.OnlineCount = c.Redis.GetOnlineCount()
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
			break
		}


	}
}

func (c *Connection) WriteToSocket(g *Group,p *ConnectionPool)  {
	for {
		msg := <-MsgChannel
		log.Println("得到")
		log.Println(msg)
		if len(msg.Text) == 0 {
			continue
		}
		if msg.UserId == c.UserId {
			continue
		}
		if _,ok := g.Clients.Load(c.UserId);!ok {
			continue
		}


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
			//p.Left <- c.UserId
			//redisErr := c.Redis.SetMessage(&msg)
			//if redisErr != nil {
			//	zLog.Error("存入发送消息到redis失败",
			//		zap.String("userid:", c.UserId),
			//		zap.Error(redisErr),
			//	)
			//}
			break
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

//func (p *ConnectionPool)GetSocketById(userId string) *websocket.Conn {
//
//
//}