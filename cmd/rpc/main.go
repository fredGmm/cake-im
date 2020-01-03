package main

import (
	"cake-im/pkg/server"
	pb "cake-im/rpc/chat/grpc"
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
)

//var addr = "192.168.0.43:9527"
var addr = "127.0.0.1:9527"
var count = 1
func main() {
	fmt.Println("开始监听 grpc 服务端口,", addr)
	listen, err3 := net.Listen("tcp", addr)
	if err3 != nil {
		fmt.Println("发生错误")
	}
	s := grpc.NewServer()
	pb.RegisterChatServer(s, &server.RpcServer{})
	//go registerServer()
	if err := s.Serve(listen); err != nil {
		fmt.Println("发生错误", err)
	}
}

func registerServer() {
	config := consulapi.DefaultConfig()
	config.Address = "192.168.0.7:8500"
	client, err := consulapi.NewClient(config)
	if err != nil {
		log.Fatal("consul client error : ", err)
	}
	checkPort := 8080
	registration := new(consulapi.AgentServiceRegistration)
	//registration.ID = "chong"

	registration.Name = "chongliao-grpc"
	registration.Port = 9527
	registration.Tags = []string{"grpc"}
	registration.Address = "192.168.0.43"

	registration.Check = &consulapi.AgentServiceCheck{
		HTTP:                           fmt.Sprintf("http://%s:%d%s", registration.Address, checkPort, "/check"),
		Timeout:                        "3s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "30s", // 指定时间不健康则注销
	}

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		log.Fatal("register server error : ", err)
	}
	client.Agent().Join("192.168.0.7", true)
	http.HandleFunc("/check", consulCheck)
	http.ListenAndServe(fmt.Sprintf(":%d", checkPort), nil)
}

func consulCheck(w http.ResponseWriter, r *http.Request) {
	s := "consulCheck" + fmt.Sprint(count) + "remote:" + r.RemoteAddr + " " + r.URL.String()
	fmt.Println(s)
	fmt.Fprintln(w, s)
	count++
}

func localIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
