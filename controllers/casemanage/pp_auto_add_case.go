package controllers

import (
	"encoding/json"
	"github.com/beego/beego/v2/core/logs"
	"go_autoapi/constants"
	"go_autoapi/models"
	"go_autoapi/utils"
	"strings"
	"time"

	//"go_autoapi/constants"
	//"go_autoapi/models"
	//"strings"
	//"time"
)


func (c *CaseManageController) PPAutoAddOneCase() {
	now := time.Now().Format(constants.TimeFormat)
	acm := models.TestCaseMongo{}
	//dom := models.Domain{}
	v := make(map[string]interface{})
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &v)
	if err != nil{
		logs.Error(err)
	}
	reqMethod := v["method"]
	reqName := v["title"]
	var reqPath string
	var paraMeter []byte
	m := make(map[string]interface{})
	if reqMethod == "POST"{
		reqPath = v["path"].(string)
		reqBody:= v["req_body_other"].(string)
		//reqBodyR := strings.Replace(reqBody,"\\","",-1)
		i := make(map[string]interface{})
		err := json.Unmarshal([]byte(reqBody), &i)
		if err != nil{
			logs.Error("json格式化失败")
		}
		//var m map[string]interface{}
		for k,s := range i["properties"].(map[string]interface{}){
			sm := s.(map[string]interface{})
			mock := sm["mock"]
			if mock == nil{
				mock = ""
			}else{
				mock = mock.(map[string]interface{})["mock"]
			}
			types := sm["type"]
			if (types == "integer" || types == "number") && mock !=""{
					m[k] = mock.(int)
			}else{
				m[k] = mock.(string)
			}
		}
	paraMeter,err = json.Marshal(m)
	}else{
		reqPath = v["path"].(string)
	}
	r := utils.GetRedis()
	testCaseId, err := r.Incr(constants.TEST_CASE_PRIMARY_KEY).Result()
	acm.Id = testCaseId
	acm.ApiName = "apiname"
	acm.Description = "皮皮自动添加用例"
	acm.CaseName = reqName.(string)
	acm.CreatedAt = now
	acm.UpdatedAt = now
	acm.Status = 0
	acm.BusinessCode = "1"
	acm.BusinessName = "皮皮"
	acm.ServiceName = "网关"
	acm.Domain = "http://testapi.ippzone.com"
	acm.ApiUrl = strings.TrimSpace(reqPath)
	acm.Parameter = string(paraMeter)
	acm.Author = "皮皮自动添加"
	acm.RequestMethod = reqMethod.(string)
	acm.Checkpoint = "{\"$.ret\":{\"eq\":1}}"
	if err := acm.AddCase(acm); err != nil {
		logs.Error("保存Case报错，err: ", err)
		c.ErrorJson(-1, "保存Case出错啦", nil)
	}
	c.SuccessJson("用例存储成功")
}


