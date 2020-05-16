package server

import (
	"cake-im/config"
	"cake-im/pkg"
	"cake-im/pkg/db"
	"cake-im/pkg/util"
	"encoding/json"
	"fmt"
	"strconv"

	"cake-im/pkg/zapLog"
	"context"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"log"
	"time"

	//"go.uber.org/zap"
	"net/http"
)

var (
	upgrader = websocket.Upgrader{
		HandshakeTimeout:time.Second * 30,
		ReadBufferSize:1024,
		WriteBufferSize:1024,
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	Pool = &ConnectionPool{
		Connections:make(map[string]*Connection),
		Join:make(chan *Connection),
		Left:make(chan string),
	}

	GlobalGroup = &MapGroup{
		Groups:make(map[string]*Group),
	}

	MessageChannel = make(chan *pkg.Message)

	Conns = &ConnPool {

	}
	zLog = zapLog.Logger()
)
//router 简单的路由
func router(conf *config.Config) (router *mux.Router) {
	r := mux.NewRouter()
	//r.HandleFunc("/group/{groupId}/{userId}", UserConn)
	//r.HandleFunc("/user/{userId}/{toUserId}", UserConn)
	//r.HandleFunc("/distribute", Distribute)
	//r.HandleFunc("/group/{groupId}/{userId}", func(w http.ResponseWriter, r *http.Request) {
	//	GroupConn(w, r, conf)
	//})

	r.HandleFunc("/group/{groupId}/{userId}", func(w http.ResponseWriter, r *http.Request) {
		HandleConn(w, r, conf)
	})

	r.HandleFunc("/distribute", func(w http.ResponseWriter, r *http.Request) {
		Distribute(w,r,conf)
	})
	http.Handle("/", r)
	return r
}

func InitWebsocket_old(conf *config.Config) error{
	_,cancel := context.WithCancel(context.Background()) //ctx
	defer cancel()
	r := router(conf)

	go GroupsHandle()
	go Pool.Handle()

	serve := &http.Server{
		Addr:conf.Websocket.Bind[0],
		Handler:r,
		ReadTimeout: time.Second * 3,
		WriteTimeout:time.Second *3,
	}
	err := serve.ListenAndServe()
	if err != nil {
		zLog.Fatal("ListenAndServe fail", zap.Error(err))
		return err
	}
	return nil
}

func InitWebsocket(conf *config.Config) error{
	_,cancel := context.WithCancel(context.Background()) //ctx
	defer cancel()

	r := router(conf)
	serve := &http.Server{
		Addr:conf.Websocket.Bind[0],
		Handler:r,
		ReadTimeout: time.Second * 3,
		WriteTimeout:time.Second *3,
	}
	go dispatchMsg(conf)
	err := serve.ListenAndServe()
	if err != nil {
		zLog.Fatal("ListenAndServe fail", zap.Error(err))
		return err
	}
	return nil
}

func GroupConn(w http.ResponseWriter, r *http.Request, conf *config.Config)  {
	// 超时控制
	//ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	//defer cancel()
	vars := mux.Vars(r)
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		//ZLogger.Fatal("https升级为wss失败", zap.Error(err))
	}
	userId := vars["userId"] //r.URL.Query().Get("user_id")
	groupId := vars["groupId"]

	zLog.Info("新的连接", zap.String("组_用户id", groupId + "-"+ userId), )
	defer func() {
		err := ws.Close()
		if err != nil {
			zLog.Warn("链接断开",
				zap.String("type", "group"),
				zap.String("GroupId", groupId),
				zap.String("userId", userId),
				zap.Error(err),
			)
			return
		}
	}()

	if len(userId) > 0 && userId != "0" {
		//var group *Group
		//本机上全局群组
		group,err := GlobalGroup.Get(groupId)
		if err != nil {
			redisContain := db.NewRedisContain(conf.Redis.GetAddr(), conf.Redis.Password)
			group = NewGroup()
			group.Id,_ = strconv.Atoi(groupId)
			group.Redis = redisContain
			GlobalGroup.set(groupId, group)
			GroupChannel <- group
		}
		//group, ok = Groups[groupId]
		//if !ok {
		//	zLog.Info("新建组",zap.String("group_id", groupId))
		//	log.Printf("新建组,group_id:%s", groupId)
		//	// 不存在这个群组
		//	group = NewGroup()
		//	Groups[groupId] = group
		//	GroupChannel <- group
		//}
		//group.Clients[userId] = client
		//group.Clients.Store(userId, c)

		c := NewConnection(ws, userId, conf.Xrpc.Addr)
		group.Join <- c

		//redisContain := db.NewRedisContain(conf.Redis.GetAddr(), conf.Redis.Password)
		//c.Redis = redisContain
		//redisContain.JoinMember(userId, groupId)

		//go c.Read(conf, group, Pool)
		//go c.Write(Pool)

		go c.ReadSocket(conf, group, Pool)
		//go c.WriteToSocket(group, Pool)

		for {
			msg := <- MessageChannel
			log.Println("得到消息")
			log.Println(msg)
			if len(msg.Text) == 0 {
				continue
			}
			if msg.UserId == c.UserId {
				log.Println("本人",c.UserId)
				continue
			}
			groupId := msg.GroupId

			currentGroup,err := GlobalGroup.Get(groupId)
			if err != nil {
				log.Println(err)
				continue
			}
			currentGroup.Clients.Range(func(key, c interface{}) bool {
				fmt.Print(key,"user_id \n")
				if err := c.(*Connection).Conn.WriteJSON(msg); err != nil {
					zLog.Warn("写入消息失败",
						zap.String("userId", "fdasf"),
						zap.Error(err),
					)
				}
				return true
			})
		}

		<-conf.Exit
	}

	return
}

//HandleConn get request then upgrade to websocket
func HandleConn(w http.ResponseWriter, r *http.Request, conf *config.Config)  {
	vars := mux.Vars(r)
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		//ZLogger.Fatal("https升级为wss失败", zap.Error(err))
	}
	userId := vars["userId"] //r.URL.Query().Get("user_id")
	groupId := vars["groupId"]


	zLog.Info("新的连接", zap.String("组_用户id", groupId + "-"+ userId), )
	defer func() {
		err := ws.Close()
		if err != nil {
			zLog.Warn("链接断开",
				zap.String("type", "group"),
				zap.String("GroupId", groupId),
				zap.String("userId", userId),
				zap.Error(err),
			)
			return
		}
	}()



	if len(userId) > 0 && userId != "0" {
		//链接加入到全局管理池
		Conns.Connections.Store(userId, ws)
		//对应 用户的加入到组
		//本机上全局群组
		group,err := GlobalGroup.Get(groupId)
		if err != nil {
			redisContain := db.NewRedisContain(conf.Redis.GetAddr(), conf.Redis.Password)
			group = NewGroup()
			group.Id,_ = strconv.Atoi(groupId)
			group.Redis = redisContain
			GlobalGroup.set(groupId, group)
		}

		//group.Clients.Store()
		c := NewConnection(ws, userId, conf.Xrpc.Addr)
		//group.Join <- c

		redisContain := db.NewRedisContain(conf.Redis.GetAddr(), conf.Redis.Password)
		//redisContain.JoinMember(userId, groupId)
		err = redisContain.SaveUserRpcAddr(groupId, userId, conf.Xrpc.Addr)
		log.Println(err)

		//go c.Read(conf, group, Pool)
		//go c.Write(Pool)
		go c.ReadSocket(conf, group, Pool)
		//go c.WriteToSocket(group, Pool)

		//for {
		//	msg := <- MessageChannel
		//	log.Println("得到消息")
		//	log.Println(msg)
		//	if len(msg.Text) == 0 {
		//		continue
		//	}
		//	if msg.UserId == c.UserId {
		//		log.Println("本人",c.UserId)
		//		continue
		//	}
		//	groupId := msg.GroupId
		//
		//	currentGroup,err := GlobalGroup.Get(groupId)
		//	if err != nil {
		//		log.Println(err)
		//		continue
		//	}
		//	currentGroup.Clients.Range(func(key, c interface{}) bool {
		//		fmt.Print(key,"user_id \n")
		//		if err := c.(*Connection).Conn.WriteJSON(msg); err != nil {
		//			zLog.Warn("写入消息失败",
		//				zap.String("userId", "fdasf"),
		//				zap.Error(err),
		//			)
		//		}
		//		return true
		//	})
		//}

		<-conf.Exit
	}

	return
}


func GroupsHandle() {
	for {
		g := <-GroupChannel
		zLog.Info("启动一个组协程，监听消息发生")
		go g.Handle()
	}
}


// Distribute 分发一个合适的服务器地址
func Distribute(w http.ResponseWriter, r *http.Request, c *config.Config){
	nodeWeight := make(map[string]int)
	nodeWeight["cakeIm_1:192.168.0.120:3102"] = 1
	//nodeWeight["cakeIm_2:192.168.0.120:3102"] = 2
	//nodeWeight["cakeIm_3:192.168.0.44:3102"] = 3
	virtualSpots := 100
	ring := util.NewHashRing(virtualSpots)
	ring.AddNodes(nodeWeight)

	//获取当前请求的ip
	ip := r.RemoteAddr
	target := ring.GetNode(ip)
	log.Printf("ip:%s", ip)

	//key := []byte(util.Key)
	//str,_ := util.Encrypt([]byte(target), key)
	data := pkg.Result{
		Code:200,
		Data: target,
	}
	result,_ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Write(result)
}

func dispatchMsg(conf *config.Config)  {
	redisContain := db.NewRedisContain(conf.Redis.GetAddr(), conf.Redis.Password)
	for {
		msg := <-MessageChannel

		log.Println("消息来啦")
		log.Println(msg)

		groupId := msg.GroupId

		//userIds := redisContain.GetUserIdsByGroupId(groupId)
		userIds := redisContain.GetUserMap(groupId)

		log.Println("所有用户",userIds)

		for uid := range userIds {

			if uid == msg.UserId {
				continue
			}
			if connection,ok := Conns.Connections.Load(uid);ok {
				socket := connection.(*websocket.Conn)
				if err := socket.WriteJSON(msg); err != nil {
					zLog.Warn("写入消息失败",
						zap.String("userId", uid),
						zap.Error(err),
					)
					err := socket.Close()
					if err != nil {
						zLog.Warn("关闭链接失败",
							zap.String("userId", uid),
							zap.Error(err),
						)
					}
					break
				}
			}else {
				//该用户已经离线了
				log.Println("该用户已经掉线 \n",uid)
			}

		}


	}
}