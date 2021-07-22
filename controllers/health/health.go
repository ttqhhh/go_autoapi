package health

import (
	"github.com/astaxie/beego/logs"
	"go_autoapi/libs"
)

type HealthController struct {
	libs.BaseController
}

func (c *HealthController) Get() {
	do := c.GetMethodName()
	switch do {
	case "check":
		c.check()

	default:
		logs.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持", nil)
	}
}

func (c *HealthController) check() {
	c.SuccessJson(nil)
}
