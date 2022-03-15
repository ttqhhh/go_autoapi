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
		apiList, err := mongo.QueryByBusiness(business_code)
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
			f.SetCellValue("Sheet1", "D1", "是否使用")

			line++
			f.SetCellValue("Sheet1", fmt.Sprintf("A%d", line), oneApi.Id)
			f.SetCellValue("Sheet1", fmt.Sprintf("B%d", line), oneApi.BusinessName)
			f.SetCellValue("Sheet1", fmt.Sprintf("C%d", line), oneApi.ApiName)
			f.SetCellValue("Sheet1", fmt.Sprintf("D%d", line), oneApi.Use)

			// 保存文件

		}
		file := strconv.FormatInt(business_code, 10)
		filename := file + ".xlsx"
		if err := f.SaveAs("/Users/tangtianqing/Desktop/active_api/" + filename); err != nil {
			fmt.Println(err)
		}

	}

}
