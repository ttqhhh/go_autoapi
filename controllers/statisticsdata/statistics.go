package statisticsdata

import (
	"github.com/astaxie/beego/logs"
	"go_autoapi/libs"
	"go_autoapi/models"
	"strings"
	"time"
)

type StatisticsController struct {
	libs.BaseController
}

const Layout = "2006-01-02 15:04:05" //时间常量

func (c *StatisticsController) Get() {
	do := c.GetMethodName()
	switch do {
	case "get_api_by_business_in_platform":
		c.getApiByBusinessInPlatform()
	case "get_api_this_week_by_create_time":
		c.getApiThisWeekByCreateAt()

	default:
		logs.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持", nil)
	}
}

func (c *StatisticsController) Post() {
	do := c.GetMethodName()
	switch do {
	case "get_all_api_by_business":
		c.getAllApiByBusiness()

	default:
		logs.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持", nil)
	}
}

func (c *StatisticsController) getAllApiByBusiness() {
	//todo 可能需要去请求别的接口 获取业务线下用到的全部接口

}

func (c *StatisticsController) getApiByBusinessInPlatform() {
	//todo 取得平台自动化所用的全部接口by business
	business := c.GetString("business")
	mongo := models.TestCaseMongo{}
	result, err := mongo.GetAllCasesByBusiness(business)
	if err != nil {
		logs.Error("统计case时通过业务线获取全部case出错")
	}
	var apiList []string
	for _, one := range result {
		api := strings.Split(one.ApiUrl, "?")[0]
		//apiList:=append(apiList, api)
		apiList = append(apiList, api)
	}
	noRepeatApiList := RemoveRepeatedElement(apiList)
	count := len(noRepeatApiList)
	c.FormSuccessJson(int64(count), noRepeatApiList)
}

func (c *StatisticsController) getApiThisWeekByCreateAt() {
	//todo 取得平台自动化所用的全部接口by business
	business := c.GetString("business")
	mongo := models.TestCaseMongo{}
	result, err := mongo.GetAllCasesByBusiness(business)
	if err != nil {
		logs.Error("统计case时通过业务线获取全部case出错")
	}
	var apiList []string
	for _, one := range result {
		CreateTime, _ := time.ParseInLocation(Layout, one.CreatedAt, time.Local) //转换为date类型

		if CreateTime.After(time.Now().AddDate(0, 0, -7)) { //如果报告时间超过一周
			api := strings.Split(one.ApiUrl, "?")[0]
			apiList = append(apiList, api)
		}

	}
	noRepeatApiList := RemoveRepeatedElement(apiList)
	count := len(noRepeatApiList)
	c.FormSuccessJson(int64(count), noRepeatApiList)

}

func RemoveRepeatedElement(arr []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}
