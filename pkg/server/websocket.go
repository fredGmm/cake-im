package server

import (
	"cake-im/pkg/client"
	"cake-im/pkg/db"
	"cake-im/pkg/zapLog"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"log"

	//"go.uber.org/zap"
	"net/http"
)

var (
	upgrader = websocket.Upgrader{}
	pool = &ConnectionPool{
		Connections:make(map[string]*Connection),
		Join:make(chan *Connection),
		Left:make(chan string),
	}
	zLog = zapLog.Logger()
)

func InitWebsocket(gClient *client.GrpcClient) (err error){
	r := mux.NewRouter()
	// 加入 群组的路由
	r.HandleFunc("/group/{groupId}/{userId}", func(w http.ResponseWriter, r *http.Request) {
		GroupConn(w, r, gClient)
	})
	//r.HandleFunc("/group/{groupId}/{userId}", UserConn)
	//r.HandleFunc("/user/{userId}/{toUserId}", UserConn)
	http.Handle("/", r)
	go GroupsHandle()
	go pool.Handle()

	//ZLogger.Info("http listen success :", zap.String("addr", addr))
	err2 := http.ListenAndServe(gClient.Config.Websocket.Bind[0], nil)
	if err2 != nil {
		zLog.Fatal("ListenAndServe fail", zap.Error(err2))
		return err2
	}
	return
}

func GroupConn(w http.ResponseWriter, r *http.Request, gClient *client.GrpcClient)  {
	vars := mux.Vars(r)

	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		//ZLogger.Fatal("https升级为wss失败", zap.Error(err))
	}
	userId := vars["userId"] //r.URL.Query().Get("user_id")
	groupId := vars["groupId"]

	zLog.Info("新的连接",
		zap.String("组-用户id", groupId + "-"+ userId),
		)

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

	//链接
	if len(userId) > 0 && userId != "0" {
		c := NewConnection(ws, userId)
		var group *Group
		var ok bool
		group, ok = Groups[groupId]
		if !ok {
			log.Printf("新建组,group_id:%s", groupId)
			// 不存在这个群组
			group = NewGroup()
			Groups[groupId] = group
			GroupChannel <- group
		}
		//group.Clients[userId] = client
		group.Clients.Store(userId, c)

		redisContain := db.NewRedisContain(gClient.Config.Redis.GetAddr(), gClient.Config.Redis.Password)
		c.Redis = redisContain
		redisContain.JoinMember(userId, groupId)

		go c.Read(gClient, group)
		go c.Write(pool)

		<-gClient.Exit
	}
}

func GroupsHandle() {
	for {
		g := <-GroupChannel
		zLog.Info("启动一个组协程，监听消息发生")
		go g.Handle()
	}
}
