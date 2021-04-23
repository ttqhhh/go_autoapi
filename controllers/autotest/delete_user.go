package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"go_autoapi/models"
	"gopkg.in/mgo.v2"
)

type DeleteUserParam struct {
	Id int64 `json:"id" bson:"_id"`
}

// 获取用户列表 登录
func (c *AutoTestController) deleteUser() {
	au := models.AutoUser{}
	dup := UpdatUserParam{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &dup); err != nil {
		logs.Error("update param parse error", err)
		c.ErrorJson(-1, "请求参数错误", nil)
	}
	_, err := au.GetUserInfoById(dup.Id)
	if err == mgo.ErrNotFound {
		logs.Error("the user is not exist", err)
		c.ErrorJson(-1, "用户不存在", nil)
	}
	err = au.DeleteUserById(dup.Id)
	if err != nil {
		logs.Error("failed to delete user info")
		c.ErrorJson(-1, "系统错误", nil)
	}
	c.SuccessJsonWithMsg(nil, "删除成功")
}
