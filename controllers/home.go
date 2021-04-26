package controllers

import (
	"github.com/astaxie/beego/logs"
	"go_autoapi/libs"
)

type HomeController struct {
	libs.BaseController
}

func (c *HomeController) Get() {
	do := c.GetMethodName()
	switch do {
	case "goHome":
		c.goHome()
	default:
		logs.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持", nil)
	}
}

func (c *HomeController) goHome() {
	//c.TplName = "home.tpl"
	c.TplName = "login.tpl"
}
