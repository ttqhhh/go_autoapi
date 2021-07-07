package monitor

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"go_autoapi/libs"
	"io/ioutil"
	"net/http"
)

type ZYMonitorController struct {
	libs.BaseController
}

func (c *ZYMonitorController) Get() {
	do := c.GetMethodName()
	switch do {
	case "mvp":
		c.mvp()
	default:
		logs.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持", nil)
	}
}

const url = "http://172.16.3.127:1090/api/v1/query_range?query=xmcs_gateway_acnt_http_latency_quantile%7Bquantile%3D%22p99%22%7D&start=1625558640&end=1625559540&step=15"

func (c *ZYMonitorController) mvp() {
	resp, err := http.Get(url)
	if err != nil {
		logs.Error("发送get请求报错, err: ", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error("发送get请求报错, err: ", err)
	}

	//res := string(body)
	res := make(map[string]interface{})
	json.Unmarshal(body, res)
	status := res["status"].(string)
	if status != "success" {
		fmt.Printf("请求结果不是success")
	} else {
		data := make(map[string]interface{})
		data = res["data"].(map[string]interface{})
		result := []interface{}{}
		result = data["result"].([]interface{})
		for _, r := range result {
			// 当有NaN时，该条数据不进行统计
			unTJ := false
			mertric := make(map[string]interface{})
			mertric = r.(map[string]interface{})
			uri := mertric["uri"].(string)
			values := []interface{}{}
			values = mertric["values"].([]interface{})
			for _, val := range values {
				arr := []interface{}{}
				arr = val.([]interface{})
				rt := arr[1].(string)
				if rt == "NaN" {
					unTJ = true
					break
				}
			}
			if unTJ {
				break
			}
		}
		// todo

	}

	fmt.Println(res)
	// 将body总结
	c.SuccessJson(nil)
}
