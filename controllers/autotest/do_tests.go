package controllers

import (
	"encoding/json"
	"go_autoapi/libs"
)

// 请求demo，如何传递jsonpath
//{
//"url": "http://127.0.0.1:8080/auto/login",
//"uuid": "xxxxx",
//"param": {
//"user_name": "liuweiqiang",
//"password": "OneixahDooquae1"
//},
//"check_point": {
//"$.code": {
//"eq": 200
//},
// 两个..email，会递归查找所有key=email元素，返回一个list
//"$..email": {
//"in": "liuweiqiang2014"
//}
//}
//}
type CheckOut struct {
	Url   string                            `json:"url"`
	Param map[string]interface{}            `json:"param"`
	Check map[string]map[string]interface{} `json:"check_point"`
}

type CaseList struct {
	CaseList []int `json:"ids"`
}

// 获取用户列表 登录
func (c *AutoTestController) performTests() {
	uuid, _ := c.GenUUid()
	var caseId int64
	caseId = 10
	u := CheckOut{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &u); err != nil {
		c.ErrorJson(-1, "请求错误", nil)
	}
	var caseList []CheckOut
	for i := 0; i < 1; i++ {
		caseList = append(caseList, CheckOut{"http://127.0.0.1:8080/auto/login", u.Param, u.Check})
	}
	//fmt.Println("case list is", caseList)
	for _, val := range caseList {
		go func(url string, uuid string, param map[string]interface{}, checkout map[string]map[string]interface{}) {
			libs.DoRequest(url, uuid, param, checkout, caseId)
		}(val.Url, uuid, val.Param, val.Check)
	}
	c.SuccessJson(map[string]interface{}{"uuid": uuid}, "OK")
}
