package report

import (
	"fmt"
	"github.com/astaxie/beego/logs"
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

/**
获取执行记录列表
*/
func (c *ReportController) runRecordList() {
	var rp = models.RunReportMongo{}
	//var ids = models.Ids{}
	page, _ := strconv.Atoi(c.GetString("page"))
	limit, _ := strconv.Atoi(c.GetString("limit"))
	//c.GetString("business")
	//c.GetString("service_name")

	//count := ids.GetCollectionLength("result")
	// 默认暂不支持business和serviceName条件查询
	result, count, err := rp.QueryByPage(-1, "", page, limit)
	if err != nil {
		c.FormErrorJson(-1, "获取报告列表数据失败")
	}
	c.FormSuccessJson(count, result)
}

func (c *ReportController) runReportDetail() {
	id, _ := c.GetInt64("id")

	type TestResult struct {
		BusinessName string   `json:"businessName"`
		ServiceName  string   `json:"serviceName"`
		CaseName     string   `json:"caseName"`
		SpendTime    string   `json:"spendTime"`
		Status       string   `json:"status"`
		Log          []string `json:"log"`
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
	testCaseMongo := models.TestCaseMongo{}
	for _, autoResult := range autoResultList {
		result := "成功"
		if autoResult.Result == models.AUTO_RESULT_FAIL {
			result = "失败"
		}
		reasons := []string{"完美Case!!!"}
		if autoResult.Reason != "" {
			reasons = strings.Split(autoResult.Reason, ";")
		}
		testCaseMongo = testCaseMongo.GetOneCase(autoResult.CaseId)

		testResult := &TestResult{
			BusinessName: testCaseMongo.BusinessName,
			ServiceName:  testCaseMongo.ServiceName,
			CaseName:     testCaseMongo.CaseName,
			SpendTime:    "-",
			Status:       result,
			Log:          reasons,
		}

		testResultList = append(testResultList, *testResult)
	}

	resp := &TemplateResp{
		TestPass:   testPass,
		TestName:   runReport.RunId,
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
