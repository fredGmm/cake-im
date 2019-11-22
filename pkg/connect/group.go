package connect

import (
	"go.uber.org/zap"
	"github.com/fredGmm/imChong/pkg/transfer"
	"sync"
)

var (
	Groups       = make(map[string]*Group)
	GroupChannel = make(chan *Group)
)

type Group struct {
	//Connections map[string]*websocket.Conn
	//Clients map[string]*Client
	Clients sync.Map
	Send    chan transfer.Message
	Join    chan *Client
	Left    chan string
}

func NewGroup() *Group {
	var sMap sync.Map
	return &Group{
		//Clients: make(map[string]*Client),
		Clients: sMap,
		Send:    make(chan transfer.Message),
		Join:    make(chan *Client),
		Left:    make(chan string),
	}
}

//群组的消息处理
func (group *Group) Handle() {
	for {
		select {
		case c := <-group.Join:
			log.Info("这个组有新用户加入",
				zap.String("userId", c.UserId),
			)
			//group.Clients[c.UserId] = c
			group.Clients.Store(c.UserId, c)

		case u := <-group.Left:
			log.Info("这个组有新用户离开",
				//zap.String("groupId", group.)
				zap.String("userId", u),
			)
			if _, ok := group.Clients.Load(u); ok {
				group.Clients.Delete(u)
			}
		case m := <-group.Send:
			log.Info("有消息")
			// 给群组每个人广播
			group.Clients.Range(func(_, c interface{}) bool {
				h := c.(*Client)
				if h.UserId != m.UserId {
					select {
					case h.Send <- m:
					default:
						break
					}
				}
				return true
			})
		}
	}
}
