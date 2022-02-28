package case_set

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	jsonpath "github.com/spyzhov/ajson"
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
	case "one_case":
		c.oneCase()
	case "page":
		c.page()
	case "get_case_set_by_id":
		c.getCaseSetById()
	case "get_set_case_by_id":
		c.getSetCaseById()
	case "get_set_case_list_by_case_set_id":
		c.getSetCaseListByCaseSetId()
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
	case "delete_set_case_by_id":
		c.DelSetCaseByID()
	default:
		logs.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持", nil)
	}
}

// ==================================== 用例集 接口 ==========================================

// 页面跳转 -- Done
func (c *CaseSetController) index() {
	business := c.GetString("business")
	c.Data["business"] = business
	c.TplName = "case_set_page.html"
}

// CaseSet列表页-分页查询 --Done

func (c *CaseSetController) oneCase() {
	business := c.GetString("business")
	id, err := c.GetInt64("id")
	if err != nil {
		logs.Error("从前台获取数据id出错，err", err)
	}
	c.Data["id"] = id
	c.Data["business"] = business
	c.TplName = "case_one_set.html"
}

// CaseSet列表页-分页查询

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
	userId, _ := c.GetSecureCookie(constants.CookieSecretKey, "user_id")
	//todo 获取author
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
	//todo xueyibing 返回增加author字段
	caseSet.Author = userId
	caseSet.Id = caseSetId
	caseSet.CreatedAt = now
	caseSet.UpdatedAt = now
	caseSet.Status = 0
	business := caseSet.BusinessCode

	businessCode, _ := strconv.Atoi(business)
	businessName := controllers.GetBusinessNameByCode(businessCode)
	caseSet.BusinessName = businessName
	//
	//// todo 千万不要删，用于处理json格式化问题（删了后某些服务会报504问题）
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
	c.Ctx.Redirect(302, "/case_set/index?business="+business)
}

type runparam struct {
	Id       int64 `json:"id"`                       // 必填
	Business int8  `json:"business" form:"business"` // 必填
}

// 通过Id运行指定CaseSet（application/json） -- Doing
func (c *CaseSetController) runById() {
	runBy, _ := c.GetSecureCookie(constants.CookieSecretKey, "user_id")
	runparam := runparam{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &runparam)
	if err != nil {
		logs.Error("解析运行指定测试用例集入参报错, err: ", err)
		c.ErrorJson(-1, "请求参数错误", nil)
	}

	// 1、获取该测试用例集的相关数据
	caseSetMongo := models.CaseSetMongo{}
	caseSet, err := caseSetMongo.CaseSetById(runparam.Id)
	if err != nil {
		c.ErrorJson(-1, err.Error(), nil)
	}
	// 将CaseSet中的公共参数，读取至当前协程内存中，后续继续加入响应中提取的值，并且依据其中的值替换caseParam中的参数值
	setParam := caseSet.Parameter
	setParamMap := map[string]interface{}{}
	err = json.Unmarshal([]byte(setParam), &setParamMap)
	if err != nil {
		logs.Error(-1, "解析测试用例集中的公共参数报错, err: ", err)
		c.ErrorJson(-1, err.Error(), nil)
	}

	setCaseMongo := models.SetCaseMongo{}
	setCaseList, err := setCaseMongo.GetSetCaseListByCaseSetId(runparam.Id)
	if err != nil {
		c.ErrorJson(-1, err.Error(), nil)
	}

	// 2、新建执行记录并入库
	uuid, _ := c.GenUUid()
	// 报告的名字：业务线-执行人-时间戳（日期）
	business := runparam.Business
	businessMap := controllers.GetBusinesses(runBy)
	businessName := "未知"
	for _, v := range businessMap {
		if int8(v["code"].(int)) == business {
			businessName = v["name"].(string)
			break
		}
	}
	runReport := models.RunReportMongo{}
	format := "20060102/150405"
	runReport.Name = businessName + "-" + runBy + "-" + time.Now().Format(format)
	runReport.CreateBy = runBy
	runReport.RunId = uuid
	runReport.IsPass = models.RUNNING
	runReport.Business = business
	// 设置报告中的用例总条数
	runReport.TotalCases = len(setCaseList)
	runReport.ServiceName = "【场景测试】: " + caseSet.CaseSetName
	id, err := runReport.Insert(runReport)
	if err != nil {
		logs.Error("插入执行记录失败", err)
		c.ErrorJson(-1, "插入执行记录失败，请呼叫本平台相关负责同学", nil)
	}

	// todo 核心逻辑
	// 起一个协程异步去串行执行CaseSet中的Case
	go func(runReportId int64) {
		// 3、运行测试用例集
		for _, setCase := range setCaseList {
			caseParam := setCase.Parameter

			// 从caseParam中，取出带有$字符的参数进行替换
			// todo 现阶段仅支持json为深度为1的参数值替换
			caseParamMap := map[string]interface{}{}
			json.Unmarshal([]byte(caseParam), &caseParamMap)
			for key, value := range caseParamMap {
				strValue, ok := value.(string)
				// 只有当value为字符串时，才考虑进行参数值替换
				if ok {
					// 当前value为"${"开头，且为"}"结尾
					if strings.HasPrefix(strValue, "${") && strings.HasSuffix(strValue, "}") {
						valueInSetParamMap, ok := setParamMap[key]
						// 当公共参数setParamMap中存在要替换的key时，进行替换；不存在时，
						if ok {
							caseParamMap[key] = valueInSetParamMap
						} else {
							// todo 公共参数中不存在该key
							reason := "公共参数中未找到指定的key, key=" + key
							libs.SaveTestResult(uuid, setCase.Id, models.NOT_INSPECTION, models.AUTO_RESULT_FAIL, reason, runBy, "", 0)
							break
						}
					}
				}
			}
			caseParamStr, err := json.Marshal(caseParamMap)
			if err != nil {
				logs.Error("场景自动化测试时, 参数替换后json字符串转换报错, err: ", err)
				reason := "场景自动化测试时, 参数替换后json字符串转换报错"
				// statusCode 为0时，表示为发送请求，前置校验逻辑直接未通过。
				libs.SaveTestResult(uuid, setCase.Id, models.NOT_INSPECTION, models.AUTO_RESULT_FAIL, reason, runBy, "", 0)
				break
			}
			caseSet.Parameter = string(caseParamStr)

			// case执行
			isOk, resp := libs.DoRequest(setCase.Domain, setCase.ApiUrl, uuid, setCase.Parameter, setCase.Checkpoint, setCase.Id, models.INSPECTION, runBy)

			// 当Case集合中某条Case不通过时，不再继续往下执行该场景测试
			if !isOk {
				break
			}
			// 通过jsonpath路径去响应中提取值，并放入setParamMap公共参数中
			extractRespMap := map[string]string{}
			err = json.Unmarshal([]byte(setCase.ExtractResp), &extractRespMap)
			if err != nil {
				logs.Error("场景自动化测试时，从响应中提取数据的配置转换json报错, err: ", err)
				reason := "场景自动化测试时，从响应中提取数据的配置转换json报错"
				// statusCode 为0时，表示为发送请求，前置校验逻辑直接未通过。
				libs.SaveTestResult(uuid, setCase.Id, models.NOT_INSPECTION, models.AUTO_RESULT_FAIL, reason, runBy, "", 0)
				break
			}
			// value为jsonpath
			for key, value := range extractRespMap {
				valueInResp, err := jsonpath.JSONPath([]byte(resp), value)
				if err != nil {
					logs.Error("根据jsonpath从响应中提取value时报错, err: ", err)
					reason := "根据jsonpath从响应中提取value时报错"
					// statusCode 为0时，表示为发送请求，前置校验逻辑直接未通过。
					libs.SaveTestResult(uuid, setCase.Id, models.NOT_INSPECTION, models.AUTO_RESULT_FAIL, reason, runBy, "", 0)
					break
				}
				// 将提取出来的值，放入setParamMap公共区域，提供后续接口使用
				setParamMap[key] = valueInResp
			}

		}

		// 4、执行记录结果状态处理
		autoResult, _ := models.GetResultByRunId(uuid)
		var isPass int8 = models.SUCCESS
		// 判断case执行结果集合中是否有失败的case，有则认为本次执行操作状态为FAIL
		for _, result := range autoResult {
			if result.Result == models.AUTO_RESULT_FAIL {
				isPass = models.FAIL
				break
			}
		}
		// 更新失败个数和本次执行记录状态
		autoResultMongo := &models.AutoResult{}
		failCount, _ := autoResultMongo.GetFailCount(uuid)
		runReport.UpdateIsPass(runReportId, isPass, failCount, runBy)
	}(id)

	c.SuccessJson(nil)
}

type delparam struct {
	Id int64 `json:"id"`
}

// 删除指定CaseSet（application/json） -- Done
func (c *CaseSetController) deleteById() {
	ids := c.GetString("id")
	id, err := strconv.ParseInt(ids, 10, 64)
	if err != nil {
		logs.Error("获取caseset id 出错，err:", err)
	}
	//delparam := delparam{}
	//err := json.Unmarshal(c.Ctx.Input.RequestBody, &delparam)
	//if err != nil {
	//	logs.Error("解析删除指定测试用例集入参报错, err: ", err)
	//	c.ErrorJson(-1, "请求参数错误", nil)
	//}
	caseSet := models.CaseSetMongo{}
	err = caseSet.DelCaseSet(id)
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
	name, _ := c.GetSecureCookie(constants.CookieSecretKey, "user_id")

	if err := c.ParseForm(&csm); err != nil { //传入user指针
		c.Ctx.WriteString("出错了！")
	}
	csm.Author = name
	business := csm.BusinessCode
	businessCode, _ := strconv.Atoi(business)
	businessName := controllers.GetBusinessNameByCode(businessCode)
	csm.BusinessName = businessName

	// todo 千万不要删，用于处理json格式化问题（删了后某些服务会报504问题）
	//param := csm.Parameter
	//v := make(map[string]interface{})
	//err := json.Unmarshal([]byte(strings.TrimSpace(param)), &v)
	//if err != nil {
	//	logs.Error("发送冒烟请求前，解码json报错，err：", err)
	//	return
	//}
	//paramByte, err := json.Marshal(v)
	//if err != nil {
	//	logs.Error("更新Case时，处理请求json报错， err:", err)
	//	c.ErrorJson(-1, "保存Case出错啦", nil)
	//}
	//csm.Parameter = string(paramByte)
	//todo 暂时不放公共参数
	csm, err := csm.UpdateCaseSet(csm.Id, csm)
	if err != nil {
		c.ErrorJson(-1, "更新测试用例集失败", nil)
	}
	c.Ctx.Redirect(302, "/case_set/index?business="+business)

}

// ==================================== 用例 接口 ==========================================
// 源Case筛选接口: /case/get_case_by_condition
// 筛选出来源Case后，调起编辑源Case的页面接口为: /case/show_copy_case?id=750&business=0

// 从一条caseset中获取全部case // todo xueyibing 分页
func (c *CaseSetController) getSetCaseListByCaseSetId() {
	caseSetId, err := c.GetInt64("case_set_id")
	if err != nil {
		logs.Error("/case_set/getSetCaseByCaseSetId接口 参数异常, err: %v", err)
		c.ErrorJson(-1, "参数异常", nil)
	}

	setCaseMongo := models.SetCaseMongo{}
	caseSet, err := setCaseMongo.GetSetCaseListByCaseSetId(caseSetId)
	if err != nil {
		c.ErrorJson(-1, "服务查询数据异常", nil)
	}
	count := len(caseSet)
	c.FormSuccessJson(int64(count), caseSet)
}

// 向CaseSet新增Case
func (c *CaseSetController) addSetCase() {
	userId, _ := c.GetSecureCookie(constants.CookieSecretKey, "user_id")
	now := time.Now().Format(constants.TimeFormat)
	scm := models.SetCaseMongo{}
	dom := models.Domain{}
	if err := c.ParseForm(&scm); err != nil { // 传入user指针
		c.Ctx.WriteString("出错了！")
	}
	// 获取域名并确认是否执行
	dom.Author = userId
	intBus, _ := strconv.Atoi(scm.BusinessCode)
	dom.Business = int8(intBus)
	dom.DomainName = scm.Domain
	if err := dom.DomainInsert(dom); err != nil {
		logs.Error("添加case的时候 domain 插入失败")
	}
	// service_id 和 service_name 在一起,需要分割后赋值
	arr := strings.Split(scm.ServiceName, ";")
	scm.ServiceName = arr[1]
	id64, _ := strconv.ParseInt(arr[0], 10, 64)
	scm.ServiceId = id64
	r := utils.GetRedis()
	testCaseId, err := r.Incr(constants.TEST_CASE_PRIMARY_KEY).Result()
	if err != nil {
		logs.Error("保存Case时，获取从redis获取唯一主键报错，err: ", err)
		c.ErrorJson(-1, "保存Case出错啦", nil)
	}
	scm.Id = testCaseId
	scm.CreatedAt = now
	scm.UpdatedAt = now
	scm.Status = 0
	business := scm.BusinessCode
	businessCode, _ := strconv.Atoi(business)
	businessName := controllers.GetBusinessNameByCode(businessCode)
	scm.BusinessName = businessName
	// 去除请求路径前后的空格
	apiUrl := scm.ApiUrl
	scm.ApiUrl = strings.TrimSpace(apiUrl)
	// todo 千万不要删，用于处理json格式化问题（删了后某些服务会报504问题）
	param := scm.Parameter
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
	scm.Parameter = string(paramByte)
	if err := scm.AddSetCase(scm); err != nil {
		c.ErrorJson(-1, err.Error(), nil)
	}
	//c.SuccessJson("添加成功")
	c.Ctx.Redirect(302, "/case_set/one_case?business="+business+"&id="+strconv.FormatInt(scm.CaseSetId, 10))
}

func (c *CaseSetController) DelSetCaseByID() {
	//delparam := delparam{}
	//err := json.Unmarshal(c.Ctx.Input.RequestBody, &delparam)
	//if err != nil {
	//	logs.Error("解析删除指定测试用例集入参报错, err: ", err)
	//	c.ErrorJson(-1, "请求参数错误", nil)
	//}
	ids := c.GetString("id")
	id, err := strconv.ParseInt(ids, 10, 64)
	if err != nil {
		logs.Error("获取set_case id 出错，err:", err)
	}
	scm := models.SetCaseMongo{}
	err = scm.DelSetCase(id)
	if err != nil {
		c.ErrorJson(-1, "删除用例失败", nil)
	}
	c.SuccessJson(nil)
}

// 获取指定CaseSet,初始化编辑页面（根据id）
func (c *CaseSetController) getSetCaseById() {
	id, err := c.GetInt64("id")
	if err != nil {
		logs.Error("/case_set/getSetCaseById接口 参数异常, err: %v", err)
		c.ErrorJson(-1, "参数异常", nil)
	}

	scm := models.SetCaseMongo{}
	mongo, err := scm.GetSetCaseById(id)
	if err != nil {
		c.ErrorJson(-1, "服务查询数据异常", nil)
	}
	c.SuccessJson(mongo)
}

// 编辑SetCase
func (c *CaseSetController) saveEditSetCase() {
	scm := models.SetCaseMongo{}
	dom := models.Domain{}
	if err := c.ParseForm(&scm); err != nil { //传入user指针
		c.Ctx.WriteString("出错了！")
	}
	// 获取域名并确认是否执行
	dom.Author = scm.Author
	intBus, _ := strconv.Atoi(scm.BusinessCode)
	dom.Business = int8(intBus)
	dom.DomainName = scm.Domain
	if err := dom.DomainInsert(dom); err != nil {
		logs.Error("添加case的时候 domain 插入失败")
	}
	// todo service_id 和 service_name 在一起,需要分割后赋值
	arr := strings.Split(scm.ServiceName, ";")
	scm.ServiceName = arr[1]
	id64, _ := strconv.ParseInt(arr[0], 10, 64)
	scm.ServiceId = id64
	caseId := scm.Id
	business := scm.BusinessCode
	businessCode, _ := strconv.Atoi(business)
	businessName := controllers.GetBusinessNameByCode(businessCode)
	scm.BusinessName = businessName
	// 去除请求路径前后的空格
	apiUrl := scm.ApiUrl
	scm.ApiUrl = strings.TrimSpace(apiUrl)
	// todo 千万不要删，用于处理json格式化问题（删了后某些服务会报504问题）
	param := scm.Parameter
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
	scm.Parameter = string(paramByte)
	scm, err = scm.UpdateSetCase(caseId, scm)
	if err != nil {
		c.ErrorJson(-1, err.Error(), nil)
	}
	//c.SuccessJson("更新成功")
	//c.Ctx.Redirect(302, "/case/show_cases?business="+business)
	c.Ctx.Redirect(302, "/case/close_windows")
}
