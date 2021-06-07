package inspection_strategy

import (
	"fmt"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
	uuid "github.com/satori/go.uuid"
	constant "go_autoapi/constants"
	controllers "go_autoapi/controllers/autotest"
	"go_autoapi/libs"
	"go_autoapi/models"
	"go_autoapi/utils"
	"io/ioutil"
	"sync"
	"time"
)

// 「测试效率团队」群web_hook-用来测试
const XIAO_NENG_QUN = "https://oapi.dingtalk.com/robot/send?access_token=6f35268d9dcb74b4b95dd338eb241832781aeaaeafd90aa947b86936f3343dbb"

// todo 临时关闭状态，上线or正常使用时，需要设为true进行开启
const IS_OPEN_SENDDING_MSG = false

// 执行线上巡检Case
func PerformInspection(businessId int8, serviceId int64, msgChannel chan string, strategy int64) (err error) {
	userId := "线上巡检"
	u2 := uuid.NewV4()
	uuid := u2.String()

	mongo := models.InspectionCaseMongo{}
	// 根据不同的执行维度，聚合需要执行的所有Case集合
	caseList := []*models.InspectionCaseMongo{}

	// 查询指定服务下所有的Case
	caseList, err = mongo.GetAllInspectionCasesByServiceAndStrategy(serviceId, strategy)
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
	serviceMongo := &models.ServiceMongo{}
	serviceMongo, _ = serviceMongo.QueryById(serviceId)
	serviceName := serviceMongo.ServiceName
	format := "20060102/150405"
	runReport.Name = businessName + "/" + serviceName + "-" + userId + "-" + time.Now().Format(format)
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

	wgOuter := sync.WaitGroup{}
	wgOuter.Add(1)
	go func() {
		wgInner := sync.WaitGroup{}
		wgInner.Add(len(caseList))
		for _, val := range caseList {
			go func(domain string, url string, uuid string, param string, checkout string, caseId int64, runBy string) {
				defer func() {
					if err := recover(); err != nil {
						logs.Error("完犊子了，大概率又特么的有个童鞋写了个垃圾Case, 去执行记录页面瞧瞧，他的执行记录会一直处于运行中的状态。。。")
						// todo 可以往外推送一个钉钉消息，通报一下这个不会写Case的同学
					}
				}()
				libs.DoRequestV2(domain, url, uuid, param, checkout, caseId, models.INSPECTION, runBy)
				// 获取用例执行进度时使用
				r := utils.GetRedis()
				r.Incr(constant.RUN_RECORD_CASE_DONE_NUM + uuid)
				wgInner.Done()
			}(val.Domain, val.ApiUrl, uuid, val.Parameter, val.Checkpoint, val.Id, userId)
		}
		wgInner.Wait()

		autoResult, _ := models.GetResultByRunId(uuid)
		var isPass int8 = models.SUCCESS
		// 判断case执行结果集合中是否有失败的case，有则认为本次执行操作状态为FAIL
		for _, result := range autoResult {
			if result.Result == models.AUTO_RESULT_FAIL {
				isPass = models.FAIL
				// todo 某个服务的巡检任务存在失败Case时，认定为本次巡检任务失败，对外发送钉钉消息通知到相关同学
				// todo 发送钉钉消息时，注意频次，预防被封群
				//logs.Warn("巡检任务失败，发送一条钉钉通知消息")
				msg := fmt.Sprintf("【业务线】: %s, 【服务】: %s。 报告链接: http://172.20.20.86:8080/report/run_report_detail?id=%d;\n", businessName, serviceName, id)
				// 将报告错误消息写进channel
				msgChannel <- msg
				break
			}
		}
		// 更新失败个数和本次执行记录状态
		autoResultMongo := &models.AutoResult{}
		failCount, _ := autoResultMongo.GetFailCount(uuid)
		runReport.UpdateIsPass(id, isPass, failCount, userId)
		wgOuter.Done()
	}()
	wgOuter.Wait()
	return
}

type ReqBody struct {
	//at      interface{}
	text    interface{}
	msgtype string
}

func DingSend(content string) {
	req := httplib.Post(XIAO_NENG_QUN).Debug(true)
	req.Header("Content-Type", "application/json;charset=utf-8")
	param := "{\"at\":{\"atMobiles\":[],\"atUserIds\":[],\"isAtAll\":false},\"text\":{\"content\":\"" + content + "\"},\"msgtype\":\"text\"}"
	req.Body(param)
	resp, err := req.Response()
	if err != nil {
		logs.Error("巡检任务失败时，发送钉钉消息报错，err: ", err)
	}
	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error("读取钉钉响应结果报错， err:", err)
	}
	logs.Info("调用钉钉发送通知接口返回: res:", string(res))
}
