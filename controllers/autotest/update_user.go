package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"go_autoapi/models"
	"gopkg.in/mgo.v2"
)

type UpdatUserParam struct {
	Id     int64  `json:"id" bson:"_id"`
	Mobile string `json:"mobile" bson:"mobile"`
	// 0：最右，1：皮皮，2：海外，3：中东，4：妈妈
	Business int `json:"business" bson:"business" valid:"Range(0, 4)"`
	//0：正常，1：删除
}

// 获取用户列表 登录
func (c *AutoTestController) updateUser() {
	uup := UpdatUserParam{}
	au := models.AutoUser{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &uup); err != nil {
		logs.Error("update param parse error", err)
		c.ErrorJson(-1, "请求参数错误", nil)
	}
	_, err := au.GetUserInfoById(uup.Id)
	if err == mgo.ErrNotFound {
		logs.Error("the user is not exist", err)
		c.ErrorJson(-1, "用户不存在", nil)
	}
	err = au.UpdateUserById(uup.Id, uup.Mobile, uup.Business)
	if err != nil {
		logs.Error("failed to update user info")
		c.ErrorJson(-1, "系统错误", nil)
	}
	c.SuccessJsonWithMsg(nil, "更新成功")
}
