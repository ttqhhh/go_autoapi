package routers

import (
	beego "github.com/beego/beego/v2/server/web"
	"go_autoapi/api"
	"go_autoapi/controllers"
	auto "go_autoapi/controllers/autotest"
	"go_autoapi/controllers/case_set"
	casemanage "go_autoapi/controllers/casemanage"
	"go_autoapi/controllers/h5listen"
	"go_autoapi/controllers/h5report"
	"go_autoapi/controllers/health"
	"go_autoapi/controllers/inspection"
	"go_autoapi/controllers/monitor"
	"go_autoapi/controllers/presstest"
	"go_autoapi/controllers/report"
	"go_autoapi/controllers/tuijian"
	"go_autoapi/controllers/web_report"
	_ "go_autoapi/controllers/web_report"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/api/*", &api.ApiController{}) // 对外暴露Api接口
	beego.Router("/user/*", &controllers.UserController{})
	beego.Router("/auto/*", &auto.AutoTestController{})
	beego.Router("/case/*", &casemanage.CaseManageController{})
	beego.Router("/service/*", &auto.ServiceController{})
	beego.Router("/business/*", &auto.BusinessController{})
	beego.Router("/report/*", &report.ReportController{})
	beego.Router("/inspection/*", &inspection.CaseController{})
	beego.Router("/web_report/*", &web_report.WebreportController{})
	beego.Router("/flowreplay/*", &tuijian.FlowReplayController{})
	beego.Router("/monitor/*", &monitor.ZYMonitorController{})
	beego.Router("/health/*", &health.HealthController{})
	beego.Router("/presstest/*", &presstest.PressTestController{})
	beego.Router("/h5listen/*", &h5listen.H5ListenController{})
	beego.Router("/h5report/*", &h5report.H5ReportController{})
	beego.Router("/allview/*", &auto.AllviewController{})
	beego.Router("/case_set/*", &case_set.CaseSetController{})
}
