package inspection_strategy

import (
	"github.com/astaxie/beego/logs"
	controllers "go_autoapi/controllers/autotest"
	"go_autoapi/controllers/inspection"
	"go_autoapi/models"
)

func Strategy1Day() error {
	logs.Info("【天】级别定时任务启动执行...")
	//msgList := make([]string, 5)
	msgChannel := make(chan string)
	restrainMsgChannel := make(chan string)
	var msgList []string
	var restrainMsgList []string
	// 起一个协程用于承接msg
	go func() {
		for true {
			msg := <-msgChannel
			msgList = append(msgList, msg)
		}
	}()
	go func() {
		for true {
			msg := <-restrainMsgChannel
			restrainMsgList = append(restrainMsgList, msg)
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
			PerformInspection(businessId, service.Id, msgChannel, restrainMsgChannel, inspection.ONE_DAY_CODE)
		}

	}
	// dingMsg中的「线上巡检」为消息关键字，不可变更
	dingMsg := "小钻风【天】级别线上巡检发现异样, 快去排查一下吧。\n"
	for _, msg := range msgList {
		dingMsg += msg
	}
	if len(msgList) > 0 {
		logs.Info("打印钉钉消息日志：\n" + dingMsg)
		if IS_OPEN_SENDDING_MSG {
			logs.Info("开始发送叮叮消息，【巡检次数一天/次】")
			DingSend(dingMsg)
		}
	}
	// dingMsg中的「线上巡检」为消息关键字，不可变更
	// dingMsg := fmt.Sprintf("【报警抑制-线上巡检】：当前Case报警次数已达3次且未得到有效解决，系统已默认将Case巡检状态置为关闭。请相关同学尽快处理！编辑或打开该Case均可恢复巡检状态。\n" + "\n" + "【业务线】：" + icm.BusinessName + "\n" + "【网关服务】：" + serviceName + "\n" + "【URI】:" + dingUri + "\n" + "【Caseid】:" + caseId)
	dingRestrainMsg := "【报警抑制-线上巡检】：当前Case报警次数已达3次且未得到有效解决，系统已默认将Case巡检状态置为关闭。请相关同学尽快处理！编辑或打开该Case均可恢复巡检状态。\n\n"
	for _, restrainMsg := range restrainMsgList {
		dingRestrainMsg += restrainMsg
	}
	if len(restrainMsgList) > 0 {
		logs.Info("打印钉钉消息日志：\n" + dingRestrainMsg)
		if IS_OPEN_SENDDING_MSG {
			DingSend(dingRestrainMsg)
		}
	}
	return nil
}
