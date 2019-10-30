package main

import (
	"go.uber.org/zap"
	"github.com/fredGmm/imChong/config"
	"github.com/fredGmm/imChong/log"
	"github.com/fredGmm/imChong/pkg/connect"
	"github.com/fredGmm/imChong/log"
)

var ZLogger = Zlog.Logger()
var exit chan int

func main() {

	c, err := config.Init()
	if err != nil {
		ZLogger.Fatal("config fail", zap.Error(err))
	}

	srv := connect.NewServer(c)

	//runtime.NumCPU()
	if err2 := connect.InitTCP(srv, c.TCP.Bind[0]); err2 != nil {
		ZLogger.Fatal("tcp init fail, err2:", zap.Error(err2))
	}

	if err3 := connect.StartWebsocket(srv, c.Websocket.Bind[0]); err3 != nil {
		ZLogger.Fatal("websocket start fail", zap.Error(err3))
	}
	<-exit
}

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
