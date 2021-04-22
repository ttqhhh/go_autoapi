package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"go_autoapi/models"
)

// 获取用户列表 登录
func (c *AutoTestController) updateUser() {
	au := models.AutoUser{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &au); err != nil {
		logs.Error("update param parse error", err)
		c.ErrorJson(-1, "请求参数错误", nil)
	}
	err := au.UpdateUserById(au.Id, au.Mobile, au.Business)
	if err != nil {
		logs.Error("failed to update user info")
		c.ErrorJson(-1, "系统错误", nil)
	}
	c.SuccessJson("更新成功")
}
