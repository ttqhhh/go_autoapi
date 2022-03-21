package h5listen

import (
	constant "go_autoapi/constants"
	"go_autoapi/models"
	"time"
)

func (c *H5ListenController) AddOneH5Date() {
	acm := models.AdH5DataMongo{}
	DataName := c.GetString("data_name")
	DataUrl := c.GetString("data_url")
	Business := c.GetString("business")
	BusinessName := ""
	switch Business {
	case "0":
		BusinessName = "zuiyou"
	case "1":
		BusinessName = "pipi"
	case "2":
		BusinessName = "haiwai"
	case "3":
		BusinessName = "zhongdong"
	case "4":
		BusinessName = "matuan"
	case "5":
		BusinessName = "business"
	case "6":
		BusinessName = "haiwai-US"
	}
	acm.Id = time.Now().Unix()
	acm.DataName = DataName
	acm.DataUrl = DataUrl
	acm.Business = Business
	acm.BusinessName = BusinessName
	nowtime := time.Now()
	acm.CreatedAt = nowtime.Format(constant.TimeFormat)

	acm.InsertAdCase(acm)
	c.Ctx.Redirect(302, "/h5listen/show_h5?business="+Business)
}
