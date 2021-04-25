package controllers

import (
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"go_autoapi/models"
	"strconv"
)


func (c *CaseManageController) ShowCases(){
	c.Data["business"] = c.GetString("business")
	c.TplName = "case_manager.html"
}

func GetServiceList(business string)(service []models.ServiceMongo){
	var busCode = int8(0)
	if business == "zuiyou"{
		busCode = int8(0)
	}else if business == "pipi"{
		busCode = int8(1)
	}else if business == "haiwai"{
		busCode = int8(2)
	}else if business == "zhongdong"{
		busCode = int8(3)
	}else if business == "mama"{
		busCode = int8(4)
	}
	serviceMongo := models.ServiceMongo{}
	services, err := serviceMongo.QueryByBusiness(busCode)
	if err != nil{
		logs.Error("find service fail")
	}
	return services
}

func (c* CaseManageController) ShowAddCase(){
	business := c.GetString("business")
	services := GetServiceList(business)
	// 获取全部service
	c.Data["services"] = services
	c.TplName = "case_add.html"
}

func (c *CaseManageController) GetAllCases(){
	acm := models.TestCaseMongo{}
	business := c.GetString("business")
	fmt.Println(business)
	page, _ := strconv.Atoi(c.GetString("page"))
	limit, _ := strconv.Atoi(c.GetString("limit"))
	result, err := acm.GetAllCases(page, limit, business)
	if err != nil {
		logs.Error("获取全部用例失败")
		logs.Error(1024, err)
	}
	c.FormSuccessJson(result)
}

func (c *CaseManageController) ShowEditCase(){
	id := c.GetString("id")
	business := c.GetString("business")
	services := GetServiceList(business)
	idInt, err := strconv.ParseInt(id,10,64)
	if err != nil{
		logs.Error("转换类型错误")
	}
	acm := models.TestCaseMongo{}
	res:= acm.GetOneCase(idInt)
	c.Data["a"] = &res
	c.Data["services"] = services
	c.TplName = "case_edit.html"
}