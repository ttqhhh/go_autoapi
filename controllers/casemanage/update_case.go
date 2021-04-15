package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"go_autoapi/models"
)

func (c *CaseManageController) updateCaseByID() {
	ac := models.TestCaseMongo{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &ac); err != nil {
		logs.Error(1024, err)
		c.ErrorJson(-1, "请求错误", nil)
	}
	caseId := ac.Id
	fmt.Println(caseId)
	acm := models.TestCaseMongo{}
	acm, err := acm.UpdateCase(caseId, ac)
	if err != nil {
		fmt.Println(err)
		c.ErrorJson(-1, "请求错误", nil)
	}
	c.SuccessJson("更新成功")
}

