package main

import (
	"github.com/WTwilight/instant-messaging-system/server/model"
	"github.com/gomodule/redigo/redis"
	"time"
)

// 定义一个全局 pool
var redisConnPool *redis.Pool

func initPool(addr string, maxIdle, maxActive int, idleTimeout time.Duration) {
	redisConnPool = &redis.Pool{
		MaxIdle: maxIdle,
		MaxActive: maxActive,
		IdleTimeout: idleTimeout,
		Dial: func() (conn redis.Conn, e error) {
			return redis.Dial("tcp", addr)
		},
	}
}

// 编写函数对UserDao初始化
func initUserDao() {
	//
	model.UserDaoHandle = model.NewUserDao(redisConnPool)
}

func init() {
	// 初始化 redis 连接池
	initPool("localhost:6379", 16,0, 300 * time.Second)
	initUserDao()
}