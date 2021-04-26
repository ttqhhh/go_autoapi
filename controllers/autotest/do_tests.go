package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
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
	CaseList []int64 `json:"ids" form:"ids" `
}

// 获取用户列表 登录
func (c *AutoTestController) performTests() {
	uuid, _ := c.GenUUid()
	cl := CaseList{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &cl); err != nil {
		c.ErrorJson(-1, "请求参数错误", nil)
	}
	caseList, err := libs.GetCasesByIds(cl.CaseList)
	if err != nil {
		logs.Error("获取测试用例列表失败", err)
		c.ErrorJson(-1, "获取测试用例失败", nil)
	}
	fmt.Println("case list is", caseList)
	if len(caseList) == 0 {
		logs.Error("没有用例", err)
		c.ErrorJson(-1, "没有用例", nil)
	}
	for _, val := range caseList {
		go func(url string, uuid string, param string, checkout string, caseId int64) {
			libs.DoRequestV2(url, uuid, val.Parameter, val.Checkpoint, val.Id)
		}(val.ApiUrl, uuid, val.Parameter, val.Checkpoint, val.Id)
	}
	c.SuccessJsonWithMsg(map[string]interface{}{"uuid": uuid}, "OK")
}
