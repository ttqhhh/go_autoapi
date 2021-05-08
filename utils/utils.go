package utils

import (
	"github.com/go-redis/redis"
	"go_autoapi/db_proxy"
	"reflect"
)

func Struct2Map(obj interface{}) map[interface{}]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[interface{}]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

func GetRedis() *redis.Client {
	_ = db_proxy.InitClient()
	return db_proxy.GetRedisObject()
}
