package controllers

import (
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"go_autoapi/models"
	"strconv"
)

func (c *CaseManageController) ShowReport() {
	c.TplName = "report.html"
}

func (c *CaseManageController) GetAllReport() {
	var rp = models.AutoResult{}
	var ids = models.Ids{}
	page, _ := strconv.Atoi(c.GetString("page"))
	limit, _ := strconv.Atoi(c.GetString("limit"))
	count := ids.GetCollectionLength("result")
	fmt.Println(page, limit)
	result,err := rp.GetAllResult(page,limit)
	if err!=nil{
		logs.Error("获取全部结果失败")
	}
	c.FormSuccessJson(count, result)
}