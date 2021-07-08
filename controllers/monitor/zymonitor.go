package monitor

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	_ "github.com/widuu/gojson"
	"go_autoapi/libs"
	"go_autoapi/models"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type ZYMonitorController struct {
	libs.BaseController
}

func (c *ZYMonitorController) Get() {
	do := c.GetMethodName()
	switch do {
	case "mvp":
		c.mvp()
	//case "mvp1":
	//	c.mvp1()
	default:
		logs.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持", nil)
	}
}

//const url = "http://172.16.3.127:1090/api/v1/query_range?query=xmcs_gateway_acnt_http_latency_quantile%7Bquantile%3D%22p99%22%7D&start=1625558640&end=1625559540&step=15"
var url = "http://172.16.3.127:1090/api/v1/query_range?query=xmcs_gateway_acnt_http_latency_quantile%7Bquantile%3D%22p99%22%7D"

func (c *ZYMonitorController) mvp() {
	st := c.GetString("start")
	et := c.GetString("end")
	step := c.GetString("step")
	start, _ := time.Parse(models.Time_format, st)
	end, _ := time.Parse(models.Time_format, et)

	startTime := start.Unix()
	endTime := end.Unix()

	startStr := strconv.Itoa(int(startTime))
	endStr := strconv.Itoa(int(endTime))

	url += "&start=" + startStr + "&end=" + endStr + "&step=" + step

	resp, err := http.Get(url)
	if err != nil {
		logs.Error("发送get请求报错, err: ", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error("发送get请求报错, err: ", err)
	}

	lastRes := map[string]interface{}{}

	//res := string(body)
	res := make(map[string]interface{})
	json.Unmarshal(body, &res)
	status := res["status"].(string)
	if status != "success" {
		fmt.Printf("请求结果不是success")
		fmt.Printf("请求结果为: %v", res)
	} else {
		data := make(map[string]interface{})
		data = res["data"].(map[string]interface{})
		results := []interface{}{}
		results = data["result"].([]interface{})
		for _, r := range results {
			result := make(map[string]interface{})
			result = r.(map[string]interface{})
			// 当有NaN时，该条数据不进行统计
			unTJ := false
			metric := make(map[string]interface{})
			metric = result["metric"].(map[string]interface{})
			uri := metric["uri"].(string)
			fmt.Printf(uri + "\n")
			rtMap := make(map[interface{}]string)
			values := []interface{}{}
			values = result["values"].([]interface{})
			for _, val := range values {
				arr := []interface{}{}
				arr = val.([]interface{})
				rt := arr[1].(string)
				if rt == "NaN" {
					unTJ = true
					break
				}
				//rt := make(map[string]interface{})
				//rt[""]
				//append(rts,rt )
				rtMap[arr[0]] = arr[1].(string)
			}
			if !unTJ {
				//break
				lastRes[uri] = rtMap
			}
		}
		// todo
		str, _ := json.Marshal(lastRes)
		fmt.Printf("打印结果为: %s", string(str))
	}

	//fmt.Println(res)
	// 将body总结
	c.SuccessJson(lastRes)
}

//func (c *ZYMonitorController) mvp1() {
//	resp, err := http.Get(url)
//	if err != nil {
//		logs.Error("发送get请求报错, err: ", err)
//	}
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		logs.Error("发送get请求报错, err: ", err)
//	}
//
//	resMap := make(map[string][]string)
//
//	bodyStr := string(body)
//	fmt.Printf("打印body: %s", bodyStr)
//	res := gojson.Json(bodyStr)
//	status := res.Get("status").Tostring()
//	if status != "success" {
//		fmt.Printf("请求结果不是success")
//	} else {
//		data := gojson.Json(bodyStr)
//		data.Get("data").Tostring()
//		fmt.Printf("打印json-data为：%s", data)
//		//result := data.Get("result").StringtoArray()
//
//		result := gojson.Json(bodyStr).Getpath("data", "result")
//		//resultList := result.Tostring()
//		//resultStr := result.Tostring()
//		//resultlist := resultStr.([]interface{})
//		//resultList = result.Getdata().([]string)
//		resultList :=  toArray(result)
//		for _, r := range resultList {
//			//当有NaN时，该条数据不进行统计
//			unTJ := false
//			mertric := gojson.Json(r)
//			uri := mertric.Get("uri").Tostring()
//			fmt.Printf(uri)
//			values := mertric.Get("values").StringtoArray()
//			for _, val := range values {
//				vs := gojson.Json(val).StringtoArray()
//				rt := vs[1]
//				if rt == "NaN" {
//					unTJ = true
//					break
//				}
//			}
//			if unTJ {
//				break
//			}
//			resMap[uri] = values
//		}
//		//todo
//
//	}
//
//	fmt.Println(resMap)
//	// 将body总结
//	c.SuccessJson(nil)
//}
