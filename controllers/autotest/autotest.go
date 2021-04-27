package controllers

import (
	"github.com/astaxie/beego/logs"
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
	case "logout":
		c.logout()
	case "get_user_list":
		c.userList()
	case "update_user":
		c.updateUser()
	case "delete_user":
		c.deleteUser()
	case "add_business":
		c.addBusiness()
	case "get_business_list":
		c.businessList()
	case "get_all_business":
		c.allBusinessList()
	case "perform_tests":
		c.performTests()
	case "get_progress":
		c.getProcess()
	case "get_result":
		c.getResult()
	default:
		logs.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持", nil)
	}
}
