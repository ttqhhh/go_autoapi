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
	case "get_case_set_by_id":
		c.getCaseSetById()
	case "get_set_case_by_id":
		c.getSetCaseById()
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
	case "save_edit_set_case":
		c.saveEditSetCase()
	default:
		logs.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持", nil)
	}
}

// ==================================== 用例集 接口 ==========================================

// 页面跳转 -- Done
func (c *CaseSetController) index() {
	c.TplName = "case_set_page.html"
}

// CaseSet列表页-分页查询 --Done
func (c *CaseSetController) page() {
	acm := models.CaseSetMongo{}
	business_code := c.GetString("business")
	page, _ := strconv.Atoi(c.GetString("page"))
	limit, _ := strconv.Atoi(c.GetString("limit"))

	result, count, err := acm.GetCaseSetByPage(page, limit, business_code)
	if err != nil {
		c.FormErrorJson(-1, "获取测试用例列表数据失败")
	}
	c.FormSuccessJson(count, result)
}

// Case集合添加（Form表单传参） -- Done
func (c *CaseSetController) addCaseSet() {
	now := time.Now().Format(constants.TimeFormat)
	caseSet := models.CaseSetMongo{}
	if err := c.ParseForm(&caseSet); err != nil { // 传入user指针
		c.Ctx.WriteString("出错了！")
	}

	r := utils.GetRedis()
	caseSetId, err := r.Incr(constants.CASE_SET_PRIMARY_KEY).Result()
	if err != nil {
		logs.Error("保存Case时，获取从redis获取唯一主键报错，err: ", err)
		c.ErrorJson(-1, "保存Case出错啦", nil)
	}
	caseSet.Id = caseSetId
	caseSet.CreatedAt = now
	caseSet.UpdatedAt = now
	caseSet.Status = 0
	business := caseSet.BusinessCode

	businessCode, _ := strconv.Atoi(business)
	businessName := controllers.GetBusinessNameByCode(businessCode)
	caseSet.BusinessName = businessName

	// todo 千万不要删，用于处理json格式化问题（删了后某些服务会报504问题）
	param := caseSet.Parameter
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
	caseSet.Parameter = string(paramByte)
	if err := caseSet.AddCaseSet(caseSet); err != nil {
		c.ErrorJson(-1, err.Error(), nil)
	}
	c.SuccessJson(nil)
	//c.Ctx.Redirect(302, "/case_set/page?page=1&limit=10&business="+business)
}

type runparam struct {
	id int64 `json:"id"`
}

// 通过Id运行指定CaseSet（application/json） -- Doing
func (c *CaseSetController) runById() {
	runparam := runparam{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &runparam)
	if err != nil {
		logs.Error("解析运行指定测试用例集入参报错, err: ", err)
		c.ErrorJson(-1, "请求参数错误", nil)
	}

	// todo 此处有非常非常重的逻辑，后面再写~~~
	// todo 此处有非常非常重的逻辑，后面再写~~~
	// todo 此处有非常非常重的逻辑，后面再写~~~

	c.SuccessJson(nil)
}

type delparam struct {
	id int64 `json:"id"`
}

// 删除指定CaseSet（application/json） -- Done
func (c *CaseSetController) deleteById() {
	delparam := delparam{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &delparam)
	if err != nil {
		logs.Error("解析删除指定测试用例集入参报错, err: ", err)
		c.ErrorJson(-1, "请求参数错误", nil)
	}
	caseSet := models.CaseSetMongo{}
	err = caseSet.DelCaseSet(delparam.id)
	if err != nil {
		c.ErrorJson(-1, err.Error(), nil)
	}
	c.SuccessJson(nil)
}

// 获取指定SetCase,初始化编辑页面（根据id）-- Done
func (c *CaseSetController) getCaseSetById() {
	id, err := c.GetInt64("id")
	if err != nil {
		logs.Warn("/service/getById接口 参数异常, err: %v", err)
		c.ErrorJson(-1, "参数异常", nil)
	}
	logs.Info("请求参数: id=%v", id)
	//serviceMongo := models.CaseSetMgo{}
	caseSetMongo := models.CaseSetMongo{}
	caseSet, err := caseSetMongo.CaseSetById(id)
	if err != nil {
		c.ErrorJson(-1, "服务查询数据异常", nil)
	}
	c.SuccessJson(caseSet)
}

// 编辑后保存CaseSet （Form表单传参） -- Done
func (c *CaseSetController) saveEditCaseSet() {
	csm := models.CaseSetMongo{}

	if err := c.ParseForm(&csm); err != nil { //传入user指针
		c.Ctx.WriteString("出错了！")
	}

	business := csm.BusinessCode
	businessCode, _ := strconv.Atoi(business)
	businessName := controllers.GetBusinessNameByCode(businessCode)
	csm.BusinessName = businessName

	// todo 千万不要删，用于处理json格式化问题（删了后某些服务会报504问题）
	param := csm.Parameter
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
	csm.Parameter = string(paramByte)

	csm, err = csm.UpdateCaseSet(csm.Id, csm)
	if err != nil {
		c.ErrorJson(-1, "更新测试用例集失败", nil)
	}
	c.SuccessJson(nil)
}

// ==================================== 用例 接口 ==========================================

// 获取指定CaseSet,初始化编辑页面（根据id）
func (c *CaseSetController) getSetCaseById() {
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

// 向CaseSet新增Case
func (c *CaseSetController) addSetCase() {

	c.SuccessJson(nil)
}

// 获取指定SetCase,初始化编辑页面（根据id）
//func (c *CaseSetController) copyCaseById() {
//	id, err := c.GetInt64("id")
//	if err != nil {
//		logs.Warn("/service/getById接口 参数异常, err: %v", err)
//		c.ErrorJson(-1, "参数异常", nil)
//	}
//	logs.Info("请求参数: id=%v", id)
//	serviceMongo := models.ServiceMongo{}
//	service, err := serviceMongo.QueryById(id)
//	if err != nil {
//		c.ErrorJson(-1, "服务查询数据异常", nil)
//	}
//	c.SuccessJson(service)
//}

// 编辑SetCase
func (c *CaseSetController) saveEditSetCase() {

	c.SuccessJson(nil)
}
