package myredis

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

var (
	pool      *redis.Pool
	redisHost = "music-02.niracler.com:6377"
	redisPass = "123456"
)

// 创建redis连接池
func newRedisPool() *redis.Pool {
	return &redis.Pool{
		Dial: func() (conn redis.Conn, e error) {
			// 1. 打开连接
			c, err := redis.Dial("tcp", redisHost)
			if err != nil {
				fmt.Println(err.Error())
				return nil, err
			}
			// 2. 访问认证
			_, err = c.Do("AUTH", redisPass)
			if err != nil {
				c.Close()
				return nil, err
			}
			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
		MaxIdle:         50,
		MaxActive:       30,
		IdleTimeout:     300 * time.Second,
		Wait:            false,
		MaxConnLifetime: 0,
	}
}

// 初始化Redis连接池
func init() {
	pool = newRedisPool()
}

// 返回Redis连接池
func RedisPool() *redis.Pool {
	return pool
}
