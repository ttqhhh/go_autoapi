package controllers

import (
	"github.com/astaxie/beego/logs"
	"go_autoapi/constants"
	"go_autoapi/models"
	"go_autoapi/utils"
	"strconv"
	"strings"
	"time"
)

func (c *CaseManageController) AddOneCase() {
	now := time.Now().Format(constants.TimeFormat)
	acm := models.TestCaseMongo{}
	if err := c.ParseForm(&acm); err != nil { //传入user指针
		c.Ctx.WriteString("出错了！")
	}
	// todo service_id 和 service_name 在一起,需要分割后赋值
	arr := strings.Split(acm.ServiceName,";")
	acm.ServiceName = arr[1]
	id64, _ := strconv.ParseInt(arr[0], 10, 64)
	acm.ServiceId = id64
	//acm.Id = models.GetId("case")
	r:=utils.GetRedis()
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
	if business == "0" {
		acm.BusinessName = "最右"
	} else if business == "1" {
		acm.BusinessName = "皮皮"
	} else if business == "2" {
		acm.BusinessName = "海外"
	} else if business == "3" {
		acm.BusinessName = "中东"
	} else if business == "4" {
		acm.BusinessName = "妈妈社区"
	} else if business == "5" {
		acm.BusinessName = "商业化"
	}
	if err := acm.AddCase(acm); err != nil {
		logs.Error("保存Case报错，err: ", err)
		c.ErrorJson(-1, "保存Case出错啦", nil)
	}
	//c.SuccessJson("添加成功")
	c.Ctx.Redirect(302, "/case/show_cases?business="+business)
}
