package libs

import (
	beego "github.com/beego/beego/v2/server/web"
	"github.com/go-redis/redis"
	"github.com/satori/go.uuid"
	_ "go_autoapi/constants"
	"go_autoapi/db_proxy"
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

//func (b *BaseController) Prepare() {
//	userId, err := b.GetSecureCookie(constant.CookieSecretKey, "user_id")
//	if err == false && b.GetMethodName() != "login" && b.GetMethodName() != "to_login" {
//		//logs.Error("not login")
//		//b.ErrorJson(-1, "not login", nil)
//		b.Redirect("/auto/to_login", 302)
//	}
//	fmt.Println(userId)
//}

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

func (b *BaseController) GetRedis() *redis.Client {
	_ = db_proxy.InitClient()
	return db_proxy.GetRedisObject()
}
