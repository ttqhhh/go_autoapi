package task

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
	"gopkg.in/mgo.v2/bson"
	"sync"
	"time"
)

const ONLINE_INSPECTION_EXPRESSION = "0 40 15 * * *"

//const ONLINE_INSPECTION_EXPRESSION = "* */2 * * * *"
const XIAO_NENG_QUN = "https://oapi.dingtalk.com/robot/send?access_token=257ff4003ced45fa2e18f230ca72a0f18bee77824f5b1a73981c25e417f19ea6"

func OnlineInspection() error {
	logs.Info("启动定时任务：Online Inspection")
	//msgList := make([]string, 5)
	msgChannel := make(chan string)
	var msgList []string
	// todo 起一个协程用于承接msg
	go func() {
		for true {
			msg := <-msgChannel
			msgList = append(msgList, msg)
		}
	}()
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
			performInspection(businessId, service.Id, msgChannel)
		}

	}
	// todo 一次巡检任务，对外发送一次失败消息；以巡检任务维度对外发送，而非服务Case维度
	dingMsg := "【 线上巡检 】报错啦, 快去排查一下。"
	for _, msg := range msgList {
		dingMsg += msg
	}
	if len(msgList) > 0 {
		logs.Info("打印钉钉消息日志：\n" + dingMsg)
		// todo 只有msgList大于1的时候，才去发送钉钉消息
	}
	return nil
}

// 执行线上巡检Case
func performInspection(businessId int8, serviceId int64, msgChannel chan string) (err error) {
	userId := "线上巡检"
	u2 := uuid.NewV4()
	uuid := u2.String()

	mongo := models.InspectionCaseMongo{}
	// 根据不同的执行维度，聚合需要执行的所有Case集合
	caseList := []*models.InspectionCaseMongo{}

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

	go func() {
		wg := sync.WaitGroup{}
		wg.Add(len(caseList))
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
				wg.Done()
			}(val.Domain, val.ApiUrl, uuid, val.Parameter, val.Checkpoint, val.Id, userId)
		}
		wg.Wait()

		//go func() {
		autoResult, _ := models.GetResultByRunId(uuid)
		var isPass int8 = models.SUCCESS
		// 判断case执行结果集合中是否有失败的case，有则认为本次执行操作状态为FAIL
		for _, result := range autoResult {
			if result.Result == models.AUTO_RESULT_FAIL {
				isPass = models.FAIL
				// todo 某个服务的巡检任务存在失败Case时，认定为本次巡检任务失败，对外发送钉钉消息通知到相关同学
				// todo 发送钉钉消息时，注意频次，预防被封群
				//logs.Warn("巡检任务失败，发送一条钉钉通知消息")
				/** 获取对应的业务线名称和服务名称 */
				//serviceName := serviceMongo.ServiceName
				//businessName := controllers.GetBusinessNameByCode(int(businessId))
				msg := fmt.Sprintf("业务线: %s, 服务: %s。 报告链接: http://localhost:8080/report/run_report_detail?id=%d;\n", businessName, serviceName, id)
				//logs.Info(msg)
				//dingSend(msg)
				// 将报告错误消息写进channel
				msgChannel <- msg
				break
			}
		}
		// 更新失败个数和本次执行记录状态
		autoResultMongo := &models.AutoResult{}
		failCount, _ := autoResultMongo.GetFailCount(uuid)
		runReport.UpdateIsPass(id, isPass, failCount, userId)
	}()
	//}()
	return
}

type ReqBody struct {
	at      interface{}
	text    interface{}
	msgtype string
}

func dingSend(content string) {
	req := httplib.Post(XIAO_NENG_QUN).Debug(true)
	req.Header("Content-Type", "application/json;charset=utf-8")
	//body := bson.M{"msgtype": "text","text": {"content":"我就是我, 是不一样的烟火"}}
	body := ReqBody{
		at:      nil,
		text:    bson.M{"content": content},
		msgtype: "text",
	}
	req.Body(body)
	resp, err := req.Response()
	if err != nil {
		logs.Error("巡检任务失败时，发送钉钉消息报错，err: ", err)
	}
	logs.Info("调用钉钉发送通知接口返回: ", resp)
}
