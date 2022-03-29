package api

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/prometheus/common/log"
	"github.com/xuri/excelize/v2"
	"go_autoapi/constants"
	"go_autoapi/libs"
	"go_autoapi/models"
	"strconv"
	"strings"
)

type GetDataController struct {
	libs.BaseController
}

func (c *GetDataController) Get() {
	do := c.GetMethodName()
	switch do {
	case "read_to_data":
		c.ReadToData() //从excal中读数据

	case "wirte_to_data":
		c.WriteToData() //将数据写入excal中

	case "wirte_to_data_fei":
		c.WriteToDataFeiLei()
	default:
		logs.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持", nil)
	}
}

func (c *GetDataController) Post() {
	do := c.GetMethodName()
	switch do {
	case "login":
	default:
		logs.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持", nil)
	}
}

func (c *GetDataController) ReadToData() {
	fileNameList := []string{
		"最右.xlsx", "皮皮.xlsx", "海外.xlsx", "中东.xlsx", "商业化.xlsx", "海外US.xlsx",
	}
	businessCode := []int64{
		0, 1, 2, 3, 5, 6,
	}
	var business = 0
	for _, fileName := range fileNameList {
		f, err := excelize.OpenFile("/Users/tangtianqing/Desktop/" + fileName)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer func() {
			if err := f.Close(); err != nil {
				fmt.Println(err)
			}
		}()
		// 获取 Sheet1 上所有单元格
		rows, err := f.GetRows("Sheet1")
		if err != nil {
			fmt.Println(err)
			return
		}
		for _, row := range rows {
			dataList := []string{}
			for _, colCell := range row {
				dataList = append(dataList, colCell)
			}
			mongo := models.AllActiveApiMongo{}
			mongo.Id, err = strconv.ParseInt(dataList[0], 10, 64)
			mongo.BusinessName = dataList[1]
			mongo.BusinessCode = businessCode[business]
			mongo.ApiName = dataList[2]
			mongo.Use, err = strconv.ParseInt(dataList[3], 10, 64)
			mongo.Insert(mongo)

		}
		business++
	}
	c.SuccessJson("success")
}

func (c *GetDataController) WriteToData() {
	var numbers = []int64{
		constants.ZuiyYou, constants.PiPi, constants.HaiWai, constants.ZhongDong, constants.ShangYeHua, constants.HaiWaiUS,
	}
	for _, business_code := range numbers {
		mongo := models.AllActiveApiMongo{}
		apiList, err := mongo.QueryByBusinessAll(business_code)
		if err != nil {
			log.Error("根据业务线获取全部接口出错")
		}
		line := 1
		f := excelize.NewFile() // 设置单元格的值
		for _, oneApi := range apiList {
			//开始创建excal 写入数据
			// 这里设置表头
			f.SetCellValue("Sheet1", "A1", "id")
			f.SetCellValue("Sheet1", "B1", "业务线")
			f.SetCellValue("Sheet1", "C1", "接口名")
			f.SetCellValue("Sheet1", "D1", "是否废弃（1=废弃）")
			f.SetCellValue("Sheet1", "E1", "是否被统计 1 = 参与统计")

			line++
			f.SetCellValue("Sheet1", fmt.Sprintf("A%d", line), oneApi.Id)
			f.SetCellValue("Sheet1", fmt.Sprintf("B%d", line), oneApi.BusinessName)
			f.SetCellValue("Sheet1", fmt.Sprintf("C%d", line), oneApi.ApiName)
			f.SetCellValue("Sheet1", fmt.Sprintf("D%d", line), oneApi.Use)
			f.SetCellValue("Sheet1", fmt.Sprintf("E%d", line), oneApi.Calculate)

			// 保存文件

		}
		file := strconv.FormatInt(business_code, 10)
		filename := file + ".xlsx"
		if err := f.SaveAs("/Users/tangtianqing/Desktop/active_api/" + filename); err != nil {
			fmt.Println(err)
		}

	}
	c.SuccessJson(nil)
}

func (c *GetDataController) WriteToDataFeiLei() {
	isUse, err := c.GetInt64("kind")
	if err != nil {
		logs.Info("类型获取错误")
	}

	var numbers = []int64{
		constants.ZuiyYou, constants.PiPi, constants.HaiWai, constants.ZhongDong, constants.ShangYeHua, constants.HaiWaiUS,
	}
	for _, business_code := range numbers {
		if isUse == 1 {
			mongo := models.AllActiveApiMongo{}
			apiList, err := mongo.QueryByBusiness(business_code)
			if err != nil {
				log.Error("根据业务线获取全部接口出错")
			}
			line := 1
			f := excelize.NewFile() // 设置单元格的值
			for _, oneApi := range apiList {
				if isUse == 1 {
					//开始创建excal 写入数据
					// 这里设置表头
					f.SetCellValue("Sheet1", "A1", "id")
					f.SetCellValue("Sheet1", "B1", "业务线")
					f.SetCellValue("Sheet1", "C1", "接口名")
					f.SetCellValue("Sheet1", "D1", "是否被统计")
					if oneApi.Calculate == 1 { //统计被统计的接口
						line++
						f.SetCellValue("Sheet1", fmt.Sprintf("A%d", line), oneApi.Id)
						f.SetCellValue("Sheet1", fmt.Sprintf("B%d", line), oneApi.BusinessName)
						f.SetCellValue("Sheet1", fmt.Sprintf("C%d", line), oneApi.ApiName)
						f.SetCellValue("Sheet1", fmt.Sprintf("D%d", line), oneApi.Calculate)

					}
				}
				// 保存文件
			}
			file := strconv.FormatInt(business_code, 10)
			filename := file + ".xlsx"
			if err := f.SaveAs("/Users/tangtianqing/Desktop/active_api/" + filename); err != nil {
				fmt.Println(err)
			}

		} else {
			line := 1
			f := excelize.NewFile() // 设置单元格的值
			listmap := GeiApi()
			acm := models.AllActiveApiMongo{}
			for _, one := range listmap[business_code] {
				_, isExist := acm.NewApiIsInDatabase(one, business_code)
				if isExist == false {
					//开始创建excal 写入数据
					// 这里设置表头
					f.SetCellValue("Sheet1", "A1", "业务线")
					f.SetCellValue("Sheet1", "B1", "接口名")
					f.SetCellValue("Sheet1", "C1", "是否被统计")
					line++
					f.SetCellValue("Sheet1", fmt.Sprintf("A%d", line), business_code)
					f.SetCellValue("Sheet1", fmt.Sprintf("B%d", line), one)
					f.SetCellValue("Sheet1", fmt.Sprintf("C%d", line), 0)

				}

			}
			file := strconv.FormatInt(business_code, 10)
			filename := file + ".xlsx"
			if err := f.SaveAs("/Users/tangtianqing/Desktop/active_api/" + filename); err != nil {
				fmt.Println(err)
			}

		}

	}
	c.SuccessJson(nil)
}

func GeiApi() map[int64][]string {
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
	noRepeatShangyehuaList := RemoveRepeatedElement(shangyehuai_list)
	noRepeatHaiwaiUSList := RemoveRepeatedElement(haiwaiUS_list)

	listMap := make(map[int64][]string)
	listMap[0] = noRepeatZuiyouList
	listMap[1] = noRepeatPipiList
	listMap[2] = noRepeatHaiwaiList
	listMap[3] = noRepeatZhongdongList
	listMap[5] = noRepeatShangyehuaList
	listMap[6] = noRepeatHaiwaiUSList
	return listMap
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
