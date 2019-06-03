package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type Message struct {
	//Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
	ToUserId string `json:"to_user_id"`
	UserId   string `json:"user_id"`
	GroupId string `json:"group_id"`
}

var clients = make(map[*websocket.Conn]bool) // 客户端连接

var userConnections = make(map[string]*websocket.Conn) // 所有人的链接

var groupMembers = make(map[string][]string)     // 群组成员
var groupChannel = make(chan Message)
var userChannel = make(chan Message)

var upgrader = websocket.Upgrader{}

func main() {

	//fs := http.FileServer(http.Dir("./src/chongChat/web/"))
	//http.Handle("/", fs)

	http.HandleFunc("/ws", handleConn)
	go sendMessageToUser()

	err := http.ListenAndServe("127.0.0.1:3000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConn(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := ws.Close()
		if err != nil {
			fmt.Println("链接断开:")
			fmt.Println(err)
			return
		}
	}()

	uid := r.URL.Query().Get("id")
	gid := r.URL.Query().Get("gid")
	//uid 校验
	if validId, err := strconv.Atoi(uid); validId > 0 && err == nil {
		// 所有在线用户链接
		fmt.Println(userConnections)
		userConnections[uid] = ws
		fmt.Printf("用户%s已经连接", uid)
		for {
			var msg Message
			err2 := ws.ReadJSON(&msg)
			if err2 != nil {
				log.Print("这个ws 已经断开")
				//log.Print(err)
				delete(userConnections, uid)
				break
			}
			fmt.Println(msg.ToUserId)
			msg.UserId = uid
			userChannel <- msg
		}
	}

	//群消息
	if validId, err := strconv.Atoi(gid); validId > 0 && err == nil {
		clients[ws] = true

		groupMembers["6"] = make([]string, 5)
		//获取群成员
		members := make([]string, 5)
		for i := 1; i < 4; i++ {
			groupMembers[gid] = append(members, strconv.Itoa(i))
		}

		print("组成员", groupMembers)
		for {
			var msg Message
			err := ws.ReadJSON(&msg)
			if err != nil {
				log.Print(err)
				delete(userConnections, uid)
				break
			}
			fmt.Println("群消息：" , msg)
			groupChannel <- msg

		}
	}
}

func sendMessageToGroup() {
	for {
		msg := <-groupChannel

		gid := msg.GroupId

		groupMembers := groupMembers[gid]
		for _,userId := range groupMembers {
			client := userConnections[userId]
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("群发消息错误: %v", err)
				err := client.Close()
				if err != nil {
					fmt.Print(err)
				}
				//delete(clients, client) ??
			}
		}
	}
}

// 一对一
func sendMessageToUser() {
	for {
		// 用户信息
		msg := <-userChannel
		toUserId := msg.ToUserId
		fmt.Println(msg.Message)
		if conn, ok := userConnections[toUserId]; ok {
			err := conn.WriteJSON(msg.Message) // 发送给那个人
			if err == nil {
				fmt.Print("发送完成！")
			} else {
				log.Printf("error: %v", err)
				err := conn.Close()
				if err != nil {
					fmt.Print(err)
				}
				fmt.Println(toUserId)
			}
		} else {
			fmt.Printf("对方不在线.toUserId:%s \n", toUserId)

			if len(msg.Message) > 0 {
				// todo 保存数据库
				resp, err := http.PostForm("http://192.168.0.66/message/create",
					url.Values{"from_id": {msg.UserId}, "to_id": {toUserId}, "content": {msg.Message}, "status": {"0"}})

				if err != nil {
					// handle error
				}

				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					// handle error
				}
				fmt.Println(string(body))
				resp.Body.Close()
			}

		}

	}
}
