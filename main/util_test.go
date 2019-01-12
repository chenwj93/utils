package main

import (
	"fmt"
	"github.com/chenwj93/utils"
	"github.com/go-redis/redis"
	"testing"
	"time"
)

func TestRandomCode(t *testing.T) {
	chars := "0123456789" //abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ
	m := make(map[string]bool, 100000)
	r, err := utils.NewRandom(chars, 10)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 外挂一个自定义混淆函数
	r.MountRearrangeFunc(func(code *string) {
		str := *code
		i := time.Now().UnixNano() % 10
		if i != 0 && i != 9 {
			str = str[i+1:] + str[:i]
		}
		code = &str
	})
	for i := 0; i < 100000; i++ {
		m[r.GenerateCode()] = true
	}
	fmt.Println(len(m))
}

func TestRedis(t *testing.T) {
	utils.Client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
		//IdleTimeout: time.Duration(utils.ParseInt(conf.GetString("redis.idle"))),
		//DialTimeout: time.Duration(utils.ParseInt(conf.GetString("redis.wait"))),
		//PoolSize:    utils.ParseInt(conf.GetString("redis.active")),
	})

	lock := utils.NewRedisLock("lock", 3)
	lock.Get()
	time.Sleep(time.Second)
	lock.Get()

}
