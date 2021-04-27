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
	//if err := json.Unmarshal(c.Ctx.Input.RequestBody, &acm); err != nil {
	//	c.ErrorJson(-1, "请求错误(json解析错误)", err)
	//}
	acm.Id = models.GetId("case")
	acm.CreatedAt = now
	acm.UpdatedAt = now
	acm.Status = 0
	business := acm.AppName
	if err := acm.AddCase(acm); err != nil {
		c.ErrorJson(-1, "请求错误", nil)
	}
	//c.SuccessJson("添加成功")
	c.Ctx.Redirect(302, "/case/show_cases?business=" + business)
}
