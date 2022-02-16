package case_set

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"go_autoapi/constants"
	controllers "go_autoapi/controllers/autotest"
	"go_autoapi/libs"
	"go_autoapi/models"
	"go_autoapi/utils"
	"strconv"
	"strings"
	"time"
)

type CaseSetController struct {
	libs.BaseController
}

func (c *CaseSetController) Get() {
	do := c.GetMethodName()
	switch do {
	case "index":
		c.index()
	case "page":
		c.page()
	case "getCaseSetById":
		c.getCaseSetById()
	case "copy_case_by_id":
		c.copyCaseById()
	case "get_case_by_id":
		c.getCaseById()
	default:
		logs.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持", nil)
	}
}

func (c *CaseSetController) Post() {
	do := c.GetMethodName()
	switch do {
	case "add_case_set":
		c.addCaseSet()
	case "save_edit_case_set":
		c.saveEditCaseSet()
	case "run_by_id":
		c.runById()
	case "delete_by_id":
		c.deleteById()
	case "add_set_case":
		c.addSetCase()
	case "copy_case_by_id":
		c.copyCaseById()
	case "save_edit_set_case":
		c.saveEditSetCase()
	default:
		logs.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持", nil)
	}
}

// 页面跳转
func (c *CaseSetController) index() {
	c.TplName = "case_set_page.html"
}

// CaseSet列表页-分页查询
func (c *CaseSetController) page() {
	acm := models.TestCaseMongo{}
	business := c.GetString("business")
	url := c.GetString("url")
	service := c.GetString("service")
	page, _ := strconv.Atoi(c.GetString("page"))
	limit, _ := strconv.Atoi(c.GetString("limit"))
	result, count, err := acm.GetCasesByConfusedUrl(page, limit, business, url, service)
	if err != nil {
		c.FormErrorJson(-1, "获取测试用例列表数据失败")
	}
	c.FormSuccessJson(count, result)
}

// Case集合添加
func (c *CaseSetController) addCaseSet() {
	now := time.Now().Format(constants.TimeFormat)
	acm := models.TestCaseMongo{}
	dom := models.Domain{}
	if err := c.ParseForm(&acm); err != nil { // 传入user指针
		c.Ctx.WriteString("出错了！")
	}
	// 获取域名并确认是否执行
	dom.Author = acm.Author
	intBus, _ := strconv.Atoi(acm.BusinessCode)
	dom.Business = int8(intBus)
	dom.DomainName = acm.Domain
	if err := dom.DomainInsert(dom); err != nil{
		logs.Error("添加case的时候 domain 插入失败")
	}
	// service_id 和 service_name 在一起,需要分割后赋值
	arr := strings.Split(acm.ServiceName, ";")
	acm.ServiceName = arr[1]
	id64, _ := strconv.ParseInt(arr[0], 10, 64)
	acm.ServiceId = id64
	//acm.Id = models.GetId("case")
	r := utils.GetRedis()
	testCaseId, err := r.Incr(constants.TEST_CASE_PRIMARY_KEY).Result()
	if err != nil {
		logs.Error("保存Case时，获取从redis获取唯一主键报错，err: ", err)
		c.ErrorJson(-1, "保存Case出错啦", nil)
	}
	acm.Id = testCaseId
	acm.CreatedAt = now
	acm.UpdatedAt = now
	acm.Status = 0
	business := acm.BusinessCode

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
	//c.SuccessJson("添加成功")
	c.Ctx.Redirect(302, "/case/show_cases?business="+business)
}

// 编辑后保存CaseSet
func (c *CaseSetController) saveEditCaseSet() {

	c.SuccessJson(nil)
}

// 获取指定CaseSet,初始化编辑页面（根据id）
func (c *CaseSetController) getCaseSetById() {
	id, err := c.GetInt64("id")
	if err != nil {
		logs.Warn("/service/getById接口 参数异常, err: %v", err)
		c.ErrorJson(-1, "参数异常", nil)
	}
	logs.Info("请求参数: id=%v", id)
	serviceMongo := models.ServiceMongo{}
	service, err := serviceMongo.QueryById(id)
	if err != nil {
		c.ErrorJson(-1, "服务查询数据异常", nil)
	}
	c.SuccessJson(service)
}

// 通过Id运行指定CaseSet
func (c *CaseSetController) runById() {

	c.SuccessJson(nil)
}

// 删除指定CaseSet
func (c *CaseSetController) deleteById() {

	c.SuccessJson(nil)
}

// 向CaseSet新增Case
func (c *CaseSetController) addSetCase() {

	c.SuccessJson(nil)
}

// 获取指定SetCase,初始化编辑页面（根据id）
func (c *CaseSetController) copyCaseById() {
	id, err := c.GetInt64("id")
	if err != nil {
		logs.Warn("/service/getById接口 参数异常, err: %v", err)
		c.ErrorJson(-1, "参数异常", nil)
	}
	logs.Info("请求参数: id=%v", id)
	serviceMongo := models.ServiceMongo{}
	service, err := serviceMongo.QueryById(id)
	if err != nil {
		c.ErrorJson(-1, "服务查询数据异常", nil)
	}
	c.SuccessJson(service)
}

// 获取指定SetCase,初始化编辑页面（根据id）
func (c *CaseSetController) getCaseById() {
	id, err := c.GetInt64("id")
	if err != nil {
		logs.Warn("/service/getById接口 参数异常, err: %v", err)
		c.ErrorJson(-1, "参数异常", nil)
	}
	logs.Info("请求参数: id=%v", id)
	serviceMongo := models.ServiceMongo{}
	service, err := serviceMongo.QueryById(id)
	if err != nil {
		c.ErrorJson(-1, "服务查询数据异常", nil)
	}
	c.SuccessJson(service)
}

// 编辑SetCase
func (c *CaseSetController) saveEditSetCase() {

	c.SuccessJson(nil)
}