package server

import (
	//"go.uber.org/zap"

	"cake-im/pkg"
	"cake-im/pkg/db"
	"context"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"log"
	"strconv"
	"sync"
	RpcClient "github.com/smallnest/rpcx/client"
)

var (
	Groups       = make(map[string]*Group)
	GroupChannel = make(chan *Group)
)

type MapGroup struct {
	Groups  map[string]*Group
	Channel chan *Group
	//Msg  chan *pkg.Message
	lock sync.RWMutex
}

func (g *MapGroup) Get(groupId string) (group *Group, err error) {
	g.lock.RLock()
	defer g.lock.RUnlock()
	 group,ok := g.Groups[groupId]
	 if ok {
	 	return group,nil
	 }else {
	 	return nil,errors.New("no found group by id")
	 }
}

func (g *MapGroup) set(groupId string, group *Group) {
	g.lock.Lock()
	defer g.lock.Unlock()
	g.Groups[groupId] = group
}

//Connections map[string]*websocket.Conn
//Clients map[string]*Client
type Group struct {
	Id      int
	Clients sync.Map

	Send    chan pkg.Message
	Join    chan *Connection
	Left    chan string
	Redis *db.RedisContain
}

func NewGroup() *Group {
	var sMap sync.Map
	return &Group{
		//Clients: make(map[string]*Client),
		Id:0,
		Clients: sMap,
		Send:    make(chan pkg.Message, 10),
		Join:    make(chan *Connection, 10),
		Left:    make(chan string, 10),
	}
}

//群组的消息处理
func (group *Group) Handle() {
	var xclient RpcClient.XClient
	var ok bool
	for {
		select {
		case c := <-group.Join:
			zLog.Info("这个组有新用户加入",
				zap.String("userId", c.UserId),
			)
			//group.Clients[c.UserId] = c
			group.Clients.Store(c.UserId, c)
			//记录到redis
			err :=group.Redis.SaveUserRpcAddr(strconv.Itoa(group.Id),c.UserId,c.WsAddr)
			if err != nil {
				zLog.Fatal("用户信息映射失败",
					zap.Int("groupId", group.Id),
					zap.String("userId", c.UserId),
				)
			}
		case u := <-group.Left:
			zLog.Info("这个组有新用户离开",
				//zap.String("groupId", group.)
				zap.String("userId", u),
			)
			if _, ok := group.Clients.Load(u); ok {
				group.Clients.Delete(u)
			}
		case m := <-group.Send:
			// 给群组每个人广播(除了自己)
			//group.Clients.Range(func(_, c interface{}) bool {
			//	h := c.(*Connection) // 群组中的连接
			//	if h.UserId != m.UserId {
			//		select {
			//		case h.Send <- m:
			//		default:
			//			break
			//		}
			//	}
			//	return true
			//})

			// 为群组每个人广播，ip:host;groupId;userId
			users := group.Redis.GetUserMap(strconv.Itoa(group.Id))
			xclients := make(map[string]RpcClient.XClient)
			for _,addr := range users {
				if xclient,ok = xclients[addr]; !ok{
					log.Println(addr)
					d := RpcClient.NewPeer2PeerDiscovery("tcp@" + addr, "")
					xclient = RpcClient.NewXClient("Msg", RpcClient.Failtry, RpcClient.RandomSelect,d,RpcClient.DefaultOption)
					//defer xclient.Close()
					xclients[addr] = xclient
				}
				args := &Req{
					UserId:m.UserId,
					Name:m.Name,
					GroupId:m.GroupId,
					MessageId:"no",
					Text:m.Text,
					Type:1,
					Time:m.Time,
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
}

