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
	case "get_domains":
		c.GetDomains()
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
	case "add_one_domain":
		c.AddOneDomain()
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

// todo 调试用的domain插入接口目前已关闭
func (c *CaseManageController) AddOneDomain() {
	//c.SuccessJson("domain 调回接口关闭")
	Dom := models.Domain{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &Dom); err != nil {
		c.ErrorJson(-1, "请求错误", nil)
	}
	err := Dom.DomainInsert(Dom)
	if err !=nil{
		c.ErrorJson(-1,"插入域名失败",err)
	}
	c.SuccessJson("成功插入域名数据")
}

func (c *CaseManageController) GetDomains(){
	business ,err:= c.GetInt8("business")
	if err !=nil{
		logs.Error("获取域名的business可能不是int8类型",err)
		c.ErrorJson(-1,"获取域名的business可能不是int8类型",nil)
	}
	Dom := models.Domain{}
	domains, err := Dom.GetDomainByBusiness(business)
	if err != nil{
		logs.Error("获取domains失败")
	}
	c.SuccessJson(domains)
}














