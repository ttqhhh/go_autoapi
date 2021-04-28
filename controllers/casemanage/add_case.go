package controllers

import (
	"go_autoapi/constants"
	"go_autoapi/models"
	"time"
)

func (c *CaseManageController) AddOneCase() {
	now := time.Now().Format(constants.TimeFormat)
	acm := models.TestCaseMongo{}
	if err := c.ParseForm(&acm); err != nil { //传入user指针
		c.Ctx.WriteString("出错了！")
	}
	acm.Id = models.GetId("case")
	acm.CreatedAt = now
	acm.UpdatedAt = now
	acm.Status = 0
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
	if err := acm.AddCase(acm); err != nil {
		c.ErrorJson(-1, "请求错误", nil)
	}
	//c.SuccessJson("添加成功")
	c.Ctx.Redirect(302, "/case/show_cases?business="+business)
}
