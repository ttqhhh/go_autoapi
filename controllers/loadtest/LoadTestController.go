package loadtest

import (
	"github.com/astaxie/beego/logs"
	"go_autoapi/libs"
)

type LoadTestController struct {
	libs.BaseController
}

func (c *LoadTestController) Get() {
	do := c.GetMethodName()
	switch do {
	case "index":
		c.index()
	default:
		logs.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持", nil)
	}
}

func (c *LoadTestController) index() {
	c.TplName = "service.html"
}