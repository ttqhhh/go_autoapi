package controllers

import (
	"encoding/json"
	"github.com/beego/beego/v2/core/logs"
	"go_autoapi/libs"
	"go_autoapi/models"
	"gopkg.in/mgo.v2/bson"
)

type FindData struct {
	key   string
	value string
}

func (c *CaseManageController) GetCasesByQuery() {
	fd := FindData{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &fd); err != nil {
		logs.Error("获取用例，解析json数据")
	}
	query := bson.M{fd.key: fd.value, "status": "0"}
	//que_str =
	acm := models.TestCaseMongo{}
	result, err := acm.GetCasesByQuery(query)
	if err != nil {
		logs.Error("通过" + fd.key + "获取用例失败")
		logs.Error(1024, err)
	}
	c.SuccessJson(result)
}

func (c *CaseManageController) GetCaseIdByService() {
	services := c.GetStrings("service")
	ids := libs.GetCasesByServices(services)
	c.SuccessJson(ids)
}

// CaseSet添加Case时，筛选Case接口
func (c *CaseManageController) GetCaseByCondition() {
	business_code := c.GetString("business_code") // 必填
	service := c.GetString("service")
	case_name := c.GetString("case_name")

	mongo := models.TestCaseMongo{}
	caseList, err := mongo.GetCasesByCondition(business_code, service, case_name)
	if err != nil {
		c.ErrorJson(-1, "指定条件获取测试用例失败", nil)
	}

	c.SuccessJson(caseList)
}
