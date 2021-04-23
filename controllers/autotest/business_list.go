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
func (c *AutoTestController) businessList() {
	bl := BusinessList{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &bl); err != nil {
		logs.Error("parse business list param error")
		c.ErrorJson(-1, "请求参数错误", nil)
	}
	if bl.Page <= 0 || bl.Offset <= 0 {
		logs.Error("parse business list param error page %v,offset %v", bl.Page, bl.Offset)
		c.ErrorJson(-1, "请求参数错误", nil)
	}
	ab := models.AutoBusiness{}
	abs, err := ab.GetBusinessList(bl.Offset, bl.Page-1)
	if err != nil {
		logs.Error("failed to get business list")
		c.ErrorJson(-1, "系统错误", nil)
	}
	c.SuccessJsonWithMsg(abs, "OK")
}

// 获取用户列表 登录
func (c *AutoTestController) allBusinessList() {
	ab := models.AutoBusiness{}
	abs, err := ab.GetAllBusiness()
	if err != nil {
		logs.Error("failed to get all business list")
		c.ErrorJson(-1, "系统错误", nil)
	}
	logs.Error("all business is ", abs)
	c.SuccessJsonWithMsg(abs, "OK")
}
