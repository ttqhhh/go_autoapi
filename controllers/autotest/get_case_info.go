package controllers

import (
	"fmt"
	"go_autoapi/models"
)

//用来解析参数
type user struct {
	A int `form:"a"`
}

func (c *AutoTestController) get_case_info() {
	u := user{}
	if err := c.ParseForm(&u); err != nil {
		c.ErrorJson(-1, "请求错误", nil)
	}
	fmt.Println(u.A)
	uid := int64(u.A)
	acm := models.AdMockCaseMongo{}
	acm, err := acm.GetAdCase(uid)
	if err != nil {
		c.ErrorJson(-1, "请求错误", nil)
	}
	c.SuccessJson(acm)
}
