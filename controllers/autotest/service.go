package controllers

import (
	"encoding/json"
	"github.com/beego/beego/v2/core/logs"
	"go_autoapi/libs"
	"go_autoapi/models"
)

type ServiceController struct {
	libs.BaseController
}

func (c *ServiceController) Get() {
	do := c.GetMethodName()
	switch do {
	case "jump":
		c.jump()
	case "page":
		c.page()
	case "list":
		c.list()
	case "getById":
		c.getById()
	default:
		logs.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持", nil)
	}
}

func (c *ServiceController) Post() {
	do := c.GetMethodName()
	switch do {
	case "add":
		c.add()
	case "remove":
		c.remove()
	case "update":
		c.update()
	default:
		logs.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持", nil)
	}
}

// 页面跳转
func (c *ServiceController) jump() {

}

// 分页查询
func (c *ServiceController) page() {
	serviceName := c.GetString("service_name")
	business, err := c.GetInt8("business")
	if err != nil {
		logs.Warn("/service/page接口 参数异常, err: %v", err)
		c.ErrorJson(-1, "参数异常", nil)
	}
	pageNo, err := c.GetInt("page_no")
	if err != nil {
		logs.Warn("/service/page接口 参数异常, err: %v", err)
		c.ErrorJson(-1, "参数异常", nil)
	}
	pageSize, err := c.GetInt("page_size", 10)
	if err != nil {
		logs.Warn("/service/page接口 参数异常, err: %v", err)
		c.ErrorJson(-1, "参数异常", nil)
	}
	logs.Info("请求参数: service_name=%v, business=%v, page_no=%v, page_size=%v", serviceName, business, pageNo, pageSize)
	serviceMongo := models.ServiceMongo{}
	services, total, err := serviceMongo.QueryByPage(business, serviceName, pageNo, pageSize)
	if err != nil {
		c.ErrorJson(-1, "服务查询数据异常", nil)
	}
	result := make(map[string]interface{})
	result["total"] = total
	result["data"] = services

	c.SuccessJson(result)
}

// 获取服务列表（可根据业务线）
func (c *ServiceController) list() {
	business, err := c.GetInt8("business")
	if err != nil {
		logs.Warn("/service/list接口 参数异常, err: %v", err)
		c.ErrorJson(-1, "参数异常", nil)
	}
	logs.Info("请求参数: business=%v", business)
	serviceMongo := models.ServiceMongo{}
	services, err := serviceMongo.QueryByBusiness(business)
	if err != nil {
		c.ErrorJson(-1, "服务查询数据异常", nil)
	}
	c.SuccessJson(services)
}

// 获取指定服务（根据id）
func (c *ServiceController) getById() {
	id, err := c.GetInt64("id")
	if err != nil {
		logs.Warn("/service/getById接口 参数异常, err: %v", err)
		c.ErrorJson(-1, "参数异常", nil)
	}
	logs.Info("请求参数: id=%v", id)
	serviceMongo := models.ServiceMongo{}
	service, err := serviceMongo.QueryById(id)
	if err != nil {
		c.ErrorJson(-1, "服务查询数据异常", nil)
	}
	c.SuccessJson(service)
}

// 添加
func (c *ServiceController) add() {
	service := &models.ServiceMongo{}
	err := c.ParseForm(service)
	if err != nil {
		logs.Warn("/service/getById接口 参数异常, err: %v", err)
		c.ErrorJson(-1, "参数异常", nil)
	}
	logs.Info("请求参数：%v", service)

	//todo 添加人字段待处理
	err = service.Insert(*service)
	if err != nil {
		c.ErrorJson(-1, "服务添加数据异常", nil)
	}
	c.SuccessJson(nil)
}

type RemoveParam struct {
	Id int64 `form:"id" json:"id"`
}

// 删除
func (c *ServiceController) remove() {
	param := &RemoveParam{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, param)
	if err != nil {
		logs.Warn("/service/remove接口 参数异常, err: %v", err)
		c.ErrorJson(-1, "参数异常", nil)
	}
	logs.Info("请求参数: id=%v", param)
	serviceMongo := models.ServiceMongo{}
	err = serviceMongo.Delete(param.Id)
	if err != nil {
		c.ErrorJson(-1, "服务删除数据异常", nil)
	}
	c.SuccessJson(nil)
}

// 更新
func (c *ServiceController) update() {
	service := &models.ServiceMongo{}
	err := c.ParseForm(service)
	if err != nil {
		logs.Warn("/service/update接口 参数异常, err: %v", err)
		c.ErrorJson(-1, "参数异常", nil)
	}
	logs.Info("请求参数: %v", service)

	// todo 更新人字段待处理
	err = service.Update(*service)
	if err != nil {
		c.ErrorJson(-1, "服务更新数据异常", nil)
	}
	c.SuccessJson(nil)
}
