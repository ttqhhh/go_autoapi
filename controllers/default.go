package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
)

type MainController struct {
	beego.Controller
}

func (a *MainController) Get() {
	a.Data["Website"] = "beego.me"
	a.Data["Email"] = "astaxie@gmail.com"
	a.TplName = "index.tpl"
}
