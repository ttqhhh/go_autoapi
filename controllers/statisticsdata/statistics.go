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

type respData struct {
	BusinessName               string  `form:"business_name" json:"business_name"`
	AllApiCount                float64 `form:"all_api_count" json:"all_api_count"`
	NewApiConut                float64 `form:"new_api_count" json:"new_api_count"`
	AllCaseCount               int     `form:"all_case_count" json:"all_case_count"`
	NewCaseConut               int     `form:"new_case_count" json:"new_case_count"`
	AllApi                     int     `form:"all_api" json:"all_api"`                                               //完成度
	DegreeOfCompletion         string  `form:"degree_of_completion" json:"degree_of_completion"`                     //完成度
	LastWeekDegreeOfCompletion string  `form:"last_week_degree_of_completion" json:"last_week_degree_of_completion"` //上周完成度

}

func (c *StatisticsController) Get() {
	do := c.GetMethodName()
	switch do {
	case "show_statistics_data":
		c.showStatisticsData()
	case "get_all_api_group": //获取对应业务线全接口
		c.GetAllApiGroupByBusiness()
	case "get_api_group_new_add": //判断每周新增
		c.getApiByBusinessNewAdd()
	case "get_all_data":
		c.getAllQuery()

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
	nowData := getFridayTime(nowTime)
	nowDatanew := nowData.AddDate(0, 0, -1)
	FridayTime := nowData.AddDate(0, 0, -7)
	nowDataStr := nowDatanew.Format("2006-01-02 15:04:05")
	showTimeStr := strings.Split(nowDataStr, " ")[0] + " " + "23:59:59"
	FridayTimeStr := FridayTime.Format("2006-01-02 15:04:05")
	c.Data["now_time"] = showTimeStr
	c.Data["check_time"] = FridayTimeStr
	c.TplName = "show_statistics_data.html"

}
func (c *StatisticsController) getAllApiByBusiness() {
	//todo 可能需要去请求别的接口 获取业务线下用到的全部接口

}

func (c *StatisticsController) GetAllApiGroupByBusiness() []respData {
	//todo 取得平台自动化所用的全部接口by business
	data := getAllApi()
	var zuiyouAllCount float64
	var pipiAllCount float64
	var haiwaiAllCount float64
	var zhongdongAllCount float64
	var matuanAllConut float64
	var shangyehuaAllCount float64
	var haiwaiUSAllCount float64
	zuiyouAllCount = data["0"]
	pipiAllCount = data["1"]
	haiwaiAllCount = data["2"]
	zhongdongAllCount = data["3"]
	matuanAllConut = data["4"]
	shangyehuaAllCount = data["5"]
	haiwaiUSAllCount = data["6"]

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
	str := strconv.FormatFloat(float64(float64(respData1.AllCaseCount)/zuiyouAllCount)*100, 'f', 2, 64) //本周完成度
	respData1.AllApi = int(zuiyouAllCount)
	if zuiyouAllCount == 0 { //获取的接口总数为0 兼容 大于1 兼容
		respData1.DegreeOfCompletion = "0%"
	} else if (float64(respData1.AllCaseCount) / zuiyouAllCount) > 1 {
		respData1.DegreeOfCompletion = "100%"
	} else {
		respData1.DegreeOfCompletion = str + "%"
	}
	strLast := strconv.FormatFloat(float64(float64(respData1.AllCaseCount-respData1.NewCaseConut)/zuiyouAllCount)*100, 'f', 2, 64) //上周完成度  全部接口-这周新增/活跃接口
	respData1.AllApi = int(zuiyouAllCount)
	if zuiyouAllCount == 0 { //获取的接口总数为0 兼容 大于1 兼容
		respData1.LastWeekDegreeOfCompletion = "0%"
	} else if (float64(respData1.AllCaseCount) / zuiyouAllCount) > 1 {
		respData1.LastWeekDegreeOfCompletion = "100%"
	} else {
		respData1.LastWeekDegreeOfCompletion = strLast + "%"
	}

	respDataList = append(respDataList, respData1)

	respData2 := respData
	respData2.BusinessName = "皮皮"
	respData2.AllApiCount = float64(len(noRepeatPipiList))
	respData2.NewApiConut = float64(resp2["pipi_count_new"])
	respData2.AllCaseCount = len(pipi_list)
	respData2.NewCaseConut = resp2["pipi_new_case"]
	str2 := strconv.FormatFloat(float64(float64(respData2.AllCaseCount)/pipiAllCount)*100, 'f', 2, 64)
	respData2.DegreeOfCompletion = str2 + "%"
	respData2.AllApi = int(pipiAllCount)
	if pipiAllCount == 0 { //获取的接口总数为0 兼容 大于1 兼容
		respData2.DegreeOfCompletion = "0%"
	} else if (float64(respData2.AllCaseCount) / pipiAllCount) > 1 {
		respData2.DegreeOfCompletion = "100%"
	} else {
		respData2.DegreeOfCompletion = str2 + "%"
	}

	strLast2 := strconv.FormatFloat(float64(float64(respData2.AllCaseCount-respData2.NewCaseConut)/pipiAllCount)*100, 'f', 2, 64) //上周完成度  全部接口-这周新增/活跃接口
	respData2.AllApi = int(pipiAllCount)
	if pipiAllCount == 0 { //获取的接口总数为0 兼容 大于1 兼容
		respData2.LastWeekDegreeOfCompletion = "0%"
	} else if (float64(respData2.AllCaseCount) / pipiAllCount) > 1 {
		respData2.LastWeekDegreeOfCompletion = "100%"
	} else {
		respData2.LastWeekDegreeOfCompletion = strLast2 + "%"
	}
	respDataList = append(respDataList, respData2)

	respData3 := respData
	respData3.BusinessName = "海外"
	respData3.AllApiCount = float64(len(noRepeatHaiwaiList))
	respData3.NewApiConut = float64(resp2["haiwai_count_new"])
	respData3.AllCaseCount = len(haiwai_list)
	respData3.NewCaseConut = resp2["haiwai_new_case"]
	str3 := strconv.FormatFloat(float64(float64(respData3.AllCaseCount)/haiwaiAllCount)*100, 'f', 2, 64)
	respData3.DegreeOfCompletion = str3 + "%"
	respData3.AllApi = int(haiwaiAllCount)
	if haiwaiAllCount == 0 { //获取的接口总数为0 兼容 大于1 兼容
		respData3.DegreeOfCompletion = "0%"
	} else if (float64(respData3.AllCaseCount) / haiwaiAllCount) > 1 {
		respData3.DegreeOfCompletion = "100%"
	} else {
		respData3.DegreeOfCompletion = str3 + "%"
	}

	strLast3 := strconv.FormatFloat(float64(float64(respData3.AllCaseCount-respData3.NewCaseConut)/haiwaiAllCount)*100, 'f', 2, 64) //上周完成度  全部接口-这周新增/活跃接口
	respData3.AllApi = int(haiwaiAllCount)
	if haiwaiAllCount == 0 { //获取的接口总数为0 兼容 大于1 兼容
		respData3.LastWeekDegreeOfCompletion = "0%"
	} else if (float64(respData3.AllCaseCount) / haiwaiAllCount) > 1 {
		respData3.LastWeekDegreeOfCompletion = "100%"
	} else {
		respData3.LastWeekDegreeOfCompletion = strLast3 + "%"
	}

	respDataList = append(respDataList, respData3)

	respData4 := respData
	respData4.BusinessName = "中东"
	respData4.AllApiCount = float64(len(noRepeatZhongdongList))
	respData4.NewApiConut = float64(resp2["zhongdong_count_new"])
	respData4.AllCaseCount = len(zhongdong_list)
	respData4.NewCaseConut = resp2["zhongdong_new_case"]
	str4 := strconv.FormatFloat(float64(float64(respData4.AllCaseCount)/zhongdongAllCount)*100, 'f', 2, 64)
	respData4.DegreeOfCompletion = str4 + "%"
	respData4.AllApi = int(zhongdongAllCount)
	if zhongdongAllCount == 0 { //获取的接口总数为0 兼容 大于1 兼容
		respData4.DegreeOfCompletion = "0%"
	} else if (float64(respData4.AllCaseCount) / zhongdongAllCount) > 1 {
		respData4.DegreeOfCompletion = "100%"
	} else {
		respData4.DegreeOfCompletion = str4 + "%"
	}

	strLast4 := strconv.FormatFloat(float64(float64(respData4.AllCaseCount-respData4.NewCaseConut)/zhongdongAllCount)*100, 'f', 2, 64) //上周完成度  全部接口-这周新增/活跃接口
	respData4.AllApi = int(zhongdongAllCount)
	if zhongdongAllCount == 0 { //获取的接口总数为0 兼容 大于1 兼容
		respData4.LastWeekDegreeOfCompletion = "0%"
	} else if (float64(respData4.AllCaseCount) / zhongdongAllCount) > 1 {
		respData4.LastWeekDegreeOfCompletion = "100%"
	} else {
		respData4.LastWeekDegreeOfCompletion = strLast4 + "%"
	}

	respDataList = append(respDataList, respData4)

	respData5 := respData
	respData5.BusinessName = "麻团"
	respData5.AllApiCount = float64(len(noRepeatMatuanList))
	respData5.NewApiConut = float64(resp2["matuan_count_new"])
	respData5.AllCaseCount = len(matuan_list)
	respData5.NewCaseConut = resp2["matuan_new_case"]
	str5 := strconv.FormatFloat(float64(float64(respData5.AllCaseCount)/matuanAllConut)*100, 'f', 2, 64)
	respData5.DegreeOfCompletion = str5 + "%"
	respData5.AllApi = int(matuanAllConut)

	//respDataList = append(respDataList, respData5)

	respData6 := respData
	respData6.BusinessName = "商业化"
	respData6.AllApiCount = float64(len(noRepeatShangyehuaList))
	respData6.NewApiConut = float64(resp2["shangyehua_count_new"])
	respData6.AllCaseCount = len(shangyehuai_list)
	respData6.NewCaseConut = resp2["shangyehua_new_case"]
	str6 := strconv.FormatFloat(float64(float64(respData6.AllCaseCount)/shangyehuaAllCount)*100, 'f', 2, 64)
	respData6.DegreeOfCompletion = str6 + "%"
	respData6.AllApi = int(shangyehuaAllCount)
	if shangyehuaAllCount == 0 { //获取的接口总数为0 兼容 大于1 兼容
		respData6.DegreeOfCompletion = "0%"
	} else if (float64(respData6.AllCaseCount) / shangyehuaAllCount) > 1 {
		respData6.DegreeOfCompletion = "100%"
	} else {
		respData6.DegreeOfCompletion = str6 + "%"
	}

	strLast6 := strconv.FormatFloat(float64(float64(respData6.AllCaseCount-respData6.NewCaseConut)/shangyehuaAllCount)*100, 'f', 2, 64) //上周完成度  全部接口-这周新增/活跃接口
	respData6.AllApi = int(shangyehuaAllCount)
	if shangyehuaAllCount == 0 { //获取的接口总数为0 兼容 大于1 兼容
		respData6.LastWeekDegreeOfCompletion = "0%"
	} else if (float64(respData6.AllCaseCount) / shangyehuaAllCount) > 1 {
		respData6.LastWeekDegreeOfCompletion = "100%"
	} else {
		respData6.LastWeekDegreeOfCompletion = strLast6 + "%"
	}

	respDataList = append(respDataList, respData6)

	respData7 := respData
	respData7.BusinessName = "海外-US"
	respData7.AllApiCount = float64(len(noRepeatHaiwaiUSList))
	respData7.NewApiConut = float64(resp2["haiwaiUS_count_new"])
	respData7.AllCaseCount = len(haiwaiUS_list)
	respData7.NewCaseConut = resp2["haiwaiUS_new_case"]
	str7 := strconv.FormatFloat(float64(float64(respData7.AllCaseCount)/haiwaiUSAllCount)*100, 'f', 2, 64)
	respData7.DegreeOfCompletion = str7 + "%"
	respData7.AllApi = int(haiwaiUSAllCount)
	if haiwaiUSAllCount == 0 { //获取的接口总数为0 兼容 大于1 兼容
		respData7.DegreeOfCompletion = "0%"
	} else if (float64(respData7.AllCaseCount) / haiwaiUSAllCount) > 1 {
		respData7.DegreeOfCompletion = "100%"
	} else {
		respData7.DegreeOfCompletion = str7 + "%"
	}

	strLast7 := strconv.FormatFloat(float64(float64(respData7.AllCaseCount-respData7.NewCaseConut)/haiwaiUSAllCount)*100, 'f', 2, 64) //上周完成度  全部接口-这周新增/活跃接口
	respData7.AllApi = int(haiwaiUSAllCount)
	if haiwaiUSAllCount == 0 { //获取的接口总数为0 兼容 大于1 兼容
		respData7.LastWeekDegreeOfCompletion = "0%"
	} else if (float64(respData7.AllCaseCount) / haiwaiUSAllCount) > 1 {
		respData7.LastWeekDegreeOfCompletion = "100%"
	} else {
		respData7.LastWeekDegreeOfCompletion = strLast7 + "%"
	}

	respDataList = append(respDataList, respData7)
	return respDataList

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
		userTimeF := getFridayTime(useTime)
		for _, one := range result {
			caseCreateTime, _ := time.ParseInLocation(Layout, one.CreatedAt, time.Local) //获取报告的创建时间 转换为data
			api := strings.Split(one.ApiUrl, "?")[0]
			switch one.BusinessCode {
			case "0": //最右
				if caseCreateTime.After(userTimeF.AddDate(0, 0, -7)) { //如果是周五，则获取前7天报告
					zuiyou_list_new = append(zuiyou_list_new, api)
				}
			case "1": //皮皮
				if caseCreateTime.After(userTimeF.AddDate(0, 0, -7)) { //如果是周五，则获取前7天报告
					pipi_list_new = append(pipi_list_new, api)
				}
			case "2": //海外
				if caseCreateTime.After(userTimeF.AddDate(0, 0, -7)) { //如果是周五，则获取前7天报告
					haiwai_list_new = append(haiwai_list_new, api)
				}
			case "3": //中东
				if caseCreateTime.After(userTimeF.AddDate(0, 0, -7)) { //如果是周五，则获取前7天报告
					zhongdong_list_new = append(zhongdong_list_new, api)
				}
			case "4": //麻团
				if caseCreateTime.After(userTimeF.AddDate(0, 0, -7)) { //如果是周五，则获取前7天报告
					shangyehuai_list_new = append(shangyehuai_list_new, api)
				}
			case "5": //商业化
				if caseCreateTime.After(userTimeF.AddDate(0, 0, -7)) { //如果是周五，则获取前7天报告
					matuan_list_new = append(matuan_list_new, api)
				}
			case "6": //海外-us
				if caseCreateTime.After(userTimeF.AddDate(0, 0, -7)) { //如果是周五，则获取前7天报告
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

func (c *StatisticsController) getAllQuery() {
	mongo := models.StatisticsMongo{}
	list, count, err := mongo.QueryAll()
	if err != nil {
		logs.Error("查询出错", err)
	}
	c.FormSuccessJson(int64(count), list)

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
		nowTimeStr := nowTime.Format("2006-01-02")
		nowTime0, _ := time.Parse("2006-01-02", nowTimeStr)
		return nowTime0
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

func getAllApi() map[string]float64 {
	data := make(map[string]float64)
	//data["0"] = 200
	//data["1"] = 200
	//data["2"] = 200
	//data["3"] = 200
	//data["4"] = 200
	//data["5"] = 200
	//data["6"] = 200

	return data

}
