package db

import "testing"

func TestNewRedisPool(t *testing.T) {
	addr := "192.168.0.7:6379"

	redisPool := NewRedisPool(addr)
	redis := redisPool.Get()
	redis.Do("ser")

}