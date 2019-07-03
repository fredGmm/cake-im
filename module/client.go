package module

import (
	"github.com/gorilla/websocket"
	"log"
	"time"
)

//var connectionPool map[string]*websocket.Conn

type ConnectionPool struct {
	Connections map[string]*Client
	Join chan *Client
	Left chan string
}

type Client struct {
	UserId string
	Connection *websocket.Conn
	Send chan Message
}

func NewClient(conn *websocket.Conn, userId string) *Client  {
	return &Client{
		UserId:userId,
		Connection:conn,
		Send:make(chan Message, 10), // 发送通道
	}
}
//
func (c *Client) ReadMessage(g *Group, to *Client) {
	preMessageTime := int64(0)
	for {
		msg := Message{}
		err := c.Connection.ReadJSON(&msg)
		if err != nil {
			log.Print("这个ws读取消息时错误，将连接移除")
			log.Print(err)
		}
		curMessageTime := time.Now().Unix()
		if curMessageTime - preMessageTime < 1 {
			continue
		}
		log.Println(msg)
		if g == nil {
			to.Send <- msg
		}else {
			// 读取到了客户端消息，写进对应群的消息
			g.Send <- msg
		}

	}
}

func (c *Client)ReadMessageToUser(p *ConnectionPool)  {
	preMessageTime := int64(0)
	for {
		msg := Message{}
		err := c.Connection.ReadJSON(&msg)
		if err != nil {
			log.Print("这个ws读取消息时错误，将连接移除")
			log.Print(err)
			p.Left <- c.UserId
			break
		}
		curMessageTime := time.Now().Unix()
		if curMessageTime - preMessageTime < 1 {
			continue
		}
		log.Println(msg)
		toUserId := msg.ToUserId
		if toClient, ok := p.Connections[toUserId]; ok {
			toClient.Send <- msg
		}else {
			// 离线消息
			log.Print("离线消息读取")
			if err != nil {
				log.Println(err)
			}
			log.Println(msg)
		}


	}
}

func (c *Client) WriteMessage(p *ConnectionPool) {
	for {
		msg := <- c.Send
		log.Print("读取")
		log.Println(msg)
		// 发送消息给 链接
		if err := c.Connection.WriteJSON(msg); err != nil {
			log.Print("关闭链接")
			log.Print(err)
			err := c.Connection.Close()
			if err != nil {
				log.Print("写消息时关闭链接失败")
				log.Print(err)
			}
			p.Left <- c.UserId
		}
	}
}

func (p *ConnectionPool) Handle() {
	for {
		select {
		case c := <-p.Join:
			log.Println("有人链接")
			p.Connections[c.UserId] = c
		case uid := <-p.Left:
			log.Print("有新用户离开")
			if _, ok := p.Connections[uid]; ok {
				delete(p.Connections, uid)
			}
		}
	}

}