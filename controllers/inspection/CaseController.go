package inspection

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/prometheus/common/log"
	"go_autoapi/constants"
	"go_autoapi/libs"
	"go_autoapi/models"
	"go_autoapi/utils"
	"strconv"
	"strings"
	"time"
)

type CaseController struct {
	libs.BaseController
}

func (c *CaseController) Get() {
	do := c.GetMethodName()
	switch do {
	case "show_cases":
		c.ShowCases()
	//case "show_add_case":
	//	c.ShowAddCase()
	//case "show_edit_case":
	//	c.ShowEditCase()
	//case "show_copy_case":
	//	c.ShowCopyCase()
	case "get_all_cases":
		c.GetAllCases()
	//case "show_report":
	//	c.ShowReport()
	//case "get_all_report":
	//	c.GetAllReport()
	//case "get_domains":
	//	c.GetDomains()
	default:
		log.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持", nil)
	}
}

func (c *CaseController) Post() {
	do := c.GetMethodName()
	switch do {
	//case "get_one_case":
	//	c.GetCasesByQuery()
	//case "update_one_case":
	//	c.updateCaseByID()
	case "add_one_case":
		c.AddOneCase()
	//case "del_one_case":
	//	c.DelCaseByID()
	//case "get_service_by_business":
	//	c.GetServiceByBusiness()
	//case "get_caseId_by_service":
	//	c.GetCaseIdByService()
	//case "do_test":
	//	c.performTests()
	//case "add_one_domain":
	//	c.AddOneDomain()
	//case "set_inspection":
	//	c.SetInspection()
	default:
		log.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持", nil)
	}
}

/** 跳转到巡检Case列表页 */
func (c *CaseController) ShowCases() {
	business := c.GetString("business")
	c.Data["business"] = business
	c.TplName = "inspection_case_manager.html"
}

func (c *CaseController) GetAllCases() {
	acm := models.InspectionCaseMongo{}
	business := c.GetString("business")
	page, _ := strconv.Atoi(c.GetString("page"))
	limit, _ := strconv.Atoi(c.GetString("limit"))
	result, count, err := acm.GetAllCases(page, limit, business)
	if err != nil {
		c.FormErrorJson(-1, "获取测试用例列表数据失败")
	}
	c.FormSuccessJson(count, result)
}

/** 添加线上巡检Case */
func (c *CaseController) AddOneCase() {
	now := time.Now().Format(constants.TimeFormat)
	acm := models.InspectionCaseMongo{}
	dom := models.Domain{}
	if err := c.ParseForm(&acm); err != nil { //传入user指针
		c.Ctx.WriteString("出错了！")
	}
	// 获取域名并确认是否执行
	dom.Author = acm.Author
	intBus, _ := strconv.Atoi(acm.BusinessCode)
	dom.Business = int8(intBus)
	dom.DomainName = acm.Domain
	if err := dom.DomainInsert(dom); err != nil {
		logs.Error("添加case的时候 domain 插入失败")
	}
	// service_id 和 service_name 在一起,需要分割后赋值
	arr := strings.Split(acm.ServiceName, ";")
	acm.ServiceName = arr[1]
	id64, _ := strconv.ParseInt(arr[0], 10, 64)
	acm.ServiceId = id64
	//acm.Id = models.GetId("case")
	r := utils.GetRedis()
	testCaseId, err := r.Incr(constants.INSPECTION_CASE_PRIMARY_KEY).Result()
	if err != nil {
		logs.Error("保存Case时，获取从redis获取唯一主键报错，err: ", err)
		c.ErrorJson(-1, "保存Case出错啦", nil)
	}
	acm.Id = testCaseId
	acm.CreatedAt = now
	acm.UpdatedAt = now
	acm.Status = 0
	business := acm.BusinessCode
	if business == "0" {
		acm.BusinessName = "最右"
	} else if business == "1" {
		acm.BusinessName = "皮皮"
	} else if business == "2" {
		acm.BusinessName = "海外"
	} else if business == "3" {
		acm.BusinessName = "中东"
	} else if business == "4" {
		acm.BusinessName = "妈妈社区"
	} else if business == "5" {
		acm.BusinessName = "商业化"
	}
	// 去除请求路径前后的空格
	apiUrl := acm.ApiUrl
	acm.ApiUrl = strings.TrimSpace(apiUrl)
	// todo 千万不要删，用于处理json格式化问题（删了后某些服务会报504问题）
	param := acm.Parameter
	v := make(map[string]interface{})
	err = json.Unmarshal([]byte(strings.TrimSpace(param)), &v)
	if err != nil {
		logs.Error("发送冒烟请求前，解码json报错，err：", err)
		return
	}
	paramByte, err := json.Marshal(v)
	if err != nil {
		logs.Error("保存Case时，处理请求json报错， err:", err)
		c.ErrorJson(-1, "保存Case出错啦", nil)
	}
	acm.Parameter = string(paramByte)
	if err := acm.AddCase(acm); err != nil {
		logs.Error("保存Case报错，err: ", err)
		c.ErrorJson(-1, "保存Case出错啦", nil)
	}
	// 保存成功后，将该条线上巡检Case关联的测试Case巡检状态切换为开启
	testCaseMongo := models.TestCaseMongo{}
	testCaseMongo.SetInspection(acm.TestCaseId, models.INSPECTION)
	c.Ctx.Redirect(302, "/inspection/show_cases?business="+business)
}
