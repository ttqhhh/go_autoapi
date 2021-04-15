package controllers

import (
	"github.com/siddontang/go/log"
	"go_autoapi/libs"
)

type CaseManageController struct {
	libs.BaseController
}

func (c *CaseManageController) Get() {
	do := c.GetMethodName()
	switch do {
	case "show_cases":
		c.showCases()
	case "get_all_cases":
		c.getAllCase()
	default:
		log.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持", nil)
	}
}

func (c *CaseManageController) Post() {
	do := c.GetMethodName()
	switch do {
	case "get_one_case":
		c.GetCasesByQuery()
	case "update_one_case":
		c.updateCaseByID()
	case "add_one_case":
		c.AddOneCase()
	default:
		log.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持", nil)
	}
}
