package controllers

import (
	"encoding/json"
	"go_autoapi/libs"
)

type CheckOut struct {
	Url   string                            `json:"url"`
	Uuid  string                            `json:"uuid"`
	Param map[string]interface{}            `json:"param"`
	Check map[string]map[string]interface{} `json:"check_point"`
}

// 获取用户列表 登录
func (c *AutoTestController) performTests() {
	u := CheckOut{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &u); err != nil {
		c.ErrorJson(-1, "请求错误", nil)
	}
	var caseList []CheckOut
	for i := 0; i < 1; i++ {
		caseList = append(caseList, CheckOut{"http://127.0.0.1:8080/auto/login", "xxxxxxx", u.Param, u.Check})
	}
	//fmt.Println("case list is", caseList)
	for _, val := range caseList {
		go func(url string, uuid string, param map[string]interface{}, checkout map[string]map[string]interface{}) {
			libs.DoRequest(url, uuid, param, checkout)
		}(val.Url, val.Uuid, val.Param, val.Check)
	}
	c.SuccessJson(nil)
}
