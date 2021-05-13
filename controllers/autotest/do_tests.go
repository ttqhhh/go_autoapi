package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	uuid "github.com/satori/go.uuid"
	constant "go_autoapi/constants"
	"go_autoapi/libs"
	"go_autoapi/models"
	"go_autoapi/utils"
	"strconv"
	"sync"
	"time"
)

// 请求demo，如何传递jsonpath
//{
//"url": "http://127.0.0.1:8080/auto/login",
//"uuid": "xxxxx",
//"param": {
//"user_name": "liuweiqiang",
//"password": "OneixahDooquae1"
//},
//"check_point": {
//"$.code": {
//"eq": 200
//},
// 两个..email，会递归查找所有key=email元素，返回一个list
//"$..email": {
//"in": "liuweiqiang2014"
//}
//}
//}
type CheckOut struct {
	Url   string                            `json:"url"`
	Param map[string]interface{}            `json:"param"`
	Check map[string]map[string]interface{} `json:"check_point"`
}
type SmokeParam struct {
	ApiUrl    string `form:"api_url" json:"api_url"`
	Parameter string `form:"parameter" json:"parameter"`
}

// 接口case冒烟
func (c *AutoTestController) performSmoke() {
	param := SmokeParam{}
	if err := c.ParseForm(&param); err != nil {
		logs.Error("请求参数解析报错, err:", err)
		c.ErrorJson(-1, "请求参数错误", nil)
	}
	// todo 后续对参数进行校验
	apiUrl := param.ApiUrl
	parameter := param.Parameter

	httpStatus, body, err := libs.DoRequestWithNoneVerify(apiUrl, parameter)
	if err != nil {
		c.ErrorJson(-1, "冒烟请求内部报错", nil)
	}
	result := make(map[string]interface{})
	result["httpCode"] = httpStatus
	result["body"] = string(body)

	c.SuccessJson(result)
}

type performParam struct {
	Type        int8    `json:"type" form:"type"`             // 必填。执行Case的维度：0-业务线维度，1-服务维度，2-Case维度
	Business    int8    `json:"business" form:"business"`     // 必填
	ServiceList []int64 `json:"serviceIds" form:"serviceIds"` // 非必填。type=1时，必填
	CaseList    []int64 `json:"ids" form:"ids" `              // 非必填
}

// 执行case的维度类型，performTests接口使用
const (
	BUSINESS_TYPE = iota
	SERVICE_TYPE
	CASE_TYPE
)

// 执行case
func (c *AutoTestController) performTests() {
	userId, _ := c.GetSecureCookie(constant.CookieSecretKey, "user_id")
	uuid, _ := c.GenUUid()
	param := performParam{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &param); err != nil {
		logs.Error("请求参数解析错误， err: ", err)
		c.ErrorJson(-1, "请求参数错误", nil)
	}
	// 进行必要的参数验证
	performType := param.Type
	business := param.Business
	if performType != BUSINESS_TYPE && performType != SERVICE_TYPE && performType != CASE_TYPE {
		c.ErrorJson(-1, "请求参数中type为不支持的类型", nil)
	}
	if performType == 1 && len(param.ServiceList) == 0 {
		c.ErrorJson(-1, "以服务维度执行Case时，请至少选择一个服务", nil)
	}

	mongo := models.TestCaseMongo{}
	// 根据不同的执行维度，聚合需要执行的所有Case集合
	caseList := []*models.TestCaseMongo{}
	if performType == BUSINESS_TYPE {
		var err error
		// 查询该业务线下所有的Case
		caseList, err = mongo.GetAllCasesByBusiness(strconv.Itoa(int(business)))
		if err != nil {
			logs.Error("获取测试用例列表失败, err: ", err)
			c.ErrorJson(-1, "业务线维度执行Case时，获取测试用例失败", nil)
		}
	} else if performType == SERVICE_TYPE {
		var err error
		// 查询指定服务集合下所有的Case
		caseList, err = mongo.GetAllCasesByServiceList(param.ServiceList)
		if err != nil {
			logs.Error("获取测试用例列表失败, err: ", err)
			c.ErrorJson(-1, "服务维度执行Case时，获取测试用例失败", nil)
		}
	} else if performType == CASE_TYPE {
		var err error
		caseList, err = libs.GetCasesByIds(param.CaseList)
		if err != nil {
			logs.Error("获取测试用例列表失败, err: ", err)
			c.ErrorJson(-1, "Case维度执行Case时，获取测试用例失败", nil)
		}
	}
	count := len(caseList)
	fmt.Println("case list is", caseList)
	if len(caseList) == 0 {
		logs.Error("没有用例")
		c.ErrorJson(-1, "没有用例", nil)
	}
	// 对本次执行操作记录进行保存
	totalCases := len(caseList)
	runReport := models.RunReportMongo{}
	// 报告的名字：业务线-执行人-时间戳（日期）
	businessMap := GetBusinesses(userId)
	businessName := "未知"
	for _, v := range businessMap {
		if int8(v["code"].(int)) == business {
			businessName = v["name"].(string)
			break
		}
	}
	format := "20060102/150405"
	runReport.Name = businessName + "-" + userId + "-" + time.Now().Format(format)
	runReport.CreateBy = userId
	runReport.RunId = uuid
	runReport.TotalCases = totalCases
	runReport.IsPass = models.RUNNING
	runReport.Business = business
	runReport.ServiceName = caseList[0].ServiceName

	id, err := runReport.Insert(runReport)
	if err != nil {
		logs.Error("插入执行记录失败", err)
		c.ErrorJson(-1, "插入执行记录失败，请呼叫本平台相关负责同学", nil)
	}

	go func() {
		wg := sync.WaitGroup{}
		wg.Add(len(caseList))
		for _, val := range caseList {
			go func(url string, uuid string, param string, checkout string, caseId int64, runBy string) {
				//libs.DoRequestV2(url, uuid, val.Parameter, val.Checkpoint, val.Id) bug?
				libs.DoRequestV2(url, uuid, param, checkout, caseId, runBy)
				// 获取用例执行进度时使用
				r := utils.GetRedis()
				r.Incr(constant.RUN_RECORD_CASE_DONE_NUM + uuid)
				wg.Done()
			}(val.ApiUrl, uuid, val.Parameter, val.Checkpoint, val.Id, userId)
		}
		wg.Wait()

		go func() {
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
			runReport.UpdateIsPass(id, isPass, failCount, userId)
		}()
	}()
	c.SuccessJsonWithMsg(map[string]interface{}{"uuid": uuid, "count": count}, "OK")
}

// 执行线上巡检Case
func PerformInspection(businessId int8, serviceId int64) (err error) {
	userId := "系统定时任务"
	u2 := uuid.NewV4()
	uuid := u2.String()

	mongo := models.TestCaseMongo{}
	// 根据不同的执行维度，聚合需要执行的所有Case集合
	caseList := []*models.TestCaseMongo{}

	// 查询指定服务下所有的Case
	caseList, err = mongo.GetAllInspectionCasesByService(serviceId)
	if err != nil {
		logs.Error("获取测试用例列表失败, err: ", err)
		return
	}
	if len(caseList) == 0 {
		logs.Info("当前服务没有线上巡检用例, serviceId: ", serviceId)
		return
	}
	// 对本次执行操作记录进行保存
	totalCases := len(caseList)
	runReport := models.RunReportMongo{}
	// 报告的名字：业务线-执行人-时间戳（日期）
	businessMap := GetAllBusinesses()
	businessName := "未知"
	for _, v := range businessMap {
		if int8(v["code"].(int)) == businessId {
			businessName = v["name"].(string)
			break
		}
	}
	format := "20060102/150405"
	runReport.Name = businessName + "-" + userId + "-" + time.Now().Format(format)
	runReport.CreateBy = userId
	runReport.RunId = uuid
	runReport.TotalCases = totalCases
	runReport.IsPass = models.RUNNING
	runReport.Business = businessId
	runReport.ServiceName = caseList[0].ServiceName

	id, err := runReport.Insert(runReport)
	if err != nil {
		logs.Error("线上巡检时， 插入执行记录失败", err)
		return
	}

	go func() {
		wg := sync.WaitGroup{}
		wg.Add(len(caseList))
		for _, val := range caseList {
			go func(url string, uuid string, param string, checkout string, caseId int64, runBy string) {
				libs.DoRequestV2(url, uuid, param, checkout, caseId, runBy)
				// 获取用例执行进度时使用
				r := utils.GetRedis()
				r.Incr(constant.RUN_RECORD_CASE_DONE_NUM + uuid)
				wg.Done()
			}(val.ApiUrl, uuid, val.Parameter, val.Checkpoint, val.Id, userId)
		}
		wg.Wait()

		go func() {
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
			runReport.UpdateIsPass(id, isPass, failCount, userId)
		}()
	}()
	return
}
