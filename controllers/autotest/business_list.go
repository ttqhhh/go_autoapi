package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"go_autoapi/models"
)

type BusinessList struct {
	Offset int `json:"offset"`
	Page   int `json:"page"`
}

// 获取用户列表 登录
func (c *AutoTestController) BusinessList() {
	bl := BusinessList{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &bl); err != nil {
		logs.Error("user list parse param error")
		c.ErrorJson(-1, "请求参数错误", nil)
	}
	if bl.Page <= 0 || bl.Offset <= 0 {
		logs.Error(1024, "param error", bl.Offset, bl.Page)
		c.ErrorJson(-1, "请求参数错误", nil)
	}
	ab := models.AutoBusiness{}
	abs, err := ab.GetBusinessList(bl.Offset, bl.Page)
	if err != nil {
		logs.Error("failed to get user list")
		c.ErrorJson(-1, "系统错误", nil)
	}
	c.SuccessJson(abs)
}
