package controllers

import (
	"github.com/beego/beego/v2/core/logs"
	constant "go_autoapi/constants"
	"go_autoapi/models"
	"strconv"
)

func (c *CaseManageController) ShowCases() {
	c.Data["business"] = c.GetString("business")
	c.TplName = "case_manager.html"
}

func GetServiceList(business string) (service []models.ServiceMongo) {
	var busCode = int8(0)
	if business == "0" {
		busCode = int8(0)
	} else if business == "1" {
		busCode = int8(1)
	} else if business == "2" {
		busCode = int8(2)
	} else if business == "3" {
		busCode = int8(3)
	} else if business == "4" {
		busCode = int8(4)
	} else if business == "5" {
		busCode = int8(5)
	}
	serviceMongo := models.ServiceMongo{}
	services, err := serviceMongo.QueryByBusiness(busCode)
	if err != nil {
		logs.Error("find service fail")
	}
	return services
}

func (c *CaseManageController) GetServiceByBusiness() {
	business := c.GetString("business")
	services := GetServiceList(business)
	c.SuccessJson(services)
}

func (c *CaseManageController) ShowAddCase() {
	userId, _ := c.GetSecureCookie(constant.CookieSecretKey, "user_id")
	business := c.GetString("business")
	services := GetServiceList(business)
	// 获取全部service
	c.Data["Author"] = userId
	c.Data["services"] = services
	c.TplName = "case_add.html"
}

func (c *CaseManageController) GetAllCases() {
	acm := models.TestCaseMongo{}
	//ids := models.Ids{}
	//count := ids.GetCollectionLength("case")
	business := c.GetString("business")
	page, _ := strconv.Atoi(c.GetString("page"))
	limit, _ := strconv.Atoi(c.GetString("limit"))
	result, count, err := acm.GetAllCases(page, limit, business)
	if err != nil {
		c.FormErrorJson(-1, "获取测试用例列表数据失败")
	}
	c.FormSuccessJson(count, result)
}

func (c *CaseManageController) ShowEditCase() {
	id := c.GetString("id")
	business := c.GetString("business")
	services := GetServiceList(business)
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		logs.Error("转换类型错误")
	}
	acm := models.TestCaseMongo{}
	res := acm.GetOneCase(idInt)
	c.Data["a"] = &res
	c.Data["services"] = services
	c.TplName = "case_edit.html"
}
