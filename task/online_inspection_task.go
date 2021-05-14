package task

import (
	"github.com/astaxie/beego/logs"
	uuid "github.com/satori/go.uuid"
	constant "go_autoapi/constants"
	controllers "go_autoapi/controllers/autotest"
	"go_autoapi/libs"
	"go_autoapi/models"
	"go_autoapi/utils"
	"sync"
	"time"
)

const ONLINE_INSPECTION_EXPRESSION = "0 0 9 * * ? *"

func OnlineInspection() error {
	logs.Info("启动定时任务：Online Inspection")
	// 获取所以业务线，并进行遍历所有的业务线
	businesses := controllers.GetAllBusinesses()
	for _, business := range businesses {
		// 遍历业务线下的所有所有服务
		serviceMongo := models.ServiceMongo{}
		businessId := int8(business["code"].(int))
		serviceMongos, err := serviceMongo.QueryByBusiness(businessId)
		if err != nil {
			logs.Error("执行线上巡检定时任务时，查询指定业务线下的服务时报错， err: ", err)
			return err
		}
		// 遍历服务下边所有的巡检Case
		for _, service := range serviceMongos {
			performInspection(businessId, service.Id)
		}
	}
	return nil
}

// 执行线上巡检Case
func performInspection(businessId int8, serviceId int64) (err error) {
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
	businessMap := controllers.GetAllBusinesses()
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
					// 某个服务的巡检任务存在失败Case时，认定为本次巡检任务失败，对外发送钉钉消息通知到相关同学
					logs.Warn("巡检任务失败，发送一条钉钉通知消息")
					//req := httplib.Post("http://beego.me/").Debug(true)
					//resp, err := req.Response()
					//if err != nil {
					//    logs.Error("巡检任务失败时，发送钉钉消息报错，err: ", err)
					//}
					//logs.Info("调用钉钉发送通知接口返回: ", resp)
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
