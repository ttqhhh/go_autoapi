package controllers

import (
	"encoding/json"
	"fmt"
	"go_autoapi/models"
	"time"
)

func (c *CaseManageController) AddOneCase() {
	now := time.Now()
	acm := models.TestCaseMongo{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &acm); err != nil {
		c.ErrorJson(-1, "请求错误", nil)
	}
	caseName := acm.CaseName
	id := models.GetId(caseName)
	acm.Id = id
	acm.CreatedAt = now
	acm.UpdatedAt = now
	if err := acm.AddCase(acm); err != nil {
		fmt.Println(err)
		c.ErrorJson(-1, "请求错误", nil)
	}
	c.SuccessJson("添加成功")
}
