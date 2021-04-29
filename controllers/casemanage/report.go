package controllers

import (
	"go_autoapi/models"
	"strconv"
)

func (c *CaseManageController) ShowReport() {
	c.TplName = "report.html"
}

func (c *CaseManageController) GetAllReport() {
	var rp = models.AutoResult{}
	//var ids = models.Ids{}
	page, _ := strconv.Atoi(c.GetString("page"))
	limit, _ := strconv.Atoi(c.GetString("limit"))
	//count := ids.GetCollectionLength("result")
	result, count, err := rp.GetAllResult(page, limit)
	if err != nil {
		c.FormErrorJson(-1, "获取报告列表数据失败")
	}
	c.FormSuccessJson(count, result)
}
