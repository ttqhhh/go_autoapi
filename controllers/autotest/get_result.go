package controllers

import (
	"encoding/json"
	"github.com/shopspring/decimal"
	constant "go_autoapi/constants"
	"go_autoapi/models"
	"go_autoapi/utils"
)

type DoProcess struct {
	Uuid  string `json:"uuid"`
	Count int64  `json:"count"`
}

// 获取测试执行进度
func (c *AutoTestController) getProcess() {
	dp := DoProcess{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &dp); err != nil {
		c.ErrorJson(-1, "请求参数错误", nil)
	}
	r := utils.GetRedis()
	defer r.Close()
	hasCount, _ := r.Get(constant.RUN_RECORD_CASE_DONE_NUM + dp.Uuid).Int64()
	progress := decimal.NewFromFloat(float64(hasCount)).Div(decimal.NewFromFloat(float64(dp.Count)))
	c.SuccessJsonWithMsg(map[string]interface{}{"progress": progress}, "OK")
}

func (c *AutoTestController) getResult() {
	dp := DoProcess{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &dp); err != nil {
		c.ErrorJson(-1, "请求参数错误", nil)
	}
	res, _ := models.GetResultByRunId(dp.Uuid)
	c.SuccessJsonWithMsg(res, "OK")
}
