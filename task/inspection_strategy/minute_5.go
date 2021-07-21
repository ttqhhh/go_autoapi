package inspection_strategy

import "fmt"

func Strategy5Min() error {
	//logs.Info("【5分钟】级别的定时任务启动执行...")
	////msgList := make([]string, 5)
	//msgChannel := make(chan string)
	//var msgList []string
	//// 起一个协程用于承接msg
	//go func() {
	//	for true {
	//		msg := <-msgChannel
	//		msgList = append(msgList, msg)
	//	}
	//}()
	//// 获取所以业务线，并进行遍历所有的业务线
	//businesses := controllers.GetAllBusinesses()
	//for _, business := range businesses {
	//	// 遍历业务线下的所有所有服务
	//	serviceMongo := models.ServiceMongo{}
	//	businessId := int8(business["code"].(int))
	//	serviceMongos, err := serviceMongo.QueryByBusiness(businessId)
	//	if err != nil {
	//		logs.Error("执行线上巡检定时任务时，查询指定业务线下的服务时报错， err: ", err)
	//		return err
	//	}
	//	// 遍历服务下边所有的巡检Case
	//	for _, service := range serviceMongos {
	//		PerformInspection(businessId, service.Id, msgChannel, inspection.FIVE_MIN_CODE)
	//	}
	//
	//}
	//// dingMsg中的「线上巡检」为消息关键字，不可变更
	//dingMsg := "小钻风【5分钟】级别线上巡检发现异样, 快去排查一下吧。\n"
	//for _, msg := range msgList {
	//	dingMsg += msg
	//}
	//if len(msgList) > 0 {
	//	logs.Info("打印钉钉消息日志：\n" + dingMsg)
	//	if IS_OPEN_SENDDING_MSG {
	//		DingSend(dingMsg)
	//	}
	//}
	fmt.Printf("打印一次日志验证tag部署")
	return nil
}
