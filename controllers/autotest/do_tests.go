package controllers

import (
	"fmt"
	"go_autoapi/libs"
)

type CheckOut struct {
	Url   string                 `json:"url"`
	Uuid  string                 `json:"uuid"`
	Param map[string]interface{} `json:"param"`
	Check interface{}            `json:"check_point"`
}

// 获取用户列表 登录
func (c *AutoTestController) performTests() {
	p := map[string]interface{}{"user_name": "liuweiqiang", "password": "OneixahDooquae1"}
	a := map[string]interface{}{"ret": `{"equal":"1"}`, "password": "OneixahDooquae1"}
	var caseList []CheckOut
	for i := 0; i < 10; i++ {
		caseList = append(caseList, CheckOut{"http://127.0.0.1:8080/auto/login", "xxxxxxx", p, a})
	}
	fmt.Println("case list is", caseList)
	for _, val := range caseList {
		go func(url string, uuid string, param map[string]interface{}, checkout interface{}) {
			//fmt.Println("%s 次执行", key,val)
			libs.DoRequest(url, uuid, param, checkout)
		}(val.Url, val.Uuid, val.Param, val.Check)
	}
	c.SuccessJson(nil)
}
