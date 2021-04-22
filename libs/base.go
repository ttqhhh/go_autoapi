package libs

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/satori/go.uuid"
	_ "go_autoapi/constants"
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

//func (b *BaseController) Prepare() {
//	userId, err := b.GetSecureCookie(constant.CookieSecretKey, "userid")
//	if err == false && b.GetMethodName() != "login" {
//		logs.Error("not login")
//		b.ErrorJson(-1, "not login", nil)
//	}
//	fmt.Println(userId)
//}

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
	return strings.Split(do, "/")[2]
}

func (b *BaseController) FormSuccessJson(data interface{}) {

	res := ReturnMsg{
		0, "success", data,
	}
	b.Data["json"] = res
	b.ServeJSON() //对json进行序列化输出
	b.StopRun()
}

func (b *BaseController) GenUUid() (string, error) {
	u2 := uuid.NewV4()
	return u2.String(), nil
}

