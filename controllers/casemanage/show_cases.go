package controllers

import (
	"github.com/beego/beego/v2/core/logs"
	"go_autoapi/models"
	"strconv"
)

//string转成int：
//int, err := strconv.Atoi(string)
//string转成int64：
//int64, err := strconv.ParseInt(string, 10, 64)
//int转成string：
//string := strconv.Itoa(int)
//int64转成string：
//string := strconv.FormatInt(int64,10)

func (c *CaseManageController) ShowCases(){
	//c.Data["Website"] = "beego.me"
	//c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "case_manager.html"
}

func (c* CaseManageController) ShowAddCase(){
	c.TplName = "case_add.html"
}


func (c *CaseManageController) GetAllCases(){
	acm := models.TestCaseMongo{}
	result, err := acm.GetAllCases()
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
	res, err := acm.GetOneCase(idInt)
	c.Data["a"] = &res
	c.TplName = "case_edit.html"
}