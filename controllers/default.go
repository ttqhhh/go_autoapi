package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
	constant "go_autoapi/constants"
	controllers "go_autoapi/controllers/autotest"
	"strconv"
)

type MainController struct {
	beego.Controller
}

// todo 勿动~~
// todo 勿动~~
// todo 勿动~~
// todo 勿动~~
func (a *MainController) Get() {
	userId, isLogined := a.GetSecureCookie(constant.CookieSecretKey, "user_id")
	if isLogined != false {
		// 当前用户为登录状态时，直接跳转Case列表页
		// 默认跳转到第一个有权限的业务线case页面
		businesses := controllers.GetBusinesses(userId)
		business := businesses[0]
		code := business["code"]
		codeInt := code.(int)
		redirectUrl := "/case/show_cases?business=" + strconv.Itoa(codeInt)
		a.Redirect(redirectUrl, 302)
	}
	a.TplName = "login.html"
}
