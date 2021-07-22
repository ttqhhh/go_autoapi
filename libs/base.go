package libs

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/satori/go.uuid"
	_ "go_autoapi/constants"
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

type ReturnMsgPage struct {
	Code  int         `json:"code"`
	Count int64       `json:"count"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

// 不需要登录的接口
var withoutLoginApi = []string{
	"/case/pp_auto_add_one_case",
	"/monitor/excute_at_first_time",
	"/monitor/excute_one_time",
	"/health/check",
	"/flowreplay/collect_flow_file",
}

// todo 该方法不要注释掉，不然登录功能就没意义了
func (b *BaseController) Prepare() {
	_, err := b.GetSecureCookie(constant.CookieSecretKey, "user_id")
	methodName := b.GetMethodName()
	if err == false && methodName != "login" && methodName != "to_login" {
		//logs.Error("not login")
		//b.ErrorJson(-1, "not login", nil)
		isNeedCheck := true
		uri := b.Ctx.Request.RequestURI
		for _, api := range withoutLoginApi {
			if api == uri {
				isNeedCheck = false
				break
			}
		}
		if isNeedCheck {
			b.Redirect("/auto/to_login", 302)
		}
	}
	//fmt.Println(userId)
}

func (b *BaseController) SuccessJsonWithMsg(data interface{}, msg string) {

	res := ReturnMsg{
		200, msg, data,
	}
	b.Data["json"] = res
	b.ServeJSON() //对json进行序列化输出
	b.StopRun()
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
	return strings.Split(do, "/")[2]
}

func (b *BaseController) FormSuccessJson(count int64, data interface{}) {

	res := ReturnMsgPage{
		0, count, "success", data,
	}
	b.Data["json"] = res
	b.ServeJSON() //对json进行序列化输出
	b.StopRun()
}

func (b *BaseController) FormErrorJson(code int, msg string) {

	res := ReturnMsgPage{
		code, 0, msg, nil,
	}
	b.Data["json"] = res
	b.ServeJSON() //对json进行序列化输出
	b.StopRun()
}

func (b *BaseController) GenUUid() (string, error) {
	u2 := uuid.NewV4()
	return u2.String(), nil
}
