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
	"fengmanlong": {constant.ZuiyYou, constant.PiPi, constant.HaiWai, constant.ZhongDong, constant.Matuan, constant.ShangYeHua, constant.HaiWaiUS},
	"xueyibing":   {constant.ZuiyYou, constant.PiPi, constant.HaiWai, constant.ZhongDong, constant.Matuan, constant.ShangYeHua, constant.HaiWaiUS},
	"wangzhen01":  {constant.ZuiyYou},
	// 普通用户
	"wangjun":       {constant.HaiWai, constant.ZhongDong, constant.Matuan, constant.HaiWaiUS},
	"sunzhiying":    {constant.HaiWai},
	"zhangjuan":     {constant.ShangYeHua, constant.ZuiyYou},
	"liuyinan":      {constant.PiPi},
	"sunmingyao":    {constant.PiPi},
	"zhangdanbing":  {constant.PiPi},
	"xufei":         {constant.HaiWai, constant.ZhongDong, constant.HaiWaiUS},
	"suhuimin":      {constant.ZhongDong},
	"wanglanjin":    {constant.ZhongDong, constant.Matuan},
	"yangjingfang":  {constant.ShangYeHua, constant.ZuiyYou},
	"chengxiaojing": {constant.ZuiyYou},
	"zhaoxiaodong":  {constant.ShangYeHua},
}

var businessCodeNameMap = map[int]string{
	constant.ZuiyYou:    zuiyou_name,
	constant.PiPi:       pipi_name,
	constant.HaiWai:     haiwai_name,
	constant.ZhongDong:  zhongdong_name,
	constant.Matuan:     mama_name,
	constant.ShangYeHua: shangyehua_name,
	constant.HaiWaiUS:   haiwaius_name,
}

func GetBusinessNameByCode(code int) string {
	businessName := businessCodeNameMap[code]
	return businessName
}

func GetUserBusinessesList(username string) []int {
	return userBusinessMap[username]
}

func (c *BusinessController) getUserBusinesses() {
	userId, _ := c.GetSecureCookie(constant.CookieSecretKey, "user_id")

	result := make(map[string]interface{})
	result["username"] = userId
	result["businesses"] = GetBusinesses(userId)
	c.SuccessJson(result)
}

func GetBusinesses(username string) []map[string]interface{} {
	businessResp := [](map[string]interface{}){}
	businesses, ok := userBusinessMap[username]
	if !ok {
		// todo 目前只对测试同学进行了限制，其他角色同学暂未进行处理
		businesses = []int{constant.ZuiyYou, constant.PiPi, constant.HaiWai, constant.ZhongDong, constant.Matuan, constant.ShangYeHua, constant.HaiWaiUS}
	} else {
		// todo 0609 测试同学现在没有限制了
		businesses = []int{constant.ZuiyYou, constant.PiPi, constant.HaiWai, constant.ZhongDong, constant.Matuan, constant.ShangYeHua, constant.HaiWaiUS}
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

// 返回所有业务线（没有任何条件限制）
func GetAllBusinesses() []map[string]interface{} {
	businessResp := [](map[string]interface{}){}
	//businesses, ok := userBusinessMap[username]
	//if !ok {
	// todo 目前只对测试同学进行了限制，其他角色同学暂未进行处理
	//businesses = []int{ZuiyYou, PiPi, HaiWai, ZhongDong, Matuan, ShangYeHua, HaiWaiUS}
	//}
	for k, v := range businessCodeNameMap {
		temp := make(map[string]interface{})
		temp["code"] = k
		temp["name"] = v
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

const (
	zuiyou_name     = "最右"
	pipi_name       = "皮皮"
	haiwai_name     = "海外"
	zhongdong_name  = "中东"
	mama_name       = "妈妈社区"
	shangyehua_name = "商业化"
	haiwaius_name   = "海外-US"
)
