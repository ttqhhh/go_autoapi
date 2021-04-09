package controllers

import (
	"encoding/json"
	"fmt"
	"go_autoapi/models"
)

func (c *AutoTestController) addCase() {
	id := models.GetId("case")
	acm := models.AutoCaseMongo{Id: id}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &acm); err != nil {
		c.ErrorJson(-1, "请求错误", nil)
	}
	err := acm.InsertCase(acm)
	if err != nil {
		fmt.Println(err)
		c.ErrorJson(-1, "请求错误2", nil)
	}
	c.SuccessJson("添加成功")
}
