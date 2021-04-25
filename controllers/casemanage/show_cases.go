package controllers

import (
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"go_autoapi/models"
	"strconv"
)

func (c *CaseManageController) ShowCases(){
	//c.Data["Website"] = "beego.me"
	//c.Data["Email"] = "astaxie@gmail.com"
	c.Data["business"] = c.GetString("business")
	c.TplName = "case_manager.html"
}

func (c* CaseManageController) ShowAddCase(){
	business,err := c.GetInt8("business", -1)
	if err!=nil{logs.Error("parseint fail")}
	serviceMongo := models.ServiceMongo{}
	services, err := serviceMongo.QueryByBusiness(business)
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
	idInt, err := strconv.ParseInt(id,10,64)
	if err != nil{
		logs.Error("转换类型错误")
	}
	acm := models.TestCaseMongo{}
	res:= acm.GetOneCase(idInt)
	c.Data["a"] = &res
	c.TplName = "case_edit.html"
}