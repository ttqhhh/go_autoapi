package controllers

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/siddontang/go/log"
	"go_autoapi/libs"
	"go_autoapi/models"
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
	case "show_copy_case":
		c.ShowCopyCase()
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
	//case "get_service_by_business":
	//	c.GetServiceByBusiness()
	case "get_caseId_by_service":
		c.GetCaseIdByService()
	//case "do_test":
	//	c.performTests()
	case "set_inspection":
		c.SetInspection()
	default:
		log.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持", nil)
	}
}

func (c *CaseManageController) SetInspection() {
	//userId, _ := c.GetSecureCookie(constant.CookieSecretKey, "user_id")
	type paramStruct struct {
		Id            int64
		Is_inspection int8
	}
	param := &paramStruct{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, param)
	if err != nil {
		logs.Error("Case设置巡检状态接口解析参数错误， err: ", err)
		c.ErrorJson(-1, "参数解析错误", nil)
	}
	// 巡检字段参数值有效性验证
	if param.Is_inspection != models.NOT_INSPECTION && param.Is_inspection != models.INSPECTION {
		c.ErrorJson(-1, "不支持的请求参数值", nil)
	}
	model := &models.TestCaseMongo{}
	err = model.SetInspection(param.Id, param.Is_inspection)
	if err != nil {
		c.ErrorJson(-1, "设置接口为巡检状态出错", nil)
	}
	c.SuccessJson(nil)
}
