package h5listen

import (
	"go_autoapi/models"
)

const XIAO_NENG_QUN_TOKEN = "5d8a1de9fbf4b87ef01ccf38b65fc4bbf5029cf1e403411fec0cb0baef057503"

func (c *H5ListenController) ShowH5() {
	acm := models.AdH5DataMongo{}
	business := c.GetString("business", "1")
	result, err := acm.GetDataByBusiness(business)
	if err != nil {
		c.FormErrorJson(-1, "获取h5列表数据失败")
	}
	c.FormSuccessJson(0, result)
}

//func (c *H5ListenController) StrategyH5() error {
//	logs.Info("【h5】定时任务启动执行...")
//	h5Mongo := models.AdH5DataMongo{}
//	h5Mongos, err := h5Mongo.GetAllData()
//	if err != nil {
//		//logs.Error("执行线上巡检定时任务时，查询指定业务线下的服务时报错， err: ", err)
//		return err
//	}
//	// 遍历服务下边所有的巡检Case
//	for _, h5data := range h5Mongos {
//		//PerformInspection(businessId, service.Id, msgChannel, restrainMsgChannel, inspection.ONE_MIN_CODE)
//		h5_url := h5data.DataUrl
//		resp, _ := http.Get(h5_url)
//		defer resp.Body.Close()
//		// 200 OK
//		status := resp.Status
//		logs.Info("返回信息是：",status)
//		if status != "200 OK" {
//			msg := fmt.Sprintf("业务名称：" + h5data.DataName + "       错误码:" + status + "     链接：" + h5_url)
//			fmt.Println(resp.Status)
//			DingSend(msg)
//			fmt.Println(status)
//			acm := models.H5RunReportMongo{
//				Id:          time.Now().Unix(),
//				Business:     h5data.Business,
//				BusinessName: "",
//				DataName:     h5data.DataName,
//				DataUrl:      h5_url,
//				ErrorCode:    status,
//				CreatedAt:    time.Now().Format("2006-01-02 15:04:05"),
//				Status:       "失败",
//			}
//			switch h5data.Business {
//			case "0":
//				acm.BusinessName="zuiyou"
//			case "1":
//				acm.BusinessName="pipi"
//			case "2":
//				acm.BusinessName="haiwai"
//			case "3":
//				acm.BusinessName="zhongdong"
//			case "4":
//				acm.BusinessName="matuan"
//			case "5":
//				acm.BusinessName="business"
//			case "6":
//				acm.BusinessName="haiwai-US"
//			}
//			dao := models.H5RunReportMongo{}
//			err := dao.Insert(acm)
//			if err != nil {
//				logs.Error("插入h5报告报错，err: ", err)
//				c.ErrorJson(-1, "插入错误", nil)
//			}
//
//		}
//
//		logs.Info("h5巡检执行完毕，如果监测到问题开始发送叮叮--------")
//	}
//
//	return nil
//}
//func DingSend(content string) {
//	var dingToken = []string{XIAO_NENG_QUN_TOKEN}
//	cli := dingtalk.InitDingTalk(dingToken, "")
//	cli.SendTextMessage(content)
//	logs.Info("钉钉消息已经发出去了")
//}
