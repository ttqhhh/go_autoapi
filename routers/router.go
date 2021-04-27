package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"go_autoapi/controllers"
	auto "go_autoapi/controllers/autotest"
	casemanage "go_autoapi/controllers/casemanage"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/user/*", &controllers.UserController{})
	beego.Router("/auto/*", &auto.AutoTestController{})
	beego.Router("/case/*", &casemanage.CaseManageController{})
	beego.Router("/service/*", &auto.ServiceController{})
	beego.Router("/business/*", &auto.BusinessController{})
}
