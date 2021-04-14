package controllers

import (
	"encoding/json"
	"fmt"
	"go_autoapi/models"
	"time"
)

func (c *CaseManageController) addCase() {
	now := time.Now()
	id := models.GetId("case")
	acm := models.TestCaseMongo{Id: id, CreatedAt: now, UpdatedAt: now}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &acm); err != nil {
		c.ErrorJson(-1, "请求错误", nil)
	}
	err := acm.InsertCase(acm)
	if err != nil {
		fmt.Println(err)
		c.ErrorJson(-1, "请求错误", nil)
	}
	c.SuccessJson("添加成功")
}
