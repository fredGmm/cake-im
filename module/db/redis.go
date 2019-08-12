package db

import (
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
	"imChong/log"
	"imChong/module"
	"time"
)
var log = Zlog.Logger()

type RedisContain struct {
	ServerId int
	Pool *redis.Pool
}

func NewRedisContain(addr, password string) *RedisContain {

	pool := &redis.Pool{
			MaxIdle:10,
		IdleTimeout:240 * time.Second,
		Wait:true,
		Dial: func() (conn redis.Conn, e error) {
			conn, err := redis.Dial("tcp", addr,
				//redis.DialConnectTimeout(10),
				//redis.DialReadTimeout(10),
				//redis.DialWriteTimeout(10),
				//redis.DialPassword(""),
				)
			if _, err := conn.Do("AUTH", password); err != nil {
				conn.Close()
				return nil, err
			}
			if err != nil {
				return nil, err
			}
			return conn, nil
		},
	}

	return &RedisContain{
		ServerId:0,
		Pool:pool,
	}
}

func getGroupMessageKey(groupId string ) string {

	return "GROUP_MESSAGE_" + groupId
}

// save message
func (r *RedisContain) SetMessage(message module.Message) error {
	conn := r.Pool.Get()
	msg, _ := json.Marshal(message)
	key := getGroupMessageKey(message.GroupId)
	_,err := conn.Do("LPUSH", key, string(msg))
	log.Info("存入消息到中间件",
		zap.String("message", message.Message),
		)
	return err
}

func (r *RedisContain) JoinMember(userId , groupId string) {
	conn := r.Pool.Get()
	//在线人数
	_, err :=conn.Do("SADD", "ONLINE_COUNT", userId)
	if err != nil {
		log.Error("在线成员添加错误")
	}

	//房间人数
	key := "GROUP_MEMBER_" + groupId
	_, err2 := conn.Do("SADD", key, userId)
	if err2 != nil {
		log.Error("组群添加成员出现错误",
			zap.String("groupId", groupId),
			zap.String("userId", userId),
			)
	}
}

func (r *RedisContain) Test(key, value string) error{
	conn := r.Pool.Get()
	_, err := conn.Do("SET",key,value)
	if err != nil {
		fmt.Println(err)
	}
	return err
}

