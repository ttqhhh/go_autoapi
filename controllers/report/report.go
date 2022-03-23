package report

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	constant "go_autoapi/constants"
	controllers "go_autoapi/controllers/autotest"
	"go_autoapi/libs"
	"go_autoapi/models"
	"strconv"
	"strings"
	"time"
)

type ReportController struct {
	libs.BaseController
}

func (c *ReportController) Get() {
	do := c.GetMethodName()
	switch do {
	case "show_run_record":
		c.ShowRunReport()
	case "run_record_list":
		c.runRecordList()
	case "show_run_record_inspection":
		c.ShowRunReportInspection()
	case "run_record_list_inspection":
		c.runRecordListInspection()
	case "run_report_detail":
		c.runReportDetail()
	default:
		logs.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持", nil)
	}
}

func (c *ReportController) ShowRunReport() {
	c.TplName = "run_report.html"
}

func (c *ReportController) ShowRunReportInspection() {
	c.TplName = "run_report_inspection.html"
}

/**
获取执行记录列表
*/
func (c *ReportController) runRecordList() {
	userId, _ := c.GetSecureCookie(constant.CookieSecretKey, "user_id")
	businessList := []int{}
	businessMap := controllers.GetBusinesses(userId)
	for _, business := range businessMap {
		for k, v := range business {
			if k == "code" {
				businessList = append(businessList, v.(int))
			}
		}
	}

	var rp = models.RunReportMongo{}
	//var ids = models.Ids{}
	page, _ := strconv.Atoi(c.GetString("page"))
	limit, _ := strconv.Atoi(c.GetString("limit"))
	//c.GetString("business")
	//c.GetString("service_name")

	//count := ids.GetCollectionLength("result")
	// 默认暂不支持business和serviceName条件查询
	result, count, err := rp.QueryByPage(businessList, "", page, limit, models.ALL, models.NOT_INSPECTION)
	if err != nil {
		c.FormErrorJson(-1, "获取报告列表数据失败")
	}
	c.FormSuccessJson(count, result)
}

func (c *ReportController) runRecordListInspection() {
	userId, _ := c.GetSecureCookie(constant.CookieSecretKey, "user_id")
	businessList := []int{}
	businessMap := controllers.GetBusinesses(userId)
	for _, business := range businessMap {
		for k, v := range business {
			if k == "code" {
				businessList = append(businessList, v.(int))
			}
		}
	}

	var rp = models.RunReportMongo{}
	//var ids = models.Ids{}
	page, _ := strconv.Atoi(c.GetString("page"))
	limit, _ := strconv.Atoi(c.GetString("limit"))
	//c.GetString("business")
	//c.GetString("service_name")

	//count := ids.GetCollectionLength("result")
	// 默认暂不支持business和serviceName条件查询
	result, count, err := rp.QueryByPage(businessList, "", page, limit, models.FAIL, models.INSPECTION)
	if err != nil {
		c.FormErrorJson(-1, "获取报告列表数据失败")
	}
	c.FormSuccessJson(count, result)
}

func (c *ReportController) runReportDetail() {
	id, _ := c.GetInt64("id")

	type TestResult struct {
		BusinessName  string   `json:"businessName"`
		ServiceName   string   `json:"serviceName"`
		CaseName      string   `json:"caseName"`
		CaseUrl       string   `json:"caseUrl"`
		CaseDetailUrl string   `json:"caseDetailUrl"`
		SpendTime     string   `json:"spendTime"`
		Status        string   `json:"status"`
		Log           []string `json:"log"`
	}

	type TemplateResp struct {
		TestPass   int          `json:"testPass"`
		TestName   string       `json:"testName"`
		TestAll    int          `json:"testAll"`
		TestFail   int          `json:"testFail"`
		BeginTime  string       `json:"beginTime"`
		TotalTime  string       `json:"totalTime"`
		TestSkip   int          `json:"testSkip"`
		TestError  int          `json:"testError"`
		TestResult []TestResult `json:"testResult"`
	}

	runReport := &models.RunReportMongo{}
	runReport, _ = runReport.QueryById(id)
	// 计算出成功总条数
	testPass := runReport.TotalCases - runReport.TotalFailCases
	// 通过修改时间-创建时间，得出本次执行操作耗时
	totalTime := ""
	createTime, _ := time.ParseInLocation(models.Time_format, runReport.CreatedAt, time.Local)
	createTimeUnix := createTime.Unix()
	updateTime, _ := time.ParseInLocation(models.Time_format, runReport.UpdatedAt, time.Local)
	updateTimeUnix := updateTime.Unix()
	totalTimeUnix := updateTimeUnix - createTimeUnix
	if totalTimeUnix < 1000 {
		totalTime = fmt.Sprintf("%d ms", totalTimeUnix)
	} else if totalTimeUnix >= 1000 && totalTimeUnix < 1000*60 {
		totalTime = fmt.Sprintf("%d s", totalTimeUnix/1000)
	} else {
		totalTime = fmt.Sprintf("%d min", totalTimeUnix/1000/60)
	}
	// 获取本次执行操作的case详情log集合
	autoResultList := []*models.AutoResult{}
	runId := runReport.RunId
	autoResultList, _ = models.GetResultByRunId(runId)
	testResultList := []TestResult{}
	inspectionCaseMongo := models.InspectionCaseMongo{}
	testCaseMongo := models.TestCaseMongo{}
	setCaseMongo := models.SetCaseMongo{}
	serNameMap := map[int64]string{}
	for _, autoResult := range autoResultList {
		result := "成功"
		if autoResult.Result == models.AUTO_RESULT_FAIL {
			result = "失败"
		}
		reasons := []string{"完美的case!!!"}

		if autoResult.Reason != "" {
			reasons = strings.Split(autoResult.Reason, ";")
		}
		if autoResult.Result != models.AUTO_RESULT_SUCCESS {
			var stat string
			stat = strconv.Itoa(autoResult.StatusCode)
			reasons = strings.Split("code:"+stat+"<br>"+autoResult.Reason+"<br>"+autoResult.Response, ";")
		}

		var testResult *TestResult
		serviceMongo := models.ServiceMongo{}
		if autoResult.IsInspection == models.INSPECTION {
			inspectionCaseMongo = inspectionCaseMongo.GetOneCase(autoResult.CaseId)
			// 根据获取服务Id去服务名
			serviceName := ""
			serviceId := inspectionCaseMongo.ServiceId
			if value, ok := serNameMap[serviceId]; ok {
				serviceName = value
			} else {
				service, err := serviceMongo.QueryById(serviceId)
				if err != nil {
					serviceName = inspectionCaseMongo.ServiceName
				} else {
					serviceName = service.ServiceName
					serNameMap[serviceId] = serviceName
				}
			}
			testResult = &TestResult{
				BusinessName:  inspectionCaseMongo.BusinessName,
				ServiceName:   serviceName,
				CaseName:      inspectionCaseMongo.CaseName,
				CaseUrl:       inspectionCaseMongo.ApiUrl,
				CaseDetailUrl: "/inspection/show_edit_case?id=" + strconv.Itoa(int(inspectionCaseMongo.Id)) + "&business=" + inspectionCaseMongo.BusinessCode,
				SpendTime:     "-",
				Status:        result,
				Log:           reasons,
			}
		} else if autoResult.IsInspection == models.NOT_INSPECTION {
			testCaseMongo = testCaseMongo.GetOneCase(autoResult.CaseId)
			// 根据获取服务Id去服务名
			serviceName := ""
			serviceId := testCaseMongo.ServiceId
			if value, ok := serNameMap[serviceId]; ok {
				serviceName = value
			} else {
				service, err := serviceMongo.QueryById(testCaseMongo.ServiceId)
				if err != nil {
					serviceName = inspectionCaseMongo.ServiceName
				} else {
					serviceName = service.ServiceName
					serNameMap[serviceId] = serviceName
				}
			}
			testResult = &TestResult{
				BusinessName:  testCaseMongo.BusinessName,
				ServiceName:   serviceName,
				CaseName:      testCaseMongo.CaseName,
				CaseUrl:       testCaseMongo.ApiUrl,
				CaseDetailUrl: "/case/show_edit_case?id=" + strconv.Itoa(int(testCaseMongo.Id)) + "&business=" + testCaseMongo.BusinessCode,
				SpendTime:     "-",
				Status:        result,
				Log:           reasons,
			}
		} else if autoResult.IsInspection == models.SENCE {
			setCaseMongo, _ := setCaseMongo.GetSetCaseById(autoResult.CaseId)
			// 根据获取服务Id去服务名
			serviceName := ""
			serviceId := setCaseMongo.ServiceId
			if value, ok := serNameMap[serviceId]; ok {
				serviceName = value
			} else {
				service, err := serviceMongo.QueryById(setCaseMongo.ServiceId)
				if err != nil {
					serviceName = inspectionCaseMongo.ServiceName
				} else {
					serviceName = service.ServiceName
					serNameMap[serviceId] = serviceName
				}
			}
			testResult = &TestResult{
				BusinessName: setCaseMongo.BusinessName,
				ServiceName:  serviceName,
				CaseName:     setCaseMongo.CaseName,
				CaseUrl:      setCaseMongo.ApiUrl,
				//CaseDetailUrl: "/case_set/get_set_case_by_id?id="+strconv.Itoa(int(setCaseMongo.Id)),
				CaseDetailUrl: "/case_set/one_case?business=" + setCaseMongo.BusinessCode + "&id=" + strconv.Itoa(int(setCaseMongo.CaseSetId)),
				SpendTime:     "-",
				Status:        result,
				Log:           reasons,
			}

		}
		testResultList = append(testResultList, *testResult)
	}

	resp := &TemplateResp{
		TestPass:   testPass,
		TestName:   runReport.Name,
		TestAll:    runReport.TotalCases,
		TestFail:   runReport.TotalFailCases,
		BeginTime:  runReport.CreatedAt,
		TotalTime:  totalTime,
		TestSkip:   0,
		TestError:  0,
		TestResult: testResultList,
	}

	result := make(map[interface{}]interface{})
	result["resultData"] = resp

	c.TplName = "template.html"
	c.Data = result
}
