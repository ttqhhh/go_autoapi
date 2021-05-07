package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	constant "go_autoapi/constants"
	"go_autoapi/libs"
	"go_autoapi/models"
	"strconv"
	"sync"
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

type CaseList struct {
	CaseList []int64 `json:"ids" form:"ids" `
}

// 执行case
func (c *AutoTestController) performTests() {
	userId, _ := c.GetSecureCookie(constant.CookieSecretKey, "user_id")
	uuid, _ := c.GenUUid()
	cl := CaseList{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &cl); err != nil {
		c.ErrorJson(-1, "请求参数错误", nil)
	}
	count := len(cl.CaseList)
	caseList, err := libs.GetCasesByIds(cl.CaseList)
	if err != nil {
		logs.Error("获取测试用例列表失败", err)
		c.ErrorJson(-1, "获取测试用例失败", nil)
	}
	fmt.Println("case list is", caseList)
	if len(caseList) == 0 {
		logs.Error("没有用例", err)
		c.ErrorJson(-1, "没有用例", nil)
	}
	// 对本次执行操作记录进行保存
	totalCases := len(caseList)
	runReport := models.RunReportMongo{}
	runReport.CreateBy = userId
	runReport.RunId = uuid
	runReport.TotalCases = totalCases
	runReport.IsPass = models.RUNNING
	code, err := strconv.Atoi(caseList[0].BusinessCode)
	if err != nil {
		logs.Error("业务线代码转换异常", err)
		c.ErrorJson(-1, "业务线代码转换异常", nil)
	}
	runReport.Business = int8(code)
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
				r := libs.GetRedis()
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

	//go func() {
	//	r := libs.GetRedis()
	//
	//	count := 0
	//	flag := true
	//	autoResultMongo := &models.AutoResult{}
	//	for flag {
	//		hasCount, _ := r.Get(constant.RUN_RECORD_CASE_DONE_NUM +uuid).Int()
	//		if hasCount == totalCases {
	//			// 去处理本次执行记录的is_pass字段
	//			autoResult, _ := models.GetResultByRunId(uuid)
	//			// todo 因为当前只记录了失败的数据，所以当autoResult中有记录时，则认为本次执行操作失败
	//			var isPass int8 = models.SUCCESS
	//			if len(autoResult)> 0 {
	//				isPass = models.FAIL
	//			}
	//			r.Del(constant.RUN_RECORD_CASE_DONE_NUM+uuid) // 将redis中的key删除，避免存在大量无用key
	//			// 更新失败个数和本次执行记录状态
	//			failCount, _ := autoResultMongo.GetFailCount(uuid)
	//			runReport.UpdateIsPass(id, isPass, failCount, userId)
	//			flag = false
	//		}
	//		if flag {
	//			time.Sleep(1*time.Second)
	//			count++
	//			// todo 最多轮训10分钟，避免有死循环协程占用资源
	//			if count > 600 {
	//				flag = false
	//			}
	//		}
	//	}
	//}()
	c.SuccessJsonWithMsg(map[string]interface{}{"uuid": uuid, "count": count}, "OK")
}
