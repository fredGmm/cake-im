package main

import (
	"fmt"
	"github.com/Unknwon/goconfig"
	"github.com/alecthomas/kingpin"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"imChong/log"
	"imChong/module/client"
	"imChong/module/db"
	"net/http"
)

type Config struct {
	dev       int
	apiPrefix string
	ip        string
	port      string
	redisAddr 	  string
	redisPwd  string
}

var (
	upgrader = websocket.Upgrader{}
	cfg *goconfig.ConfigFile
	config = Config{}
	configPath string
	groupChannel = make(chan *client.Group)
	Groups = make(map[string]*client.Group)
)
// 链接存放池 （可以适当备份）
var pool = &client.ConnectionPool{
	Connections: make(map[string]*client.Client),
	Join:        make(chan *client.Client),
	Left:        make(chan string),
}
var ZLogger = Zlog.Logger()
var exit chan int

const (
	version = "0.0.1"
)


func parseFlag() {
	var defaultConfigPath string
	kingpin.HelpFlag.Short('h')
	kingpin.Version("0.0.1")
	defaultConfigPath = "E:\\code\\golang\\src\\imChong\\Config.ini"
	kingpin.Flag("configPath", "配置文件地址").Default(defaultConfigPath).StringVar(&configPath)
	kingpin.Parse() // first parse conf
	cfg, _ = goconfig.LoadConfigFile(configPath)
	config.apiPrefix, _ = cfg.GetValue("product", "apiPrefix")
	config.ip, _ = cfg.GetValue("product", "ip")
	config.port, _ = cfg.GetValue("product", "port")
	config.redisAddr, _ = cfg.GetValue("product", "redisAddr")
	config.redisPwd, _ = cfg.GetValue("product", "redisPwd")
}

func main() {
	parseFlag()
	
	r := mux.NewRouter()
	r.HandleFunc("/group/{groupId}/{userId}", GroupConn)
	r.HandleFunc("/user/{userId}/{toUserId}", UserConn)
	http.Handle("/", r)
	go GroupsHandle()
	go pool.Handle()
	ZLogger.Info(config.ip + ":" + config.port)
	err := http.ListenAndServe(config.ip+":"+config.port, nil)
	if err != nil {
		ZLogger.Fatal("ListenAndServe fail", zap.Error(err))
	}
}

// 新建连接
func GroupConn(w http.ResponseWriter, r *http.Request) {
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
		c := client.NewClient(ws, userId)
		var group *client.Group
		var ok bool
		group, ok = Groups[groupId]
		if !ok {
			// 不存在这个群组
			group = client.NewGroup()
			Groups[groupId] = group
			groupChannel <- group
		}
		//group.Clients[userId] = client
		group.Clients.Store(userId, c)
		redisContain := db.NewRedisContain(config.redisAddr, config.redisPwd)
		c.Redis = redisContain
		redisContain.JoinMember(userId, groupId)
		go c.ReadMessage(group)
		go c.WriteMessage(pool)
		<-exit
	}
}

func UserConn(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//不检查来源
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	userId := vars["userId"]
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		ZLogger.Warn("链接断开",
			zap.String("type", "user"),
			zap.String("userId", userId),
			zap.Error(err),
		)
	}

	ZLogger.Info("新用户链接", zap.String("type", "user"), zap.String("userId", userId), )
	defer func() {
		err := ws.Close()
		if err != nil {
			fmt.Println("链接时 断开:")
			fmt.Println(userId)
			fmt.Println(err)
			// 发送消息通知
			return
		}
	}()

	var c *client.Client
	//链接
	if len(userId) > 0 && userId != "0" {
		var ok bool
		if c, ok = pool.Connections[userId]; !ok {
			c = client.NewClient(ws, userId)
		}
		//pool.Left <- userId
		pool.Join <- c
		go c.ReadMessageToUser(pool)
		go c.WriteMessage(pool)
	}

	<-exit
}

func GroupsHandle() {
	for {
		g := <-groupChannel
		ZLogger.Info("启动一个组协程")
		go g.Handle()
	}
}

//func handleConn(w http.ResponseWriter, r *http.Request) {
//	//不检查来源
//	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
//	ws, err := upgrader.Upgrade(w, r, nil)
//	if err != nil {
//		log.Print("https升级为wss失败:err:")
//		log.Fatal(err)
//	}
//	defer func() {
//		err := ws.Close()
//		if err != nil {
//			fmt.Println("链接断开:")
//			fmt.Println(err)
//			return
//		}
//	}()
//
//	uid := r.URL.Query().Get("id")
//	gid := r.URL.Query().Get("gid")
//	//uid 校验
//	if validId, err := strconv.Atoi(uid); validId > 0 && err == nil {
//		// 所有在线用户链接
//		fmt.Println(userConnections)
//		userConnections[uid] = ws
//		fmt.Printf("用户%s已经连接", uid)
//		for {
//			var msg Message
//			err2 := ws.ReadJSON(&msg)
//			if err2 != nil {
//				log.Print("这个ws 已经断开")
//				//log.Print(err)
//				delete(userConnections, uid)
//				break
//			}
//			fmt.Println(msg.ToUserId)
//			msg.UserId = uid
//
//			if validToId, err := strconv.Atoi(msg.ToUserId); validToId > 0 && err == nil {
//				userChannel <- msg
//			} else if validGid, err := strconv.Atoi(gid); validGid > 0 && err == nil {
//				msg.GroupId = gid
//				fmt.Printf("群ID:%s,消息：%s：", gid, msg.Message)
//				groupChannel <- msg
//			} else {
//				fmt.Print("莫名消息")
//				fmt.Println(msg)
//			}
//		}
//	}
//}

//func sendMessageToGroup() {
//	for {
//		msg := <-groupChannel
//		gid := msg.GroupId
//		fmt.Println(msg)
//		//获取组成员
//		//groupMembers := groupMembers[gid]
//		resp, err := http.Get(config.apiPrefix + "group/members?group_id=" + gid + "&user_id=" + msg.UserId)
//		if err != nil {
//			fmt.Print(err)
//		}
//		body, err2 := ioutil.ReadAll(resp.Body)
//		if err2 != nil {
//			fmt.Print(err)
//		}
//		result := Result{}
//		err3 := json.Unmarshal(body, &result)
//		if err3 != nil {
//			fmt.Print(err3)
//		}
//		fmt.Println(result.Data)
//		for _, id := range result.Data {
//			userId := strconv.Itoa(id)
//			fmt.Println(id)
//			fmt.Println(userId)
//			if client, ok := userConnections[userId]; ok {
//				err4 := client.WriteJSON(msg)
//				if err4 != nil {
//					log.Printf("群发消息错误: %v", err4)
//					err5 := client.Close()
//					if err5 != nil {
//						fmt.Print(err)
//					}
//				} else {
//					fmt.Print("发送完毕！")
//				}
//			} else {
//				fmt.Printf("不存在这种群用户id: %s\n", userId)
//			}
//		}
//	}
//}
//
//// 一对一
//func sendMessageToUser() {
//	for {
//		// 用户信息
//		msg := <-userChannel
//		toUserId := msg.ToUserId
//		fmt.Println(msg.Message)
//		if conn, ok := userConnections[toUserId]; ok {
//			err := conn.WriteJSON(msg.Message) // 发送给那个人
//			if err == nil {
//				fmt.Print("发送完成！")
//			} else {
//				log.Printf("error: %v", err)
//				err := conn.Close()
//				if err != nil {
//					fmt.Print(err)
//				}
//				fmt.Println(toUserId)
//			}
//		} else {
//			fmt.Printf("对方不在线.toUserId:%s \n", toUserId)
//
//			if len(msg.Message) > 0 {
//				// todo 保存数据库
//				resp, err := http.PostForm(config.apiPrefix+"message/create",
//					url.Values{"from_id": {msg.UserId}, "to_id": {toUserId}, "content": {msg.Message}, "status": {"0"}})
//				if err != nil {
//
//				}
//				body, err := ioutil.ReadAll(resp.Body)
//				if err != nil {
//					fmt.Print(err)
//				}
//				fmt.Println(string(body))
//				resp.Body.Close()
//			}
//
//		}
//
//	}
//}
//
//
//func simulateConn(){
//
//	for i:=0;i<1000;i++ {
//		userConnections[strconv.Itoa(i)] = &websocket.Conn{}
//	}
//}
