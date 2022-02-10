package h5listen

import (
	"github.com/siddontang/go/log"
	"go_autoapi/libs"
)

type H5ListenController struct {
	libs.BaseController
}

func (c *H5ListenController) Get() {
	do := c.GetMethodName()
	switch do {
	case "show_h5":
		business := c.GetString("business")
		c.Data["business"] = business
		c.TplName = "h5.html"

	case "show_h5_data":
		c.ShowH5()

	default:
		log.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持", nil)
	}
}

func (c *H5ListenController) Post() {
	do := c.GetMethodName()
	switch do {
	case "update_one_h5_data":
		c.updateH5DateByID()
	case "add_one_h5_data":
		c.AddOneH5Date()
	case "del_one_h5_data":
		c.DelH5DateByID()
	case "show_h5_data":
		c.ShowH5()
	//case "StrategyH5":
	//	c.StrategyH5()

	default:
		log.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持", nil)
	}
}
