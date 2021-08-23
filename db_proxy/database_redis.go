package db_proxy

import (
	"github.com/go-redis/redis"
)

// 声明一个全局的rdb变量
var rdb *redis.Client

// 初始化连接
func InitClient() (err error) {
	rdb = redis.NewClient(&redis.Options{
		//Addr: "172.20.20.2:6379", // 测试
		Addr:     "172.16.2.86:6379", // 生产
		Password: "",                 // no password set
		DB:       0,                  // use default DB
	})

	_, err = rdb.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

func GetRedisObject() *redis.Client {
	return rdb
}
