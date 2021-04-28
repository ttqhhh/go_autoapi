package controllers

import (
	"github.com/astaxie/beego/logs"
	constant "go_autoapi/constants"
	"go_autoapi/libs"
)

type BusinessController struct {
	libs.BaseController
}

func (c *BusinessController) Get() {
	do := c.GetMethodName()
	switch do {
	case "get_user_businesses":
		c.getUserBusinesses()
	default:
		logs.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持", nil)
	}
}

var userBusinessMap = map[string][]int{
	// 4个超级管理员
	"fengmanlong": {zuiyou, pipi, haiwai, zhongdong, mama, shangyehua},
	"xueyibing":   {zuiyou, pipi, haiwai, zhongdong, mama, shangyehua},
	"liuweiqiang": {zuiyou, pipi, haiwai, zhongdong, mama, shangyehua},
	"wangzhen":    {zuiyou, pipi, haiwai, zhongdong, mama, shangyehua},
	// 普通用户
	"wangjun":       {haiwai, zhongdong, mama},
	"sunzhiying":    {haiwai},
	"zhangjuan":     {shangyehua, zuiyou},
	"liuyinan":      {pipi},
	"sunmingyao":    {pipi},
	"zhangdanbing":  {pipi},
	"xufei":         {haiwai, zhongdong},
	"suhuimin":      {zhongdong},
	"wanglanjin":    {zhongdong, mama},
	"yangjingfang":  {shangyehua, zuiyou},
	"chengxiaojing": {zuiyou},
	"zhaoxiaodong":  {shangyehua},
}

var businessCodeNameMap = map[int]string{
	zuiyou:     zuiyou_name,
	pipi:       pipi_name,
	haiwai:     haiwai_name,
	zhongdong:  zhongdong_name,
	mama:       mama_name,
	shangyehua: shangyehua_name,
}

func (c *BusinessController) getUserBusinesses() {
	userId, _ := c.GetSecureCookie(constant.CookieSecretKey, "user_id")

	result := make(map[string]interface{})
	result["username"] = userId
	result["businesses"] = getBusinesses(userId)
	c.SuccessJson(result)
}

func getBusinesses(username string) []map[string]interface{} {
	businessResp := [](map[string]interface{}){}
	businesses, ok := userBusinessMap[username]
	if !ok {
		// todo 目前只对测试同学进行了限制，其他角色同学暂未进行处理
		businesses = []int{zuiyou, pipi, haiwai, zhongdong, mama, shangyehua}
	}
	for _, v := range businesses {
		temp := make(map[string]interface{})

		temp["code"] = v
		temp["name"] = businessCodeNameMap[v]

		businessResp = append(businessResp, temp)
	}
	// 对resp中的值进行排序
	for i := 0; i < len(businessResp)-1; i++ {
		for j := 0; j < len(businessResp)-i-1; j++ {
			jv := businessResp[j]["code"]
			j1v := businessResp[j+1]["code"]
			if jv.(int) > j1v.(int) {
				swap := businessResp[j]
				businessResp[j] = businessResp[j+1]
				businessResp[j+1] = swap
			}
		}
	}
	return businessResp
}

// 0：最右，1：皮皮，2：海外，3：中东，4：妈妈，5：商业化
const (
	zuiyou = iota
	pipi
	haiwai
	zhongdong
	mama
	shangyehua
)

const (
	zuiyou_name     = "最右"
	pipi_name       = "皮皮"
	haiwai_name     = "海外"
	zhongdong_name  = "中东"
	mama_name       = "妈妈社区"
	shangyehua_name = "商业化"
)
