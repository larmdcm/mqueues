package redispool

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

//构造一个链接函数，如果没有密码，passwd为空字符串
func RedisConn(ip,port,passwd string) (redis.Conn, error) {
	c,err := redis.Dial("tcp",
		ip+":"+port,
		redis.DialConnectTimeout(5*time.Second),
		redis.DialReadTimeout(1*time.Second),
		redis.DialWriteTimeout(1*time.Second),
		redis.DialPassword(passwd),
		redis.DialKeepAlive(1*time.Second),
	)
	return c,err
}

//构造一个连接池
//url为包装了redis的连接参数ip,port,passwd
func NewPool(ip,port,passwd string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:            5,    //定义redis连接池中最大的空闲链接为5
		MaxActive:          18,    //在给定时间已分配的最大连接数(限制并发数)
		IdleTimeout:        240 * time.Second,
		MaxConnLifetime:    300 * time.Second,
		Dial:               func() (redis.Conn,error) { return RedisConn(ip,port,passwd) },
	}
}

