package controllers

import (
	"github.com/astaxie/beego/logs"
	"go_autoapi/constants"
	controllers "go_autoapi/controllers/autotest"
	"go_autoapi/models"
	"go_autoapi/utils"
	"strconv"
	"strings"
	"time"
)

func (c *CaseManageController) AddOneCase() {
	name, _ := c.GetSecureCookie(constants.CookieSecretKey, "user_id")
	now := time.Now().Format(constants.TimeFormat)
	acm := models.TestCaseMongo{}
	dom := models.Domain{}
	if err := c.ParseForm(&acm); err != nil { // 传入user指针
		c.Ctx.WriteString("出错了！")
	}
	acm.Author = name
	// 获取域名并确认是否执行
	dom.Author = acm.Author
	intBus, _ := strconv.Atoi(acm.BusinessCode)
	dom.Business = int8(intBus)
	if strings.Contains(constants.ONLINE_DOMAIN, acm.Domain) {
		c.ErrorJson(-1, "保存Case出错啦,线下自动化禁止使用线上域名", nil)
	} else {
		dom.DomainName = acm.Domain
	}
	if err := dom.DomainInsert(dom); err != nil {
		logs.Error("添加case的时候 domain 插入失败")
	}
	// service_id 和 service_name 在一起,需要分割后赋值
	arr := strings.Split(acm.ServiceName, ";")
	acm.ServiceName = arr[1]
	id64, _ := strconv.ParseInt(arr[0], 10, 64)
	acm.ServiceId = id64
	//acm.Id = models.GetId("case")
	r := utils.GetRedis()
	defer r.Close()
	testCaseId, err := r.Incr(constants.TEST_CASE_PRIMARY_KEY).Result()
	if err != nil {
		logs.Error("保存Case时，获取从redis获取唯一主键报错，err: ", err)
		c.ErrorJson(-1, "保存Case出错啦", nil)
	}
	acm.Id = testCaseId
	acm.CreatedAt = now
	acm.UpdatedAt = now
	acm.Status = 0
	business := acm.BusinessCode
	//if business == "0" {
	//	acm.BusinessName = "最右"
	//} else if business == "1" {
	//	acm.BusinessName = "皮皮"
	//} else if business == "2" {
	//	acm.BusinessName = "海外"
	//} else if business == "3" {
	//	acm.BusinessName = "中东"
	//} else if business == "4" {
	//	acm.BusinessName = "妈妈社区"
	//} else if business == "5" {
	//	acm.BusinessName = "商业化"
	//}
	businessCode, _ := strconv.Atoi(business)
	businessName := controllers.GetBusinessNameByCode(businessCode)
	acm.BusinessName = businessName
	// 去除请求路径前后的空格
	apiUrl := acm.ApiUrl
	acm.ApiUrl = strings.TrimSpace(apiUrl)
	// todo 千万不要删，用于处理json格式化问题（删了后某些服务会报504问题）
	// todo 暂时把格式化处理的相关逻辑挪到了DoRequest中
	//param := acm.Parameter
	//v := make(map[string]interface{})
	//err = json.Unmarshal([]byte(strings.TrimSpace(param)), &v)
	//if err != nil {
	//	logs.Error("发送冒烟请求前，解码json报错，err：", err)
	//	return
	//}
	//paramByte, err := json.Marshal(v)
	//if err != nil {
	//	logs.Error("保存Case时，处理请求json报错， err:", err)
	//	c.ErrorJson(-1, "保存Case出错啦", nil)
	//}
	//acm.Parameter = string(paramByte)
	if err := acm.AddCase(acm); err != nil {
		c.ErrorJson(-1, err.Error(), nil)
	}
	//c.SuccessJson("添加成功")
	//c.Ctx.Redirect(302, "/case/show_cases?business="+business)
	c.Ctx.Redirect(302, "/case/close_windows")
}
