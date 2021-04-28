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
		c.ShowCases()
	case "show_add_case":
		c.ShowAddCase()
	case "show_edit_case":
		c.ShowEditCase()
	case "get_all_cases":
		c.GetAllCases()
	case "show_report":
		c.ShowReport()
	case "get_all_report":
		c.GetAllReport()
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
	case "del_one_case":
		c.DelCaseByID()
	case "get_service_by_business":
		c.GetServiceByBusiness()
	//case "do_test":
	//	c.performTests()
	default:
		log.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持", nil)
	}
}
