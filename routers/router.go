package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"go_autoapi/api"
	"go_autoapi/controllers"
	auto "go_autoapi/controllers/autotest"
	casemanage "go_autoapi/controllers/casemanage"
	"go_autoapi/controllers/inspection"
	"go_autoapi/controllers/report"
	"go_autoapi/controllers/web_report"
	_ "go_autoapi/controllers/web_report"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/api/*", &api.ApiController{})
	beego.Router("/user/*", &controllers.UserController{})
	beego.Router("/auto/*", &auto.AutoTestController{})
	beego.Router("/case/*", &casemanage.CaseManageController{})
	beego.Router("/service/*", &auto.ServiceController{})
	beego.Router("/business/*", &auto.BusinessController{})
	beego.Router("/report/*", &report.ReportController{})
	beego.Router("/inspection/*", &inspection.CaseController{})
	beego.Router("/web_report/*", &web_report.WebreportController{})
}
