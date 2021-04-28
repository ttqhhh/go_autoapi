package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"go_autoapi/constants"
	"go_autoapi/models"
	"gopkg.in/mgo.v2"
	"time"
)

// 添加业务线
type Business struct {
	BusinessName string `json:"business_name"`
}

func (c *AutoTestController) addBusiness() {
	name, _ := c.GetSecureCookie(constants.CookieSecretKey, "user_id")
	now := time.Now().Format(constants.TimeFormat)
	id := models.GetId("business")
	b := Business{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &b); err != nil {
		logs.Error("bad param")
		c.ErrorJson(-1, "请求参数错误", nil)
	}
	ab := models.AutoBusiness{Id: id, BusinessName: b.BusinessName, CreatedAt: now, UpdatedAt: now, Author: name}
	_, err := ab.GetBusinessByName(b.BusinessName)
	if err == mgo.ErrNotFound {
		logs.Error("GetBusinessByName error", err)
		err = ab.InsertBusiness(ab)
		if err != nil {
			logs.Error("failed to add business", err)
			c.ErrorJson(-1, "系统错误", nil)
		}
	}
	c.SuccessJsonWithMsg(nil, "OK")
}
