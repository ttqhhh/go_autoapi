package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
)

type MainController struct {
	beego.Controller
}

// todo 勿动~~
// todo 勿动~~
// todo 勿动~~
// todo 勿动~~
func (a *MainController) Get() {
	a.TplName = "login.html"
}
