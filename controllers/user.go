package controllers

import (
	"fmt"
	"github.com/siddontang/go/log"
	"go_autoapi/libs"
	"go_autoapi/models"
)

type UserController struct {
	libs.BaseController
}

func (c *UserController) Get() {
	do := c.GetMethodName()
	switch do {
	case "profile":
		c.profile()
	case "add":
		c.add()
	case "add_mgo":
		c.add_mgo()
	case "get_mgo":
		c.get_mgo()
	default:
		log.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持", nil)
	}
}
func (c *UserController) add() {
	fmt.Println("xxxxxx", "add", c.Ctx.Request.PostForm)
	amc := models.AdMockCase{}
	data, _ := amc.QueryByUUid("5eb62275-9818-4101-a477-6fef0bb9c7bd")
	fmt.Println(&data)
	c.SuccessJsonWithMsg(&data, "OK")
}
