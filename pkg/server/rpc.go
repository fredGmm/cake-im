package server

import (
	"cake-im/pkg"
	"context"
	"fmt"
	"time"
	"log"
)

type Req struct {
	UserId string
	Name string
	GroupId string
	MessageId string
	Text string
	Type int
	Time string
}

type Resp struct {
	Id string
	Status string
	Result string
	Err error
}

type Msg struct {
	Text string
}

//Send 将消息发到用户的连接上
func (m *Msg) Send(ctx context.Context,req *Req,resp *Resp) error {
	log.Println("start")

	msg := &pkg.Message{
		UserId :req.UserId,
		Name  : req.Name,
		Avatar  :   "",
		Text   :    req.Text,
		GroupId :   req.GroupId,
		ToUserId :  "0",
		Type:        "BroadcastRadio",
		OnlineCount : 0,
		Time  :     req.Time,
	}
	MessageChannel <- msg


	// 信息推入到 队列
	//conn.Redis.SetMessage(&msg)
	resp.Id = "finish \n"
	return nil
}




func (m *Msg) Push(ctx context.Context,req *Req, resp *Resp) error  {

	fmt.Printf("打印req, %s", req.MessageId)
	fmt.Printf("打印resp, %s", resp.Id)
	return nil
}

func (m *Msg) Test(ctx context.Context,req *Req, resp *Resp) error {
	resp.Id = fmt.Sprintf("当前时间:%s \n", time.Now().Format("2006-01-02 15:04:05"),)

	return nil
}

