package controllers

import (
	"fmt"
	"go_autoapi/libs"
)

//type UserList struct {
//	Offset int `json:"offset"`
//	Page   int `json:"page"`
//}

// 获取用户列表 登录
func (c *AutoTestController) performTests() {
	loginMap := make(map[string]interface{})
	loginMap["user_name"] = "liuweiqiang" //账号
	loginMap["password"] = "OneixahDooquae1"
	for i := 0; i <= 100; i++ {
		go func(count int) {
			fmt.Println("%s 次执行", count)
			libs.DoRequest("http://127.0.0.1:8080/auto/login", "xxxxxxx", loginMap, map[string]interface{}{})
		}(i)
	}
	c.SuccessJson(nil)
}
