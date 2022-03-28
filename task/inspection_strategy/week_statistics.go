package inspection_strategy

import (
	"github.com/astaxie/beego/logs"
	controllers "go_autoapi/controllers/statisticsdata"
	"go_autoapi/models"
)

const Layout = "2006-01-02 15:04:05" //时间常量
const ApiCountZuiyou = 200           //每条业务线接口总数
const ApiCountPipi = 200
const ApiCountHaiwai = 200
const ApiCountZhongdong = 200
const ApiCountMatuan = 200
const ApiCountShangyehua = 200
const ApiCountHaiwaiUS = 200

func Statistics1Week() error {
	//todo 取得平台自动化所用的全部接口by business
	logs.Info("定时启动更新数据统计 周五0点")
	mongo := models.StatisticsMongo{}
	c := controllers.StatisticsController{}

	respDataList := c.GetAllApiGroupByBusiness()
	for _, one := range respDataList {
		mongo.AllApi = one.AllApi
		mongo.DegreeOfCompletion = one.DegreeOfCompletion
		mongo.NewCaseConut = one.NewCaseConut
		mongo.AllCaseCount = one.AllCaseCount
		mongo.BusinessName = one.BusinessName
		mongo.LastWeekDegreeOfCompletion = one.LastWeekDegreeOfCompletion
		mongo.UnUseApi = one.UnUseApi
		err := mongo.Insert(mongo)
		if err != nil {
			logs.Error("插入数据出错")
		}

	}

	return nil

}
