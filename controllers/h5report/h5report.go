package h5report

import (
	"github.com/siddontang/go/log"
	"go_autoapi/libs"
	"go_autoapi/models"
)

type H5ReportController struct {
	libs.BaseController
}

func (c *H5ReportController) Get() {
	do := c.GetMethodName()
	switch do {
	case "show_h5_report":
		c.show_h5_report()
	case "show_report":
		c.show_report()

	default:
		log.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持", nil)
	}
}

func (c *H5ReportController) Post() {
	do := c.GetMethodName()
	switch do {
	default:
		log.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持", nil)
	}
}

func (c *H5ReportController) show_report() {
	c.TplName = "h5_report.html"

}

func (c *H5ReportController) show_h5_report() {
	acm := models.H5RunReportMongo{}
	//business := c.GetString("business", "1")
	result, err := acm.Query()
	if err != nil {
		c.FormErrorJson(-1, "获取h5报告数据失败")
	}
	c.FormSuccessJson(0, result)
}
