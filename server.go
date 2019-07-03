package main

import (
	"fmt"
	"github.com/Unknwon/goconfig"
	"github.com/alecthomas/kingpin"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"imChong/module"
	"log"
	"net/http"
)

type Result struct {
	Code    int    `json:"code"`
	Data    []int  `json:"data"`
	Message string `json:"message"`
}

var upgrader = websocket.Upgrader{}
var cfg *goconfig.ConfigFile
var configPath string

type Config struct {
	dev       int
	apiPrefix string
	ip        string
	port      string
}
var config = Config{}

var groupChannel = make(chan *module.Group)
var Groups = make(map[string]*module.Group)
var pool = &module.ConnectionPool{
	Connections:make(map[string]*module.Client),
	Join:make(chan *module.Client),
	Left:make(chan string),

}
var exit chan int

func parseFlag() {
	kingpin.HelpFlag.Short('h')
	kingpin.Version("0.0.1")
	kingpin.Flag("dev", "是否开发").Default("1").IntVar(&config.dev)
	kingpin.Flag("configPath", "配置文件地址").Default("E:\\code\\golang\\src\\imChong\\Config.ini").StringVar(&configPath)
	kingpin.Parse() // first parse conf
	if config.dev == 1 {
		cfg, _ = goconfig.LoadConfigFile(configPath)
		config.apiPrefix, _ = cfg.GetValue("dev", "apiPrefix")
		config.ip, _ = cfg.GetValue("dev", "ip")
		config.port, _ = cfg.GetValue("dev", "port")
	} else {
		cfg, _ = goconfig.LoadConfigFile("/data/chongliao.com/imChong/Config.ini")
		config.apiPrefix, _ = cfg.GetValue("product", "apiPrefix")
		config.ip, _ = cfg.GetValue("product", "ip")
		config.port, _ = cfg.GetValue("product", "port")
	}
}

func main() {
	parseFlag()
	r := mux.NewRouter()
	r.HandleFunc("/group/{groupId}/{userId}", GroupConn)
	r.HandleFunc("/user/{userId}/{toUserId}", UserConn)
	http.Handle("/", r)
	go GroupsHandle()
	go pool.Handle()
	err := http.ListenAndServe(config.ip+":"+config.port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

// 新建连接
func GroupConn(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	//不检查来源
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("https升级为wss失败:err:")
		log.Fatal(err)
	}
	//userId := r.URL.Query().Get("id") // 用户id
	//groupId := r.URL.Query().Get("group_id") // 群组id

	userId := vars["userId"]
	groupId := vars["groupId"]
	fmt.Println(groupId,userId)
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

	//链接
	if len(userId) > 0 && userId != "0" {
		client := module.NewClient(ws, userId)
		var group *module.Group
		var ok bool
		group, ok = Groups[groupId]
		if !ok {
			// 不存在这个群组
			group = module.NewGroup()
			Groups[groupId] = group
			groupChannel <- group
		}
		group.Clients[userId] = client
		go client.ReadMessage(group, nil)
		go client.WriteMessage(pool)
		<-exit
	}
}

func UserConn(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//不检查来源
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("https升级为wss失败:err:")
		log.Fatal(err)
	}
	userId := vars["userId"]
	fmt.Println(userId)
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

	var client *module.Client
	//链接
	if len(userId) > 0 && userId != "0" {
		var ok bool
		if client,ok = pool.Connections[userId]; !ok {
			client = module.NewClient(ws, userId)
		}
		//pool.Left <- userId
		pool.Join <- client
		go client.ReadMessageToUser(pool)
		go client.WriteMessage(pool)
	}

	<-exit
}

func GroupsHandle() {
	for {
		g := <- groupChannel
		log.Print("启动")
		log.Print(g)
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