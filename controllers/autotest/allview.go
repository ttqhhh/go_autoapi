package controllers

import (
	"github.com/beego/beego/v2/core/logs"
	"go_autoapi/libs"
	"go_autoapi/models"
	_ "image/color"
	_ "net/http"
)

type AllviewController struct {
	libs.BaseController
}

func (c *AllviewController) Get() {
	do := c.GetMethodName()
	switch do {
	case "jump":
		c.jump()
	default:
		logs.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持", nil)
	}
}

// 页面跳转
func (c *AllviewController) jump() {
	c.TplName = "allview.html"
}

func (c *AllviewController) Post() { // 获取到用户传入的数据，解析到结构体中
	do := c.GetMethodName()
	switch do {
	case "save":
		c.save()
	case "check":
		c.check()
	case "remove":
		c.remove()
	default:
		logs.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持", nil)
	}
}

// 添加
func (c *AllviewController) save() {
	addlink := &models.AddLinkMongo{}
	err := c.ParseForm(addlink)
	logs.Info("请求参数：%v", addlink)
	if err != nil {
		logs.Warn("/allview/save接口 参数异常，err: %v", err)
		c.ErrorJson(-1, "参数异常", nil)
		logs.Info("请求参数：%v", addlink)
	}

	err = addlink.Insert(*addlink)
	if err != nil {
		c.ErrorJson(-1, "服务添加数据异常", nil)
	}
	c.Redirect("/allview/jump", 302)
}

// 查询
func (c *AllviewController) check() {
	addlinkMongo := models.AddLinkMongo{}
	addlinks, err := addlinkMongo.Lists()

	if err != nil {
		c.ErrorJson(-1, "服务查询数据异常", nil)
	}
	c.SuccessJson(addlinks)
}

type _RemoveParam struct {
	Name string `form:"name" json:"name"`
}

// 删除
func (c *AllviewController) remove() {
	name := c.GetString("name")
	param := &_RemoveParam{}
	param.Name = name
	mongo := &models.AddLinkMongo{}
	err := mongo.Delete(param.Name)

	if err != nil {
		logs.Warn("/allview/remove接口 参数异常, err: %v", err)
		c.ErrorJson(-1, "参数异常", nil)
	}
	c.Redirect("/allview/jump", 302)
}
