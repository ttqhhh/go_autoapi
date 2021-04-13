package libs

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	beego "github.com/beego/beego/v2/server/web"
	constant "go_autoapi/constants"
	"strings"
)

type BaseController struct {
	beego.Controller
}

type ReturnMsg struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (b *BaseController) Prepare() {
	userId, err := b.GetSecureCookie(constant.CookieSecretKey, "userid")
	if err == false && b.GetMethodName() != "login" {
		logs.Error("not login")
		b.ErrorJson(-1, "not login", nil)
	}
	fmt.Println(userId)
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

func (b *BaseController) GetMethodName() (do string) {
	do = b.Ctx.Request.URL.Path
	fmt.Println("url is ", do)
	return strings.Split(do, "/")[2]
}
