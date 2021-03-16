package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
)

type BaseController struct {
	beego.Controller
}

type ReturnMsg struct {
	Code int `json:"code"`
	Msg  string `json:"msg"`
	Data interface{} `json:"data"`
}

func (b *BaseController) SuccessJson(data interface{}) {

	res := ReturnMsg{
		200, "success", data,
	}
	b.Data["json"] = res
	b.ServeJSON() //对json进行序列化输出
	b.StopRun()
}

func (b *BaseController) ErrorJson(code int, msg string, data interface{}) {

	res := ReturnMsg{
		code, msg, data,
	}

	b.Data["json"] = res
	b.ServeJSON() //对json进行序列化输出
	b.StopRun()
}
