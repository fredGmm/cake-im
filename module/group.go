package module

import (
	"fmt"
	"log"
)

type Group struct {
	//Connections map[string]*websocket.Conn
	Clients map[string]*Client
	Send chan Message
	Join chan *Client
	Left chan string
}


func NewGroup() *Group {
	return &Group{
		Clients: make(map[string]*Client),
		Send:make(chan Message),
		Join:make(chan *Client),
		Left:make(chan string),
	}
}



//群组的消息处理
func (group *Group) Handle(){
	for {
		select {
		case c := <- group.Join:
			log.Print("这个组有新用户加入")
			group.Clients[c.UserId] = c

		case u := <- group.Left:
			log.Print("这个组有新用户离开")
			if _, ok := group.Clients[u]; ok {
				delete(group.Clients, u)
			}
		case m := <- group.Send:
			log.Print("有消息")
			fmt.Println(m)
			for _, c := range group.Clients {
				log.Print(c)
				select {
				case c.Send <- m :
				default:
					break
				}
			}
		}
	}
}
