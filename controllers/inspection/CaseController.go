package inspection

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/prometheus/common/log"
	"go_autoapi/constants"
	controllers "go_autoapi/controllers/autotest"
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

// 巡检策略常量集合，以分钟为基础单位
const (
	FIVE_MIN_CODE    = 5
	TEN_MIN_CODE     = 10
	ONE_QUARTER_CODE = 15
	HALF_HOUR_CODE   = 30
	ONE_HOUR_CODE    = 60
	SIX_HOUR_CODE    = 360
	HALF_DAY_CODE    = 720
	ONE_DAY_CODE     = 1440
)

const (
	FIVE_MIN    = "5分钟"
	TEN_MIN     = "10分钟"
	ONE_QUARTER = "15分钟"
	HALF_HOUR   = "30分钟"
	ONE_HOUR    = "1小时"
	SIX_HOUR    = "6小时"
	HALF_DAY    = "12小时"
	ONE_DAY     = "24小时"
)

const (
	//FIVE_MIN_EXPRESSION    = "0 0/5 * * * *"
	FIVE_MIN_EXPRESSION    = "0/5 * * * * *"
	TEN_MIN_EXPRESSION     = "0 0/10 * * * *"
	ONE_QUARTER_EXPRESSION = "0 0/15 * * * *"
	HALF_HOUR_EXPRESSION   = "0 0/30 * * * *"
	ONE_HOUR_EXPRESSION    = "0 0 * * * *"
	SIX_HOUR_EXPRESSION    = "0 0 0/6 * * *"
	HALF_DAY_EXPRESSION    = "0 0 0/12 * * *"
	ONE_DAY_EXPRESSION     = "0 0 0 * * *"
)

var StrategyCode2Name = map[int]string{
	FIVE_MIN_CODE:    FIVE_MIN,
	TEN_MIN_CODE:     TEN_MIN,
	ONE_QUARTER_CODE: ONE_QUARTER,
	HALF_HOUR_CODE:   HALF_HOUR,
	ONE_HOUR_CODE:    ONE_HOUR,
	SIX_HOUR_CODE:    SIX_HOUR,
	HALF_DAY_CODE:    HALF_DAY,
	ONE_DAY_CODE:     ONE_DAY,
}

var StrategyCode2Expression = map[int]string{
	FIVE_MIN_CODE:    FIVE_MIN_EXPRESSION,
	TEN_MIN_CODE:     TEN_MIN_EXPRESSION,
	ONE_QUARTER_CODE: ONE_QUARTER_EXPRESSION,
	HALF_HOUR_CODE:   HALF_HOUR_EXPRESSION,
	ONE_HOUR_CODE:    ONE_HOUR_EXPRESSION,
	SIX_HOUR_CODE:    SIX_HOUR_EXPRESSION,
	HALF_DAY_CODE:    HALF_DAY_EXPRESSION,
	ONE_DAY_CODE:     ONE_DAY_EXPRESSION,
}

func (c *CaseController) Get() {
	do := c.GetMethodName()
	switch do {
	case "show_cases":
		c.ShowCases()
	case "show_edit_case":
		c.ShowEditCase()
	case "show_case_detail":
		c.ShowCaseDeatil()
	case "get_all_cases":
		c.GetAllCases()
	case "get_all_strategy":
		c.GetAllStrategy()
	default:
		log.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持", nil)
	}
}

func (c *CaseController) Post() {
	do := c.GetMethodName()
	switch do {
	case "update_one_case":
		c.updateCaseByID()
	case "add_one_case":
		c.AddOneCase()
	case "del_one_inspection_case":
		c.DelCaseByID()
	default:
		log.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持", nil)
	}
}
func (c *CaseController) updateCaseByID() {
	icm := models.InspectionCaseMongo{}
	dom := models.Domain{}
	if err := c.ParseForm(&icm); err != nil { //传入user指针
		c.Ctx.WriteString("出错了！")
	}
	// 获取域名并确认是否执行
	dom.Author = icm.Author
	intBus, _ := strconv.Atoi(icm.BusinessCode)
	dom.Business = int8(intBus)
	dom.DomainName = icm.Domain
	if err := dom.DomainInsert(dom); err != nil {
		logs.Error("添加case的时候 domain 插入失败")
	}
	// todo service_id 和 service_name 在一起,需要分割后赋值
	arr := strings.Split(icm.ServiceName, ";")
	icm.ServiceName = arr[1]
	id64, _ := strconv.ParseInt(arr[0], 10, 64)
	icm.ServiceId = id64
	caseId := icm.Id
	business := icm.BusinessCode
	//strategy := icm.Strategy
	businessCode, _ := strconv.Atoi(business)
	businessName := controllers.GetBusinessNameByCode(businessCode)
	icm.BusinessName = businessName
	// 去除请求路径前后的空格
	apiUrl := icm.ApiUrl
	icm.ApiUrl = strings.TrimSpace(apiUrl)
	// todo 千万不要删，用于处理json格式化问题（删了后某些服务会报504问题）
	param := icm.Parameter
	v := make(map[string]interface{})
	err := json.Unmarshal([]byte(strings.TrimSpace(param)), &v)
	if err != nil {
		logs.Error("发送冒烟请求前，解码json报错，err：", err)
		return
	}
	paramByte, err := json.Marshal(v)
	if err != nil {
		logs.Error("更新Case时，处理请求json报错， err:", err)
		c.ErrorJson(-1, "保存Case出错啦", nil)
	}
	icm.Parameter = string(paramByte)
	// 查询出当前该条Case的巡检状态，并设置到将要更新的acm结构中去
	//InspectionCaseMongo := icm.GetOneCase(caseId)
	//icm.IsInspection = InspectionCaseMongo.IsInspection
	icm, err = icm.UpdateCase(caseId, icm)
	if err != nil {
		logs.Error("更新Case报错，err: ", err)
		c.ErrorJson(-1, "请求错误", nil)
	}
	//c.SuccessJson("更新成功")
	c.Ctx.Redirect(302, "/inspection/show_cases?business="+business)
}

func (c *CaseController) ShowCaseDeatil() {
	id := c.GetString("id")
	//business := c.GetString("business")
	//services := GetServiceList(business)
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		logs.Error("转换类型错误")
	}
	acm := models.InspectionCaseMongo{}
	res := acm.GetOneCase(idInt)
	c.Data["a"] = &res
	//c.Data["services"] = services
	c.TplName = "inspection_case_detail.html"
}
func (c *CaseController) ShowEditCase() {
	id := c.GetString("id")
	//business := c.GetString("business")
	//services := GetServiceList(business)
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		logs.Error("转换类型错误")
	}
	acm := models.InspectionCaseMongo{}
	res := acm.GetOneCase(idInt)
	c.Data["a"] = &res
	//c.Data["services"] = services
	c.TplName = "inspection_case_edit.html"
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
	serviceId, _ := c.GetInt64("serviceId", 0)
	uri := c.GetString("uri", "")
	strategy, _ := c.GetInt64("strategy", 0)

	result, count, err := acm.GetAllCases(page, limit, business, serviceId, uri, strategy)
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
	//if business == "0" {
	//	acm.BusinessName = "最右"
	//} else if business == "1" {
	//	acm.BusinessName = "皮皮"
	//} else if business == "2" {
	//	acm.BusinessName = "海外"
	//} else if business == "3" {
	//	acm.BusinessName = "中东"
	//} else if business == "4" {
	//	acm.BusinessName = "妈妈社区"
	//} else if business == "5" {
	//	acm.BusinessName = "商业化"
	//}
	businessCode, _ := strconv.Atoi(business)
	businessName := controllers.GetBusinessNameByCode(businessCode)
	acm.BusinessName = businessName
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
	//c.Ctx.Redirect(302, "/inspection/show_cases?business="+business)
	c.Ctx.Redirect(302, "/case/close_windows")
}

func (c *CaseController) DelCaseByID() {
	caseID := c.GetString("id")
	ac := models.InspectionCaseMongo{}
	// 先将Case表中的关联关系干掉
	id, err := strconv.Atoi(caseID)
	if err != nil {
		logs.Error("请求参数类型转换报错， err:", err)
		c.ErrorJson(-1, "请求参数转换异常", nil)
	}
	inspectionCaseMongo := ac.GetOneCase(int64(id))
	testCaseId := inspectionCaseMongo.TestCaseId
	testCaseMongo := models.TestCaseMongo{}
	testCaseMongo.SetInspection(testCaseId, models.NOT_INSPECTION)
	// 删除巡检表中的数据
	caseIDInt, err := strconv.ParseInt(caseID, 10, 64)
	if err != nil {
		logs.Error("在删除用例的时候类型转换失败")
	}
	ac.DelCase(caseIDInt)
	c.SuccessJson(nil)
}

func (c *CaseController) GetAllStrategy() {
	var result []map[string]interface{}
	for key, val := range StrategyCode2Name {
		temp := make(map[string]interface{})
		temp["code"] = key
		temp["name"] = val
		result = append(result, temp)
	}
	for i := 0; i < len(result)-1; i++ {
		for j := 0; j < len(result)-i-1; j++ {
			mappingJ := result[j]
			mappingJNext := result[j+1]
			if mappingJ["code"].(int) > mappingJNext["code"].(int) {
				temp := result[j]
				result[j] = result[j+1]
				result[j+1] = temp
			}
		}
	}
	c.SuccessJson(result)
}
