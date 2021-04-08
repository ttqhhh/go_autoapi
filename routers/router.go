package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"go_autoapi/controllers"
	auto "go_autoapi/controllers/autotest"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/user/*", &controllers.UserController{})
	beego.Router("/auto/*", &auto.AutoTestController{})
}
