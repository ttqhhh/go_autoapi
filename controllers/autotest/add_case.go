package controllers

import (
	"encoding/json"
	"fmt"
	"go_autoapi/constants"
	"go_autoapi/models"
	"time"
)

func (c *AutoTestController) addCase() {
	now := time.Now().Format(constants.TimeFormat)
	id := models.GetId("case")
	acm := models.AutoCaseMongo{Id: id, CreatedAt: now, UpdatedAt: now}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &acm); err != nil {
		c.ErrorJson(-1, "请求错误", nil)
	}
	err := acm.InsertCase(acm)
	if err != nil {
		fmt.Println(err)
		c.ErrorJson(-1, "请求错误", nil)
	}
	c.SuccessJsonWithMsg(nil, "添加成功")
}
