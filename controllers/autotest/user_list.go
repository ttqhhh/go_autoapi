package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"go_autoapi/models"
)

type UserList struct {
	Offset int `json:"offset"`
	Page   int `json:"page"`
}

// 获取用户列表 登录
func (c *AutoTestController) userList() {
	ul := UserList{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &ul); err != nil {
		logs.Error("user list parse param error")
		c.ErrorJson(-1, "请求参数错误", nil)
	}
	if ul.Page <= 0 || ul.Offset <= 0 {
		logs.Error(1024, "param error", ul.Offset, ul.Page)
		c.ErrorJson(-1, "请求参数错误", nil)
	}
	au := models.AutoUser{}
	auu, err := au.GetUserList(ul.Offset, ul.Page-1)
	if err != nil {
		logs.Error("failed to get user list")
		c.ErrorJson(-1, "系统错误", nil)
	}
	c.SuccessJsonWithMsg(auu, "OK")
}

// todo 重写了获取用户数据列表的方法接口，需要再把这块整理整理
func (c *AutoTestController) getUserList() {
	//ul := UserList{}
	//if err := json.Unmarshal(c.Ctx.Input.RequestBody, &ul); err != nil {
	//	logs.Error("user list parse param error")
	//	c.ErrorJson(-1, "请求参数错误", nil)
	//}
	page, err := c.GetUint8("page", 1)
	if err != nil {
		logs.Error("param error: %v", err)
		c.ErrorJson(-1, "请求参数错误", nil)
	}
	limit, err := c.GetInt8("limit", 10)
	if err != nil {
		logs.Error("param error: %v", err)
		c.ErrorJson(-1, "请求参数错误", nil)
	}
	if page <= 0 || limit <= 0 {
		logs.Error(1024, "param error", limit, page)
		c.ErrorJson(-1, "请求参数错误", nil)
	}
	// 获取用户列表
	au := models.AutoUser{}
	auu, err := au.GetUserList(int(limit), int(page))
	if err != nil {
		logs.Error("failed to get user list")
		c.ErrorJson(-1, "系统错误", nil)
	}
	//c.SuccessJsonWithMsg(auu, "OK")
	// 获取用户条数
	total, err := au.GetActivateUserCount()
	if err != nil {
		logs.Error("failed to get user list count")
		c.ErrorJson(-1, "系统错误", nil)
	}

	res := make(map[string]interface{})
	res["code"] = 0
	res["msg"] = "成功"
	res["count"] = total
	res["data"] = auu

	c.Data["json"] = res
	c.ServeJSON() //对json进行序列化输出
	c.StopRun()
}

func (c *AutoTestController) userIndex() {
	c.TplName = "user.html"
}
