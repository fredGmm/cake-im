package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
	ToUserId string `json:"to_user_id"`
	UserId   string `json:"user_id"`
}

var clients = make(map[*websocket.Conn]bool) // 客户端连接
var broadcast = make(chan Message)
var upgrader = websocket.Upgrader{}

var userConns = make(map[string]*websocket.Conn)
var userChannel = make(chan Message)

func main() {

	//fs := http.FileServer(http.Dir("./src/chongChat/web/"))
	//http.Handle("/", fs)

	http.HandleFunc("/ws", handleConn)

	go handleMessage()

	go sendMessageToUser()

	log.Println("绑定当前ip的8000端口")
	err := http.ListenAndServe("192.168.0.154:8000", nil)
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
	fmt.Println(r.URL)
	uid := r.URL.Query().Get("user_id")

	if len(uid) > 0 {
		userConns[uid] = ws
	}

	fmt.Println(userConns)

	defer func() {
		err := ws.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	clients[ws] = true

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Print(err)
			delete(clients, ws)
			break
		}
		fmt.Println(msg)
		fmt.Println(msg.ToUserId)

		//info := strings.Split(msg.Message, "|")
		//
		//fmt.Println(info)
		// msg 格式， userID|msg|secret
		if len(msg.ToUserId) > 0{
			// 发送给个人
			fmt.Println(msg.ToUserId)
			userChannel <- msg
		}else {
			broadcast <- msg
		}

	}
}

func handleMessage() {

	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				err := client.Close()
				if err != nil {
					fmt.Print(err)
				}
				delete(clients, client)
			}
		}
	}
}

// 一对一
func sendMessageToUser() {

	for {
		// 用户信息
		msg := <-userChannel
		toUserId = msg.ToUserId

		if userConns[to]
		for user, client := range userConns {
			err := client.WriteJSON(msg) // 发送给那个人

			if err == nil {

				fmt.Print("发送完成！")
			}else {
				log.Printf("error: %v", err)
				err := client.Close()
				if err != nil {
					fmt.Print(err)
				}
				println(user)
				delete(userConns, user)
			}
		}
	}

}
