package api

import (
	"github.com/astaxie/beego/logs"
	"go_autoapi/constants"
	"go_autoapi/libs"
	"go_autoapi/models"
	"strconv"
)

type ToolsController struct {
	libs.BaseController
}

func (c *ToolsController) Get() {
	do := c.GetMethodName()
	switch do {
	case "flush_token":
		c.flush_token()
	default:
		logs.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持", nil)
	}
}

func (c *ToolsController) Post() {
	do := c.GetMethodName()
	switch do {
	default:
		logs.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持", nil)
	}
}

func (c *ToolsController) flush_token() {
	business := c.GetString("business", "") // 枚举： 0：最右，1：皮皮，2：海外，3：中东，4：妈妈，5：商业化，6：海外-US
	token := c.GetString("token", "")
	env := c.GetString("env", "") // 枚举：prod、test
	mid := c.GetString("mid", "") // 非必填

	if business == "" {
		c.ErrorJson(-1, "business字段不可为空", nil)
	}
	if token == "" {
		c.ErrorJson(-1, "token字段不可为空", nil)
	}
	if env == "" {
		c.ErrorJson(-1, "env字段不可为空", nil)
	}
	businessInt, err := strconv.Atoi(business)
	if err != nil {
		c.ErrorJson(-1, "您输入的business参数有误", nil)
	}
	if businessInt != constants.ZuiyYou && businessInt != constants.PiPi && businessInt != constants.HaiWai && businessInt != constants.ZhongDong && businessInt != constants.Matuan && businessInt != constants.ShangYeHua && businessInt != constants.HaiWaiUS {
		c.ErrorJson(-1, "您输入的business参数必须满足0-6之间的数", nil)
	}
	if env != "prod" && env != "test" {
		c.ErrorJson(-1, "您输入的env参数只能是prod或者test", nil)
	}

	if env == "test" {
		mongo := models.TestCaseMongo{}
		err = mongo.FlushAllTokenByBusiness(business, token, mid)
		if err != nil {
			c.ErrorJson(-1, err.Error(), nil)
		}
	} else {
		mongo := models.InspectionCaseMongo{}
		err = mongo.FlushAllTokenByBusiness(business, token, mid)
		if err != nil {
			c.ErrorJson(-1, err.Error(), nil)
		}
	}

	c.SuccessJson(nil)
}
