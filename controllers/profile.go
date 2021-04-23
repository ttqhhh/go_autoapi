package controllers

import (
	"fmt"
	"go_autoapi/models"
)

//用来解析参数
type user struct {
	A int `form:"a"`
}

func (c *UserController) profile() {
	fmt.Println("xxxxxx", "profile")
	amc := models.AdMockCase{}
	data, _ := amc.QueryByUUid("5eb62275-9818-4101-a477-6fef0bb9c7bd")
	fmt.Println(&data)
	c.SuccessJsonWithMsg(&data, "OK")
}

func (c *UserController) add_mgo() {
	acm := models.AdMockCaseMongo{}
	err := acm.InsertAdCase(models.AdMockCaseMongo{Id: 1, CaseName: "aaaaaa", CaseDesc: "bbbbbbb"})
	if err != nil {
		c.ErrorJson(-1, "请求错误", nil)
	}
	c.SuccessJsonWithMsg(nil, "添加成功")
}

func (c *UserController) get_mgo() {
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
	c.SuccessJsonWithMsg(acm, "OK")
}
