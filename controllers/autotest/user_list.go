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
		logs.Error(1024, err)
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
		c.ErrorJson(-1, "登录失败", nil)
	}
	c.SuccessJson(auu)
}
