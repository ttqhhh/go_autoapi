package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
)

type autoCase struct {
	Id int `form:"id"`
}

func (c *CaseManageController) getCases(){
	fmt.Println("获取用例数据，返回并展示")
	ac := autoCase{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &ac); err != nil {
		logs.Error(1024, err)
		c.ErrorJson(-1, "请求错误", nil)
	}
}
