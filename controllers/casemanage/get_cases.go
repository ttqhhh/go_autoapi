package controllers

import (
	"encoding/json"
	"github.com/beego/beego/v2/core/logs"
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
