package inspection_strategy

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/blinkbean/dingtalk"
	uuid "github.com/satori/go.uuid"
	constant "go_autoapi/constants"
	controllers "go_autoapi/controllers/autotest"
	"go_autoapi/libs"
	"go_autoapi/models"
	"go_autoapi/utils"
	"sync"
	"time"
)

// 「测试效率团队」群web_hook-用来测试
//const XIAO_NENG_QUN = "https://oapi.dingtalk.com/robot/send?access_token=6f35268d9dcb74b4b95dd338eb241832781aeaaeafd90aa947b86936f3343dbb"
//const XIAO_NENG_QUN_TOKEN = "6f35268d9dcb74b4b95dd338eb241832781aeaaeafd90aa947b86936f3343dbb"

// 「测试管理群」群web_hook-用来测试
//const XIAO_NENG_QUN = "https://oapi.dingtalk.com/robot/send?access_token=a822521c37e0d566563652452a1fdd692f27f1746d59c4229dd91047ba52f325"
//const XIAO_NENG_QUN_TOKEN = "a822521c37e0d566563652452a1fdd692f27f1746d59c4229dd91047ba52f325"

// 「测试管理群」群web_hook-用来测试
//const XIAO_NENG_QUN = "https://oapi.dingtalk.com/robot/send?access_token=60ee4a400b625f8bb3284f12a2a5b8e6bf9eabb76fd23982359ffbb23e591a4d"
const XIAO_NENG_QUN_TOKEN = "60ee4a400b625f8bb3284f12a2a5b8e6bf9eabb76fd23982359ffbb23e591a4d"

// todo 上线or正常使用时，需要设为true进行开启
const IS_OPEN_SENDDING_MSG = true

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
				// 当巡检用例执行失败时，再进行2次补偿重试
				retryTimes := 0
				isContinue := true
				for isContinue {
					isOk := libs.DoRequestV2(domain, url, uuid, param, checkout, caseId, models.INSPECTION, runBy)
					if retryTimes >= 2 || isOk {
						isContinue = false
					}
					retryTimes++
				}
				// 获取用例执行进度时使用
				r := utils.GetRedis()
				r.Incr(constant.RUN_RECORD_CASE_DONE_NUM + uuid)
				wgInner.Done()
			}(val.Domain, val.ApiUrl, uuid, val.Parameter, val.Checkpoint, val.Id, userId)
		}
		wgInner.Wait()
		autoResult, _ := models.GetResultByRunId(uuid)
		var isPass int8 = models.SUCCESS
		//用来盛放同一个Case多次执行的结果
		case2ResultMap := make(map[int64][]*models.AutoResult)
		// 判断case执行结果集合中是否有失败的case，有则认为本次执行操作状态为FAIL
		for _, result := range autoResult {
			if result.Result == models.AUTO_RESULT_FAIL {
				caseId := result.CaseId
				autoResultList := case2ResultMap[caseId]
				if autoResultList == nil || len(autoResultList) == 0 {
					autoResultList = []*models.AutoResult{}
				}
				autoResultList = append(autoResultList, result)
				case2ResultMap[caseId] = autoResultList

				isPass = models.FAIL
				//logs.Warn("巡检任务失败，发送一条钉钉通知消息")
				//msg := fmt.Sprintf("【业务线】: %s, 【服务】: %s。 报告链接: http://172.20.20.86:8080/report/run_report_detail?id=%d;\n", businessName, serviceName, id)
				// 将报告错误消息写进channel
				//msgChannel <- msg
				//break
			}
		}
		baseMsg := fmt.Sprintf("【业务线】: %s, 【服务】: %s。 报告链接: http://172.16.2.86:8080/report/run_report_detail?id=%d;\n\n", businessName, serviceName, id)
		// 遍历case2ResultMap，哪个caseId对应的value长度为3，则该条Case为失败Case
		msg := ""
		for caseId, autoResultList := range case2ResultMap {
			if len(autoResultList) > 2 {
				//todo 此时该条巡检Case有问题，进行对外通知
				//testCaseMongo := models.TestCaseMongo{}
				//testCaseMongo = testCaseMongo.GetOneCase(caseId)
				icm := models.InspectionCaseMongo{}
				icm = icm.GetOneCase(caseId)
				caseName := icm.CaseName
				uri := icm.ApiUrl

				autoResult := autoResultList[2]
				resp := autoResult.Response
				//resp := "{\"ret\":1,\"data\":{\"banner\":[{\"name\":\"gaokaozhiyuan\",\"img\":\"file.izuiyou.com/img/png/id/1567506494\",\"url\":\"zuiyou://eventactivity?eventActivityId=330334\"},{\"name\":\"fangyandugongyue\",\"img\":\"https://file.izuiyou.com/img/png/id/1568965881\",\"url\":\"zuiyou://postdetail?id=233156321\"},{\"name\":\"MCNzhaomu\",\"img\":\"https://file.izuiyou.com/img/png/id/1564965136\",\"url\":\"https://h5.izuiyou.com/hybrid/template/smartH5?\\u0026id=329839\"},{\"name\":\"shenhezhuanqu\",\"img\":\"https://file.izuiyou.com/img/png/id/1568973427\",\"url\":\"https://h5.izuiyou.com/hybrid/censor/entry\"},{\"name\":\"maishoudian\",\"img\":\"https://file.izuiyou.com/img/png/id/1561594314\",\"url\":\"zuiyou://postdetail?id=232312632\"},{\"name\":\"wanyouxi\",\"img\":\"https://file.izuiyou.com/img/png/id/1567612281\",\"url\":\"http://www.shandw.com/auth\"}]}}"
				reason := autoResult.Reason
				statusCode := autoResult.StatusCode

				// todo 某个服务的巡检任务存在失败Case时，认定为本次巡检任务失败，对外发送钉钉消息通知到相关同学
				// todo 发送钉钉消息时，注意频次，预防被封群
				msg += fmt.Sprintf("【Case名称】: %s;\n【接口路径】: %s;\n【请求状态码】: %d;\n【失败原因】: %s;\n【响应结果】: %s;\n\n", caseName, uri, statusCode, reason, resp)
				//msg = fmt.Sprintf("%s【Case名称: %s; 接口路径: %s; 失败原因: %s; 】\n", msg, caseName, uri, reason)
				// 将报告错误消息写进channel
				//msgChannel <- msg
				//break
			}
		}
		if msg != "" {
			totalMsg := baseMsg + msg
			msgChannel <- totalMsg
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
	var dingToken = []string{XIAO_NENG_QUN_TOKEN}
	cli := dingtalk.InitDingTalk(dingToken, "")
	cli.SendTextMessage(content)
}

//func DingSend(content string) {
//	req := httplib.Post(XIAO_NENG_QUN).Debug(true)
//	req.Header("Content-Type", "application/json;charset=utf-8")
//	param := "{\"at\":{\"atMobiles\":[],\"atUserIds\":[],\"isAtAll\":false},\"text\":{\"content\":\"" + content + "\"},\"msgtype\":\"text\"}"
//	req.Body(param)
//	resp, err := req.Response()
//	if err != nil {
//		logs.Error("巡检任务失败时，发送钉钉消息报错，err: ", err)
//	}
//	res, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		logs.Error("读取钉钉响应结果报错， err:", err)
//	}
//	logs.Info("调用钉钉发送通知接口返回: res:", string(res))
//}
