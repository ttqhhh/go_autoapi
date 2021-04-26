package controllers

import (
	"encoding/json"
	"github.com/beego/beego/v2/core/logs"
	constant "go_autoapi/constants"
	"go_autoapi/libs"
	"go_autoapi/models"
)

type ServiceController struct {
	libs.BaseController
}

func (c *ServiceController) Get() {
	do := c.GetMethodName()
	switch do {
	case "index":
		c.index()
	//case "addIndex":
	//	c.addIndex()
	case "page":
		c.page()
	case "list":
		c.list()
	case "getById":
		c.getById()
	case "business":
		c.business()
	default:
		logs.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持", nil)
	}
}

func (c *ServiceController) Post() {
	do := c.GetMethodName()
	switch do {
	case "save":
		c.save()
	case "remove":
		c.remove()
	//case "update":
	//	c.update()
	default:
		logs.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持", nil)
	}
}

// 页面跳转
func (c *ServiceController) index() {
	c.TplName = "service.html"
}

//func (c *ServiceController) addIndex() {
//	id, err := c.GetInt64("id", -1)
//	if err != nil {
//		logs.Warn("/service/getById接口 参数异常, err: %v", err)
//		c.ErrorJson(-1, "参数异常", nil)
//	}
//	logs.Info("请求参数: id=%v", id)
//
//	c.Data["id"] = id
//	c.TplName = "service_add.tpl"
//}

// 分页查询
func (c *ServiceController) page() {
	serviceName := c.GetString("service_name")
	business, err := c.GetInt8("business", -1)
	if err != nil {
		logs.Warn("/service/page接口 参数异常, err: %v", err)
		c.ErrorJson(-1, "参数异常", nil)
	}
	pageNo, err := c.GetInt("page", 1)
	if err != nil {
		logs.Warn("/service/page接口 参数异常, err: %v", err)
		c.ErrorJson(-1, "参数异常", nil)
	}
	pageSize, err := c.GetInt("limit", 10)
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
	//result := make(map[string]interface{})
	//result["total"] = total
	//result["data"] = services

	//c.SuccessJson(result)
	res := make(map[string]interface{})
	res["code"] = 0
	res["msg"] = "成功"
	res["count"] = total
	res["data"] = services

	c.Data["json"] = res
	c.ServeJSON() //对json进行序列化输出
	c.StopRun()
}

// 获取服务列表（可根据业务线）
func (c *ServiceController) list() {
	business, err := c.GetInt8("business", -1)
	if err != nil {
		logs.Warn("/service/list接口 参数异常, err: %v", err)
		c.ErrorJson(-1, "参数异常", nil)
	}
	if business == -1 {
		logs.Warn("/service/list接口 未指定业务线", err)
		c.ErrorJson(-1, "参数异常，请指定业务线", nil)

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
func (c *ServiceController) save() {
	userId, _ := c.GetSecureCookie(constant.CookieSecretKey, "user_id")

	service := &models.ServiceMongo{}
	err := c.ParseForm(service)
	if err != nil {
		logs.Warn("/service/save接口 参数异常, err: %v", err)
		c.ErrorJson(-1, "参数异常", nil)
	}
	logs.Info("请求参数：%v", service)

	if string(service.Id) == "" || service.Id == -1 || service.Id == 0 {
		service.CreateBy = userId
		err = service.Insert(*service)
		if err != nil {
			c.ErrorJson(-1, "服务添加数据异常", nil)
		}
	} else {
		service.UpdateBy = userId
		err = service.Update(*service)
		if err != nil {
			c.ErrorJson(-1, "服务更新数据异常", nil)
		}
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
//func (c *ServiceController) update() {
//	service := &models.ServiceMongo{}
//	err := c.ParseForm(service)
//	if err != nil {
//		logs.Warn("/service/update接口 参数异常, err: %v", err)
//		c.ErrorJson(-1, "参数异常", nil)
//	}
//	logs.Info("请求参数: %v", service)
//
//	// todo 更新人字段待处理
//	err = service.Update(*service)
//	if err != nil {
//		c.ErrorJson(-1, "服务更新数据异常", nil)
//	}
//	c.SuccessJson(nil)
//}
// 用来mock业务线数据
func (c *ServiceController) business() {
	var data [6]map[string]interface{}

	//data[0] = "最右"
	//data[1] = "皮皮"
	//data[2] = "中东"
	//data[3] = "海外"
	//data[4] = "商业化"
	//data[5] = "妈妈社区"
	one := make(map[string]interface{})
	one["code"] = 0
	one["name"] = "最右"
	data[0] = one

	one = make(map[string]interface{})
	one["code"] = 1
	one["name"] = "皮皮"
	data[1] = one

	one = make(map[string]interface{})
	one["code"] = 2
	one["name"] = "中东"
	data[2] = one

	one = make(map[string]interface{})
	one["code"] = 3
	one["name"] = "海外"
	data[3] = one

	one = make(map[string]interface{})
	one["code"] = 4
	one["name"] = "商业化"
	data[4] = one

	one = make(map[string]interface{})
	one["code"] = 5
	one["name"] = "妈妈社区"
	data[5] = one

	c.SuccessJson(data)
}
