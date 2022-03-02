package statisticsdata

import (
	"github.com/astaxie/beego/logs"
	"go_autoapi/libs"
	"go_autoapi/models"
	"strconv"
	"strings"
	"time"
)

type StatisticsController struct {
	libs.BaseController
}

const Layout = "2006-01-02 15:04:05" //时间常量
const ApiCountZuiyou = 200           //每条业务线接口总数
const ApiCountPipi = 200
const ApiCountHaiwai = 200
const ApiCountZhongdong = 200
const ApiCountMatuan = 200
const ApiCountShangyehua = 200
const ApiCountHaiwaiUS = 200

type respData struct {
	BusinessName       string  `form:"business_name" json:"business_name"`
	AllApiCount        float64 `form:"all_api_count" json:"all_api_count"`
	NewApiConut        float64 `form:"new_api_count" json:"new_api_count"`
	AllCaseCount       int     `form:"all_case_count" json:"all_case_count"`
	NewCaseConut       int     `form:"new_case_count" json:"new_case_count"`
	AllApi             int     `form:"all_api" json:"all_api"`                           //完成度
	DegreeOfCompletion string  `form:"degree_of_completion" json:"degree_of_completion"` //完成度

}

func (c *StatisticsController) Get() {
	do := c.GetMethodName()
	switch do {
	case "show_statistics_data":
		c.showStatisticsData()
	case "get_all_api_group": //获取对应业务线全接口
		c.getAllApiGroupByBusiness()
	case "get_api_group_new_add": //判断每周新增
		c.getApiByBusinessNewAdd()

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

func (c *StatisticsController) showStatisticsData() {
	nowTime := time.Now()
	if nowTime.Weekday().String() == "Friday" {
		FridayTime := nowTime.AddDate(0, 0, -7)
		timeNowstr := nowTime.Format("2006-01-02 15:04:05")
		FridayTimeStr := FridayTime.Format("2006-01-02 15:04:05")
		c.Data["now_time"] = timeNowstr
		c.Data["check_time"] = FridayTimeStr
		c.TplName = "show_statistics_data.html"

	} else {
		FridayTime := getFridayTime(nowTime)
		timeNowstr := nowTime.Format("2006-01-02 15:04:05")
		FridayTimeStr := FridayTime.Format("2006-01-02 15:04:05")
		c.Data["now_time"] = timeNowstr
		c.Data["check_time"] = FridayTimeStr
		c.TplName = "show_statistics_data.html"
	}

}
func (c *StatisticsController) getAllApiByBusiness() {
	//todo 可能需要去请求别的接口 获取业务线下用到的全部接口

}

func (c *StatisticsController) getAllApiGroupByBusiness() {
	//todo 取得平台自动化所用的全部接口by business
	mongo := models.TestCaseMongo{}
	result, err := mongo.GetAllCasesNoBusiness()
	if err != nil {
		logs.Error("统计case时通过业务线获取全部case出错")
	}
	var zuiyou_list []string
	var pipi_list []string
	var haiwai_list []string
	var zhongdong_list []string
	var shangyehuai_list []string
	var matuan_list []string
	var haiwaiUS_list []string
	for _, one := range result {
		api := strings.Split(one.ApiUrl, "?")[0]
		switch one.BusinessCode {
		case "0": //最右
			zuiyou_list = append(zuiyou_list, api)
		case "1": //皮皮
			pipi_list = append(pipi_list, api)
		case "2": //海外
			haiwai_list = append(haiwai_list, api)
		case "3": //中东
			zhongdong_list = append(zhongdong_list, api)
		case "4": //麻团
			matuan_list = append(matuan_list, api)
		case "5": //商业化
			shangyehuai_list = append(shangyehuai_list, api)
		case "6": //海外-us
			haiwaiUS_list = append(haiwaiUS_list, api)
		default:
			logs.Warn("no business")
		}

	}
	noRepeatZuiyouList := RemoveRepeatedElement(zuiyou_list)
	noRepeatPipiList := RemoveRepeatedElement(pipi_list)
	noRepeatHaiwaiList := RemoveRepeatedElement(haiwai_list)
	noRepeatZhongdongList := RemoveRepeatedElement(zhongdong_list)
	noRepeatMatuanList := RemoveRepeatedElement(matuan_list)
	noRepeatShangyehuaList := RemoveRepeatedElement(shangyehuai_list)
	noRepeatHaiwaiUSList := RemoveRepeatedElement(haiwaiUS_list)

	resp2 := c.getApiByBusinessNewAdd()

	respDataList := []respData{}

	respData := respData{}
	respData1 := respData
	respData1.BusinessName = "最右"
	respData1.AllApiCount = float64(len(noRepeatZuiyouList))
	respData1.NewApiConut = float64(resp2["zuiyou_count_new"])
	respData1.AllCaseCount = len(zuiyou_list)
	respData1.NewCaseConut = resp2["zuiyou_new_case"]
	str := strconv.FormatFloat(float64(respData1.AllApiCount/ApiCountZuiyou)*100, 'f', 2, 64)
	respData1.AllApi = ApiCountZuiyou
	respData1.DegreeOfCompletion = str + "%"

	respDataList = append(respDataList, respData1)

	respData2 := respData
	respData2.BusinessName = "皮皮"
	respData2.AllApiCount = float64(len(noRepeatPipiList))
	respData2.NewApiConut = float64(resp2["pipi_count_new"])
	respData2.AllCaseCount = len(pipi_list)
	respData2.NewCaseConut = resp2["pipi_new_case"]
	str2 := strconv.FormatFloat(float64(respData2.AllApiCount/ApiCountZuiyou)*100, 'f', 2, 64)
	respData2.DegreeOfCompletion = str2 + "%"
	respData2.AllApi = ApiCountZuiyou

	respDataList = append(respDataList, respData2)

	respData3 := respData
	respData3.BusinessName = "海外"
	respData3.AllApiCount = float64(len(noRepeatHaiwaiList))
	respData3.NewApiConut = float64(resp2["haiwai_count_new"])
	respData3.AllCaseCount = len(haiwai_list)
	respData3.NewCaseConut = resp2["haiwai_new_case"]
	str3 := strconv.FormatFloat(float64(respData3.AllApiCount/ApiCountZuiyou)*100, 'f', 2, 64)
	respData3.DegreeOfCompletion = str3 + "%"
	respData3.AllApi = ApiCountZuiyou

	respDataList = append(respDataList, respData3)

	respData4 := respData
	respData4.BusinessName = "中东"
	respData4.AllApiCount = float64(len(noRepeatZhongdongList))
	respData4.NewApiConut = float64(resp2["zhongdong_count_new"])
	respData4.AllCaseCount = len(zhongdong_list)
	respData4.NewCaseConut = resp2["zhongdong_new_case"]
	str4 := strconv.FormatFloat(float64(respData4.AllApiCount/ApiCountZuiyou)*100, 'f', 2, 64)
	respData4.DegreeOfCompletion = str4 + "%"
	respData4.AllApi = ApiCountZuiyou

	respDataList = append(respDataList, respData4)

	respData5 := respData
	respData5.BusinessName = "麻团"
	respData5.AllApiCount = float64(len(noRepeatMatuanList))
	respData5.NewApiConut = float64(resp2["matuan_count_new"])
	respData5.AllCaseCount = len(matuan_list)
	respData5.NewCaseConut = resp2["matuan_new_case"]
	str5 := strconv.FormatFloat(float64(respData5.AllApiCount/ApiCountZuiyou)*100, 'f', 2, 64)
	respData5.DegreeOfCompletion = str5 + "%"
	respData5.AllApi = ApiCountZuiyou

	respDataList = append(respDataList, respData5)

	respData6 := respData
	respData6.BusinessName = "商业化"
	respData6.AllApiCount = float64(len(noRepeatShangyehuaList))
	respData6.NewApiConut = float64(resp2["shangyehua_count_new"])
	respData6.AllCaseCount = len(shangyehuai_list)
	respData6.NewCaseConut = resp2["shangyehua_new_case"]
	str6 := strconv.FormatFloat(float64(respData6.AllApiCount/ApiCountZuiyou)*100, 'f', 2, 64)
	respData6.DegreeOfCompletion = str6 + "%"
	respData6.AllApi = ApiCountZuiyou

	respDataList = append(respDataList, respData6)

	respData7 := respData
	respData7.BusinessName = "海外-US"
	respData7.AllApiCount = float64(len(noRepeatHaiwaiUSList))
	respData7.NewApiConut = float64(resp2["haiwaiUS_count_new"])
	respData7.AllCaseCount = len(haiwaiUS_list)
	respData7.NewCaseConut = resp2["haiwaiUS_new_case"]
	str7 := strconv.FormatFloat(float64(respData7.AllApiCount/ApiCountZuiyou)*100, 'f', 2, 64)
	respData7.DegreeOfCompletion = str7 + "%"
	respData7.AllApi = ApiCountZuiyou

	respDataList = append(respDataList, respData7)
	c.FormSuccessJson(int64(len(respDataList)), respDataList)
}

func (c *StatisticsController) getApiByBusinessNewAdd() map[string]int {
	//todo 根据create判断新增，如果now是周五，则判断当前时间前7天，如果不是周五，向前查找，直到找到周五后，判断周五的后七天
	mongo := models.TestCaseMongo{}
	result, _ := mongo.GetAllCasesNoBusiness()
	var zuiyou_list_new []string
	var pipi_list_new []string
	var haiwai_list_new []string
	var zhongdong_list_new []string
	var shangyehuai_list_new []string
	var matuan_list_new []string
	var haiwaiUS_list_new []string

	useTime := time.Now()
	if useTime.Weekday().String() == "Friday" {
		for _, one := range result {
			caseCreateTime, _ := time.ParseInLocation(Layout, one.CreatedAt, time.Local) //获取报告的创建时间 转换为data
			api := strings.Split(one.ApiUrl, "?")[0]
			switch one.BusinessCode {
			case "0": //最右
				if caseCreateTime.After(useTime.AddDate(0, 0, -7)) { //如果是周五，则获取前7天报告
					zuiyou_list_new = append(zuiyou_list_new, api)
				}
			case "1": //皮皮
				if caseCreateTime.After(useTime.AddDate(0, 0, -7)) { //如果是周五，则获取前7天报告
					pipi_list_new = append(pipi_list_new, api)
				}
			case "2": //海外
				if caseCreateTime.After(useTime.AddDate(0, 0, -7)) { //如果是周五，则获取前7天报告
					haiwai_list_new = append(haiwai_list_new, api)
				}
			case "3": //中东
				if caseCreateTime.After(useTime.AddDate(0, 0, -7)) { //如果是周五，则获取前7天报告
					zhongdong_list_new = append(zhongdong_list_new, api)
				}
			case "4": //麻团
				if caseCreateTime.After(useTime.AddDate(0, 0, -7)) { //如果是周五，则获取前7天报告
					shangyehuai_list_new = append(shangyehuai_list_new, api)
				}
			case "5": //商业化
				if caseCreateTime.After(useTime.AddDate(0, 0, -7)) { //如果是周五，则获取前7天报告
					matuan_list_new = append(matuan_list_new, api)
				}
			case "6": //海外-us
				if caseCreateTime.After(useTime.AddDate(0, 0, -7)) { //如果是周五，则获取前7天报告
					haiwaiUS_list_new = append(haiwaiUS_list_new, api)
				}
			default:
				logs.Warn("no business")
			}
		}
	} else {
		useTimeF := getFridayTime(useTime) //获取上周五的时间
		for _, one := range result {
			caseCreateTime, _ := time.ParseInLocation(Layout, one.CreatedAt, time.Local) //获取报告的创建时间 转换为data
			api := strings.Split(one.ApiUrl, "?")[0]
			switch one.BusinessCode {
			case "0": //最右
				if caseCreateTime.Before(useTimeF.AddDate(0, 0, 7)) && caseCreateTime.After(useTimeF) { //如不是是周五，找到上一个周五，往后推7天，并大于上个周五
					zuiyou_list_new = append(zuiyou_list_new, api)
				}
			case "1": //皮皮
				if caseCreateTime.Before(useTimeF.AddDate(0, 0, 7)) && caseCreateTime.After(useTimeF) {
					pipi_list_new = append(pipi_list_new, api)
				}
			case "2": //海外
				if caseCreateTime.Before(useTimeF.AddDate(0, 0, 7)) && caseCreateTime.After(useTimeF) {
					haiwai_list_new = append(haiwai_list_new, api)
				}
			case "3": //中东
				if caseCreateTime.Before(useTimeF.AddDate(0, 0, 7)) && caseCreateTime.After(useTimeF) {
					zhongdong_list_new = append(zhongdong_list_new, api)
				}
			case "4": //麻团
				if caseCreateTime.Before(useTimeF.AddDate(0, 0, 7)) && caseCreateTime.After(useTimeF) {
					matuan_list_new = append(matuan_list_new, api)
				}
			case "5": //商业化
				if caseCreateTime.Before(useTimeF.AddDate(0, 0, 7)) && caseCreateTime.After(useTimeF) {
					shangyehuai_list_new = append(shangyehuai_list_new, api)
				}
			case "6": //海外-us
				if caseCreateTime.Before(useTimeF.AddDate(0, 0, 7)) && caseCreateTime.After(useTimeF) {
					haiwaiUS_list_new = append(haiwaiUS_list_new, api)
				}
			default:
				logs.Warn("no business")
			}
		}
	}

	noRepeatZuiyouList := RemoveRepeatedElement(zuiyou_list_new)
	noRepeatPipiList := RemoveRepeatedElement(pipi_list_new)
	noRepeatHaiwaiList := RemoveRepeatedElement(haiwai_list_new)
	noRepeatZhongdongList := RemoveRepeatedElement(zhongdong_list_new)
	noRepeatMatuanList := RemoveRepeatedElement(matuan_list_new)
	noRepeatShangyehuaList := RemoveRepeatedElement(shangyehuai_list_new)
	noRepeatHaiwaiUSList := RemoveRepeatedElement(haiwaiUS_list_new)

	resp := make(map[string]int)
	resp["zuiyou_new_case"] = len(zuiyou_list_new)
	resp["pipi_new_case"] = len(pipi_list_new)
	resp["haiwai_new_case"] = len(haiwai_list_new)
	resp["zhongdong_new_case"] = len(zhongdong_list_new)
	resp["matuan_new_case"] = len(matuan_list_new)
	resp["shangyehua_new_case"] = len(shangyehuai_list_new)
	resp["haiwaiUS_new_case"] = len(haiwaiUS_list_new)
	resp["zuiyou_count_new"] = len(noRepeatZuiyouList)
	resp["pipi_count_new"] = len(noRepeatPipiList)
	resp["haiwai_count_new"] = len(noRepeatHaiwaiList)
	resp["zhongdong_count_new"] = len(noRepeatZhongdongList)
	resp["matuan_count_new"] = len(noRepeatMatuanList)
	resp["shangyehua_count_new"] = len(noRepeatShangyehuaList)
	resp["haiwaiUS_count_new"] = len(noRepeatHaiwaiUSList)

	return resp

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

func getFridayTime(nowTime time.Time) time.Time { //返回当前时间的上一个周五的0点
	if nowTime.Weekday().String() == "Friday" {
		return nowTime
	} else {
		for i := 1; i <= 7; i++ {
			if nowTime.AddDate(0, 0, -i).Weekday().String() == "Friday" {
				FridayNowTime := nowTime.AddDate(0, 0, -i)
				FridayNowTimeString := FridayNowTime.Format("2006-01-02")
				Fridaw0Time, _ := time.Parse("2006-01-02", FridayNowTimeString) //当前时间至周五0点
				return Fridaw0Time
			}

		}
		return nowTime
	}

}
