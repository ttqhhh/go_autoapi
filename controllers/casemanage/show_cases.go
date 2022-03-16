package controllers

import (
	"github.com/beego/beego/v2/core/logs"
	constant "go_autoapi/constants"
	controllers "go_autoapi/controllers/autotest"
	"go_autoapi/models"
	"strconv"
)

func (c *CaseManageController) ShowCases() {
	business := c.GetString("business")
	c.Data["business"] = business
	c.TplName = "case_manager.html"
}

//func GetServiceList(business string) (service []models.ServiceMongo) {
//	bs, err := strconv.Atoi(business)
//	if err != nil{
//		logs.Error("类型转换失败", err)
//	}
//	busCode := int8(bs)
//	serviceMongo := models.ServiceMongo{}
//	services, err := serviceMongo.QueryByBusiness(busCode)
//	if err != nil {
//		logs.Error("find service fail")
//	}
//	return services
//}

//func (c *CaseManageController) GetServiceByBusiness() {
//	business := c.GetString("business")
//	services := GetServiceList(business)
//	c.SuccessJson(services)
//}

func (c *CaseManageController) ShowAddCase() {
	userId, _ := c.GetSecureCookie(constant.CookieSecretKey, "user_id")
	business := c.GetString("business")
	//services := GetServiceList(business)
	// 获取全部service
	c.Data["Author"] = userId
	c.Data["business"] = business
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

func (c *CaseManageController) SearchCase() {
	acm := models.TestCaseMongo{}
	business := c.GetString("business")
	url := c.GetString("url")
	service := c.GetString("service", "-1")
	author := c.GetString("author")
	case_name := c.GetString("case_name")
	id := c.GetString("case_id", "-1")
	caseId, err := strconv.Atoi(id)
	if err != nil {
		logs.Error("自动化Case分页接口报错, err: ", err)
		c.ErrorJson(-1, err.Error(), nil)
	}
	serviceId, _ := strconv.Atoi(service)
	page, _ := strconv.Atoi(c.GetString("page"))
	limit, _ := strconv.Atoi(c.GetString("limit"))
	result, count, err := acm.GetCasesByConfusedUrl(page, limit, business, url, serviceId, caseId, case_name, author)
	if err != nil {
		c.FormErrorJson(-1, "获取测试用例列表数据失败")
	}
	c.FormSuccessJson(count, result)
}

func (c *CaseManageController) ShowEditCase() {
	id := c.GetString("id")
	//business := c.GetString("business")
	//services := GetServiceList(business)
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		logs.Error("转换类型错误")
	}
	acm := models.TestCaseMongo{}
	res := acm.GetOneCase(idInt)
	c.Data["a"] = &res
	//c.Data["services"] = services
	c.TplName = "case_edit.html"
}

func (c *CaseManageController) ShowCopyCase() {
	id := c.GetString("id")
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		logs.Error("转换类型错误")
	}
	acm := models.TestCaseMongo{}
	res := acm.GetOneCase(idInt)
	c.Data["a"] = &res
	//c.Data["services"] = services
	c.TplName = "case_copy.html"
}

func (c *CaseManageController) ShowCopyCaseSet() {
	id := c.GetString("id")
	caseSetId, err := c.GetInt64("case_set_id", 10)
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		logs.Error("转换类型错误")
	}
	acm := models.TestCaseMongo{}
	res := acm.GetOneCase(idInt)
	c.Data["a"] = &res
	c.Data["case_set_id"] = caseSetId
	c.TplName = "case_copy_new.html"
}

// 将填写的自动化case转换成场景测试步骤
func (c *CaseManageController) TransferToCaseSetStep() {
	//id := c.GetString("id")
	caseSetId, err := c.GetInt64("case_set_id")
	if err != nil {
		logs.Error("参数错误")
		c.ErrorJson(-1, "参数错误", nil)
	}
	businessCode := c.GetString("business_code")
	serviceId := c.GetString("service_id")
	author := c.GetString("author")
	smokeResponse := c.GetString("smoke_response")
	caseName := c.GetString("case_name")
	description := c.GetString("description")
	requestMethod := c.GetString("request_method")
	domain := c.GetString("domain")
	apiUrl := c.GetString("api_url")
	parameter := c.GetString("parameter")
	//responseParam := c.GetString("response_param")
	checkPoint := c.GetString("check_point")

	/**
	let business_code = $("[name='business_code']").val()
	let service_name = $("[name='service_name']").val()
	let author = $("[name='author']").val()
	let smoke_response = $("[name='smoke_response']").val()
	let case_name = $("[name='case_name']").val()
	let description = $("[name='description']").val()
	let request_method = $("[name='request_method']").val()
	let domain = $("[name='domain']").val()
	let api_url = $("[name='api_url']").val()
	let parameter = $("[name='parameter']").val()
	let response_param = $("[name='response_param']").val()
	let check_point = $("[name='check_point']").val()
	*/

	res := models.TestCaseMongo{}
	//res.BusinessCode = businessCode
	businessCodeInt, err := strconv.Atoi(businessCode)
	if err != nil {
		logs.Error("TransferToCaseSetStep方法字符串转int报错, err: ", err)
	}
	res.BusinessName = controllers.GetBusinessNameByCode(businessCodeInt)
	serviceIdInt64, err := strconv.ParseInt(serviceId, 10, 64)
	if err != nil {
		logs.Error("TransferToCaseSetStep方法字符串转int64报错, err: ", err)
	}
	res.ServiceId = serviceIdInt64
	res.Author = author
	res.SmokeResponse = smokeResponse
	res.CaseName = caseName
	res.Description = description
	res.RequestMethod = requestMethod
	res.Domain = domain
	res.ApiUrl = apiUrl
	res.Parameter = parameter
	res.Checkpoint = checkPoint
	c.Data["a"] = &res
	c.Data["case_set_id"] = caseSetId
	c.TplName = "case_copy_new.html"
}

func (c *CaseManageController) ShowInspectionCase() {
	id := c.GetString("id")
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		logs.Error("转换类型错误")
	}
	acm := models.TestCaseMongo{}
	res := acm.GetOneCase(idInt)
	c.Data["a"] = &res
	//c.Data["services"] = services
	c.TplName = "case_inspection.html"
}
