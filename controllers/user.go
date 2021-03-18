package controllers

import (
	"fmt"
	"github.com/siddontang/go/log"
	"go_autoapi/models"
)

type UserController struct {
	BaseController
}

func (c *UserController) Get() {
	do := c.GetMethodName()
	switch do {
	case "profile":
		c.profile()
	case "add":
		c.add()
	default:
		log.Warn("action: %s, not implemented", do)
	}
}
func (c *UserController) add() {
	fmt.Println("xxxxxx", "add", c.Ctx.Request.PostForm)
	amc := models.AdMockCase{}
	data, _ := amc.QueryByUUid("5eb62275-9818-4101-a477-6fef0bb9c7bd")
	fmt.Println(&data)
	c.SuccessJson(&data)
}
