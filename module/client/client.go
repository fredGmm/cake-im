package client

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"imChong/log"
	"imChong/module"
	"imChong/module/db"
	"time"
)
var (
	//redis = db.NewRedisContain("192.168.0.7:6379")
	log = Zlog.Logger()
)

type ConnectionPool struct {
	Connections map[string]*Client
	Join chan *Client
	Left chan string
}

type Client struct {
	UserId string
	Connection *websocket.Conn
	Send chan module.Message
	Redis *db.RedisContain
}

func NewClient(conn *websocket.Conn, userId string) *Client  {
	return &Client{
		UserId:userId,
		Connection:conn,
		Send:make(chan module.Message, 10), // 发送通道
	}
}
//
func (c *Client) ReadMessage(g *Group) {
	preMessageTime := int64(0)
	for {
		msg := module.Message{}
		err := c.Connection.ReadJSON(&msg)
		if err != nil {
			log.Warn("读取消息时发生错误，将移除！", zap.Error(err))
			g.Left <- c.UserId
			break
		}
		curMessageTime := time.Now().Unix()
		if curMessageTime - preMessageTime < 1 {
			continue
		}
		msg.UserId = c.UserId
		log.Info("发送的消息",
			zap.String("message", msg.Message),
			zap.String("to_user_id", msg.ToUserId),
			zap.String("user_id", msg.UserId),
			zap.String("group_id", msg.GroupId),
			//zap.Bool("g", g == nil),
			)
		//if (msg == Message{}) {
			g.Send <- msg
		//}
		redisErr := c.Redis.SetMessage(msg)
		if redisErr != nil {
			log.Error("存入发送消息到redis失败",
				zap.String("groupId:", msg.GroupId),
				zap.String("userId:", c.UserId),
				zap.Error(redisErr),
			)
		}

	}
}

func (c *Client)ReadMessageToUser(p *ConnectionPool)  {
	preMessageTime := int64(0)
	for {
		msg := module.Message{}
		err := c.Connection.ReadJSON(&msg)
		if err != nil {
			log.Warn("读取消息时错误，将移除！",
				zap.String("userId", c.UserId),
				zap.Error(err),
				)
			p.Left <- c.UserId
			break
		}
		msg.UserId = c.UserId
		curMessageTime := time.Now().Unix()
		if curMessageTime - preMessageTime < 1 {
			continue
		}
		log.Info("发送的消息",
			zap.String("message", msg.Message),
			zap.String("to_user_id", msg.ToUserId),
			zap.String("user_id", msg.UserId),
			zap.String("group_id", msg.GroupId),
		)

		toUserId := msg.ToUserId
		if toClient, ok := p.Connections[toUserId]; ok {
			toClient.Send <- msg
		}else {
			// 离线消息
			log.Info("离线消息")
		}


	}
}

func (c *Client) WriteMessage(p *ConnectionPool) {
	for {
		msg := <- c.Send
		if (msg == module.Message{}) {
			continue
		}

		log.Info("写入消息",
			zap.String("message", msg.Message),
			zap.String("to_user_id", msg.ToUserId),
			zap.String("user_id", msg.UserId),
			zap.String("group_id", msg.GroupId),
		)
		// 发送消息给 链接
		if err := c.Connection.WriteJSON(msg); err != nil {
			log.Warn("写入消息失败",
				zap.String("userId", c.UserId),
				zap.Error(err),
				)
			err := c.Connection.Close()
			if err != nil {
				log.Warn("关闭链接失败",
					zap.String("userId", c.UserId),
					zap.Error(err),
				)
			}
			p.Left <- c.UserId
			redisErr := c.Redis.SetMessage(msg)
			if redisErr != nil {
				log.Error("存入发送消息到redis失败",
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
			p.Connections[c.UserId] = c
		case uid := <-p.Left:
			log.Warn("有人转身离开", zap.String("userId", uid))
			if _, ok := p.Connections[uid]; ok {
				delete(p.Connections, uid)
			}
		}
	}

}