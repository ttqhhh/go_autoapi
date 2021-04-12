package controllers

import (
	"github.com/siddontang/go/log"
	"go_autoapi/libs"
)

type AutoTestController struct {
	libs.BaseController
}

func (c *AutoTestController) Post() {
	do := c.GetMethodName()
	switch do {
	case "get_case_info":
		c.getCaseInfo()
	case "add_case":
		c.addCase()
	case "update_case":
		c.updateCaseInfo()
	case "login":
		c.login()
	default:
		log.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持", nil)
	}
}
