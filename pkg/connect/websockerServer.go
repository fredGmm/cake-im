package connect

import (
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"github.com/fredGmm/imChong/pkg/db"
	"net/http"
)

var (
	upgrader = websocket.Upgrader{}
	pool     = &ConnectionPool{
		Connections: make(map[string]*Client),
		Join:        make(chan *Client),
		Left:        make(chan string),
	}
)

func StartWebsocket(s *Server, addr string) error {
	r := mux.NewRouter()
	// 加入 群组的路由
	r.HandleFunc("/group/{groupId}/{userId}", func(w http.ResponseWriter, r *http.Request) {
		GroupConn(w, r, s)
	})
	//r.HandleFunc("/group/{groupId}/{userId}", UserConn)
	//r.HandleFunc("/user/{userId}/{toUserId}", UserConn)
	http.Handle("/", r)
	go GroupsHandle()
	go pool.Handle()

	ZLogger.Info("http listen success :", zap.String("addr", addr))
	err2 := http.ListenAndServe(addr, nil)
	if err2 != nil {
		ZLogger.Fatal("ListenAndServe fail", zap.Error(err2))
		return err2
	}
	return nil
}

// 新建连接
func GroupConn(w http.ResponseWriter, r *http.Request, s *Server) {
	vars := mux.Vars(r)
	//不检查来源
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		ZLogger.Fatal("https升级为wss失败", zap.Error(err))
	}
	//userId := r.URL.Query().Get("id") // 用户id
	userId := vars["userId"]
	groupId := vars["groupId"]
	ZLogger.Info("new connect",
		zap.String("ip", r.RemoteAddr),
		zap.String("type", "group"),
		zap.String("groupId", groupId),
		zap.String("userId", userId),
		zap.String("path", r.URL.Path),
	)
	defer func() {
		err := ws.Close()
		if err != nil {
			ZLogger.Warn("链接断开",
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
		c := NewClient(ws, userId)
		var group *Group
		var ok bool
		group, ok = Groups[groupId]
		if !ok {
			// 不存在这个群组
			group = NewGroup()
			Groups[groupId] = group
			GroupChannel <- group
		}
		//group.Clients[userId] = client
		group.Clients.Store(userId, c)

		redisContain := db.NewRedisContain(s.c.Redis.GetAddr(), s.c.Redis.Password)
		c.Redis = redisContain
		redisContain.JoinMember(userId, groupId)

		go c.ReadMessage(s, group)
		go c.WriteMessage(pool)
		<-s.Exit
	}

}

//func UserConn(w http.ResponseWriter, r *http.Request) {
//	vars := mux.Vars(r)
//	//不检查来源
//	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
//	userId := vars["userId"]
//	ws, err := upgrader.Upgrade(w, r, nil)
//	if err != nil {
//		ZLogger.Warn("链接断开",
//			zap.String("type", "user"),
//			zap.String("userId", userId),
//			zap.Error(err),
//		)
//	}
//
//	ZLogger.Info("新用户链接", zap.String("type", "user"), zap.String("userId", userId), )
//	defer func() {
//		err := ws.Close()
//		if err != nil {
//			fmt.Println("链接时 断开:")
//			fmt.Println(userId)
//			fmt.Println(err)
//			// 发送消息通知
//			return
//		}
//	}()
//
//	var c *client.Client
//	//链接
//	if len(userId) > 0 && userId != "0" {
//		var ok bool
//		if c, ok = pool.Connections[userId]; !ok {
//			c = client.NewClient(ws, userId)
//		}
//		//pool.Left <- userId
//		pool.Join <- c
//		go c.ReadMessageToUser(pool)
//		go c.WriteMessage(pool)
//	}
//
//	<-exit
//}

func GroupsHandle() {
	for {
		g := <-GroupChannel
		ZLogger.Info("启动一个组协程")
		go g.Handle()
	}
}
