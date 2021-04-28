package controllers

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"go_autoapi/models"
	"strconv"
)

func (c *CaseManageController) updateCaseByID() {
	acm := models.TestCaseMongo{}
	if err := c.ParseForm(&acm); err != nil { //传入user指针
		c.Ctx.WriteString("出错了！")
	}
	caseId := acm.Id
	business := acm.BusinessCode
	if business == "0" {
		acm.BusinessName = "最右"
	} else if business == "1" {
		acm.BusinessName = "皮皮"
	} else if business == "2" {
		acm.BusinessName = "海外"
	} else if business == "3" {
		acm.BusinessName = "中东"
	} else if business == "4" {
		acm.BusinessName = "妈妈社区"
	} else if business == "5" {
		acm.BusinessName = "商业化"
	}
	acm, err := acm.UpdateCase(caseId, acm)
	if err != nil {
		fmt.Println(err)
		c.ErrorJson(-1, "请求错误", nil)
	}
	//c.SuccessJson("更新成功")
	c.Ctx.Redirect(302, "/case/show_cases?business="+business)
}

func (c *CaseManageController) DelCaseByID() {
	caseID := c.GetString("id")
	ac := models.TestCaseMongo{}
	caseIDInt, err := strconv.ParseInt(caseID, 10, 64)
	if err != nil {
		logs.Error("在删除用例的时候类型转换失败")
	}
	ac.DelCase(caseIDInt)
	logs.Info("删除成功")
}
