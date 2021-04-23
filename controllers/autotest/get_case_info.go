package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"go_autoapi/models"
)

//用来解析参数
type autoCase struct {
	Id int `form:"id"`
}

func (c *AutoTestController) getCaseInfo() {
	ac := autoCase{}
	fmt.Println(ac)
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &ac); err != nil {
		logs.Error(1024, err)
		c.ErrorJson(-1, "请求错误", nil)
	}
	caseId := int64(ac.Id)
	fmt.Println(caseId)
	acm := models.AutoCaseMongo{}
	acm, err := acm.GetCaseById(caseId)
	if err != nil {
		fmt.Println(err)
		c.ErrorJson(-1, "请求错误", nil)
	}
	c.SuccessJsonWithMsg(acm, "OK")
}
