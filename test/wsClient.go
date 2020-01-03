package main

import (
	"bufio"
	"cake-im/pkg"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"golang.org/x/net/websocket"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var exit chan int
var timer time.Time
var timeLayout = "2006-01-02 15:04:05"

func main() {
	// 使用时间串获取时间对象
	timer = time.Now()
	userId := strconv.Itoa(timer.Hour()) + strconv.Itoa(timer.Second()) + strconv.Itoa(timer.Minute())
	log.Printf("为您分配的id:%s", userId)
	conn := newWebsocketClient("1", userId)

	//log.Print("connect websocket successful")

	reader := bufio.NewReader(os.Stdin)

	go receiveMsg(conn)
	for {
		fmt.Print("我：")
		text, _ := reader.ReadString('\n')
		//text,_, _ := reader.ReadLine()
		msg := &pkg.Message{
			UserId: userId,
			//Message:time.Now().String(),
			Text:     strings.TrimRight(text, "\n"),
			GroupId:  "1",
			ToUserId: "0",
			Type:     "Broadcast",
			Name:"robot",
			Time:timer.Format("2019-02-28 10:48:21"),
			Avatar:"https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png",
		}
		sendData, _ := json.Marshal(msg)
		_, err := conn.Write(sendData)
		if err != nil {
			log.Printf("send msg fail, err: %v", err)
			return
		}
		//log.Printf("发送字节数 %d", n)
		//fmt.Println("\n ")
	}

}

func newWebsocketClient(groupId, userId string) *websocket.Conn {
	wsUrl := "ws://127.0.0.1:3102/group/" + groupId + "/" + userId
	//wsUrl := "wss://chongliao.j-book.cn/ws/group/" + groupId + "/" + userId
	//conn,err := net.DialTimeout("tcp", "127.0.0.1:3102", 3*time.Second)
	conn, err := websocket.Dial(wsUrl, "", "ws://chongliao。j-book.cn")
	if err != nil {
		log.Print("dial error")
		log.Println(err)
		os.Exit(90000)
	}
	return conn
}

// 批量连接 websocket 服务器

func batchNewWebsocketClient() {

}

func receiveMsg(conn *websocket.Conn) {

	for {
		buf := make([]byte, 2048)
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("read msg error . => %s", err.Error())
			os.Exit(90001)
		}
		msg := &pkg.Message{}
		err2 := json.Unmarshal(buf[:n], msg)
		if err2 != nil {
			log.Printf("读取字节数 %d", n)
			log.Print("消息格式错误")
			log.Printf("消息体 %s", string(buf))
			log.Printf("消息体 %s", err2.Error())
		}
		//fmt.Printf("接收到:%s \n", msg.Text)
		color.Blue(msg.Text)
	}
}
