package statisticsdata

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"
	"github.com/prometheus/common/log"
	"go_autoapi/constants"
	"go_autoapi/libs"
	"go_autoapi/models"
	"net/http"
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
	AllApi                     int     `form:"all_api" json:"all_api"`
	DegreeOfCompletion         string  `form:"degree_of_completion" json:"degree_of_completion"`                     //完成度
	LastWeekDegreeOfCompletion string  `form:"last_week_degree_of_completion" json:"last_week_degree_of_completion"` //上周完成度
	UnUseApi                   int     `form:"un_use_api" json:"un_use_api"`                                         //废弃接口数

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

func (c *StatisticsController) GetAllApiGroupByBusiness() []respData {
	// 1.去平台获取全部活跃接口数  并且增量添加
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
	//-----------------------------------------------------------------------------------
	//2.通过所有业务线 获取所有case 并且吧url切出来去重
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
	//-------------------------------------------------------------------
	//3. 获取本周新增数据（这里逻辑不要动 原来是实时修改数据 后来改为定时任务）
	resp2 := c.getApiByBusinessNewAdd()
	//--------------------------------------------------------------------
	respDataList := []respData{}      //生命一个存放对象对数组
	acm := models.AllActiveApiMongo{} //全局定义对象
	//--------------------------------------------------------------------
	//重点计算 给结构体值
	respData := respData{}
	respData1 := respData
	respData1.BusinessName = "最右"
	respData1.AllCaseCount = len(zuiyou_list)
	respData1.NewCaseConut = resp2["zuiyou_new_case"]
	var EffectiveApiZY = 0           //初始化有效接口为0
	acm = models.AllActiveApiMongo{} //实例化这个对象 使用他的方法来判断接口是否存在
	for _, one := range noRepeatZuiyouList {
		acm, isExist := acm.NewApiIsInDatabase(one, constants.ZuiyYou)
		if isExist == true {
			acm.Calculate = 0
			acm.ChangeApiCalculate(acm.Id, acm)
			fmt.Print("参与计算的接口：" + one)
			EffectiveApiZY++ //一旦该case对应的api存在数据库 证明该api为有效api 可以用来计算覆盖率
		}
	}
	respData1.UnUseApi = acm.GetAllUnUseApiCount(constants.ZuiyYou)                             //获取一个废弃数
	str := strconv.FormatFloat(float64(float64(EffectiveApiZY)/zuiyouAllCount)*100, 'f', 2, 64) //本周完成度
	respData1.AllApi = int(zuiyouAllCount)
	if zuiyouAllCount == 0 { //获取的接口总数为0 兼容 大于1 兼容
		respData1.DegreeOfCompletion = "0%"
	} else if (float64(EffectiveApiZY) / zuiyouAllCount) > 1 {
		respData1.DegreeOfCompletion = "100%"
	} else {
		respData1.DegreeOfCompletion = str + "%"
	}
	//上周覆盖率 todo 想想怎么处理
	ac := models.StatisticsMongo{}
	lastzuiyou, err := ac.QueryLast6ByBusiness("最右")
	if err != nil {
		logs.Error("从数据库查询后六条中数据出错")
	}
	respData1.LastWeekDegreeOfCompletion = lastzuiyou.DegreeOfCompletion
	respDataList = append(respDataList, respData1)

	//皮皮
	respData2 := respData
	respData2.BusinessName = "皮皮"
	respData2.AllApiCount = float64(len(noRepeatPipiList))
	respData2.NewApiConut = float64(resp2["pipi_count_new"])
	respData2.AllCaseCount = len(pipi_list)
	respData2.NewCaseConut = resp2["pipi_new_case"]
	var EffectiveApiPP = 0           //初始化有效接口为0
	acm = models.AllActiveApiMongo{} //实例化这个对象 使用他的方法来判断接口是否存在
	for _, one := range noRepeatPipiList {
		acm, isExist := acm.NewApiIsInDatabase(one, constants.PiPi)
		if isExist == true {
			acm.Calculate = 0
			acm.ChangeApiCalculate(acm.Id, acm)
			fmt.Print("参与计算的接口：" + one)
			EffectiveApiPP++ //一旦该case对应的api存在数据库 证明该api为有效api 可以用来计算覆盖率
		}
	}
	respData2.UnUseApi = acm.GetAllUnUseApiCount(constants.PiPi) //获取一个废弃数
	str2 := strconv.FormatFloat(float64(float64(EffectiveApiPP)/pipiAllCount)*100, 'f', 2, 64)
	respData2.DegreeOfCompletion = str2 + "%"
	respData2.AllApi = int(pipiAllCount)
	if pipiAllCount == 0 { //获取的接口总数为0 兼容 大于1 兼容
		respData2.DegreeOfCompletion = "0%"
	} else if (float64(EffectiveApiPP) / pipiAllCount) > 1 {
		respData2.DegreeOfCompletion = "100%"
	} else {
		respData2.DegreeOfCompletion = str2 + "%"

	}
	lastPIPI, err := ac.QueryLast6ByBusiness("皮皮")
	if err != nil {
		logs.Error("从数据库查询后六条中数据出错")
	}
	respData2.LastWeekDegreeOfCompletion = lastPIPI.DegreeOfCompletion
	respDataList = append(respDataList, respData2)
	//海外
	respData3 := respData
	respData3.BusinessName = "海外"
	respData3.AllApiCount = float64(len(noRepeatHaiwaiList))
	respData3.NewApiConut = float64(resp2["haiwai_count_new"])
	respData3.AllCaseCount = len(haiwai_list)
	respData3.NewCaseConut = resp2["haiwai_new_case"]
	var EffectiveApiHW = 0           //初始化有效接口为0
	acm = models.AllActiveApiMongo{} //实例化这个对象 使用他的方法来判断接口是否存在
	for _, one := range noRepeatHaiwaiList {
		acm, isExist := acm.NewApiIsInDatabase(one, constants.HaiWai)
		if isExist == true {
			acm.Calculate = 0
			acm.ChangeApiCalculate(acm.Id, acm)
			EffectiveApiHW++ //一旦该case对应的api存在数据库 证明该api为有效api 可以用来计算覆盖率
		}
	}
	respData3.UnUseApi = acm.GetAllUnUseApiCount(constants.HaiWai) //获取一个废弃数
	str3 := strconv.FormatFloat(float64(float64(EffectiveApiHW)/haiwaiAllCount)*100, 'f', 2, 64)
	respData3.DegreeOfCompletion = str3 + "%"
	respData3.AllApi = int(haiwaiAllCount)
	if haiwaiAllCount == 0 { //获取的接口总数为0 兼容 大于1 兼容
		respData3.DegreeOfCompletion = "0%"
	} else if (float64(EffectiveApiHW) / haiwaiAllCount) > 1 {
		respData3.DegreeOfCompletion = "100%"
	} else {
		respData3.DegreeOfCompletion = str3 + "%"
	}

	lasthaiwai, err := ac.QueryLast6ByBusiness("海外")
	if err != nil {
		logs.Error("从数据库查询后六条中数据出错")
	}
	respData3.LastWeekDegreeOfCompletion = lasthaiwai.DegreeOfCompletion

	respDataList = append(respDataList, respData3)

	respData4 := respData
	respData4.BusinessName = "中东"
	respData4.AllApiCount = float64(len(noRepeatZhongdongList))
	respData4.NewApiConut = float64(resp2["zhongdong_count_new"])
	respData4.AllCaseCount = len(zhongdong_list)
	respData4.NewCaseConut = resp2["zhongdong_new_case"]
	var EffectiveApiZD = 0           //初始化有效接口为0
	acm = models.AllActiveApiMongo{} //实例化这个对象 使用他的方法来判断接口是否存在
	for _, one := range noRepeatZhongdongList {
		acm, isExist := acm.NewApiIsInDatabase(one, constants.ZhongDong)
		if isExist == true {
			acm.Calculate = 0
			acm.ChangeApiCalculate(acm.Id, acm)
			EffectiveApiZD++ //一旦该case对应的api存在数据库 证明该api为有效api 可以用来计算覆盖率
		}
	}
	respData4.UnUseApi = acm.GetAllUnUseApiCount(constants.ZhongDong) //获取一个废弃数
	str4 := strconv.FormatFloat(float64(float64(EffectiveApiZD)/zhongdongAllCount)*100, 'f', 2, 64)
	respData4.DegreeOfCompletion = str4 + "%"
	respData4.AllApi = int(zhongdongAllCount)
	if zhongdongAllCount == 0 { //获取的接口总数为0 兼容 大于1 兼容
		respData4.DegreeOfCompletion = "0%"
	} else if (float64(EffectiveApiZD) / zhongdongAllCount) > 1 {
		respData4.DegreeOfCompletion = "100%"
	} else {
		respData4.DegreeOfCompletion = str4 + "%"

	}
	lastzd, err := ac.QueryLast6ByBusiness("中东")
	if err != nil {
		logs.Error("从数据库查询后六条中数据出错")
	}
	respData4.LastWeekDegreeOfCompletion = lastzd.DegreeOfCompletion
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
	var EffectiveApiSYH = 0          //初始化有效接口为0
	acm = models.AllActiveApiMongo{} //实例化这个对象 使用他的方法来判断接口是否存在
	for _, one := range noRepeatShangyehuaList {
		acm, isExist := acm.NewApiIsInDatabase(one, constants.ShangYeHua)
		if isExist == true {
			acm.Calculate = 0
			acm.ChangeApiCalculate(acm.Id, acm)
			EffectiveApiSYH++ //一旦该case对应的api存在数据库 证明该api为有效api 可以用来计算覆盖率
		}
	}
	respData6.UnUseApi = acm.GetAllUnUseApiCount(constants.ShangYeHua) //获取一个废弃数
	str6 := strconv.FormatFloat(float64(float64(EffectiveApiSYH)/shangyehuaAllCount)*100, 'f', 2, 64)
	respData6.DegreeOfCompletion = str6 + "%"
	respData6.AllApi = int(shangyehuaAllCount)
	if shangyehuaAllCount == 0 { //获取的接口总数为0 兼容 大于1 兼容
		respData6.DegreeOfCompletion = "0%"
	} else if (float64(EffectiveApiSYH) / shangyehuaAllCount) > 1 {
		respData6.DegreeOfCompletion = "100%"
	} else {
		respData6.DegreeOfCompletion = str6 + "%"

	}
	lastsyh, err := ac.QueryLast6ByBusiness("商业化")
	if err != nil {
		logs.Error("从数据库查询后六条中数据出错")
	}
	respData6.LastWeekDegreeOfCompletion = lastsyh.DegreeOfCompletion
	respDataList = append(respDataList, respData6)

	respData7 := respData
	respData7.BusinessName = "海外-US"
	respData7.AllApiCount = float64(len(noRepeatHaiwaiUSList))
	respData7.NewApiConut = float64(resp2["haiwaiUS_count_new"])
	respData7.AllCaseCount = len(haiwaiUS_list)
	respData7.NewCaseConut = resp2["haiwaiUS_new_case"]
	var EffectiveApiHWUS = 0         //初始化有效接口为0
	acm = models.AllActiveApiMongo{} //实例化这个对象 使用他的方法来判断接口是否存在
	for _, one := range noRepeatHaiwaiUSList {
		acm, isExist := acm.NewApiIsInDatabase(one, constants.HaiWaiUS)
		if isExist == true {
			acm.Calculate = 0
			acm.ChangeApiCalculate(acm.Id, acm)
			EffectiveApiHWUS++ //一旦该case对应的api存在数据库 证明该api为有效api 可以用来计算覆盖率
		}
	}
	respData7.UnUseApi = acm.GetAllUnUseApiCount(constants.HaiWaiUS) //获取一个废弃数
	str7 := strconv.FormatFloat(float64(float64(EffectiveApiHWUS)/haiwaiUSAllCount)*100, 'f', 2, 64)
	respData7.DegreeOfCompletion = str7 + "%"
	respData7.AllApi = int(haiwaiUSAllCount)
	if haiwaiUSAllCount == 0 { //获取的接口总数为0 兼容 大于1 兼容
		respData7.DegreeOfCompletion = "0%"
	} else if (float64(EffectiveApiHWUS) / haiwaiUSAllCount) > 1 {
		respData7.DegreeOfCompletion = "100%"
	} else {
		respData7.DegreeOfCompletion = str7 + "%"

	}
	//上周新增 在这里直接查后6条 因为这时还没入库
	lastHwus, err := ac.QueryLast6ByBusiness("海外-US")
	if err != nil {
		logs.Error("从数据库查询后六条中数据出错")
	}
	respData7.LastWeekDegreeOfCompletion = lastHwus.DegreeOfCompletion
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
	cookie := getLogin()
	cookiehaiwai := getLoginHaiWai()
	cookiehaiwaius := getLoginHaiWaiUS()
	cookiezd := getLoginZhongDong()

	zuiyou := getZyAllApiCount(cookie)
	pipi := getPPAllApiCount(cookie)
	haiwai := gethaiwaiAllApiCount(cookiehaiwai)
	zd := getzhongdongAllApiCount(cookiezd)
	shangyehua := getSyhAllApiCount(cookie)
	haiwaius := gethaiwaiUSAllApiCount(cookiehaiwaius)
	data["0"] = zuiyou
	data["1"] = pipi
	data["2"] = haiwai
	data["3"] = zd  //中东
	data["4"] = 200 //麻团
	data["5"] = shangyehua
	data["6"] = haiwaius
	return data

}

//-----------------------------------------------------------------

//用到的常量
const ZY_grafana_login_url = "http://grafana.ixiaochuan.cn/login"
const HW_grafana_login_url = "http://dashboard.icocofun.net/login"
const HWUS_grafana_login_url = "http://grafanaus.icocofun.net/login"
const ZD_grafana_login_url = "http://grafana.mehiya.com/login"
const AD_gateway_path_url = "http://grafana.ixiaochuan.cn/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_gateway_ad_http_latency_count%5B1m%5D))by(uri)&start=1646064000&end=1648742400&step=7200"

//解析json所需要的结构体
type JsonData struct {
	Status string `json:"status"`
	Data   Data   `json:"data"`
}

type Data struct {
	ResultType string   `json:"resultType"`
	Result     []Result `json:"result"`
}

type Result struct {
	Meturc Metruc          `json:"metric"`
	Values [][]interface{} `json:"values"`
}

type Metruc struct {
	Uri string `json:"uri"`
}

// 4个登录获取cookie的方法
func getLogin() *http.Cookie {
	req := httplib.Post(ZY_grafana_login_url)
	req.Param("user", "wangzhen01")
	req.Param("password", "Iepohg5go4iawoo")
	req.Param("email", "")
	resp, err := req.Response()
	if err != nil {
	}
	cookies := resp.Cookies()
	for _, cookie := range cookies {
		return cookie
	}
	return nil
}

func getLoginHaiWai() *http.Cookie {
	req := httplib.Post(HW_grafana_login_url)
	req.Param("user", "wangzhen01")
	req.Param("password", "Iepohg5go4iawoo")
	resp, err := req.Response()
	if err != nil {
	}
	cookies := resp.Cookies()
	for _, cookie := range cookies {
		return cookie
	}
	return nil
}

func getLoginHaiWaiUS() *http.Cookie {
	req := httplib.Post(HWUS_grafana_login_url)
	req.Param("user", "wangzhen01")
	req.Param("password", "Iepohg5go4iawoo")
	resp, err := req.Response()
	if err != nil {
	}
	cookies := resp.Cookies()
	for _, cookie := range cookies {
		return cookie
	}
	return nil
}

func getLoginZhongDong() *http.Cookie {
	req := httplib.Post(ZD_grafana_login_url)
	req.Param("user", "chengxiaojing")
	req.Param("password", "ls51HGb8y0MA")
	resp, err := req.Response()
	if err != nil {
	}
	cookies := resp.Cookies()
	for _, cookie := range cookies {
		return cookie
	}
	return nil
}

//获取各个业务线活跃接口的方法

func getSyhAllApiCount(cookie *http.Cookie) float64 {
	respData := doReq(AD_gateway_path_url, cookie)
	toJson := JsonData{}
	err := json.Unmarshal([]byte(respData), &toJson)
	if err != nil {

	}
	acm := models.AllActiveApiMongo{}
	count := acm.QueryAllCountByBusinessCount(constants.ShangYeHua)
	//todo tangtianqing  后续增量查询需要开放这段代码
	for _, one := range toJson.Data.Result {
		for _, ones := range one.Values {
			if ones[1] != "0" {
				acm := models.AllActiveApiMongo{}
				acm, isEsixt := acm.NewApiIsInDatabase(one.Meturc.Uri, constants.ShangYeHua) //传入 api_name business 查看是否存在
				if isEsixt == false {
					acm.BusinessName = "商业化"
					acm.BusinessCode = constants.ShangYeHua
					acm.ApiName = one.Meturc.Uri
					acm.Insert(acm)
					count++
				}
				break
			}
		}
	}
	return float64(count)
}

func getZyAllApiCount(cookie *http.Cookie) float64 {
	zuiyouURLlsit := []string{
		"http://grafana.ixiaochuan.cn/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_gateway_rec_http_latency_count%5B1m%5D))by(uri)&start=1646064000&end=1648742400&step=7200",
		"http://grafana.ixiaochuan.cn/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_gateway_topic_http_latency_count%5B1m%5D))by(uri)&start=1646064000&end=1648742400&step=7200",
		"http://grafana.ixiaochuan.cn/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_gateway_post_http_latency_count%5B1m%5D))by(uri)&start=1646064000&end=1648742400&step=7200",
		"http://grafana.ixiaochuan.cn/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_gateway_rev_http_latency_count%5B1m%5D))by(uri)&start=1646064000&end=1648742400&step=7200",
		"http://grafana.ixiaochuan.cn/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_gateway_acnt_http_latency_count%5B1m%5D))by(uri)&start=1646064000&end=1648742400&step=7200",
		"http://grafana.ixiaochuan.cn/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_gateway_cfg_http_latency_count%5B1m%5D))by(uri)&start=1646064000&end=1648742400&step=7200",
		"http://grafana.ixiaochuan.cn/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_gateway_danmaku_http_latency_count%5B1m%5D))by(uri)&start=1646064000&end=1648742400&step=7200",
		"http://grafana.ixiaochuan.cn/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_gateway_misc_http_latency_count%5B1m%5D))by(uri)&start=1646064000&end=1648742400&step=7200",
		"http://grafana.ixiaochuan.cn/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_snssrv_gateway_native_http_latency_count%5B1m%5D))by(uri)&start=1646064000&end=1648742400&step=7200",
		"http://grafana.ixiaochuan.cn/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_snssrv_gateway-nearby_http_latency_count%5B1m%5D))by(uri)&start=1646064000&end=1648742400&step=7200",
		"http://grafana.ixiaochuan.cn/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_zy_gateway_teamchat_http_latency_count%5B1m%5D))by(uri)&start=1646064000&end=1648742400&step=7200",
		"http://grafana.ixiaochuan.cn/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_chatsrv_gateway_http_latency_count%5B1m%5D))by(uri)&start=1646064000&end=1648742400&step=7200",
	}
	acm := models.AllActiveApiMongo{}
	count := acm.QueryAllCountByBusinessCount(constants.ZuiyYou)
	for _, i := range zuiyouURLlsit { //todo 增量查询开放
		respData := doReq(i, cookie)
		toJson := JsonData{}
		err := json.Unmarshal([]byte(respData), &toJson)
		if err != nil {
			log.Error("转换类型出错", err)
		}
		for _, one := range toJson.Data.Result {
			for _, ones := range one.Values {
				if ones[1] != "0" {
					acm := models.AllActiveApiMongo{}
					acm, isEsixt := acm.NewApiIsInDatabase(one.Meturc.Uri, constants.ZuiyYou) //传入 api_name business 查看是否存在
					if isEsixt == false {
						acm.BusinessName = "最右"
						acm.BusinessCode = constants.ZuiyYou
						acm.ApiName = one.Meturc.Uri
						acm.Insert(acm)
						count++
					}
					break
				}
			}
		}

	}

	return float64(count)
}

func getPPAllApiCount(cookie *http.Cookie) float64 {
	pipiURLlsit := []string{
		"http://grafana.ixiaochuan.cn/api/datasources/proxy/5/api/v1/query_range?query=sum%20by(uri)%20(irate(xcmetrics_httpsrv_qps%7Bjob%3D%22pp-gateway-misc%22%7D%5B1m%5D))&start=1646064000&end=1648742400&step=7200",
		"http://grafana.ixiaochuan.cn/api/datasources/proxy/5/api/v1/query_range?query=sum%20by(uri)%20(irate(xcmetrics_httpsrv_qps%7Bjob%3D%22pp-gateway-acnt%22%7D%5B1m%5D))&start=1646064000&end=1648742400&step=7200",
		"http://grafana.ixiaochuan.cn/api/datasources/proxy/5/api/v1/query_range?query=sum%20by(uri)%20(irate(xcmetrics_httpsrv_qps%7Bjob%3D%22pp-gateway-internal%22%7D%5B1m%5D))&start=1646064000&end=1648742400&step=7200",
		"http://grafana.ixiaochuan.cn/api/datasources/proxy/5/api/v1/query_range?query=sum%20by(uri)%20(irate(xcmetrics_httpsrv_qps%7Bjob%3D%22pp-gateway-point%22%7D%5B1m%5D))&start=1646064000&end=1648742400&step=7200",
		"http://grafana.ixiaochuan.cn/api/datasources/proxy/5/api/v1/query_range?query=sum%20by(uri)%20(irate(xcmetrics_httpsrv_qps%7Bjob%3D%22pp-gateway-post%22%7D%5B1m%5D))&start=1646064000&end=1648742400&step=7200",
		"http://grafana.ixiaochuan.cn/api/datasources/proxy/5/api/v1/query_range?query=sum%20by(uri)%20(irate(xcmetrics_httpsrv_qps%7Bjob%3D%22pp-gateway-rec%22%7D%5B1m%5D))&start=1646064000&end=1648742400&step=7200",
		"http://grafana.ixiaochuan.cn/api/datasources/proxy/5/api/v1/query_range?query=sum%20by(uri)%20(irate(xcmetrics_httpsrv_qps%7Bjob%3D%22pp-gateway-review%22%7D%5B1m%5D))&start=1646064000&end=1648742400&step=7200",
		"http://grafana.ixiaochuan.cn/api/datasources/proxy/5/api/v1/query_range?query=sum%20by(uri)%20(irate(xcmetrics_httpsrv_qps%7Bjob%3D%22pp-gateway-topic%22%7D%5B1m%5D))&start=1646064000&end=1648742400&step=7200",
		"http://grafana.ixiaochuan.cn/api/datasources/proxy/5/api/v1/query_range?query=sum%20by(uri)%20(irate(xcmetrics_httpsrv_qps%7Bjob%3D%22pp-gateway-town%22%7D%5B1m%5D))&start=1646064000&end=1648742400&step=7200",
	}
	acm := models.AllActiveApiMongo{}
	count := acm.QueryAllCountByBusinessCount(constants.PiPi)
	for _, i := range pipiURLlsit {
		//print(i)
		respData := doReq(i, cookie)
		toJson := JsonData{}
		err := json.Unmarshal([]byte(respData), &toJson)
		if err != nil {

		}
		for _, one := range toJson.Data.Result {
			for _, ones := range one.Values {
				if ones[1] != "0" {
					acm := models.AllActiveApiMongo{}
					acm, isEsixt := acm.NewApiIsInDatabase(one.Meturc.Uri, constants.PiPi) //传入 api_name business 查看是否存在
					if isEsixt == false {
						acm.BusinessName = "皮皮"
						acm.BusinessCode = constants.PiPi
						acm.ApiName = one.Meturc.Uri
						acm.Insert(acm)
						count++
					}
					break
				}
			}
		}

	}

	return float64(count)
}

func gethaiwaiAllApiCount(cookie *http.Cookie) float64 {
	acm := models.AllActiveApiMongo{}
	count := acm.QueryAllCountByBusinessCount(constants.HaiWai)
	haiwaiURLlsit := []string{
		"http://dashboard.icocofun.net/api/datasources/proxy/25/api/v1/query_range?query=sum(rate(xms_omg_gateway_http_latency_count%5B1m%5D))by(uri)&start=1644288000&end=1646880000&step=1200&timeout=300s",
		"http://dashboard.icocofun.net/api/datasources/proxy/25/api/v1/query_range?query=sum(rate(xms_chatsrv_gateway_http_latency_count%5B1m%5D))by(uri)&start=1644288000&end=1646880000&step=1200&timeout=300s",
		"http://dashboard.icocofun.net/api/datasources/proxy/25/api/v1/query_range?query=sum(rate(xms_gateway_ad_http_latency_count%5B1m%5D))by(uri)&start=1644288000&end=1646880000&step=1200&timeout=300s",
		"http://dashboard.icocofun.net/api/datasources/proxy/25/api/v1/query_range?query=sum(rate(xms_omg_gateway_acnt_http_latency_count%5B1m%5D))by(uri)&start=1644288000&end=1646880000&step=1200&timeout=300s",
		"http://dashboard.icocofun.net/api/datasources/proxy/25/api/v1/query_range?query=sum(rate(xms_omg_gateway_index_http_latency_count%5B1m%5D))by(uri)&start=1644288000&end=1646880000&step=1200&timeout=300s",
		"http://dashboard.icocofun.net/api/datasources/proxy/25/api/v1/query_range?query=sum(rate(xms_omg_gateway_post_http_latency_count%5B1m%5D))by(uri)&start=1644288000&end=1646880000&step=1200&timeout=300s",
		"http://dashboard.icocofun.net/api/datasources/proxy/25/api/v1/query_range?query=sum(rate(xms_omg_gateway_review_http_latency_count%5B1m%5D))by(uri)&start=1644288000&end=1646880000&step=1200&timeout=300s",
		"http://dashboard.icocofun.net/api/datasources/proxy/25/api/v1/query_range?query=sum(rate(xms_omg_gateway_topic_http_latency_count%5B1m%5D))by(uri)&start=1644288000&end=1646880000&step=1200&timeout=300s",
	}
	for _, i := range haiwaiURLlsit {
		//print(i)
		respData := doReq(i, cookie)
		toJson := JsonData{}
		err := json.Unmarshal([]byte(respData), &toJson)
		if err != nil {

		}
		for _, one := range toJson.Data.Result {
			for _, ones := range one.Values {
				if ones[1] != "0" {
					acm := models.AllActiveApiMongo{}
					acm, isEsixt := acm.NewApiIsInDatabase(one.Meturc.Uri, constants.HaiWai) //传入 api_name business 查看是否存在
					if isEsixt == false {
						acm.BusinessName = "海外"
						acm.BusinessCode = constants.HaiWai
						acm.ApiName = one.Meturc.Uri
						acm.Insert(acm)
						count++
					}
					break
				}
			}
		}

	}

	return float64(count)
}

func gethaiwaiUSAllApiCount(cookie *http.Cookie) float64 {
	acm := models.AllActiveApiMongo{}
	count := acm.QueryAllCountByBusinessCount(constants.HaiWaiUS)
	haiwaiUSURLlsit := []string{
		"http://grafanaus.icocofun.net/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_maga_gateway_http_latency_count%7Buri!%3D%22%2Fhealthcheck%22%7D%5B1m%5D))by(uri)&start=1646496000&end=1646668680&step=120",
		"http://grafanaus.icocofun.net/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_maga_chatsrv_gateway_http_latency_count%7Buri!%3D%22%2Fhealthcheck%22%7D%5B1m%5D))by(uri)&start=1646496000&end=1646668680&step=120",
		"http://grafanaus.icocofun.net/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_maga_gateway_account_http_latency_count%7Buri!%3D%22%2Fhealthcheck%22%7D%5B1m%5D))by(uri)&start=1646496000&end=1646668680&step=120",
	}
	for _, i := range haiwaiUSURLlsit {
		//print(i)
		respData := doReq(i, cookie)
		toJson := JsonData{}
		err := json.Unmarshal([]byte(respData), &toJson)
		if err != nil {

		}
		for _, one := range toJson.Data.Result {
			for _, ones := range one.Values {
				if ones[1] != "0" {
					acm := models.AllActiveApiMongo{}
					acm, isEsixt := acm.NewApiIsInDatabase(one.Meturc.Uri, constants.HaiWaiUS) //传入 api_name business 查看是否存在
					if isEsixt == false {
						acm.BusinessName = "海外US"
						acm.BusinessCode = constants.HaiWaiUS
						acm.ApiName = one.Meturc.Uri
						acm.Insert(acm)
						count++
					}
					break
				}
			}
		}

	}

	return float64(count)
}

func getzhongdongAllApiCount(cookie *http.Cookie) float64 {

	acm := models.AllActiveApiMongo{}
	count := acm.QueryAllCountByBusinessCount(constants.ZhongDong)
	ZhongDongURLlsit := []string{
		"http://grafana.mehiya.com/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_me_live_chat-gateway_http_latency_count%5B1m%5D))by(uri)&start=1644303600&end=1646895600&step=1800",
		"http://grafana.mehiya.com/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_me_live_gamestore_gateway_http_latency_count%5B1m%5D))by(uri)&start=1644303600&end=1646895600&step=1800",
		"http://grafana.mehiya.com/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_me_live_trade_gateway_http_latency_count%5B1m%5D))by(uri)&start=1644303600&end=1646895600&step=1800",
		"http://grafana.mehiya.com/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_me_live_gateway_http_latency_count%5B1m%5D))by(uri)&start=1644303600&end=1646895600&step=1800",
	}
	for _, i := range ZhongDongURLlsit {
		//print(i)
		respData := doReq(i, cookie)
		toJson := JsonData{}
		err := json.Unmarshal([]byte(respData), &toJson)
		if err != nil {

		}
		for _, one := range toJson.Data.Result {
			for _, ones := range one.Values {
				if ones[1] != "0" {
					acm := models.AllActiveApiMongo{}
					acm, isEsixt := acm.NewApiIsInDatabase(one.Meturc.Uri, constants.ZhongDong) //传入 api_name business 查看是否存在
					if isEsixt == false {
						acm.BusinessName = "中东"
						acm.BusinessCode = constants.ZhongDong
						acm.ApiName = one.Meturc.Uri
						acm.Insert(acm)
						count++
					}
					break
				}
			}
		}

	}

	return float64(count)
}

func doReq(url string, cookie *http.Cookie) string {
	req := httplib.Get(url)
	req.SetCookie(cookie)
	str, err := req.String()
	if err != nil {
		logs.Error("请求失败，err:", err)
	}
	return str
}
