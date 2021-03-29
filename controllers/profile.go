package controllers

import (
	"fmt"
	"go_autoapi/models"
)

func (c *UserController) profile() {
	fmt.Println("xxxxxx", "profile")
	amc := models.AdMockCase{}
	data, _ := amc.QueryByUUid("5eb62275-9818-4101-a477-6fef0bb9c7bd")
	fmt.Println(&data)
	c.SuccessJson(&data)
}

func (c *UserController) add_mgo() {
	acm := models.AdMockCaseMongo{}
	err := acm.InsertAdCase(models.AdMockCaseMongo{Id: 1, CaseName: "aaaaaa", CaseDesc: "bbbbbbb"})
	if err != nil {
		c.ErrorJson(-1, "请求错误", nil)
	}
	c.SuccessJson("添加成功")
}

func (c *UserController) get_mgo() {
	acm := models.AdMockCaseMongo{}
	acm, err := acm.GetAdCase()
	if err != nil {
		c.ErrorJson(-1, "请求错误", nil)
	}
	c.SuccessJson(acm)
}
