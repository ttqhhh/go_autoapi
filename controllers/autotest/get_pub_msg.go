package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"go_autoapi/models"
)

//用来解析参数
type formPubMsg struct {
	Id int `json:"id"`
}

func (c *AutoTestController) showPubInfo(){
	c.TplName = "show_pub_info.html"
}

func (c *AutoTestController) getPubInfo() {
	ac := formPubMsg{}
	//fmt.Println(ac)
	//fmt.Println(c.Ctx.Input.RequestBody)
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &ac); err != nil {
		logs.Error(1024, err)
		c.ErrorJson(-1, "请求错误", nil)
	}
	businessId := int8(ac.Id)
	fmt.Println(businessId)
	pm := models.PublishMsg{}
	pmb, err := pm.GetPubMsgByBusiness(businessId)
	if err != nil {
		fmt.Println(err)
		c.ErrorJson(-1, "请求错误", nil)
	}
	c.ElemJson(pmb)
}
