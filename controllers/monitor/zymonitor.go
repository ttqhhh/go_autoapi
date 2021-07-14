package monitor

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/widuu/gojson"
	"go_autoapi/libs"
	"go_autoapi/models"
	"go_autoapi/task/monitor"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ZYMonitorController struct {
	libs.BaseController
}

func (c *ZYMonitorController) Get() {
	do := c.GetMethodName()
	switch do {
	case "mvp":
		c.getRtDetailByRange()
	case "mvp1":
		c.mvp1()
	case "test":
		c.test()
	default:
		logs.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持", nil)
	}
}

//const url = "http://172.16.3.127:1090/api/v1/query_range?query=xmcs_gateway_acnt_http_latency_quantile%7Bquantile%3D%22p99%22%7D&start=1625558640&end=1625559540&step=15"
const baseUrl = "http://172.16.3.127:1090/api/v1/query_range?query=xmcs_gateway_acnt_http_latency_quantile%7Bquantile%3D%22p99%22%7D"

/**
获取接口每个整点的响应时间详情
通过起始时间、终止时间、步长（建议3600，每一个整点）
*/
func (c *ZYMonitorController) getRtDetailByRange() {
	st := c.GetString("start")
	et := c.GetString("end")
	step := c.GetString("step")
	loc, _ := time.LoadLocation("Local")
	start, _ := time.ParseInLocation(models.Time_format, st, loc)
	end, _ := time.ParseInLocation(models.Time_format, et, loc)

	startTime := start.Unix()
	endTime := end.Unix()

	startStr := strconv.Itoa(int(startTime))
	endStr := strconv.Itoa(int(endTime))

	url := baseUrl + "&start=" + startStr + "&end=" + endStr + "&step=" + step

	client := &http.Client{Timeout: 5 * time.Second}
	reqest, err := http.NewRequest("GET", url, nil)
	//resp, err := client.Get(url)
	reqest.Header.Set("Cache-Control", "no-cache")
	resp, err := client.Do(reqest)

	if err != nil {
		logs.Error("发送get请求报错, err: ", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error("发送get请求报错, err: ", err)
	}

	lastRes := make(map[string]interface{})

	//res := string(body)
	res := make(map[string]interface{})
	json.Unmarshal(body, &res)
	status := res["status"].(string)
	if status != "success" {
		fmt.Printf("请求结果不是success")
		fmt.Printf("请求结果为: %v", res)
	} else {
		fmt.Printf("打印请求结果为: %v", res)
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
			rtMap := make(map[string]string)
			values := []interface{}{}
			values = result["values"].([]interface{})
			for _, val := range values {
				//arr := []interface{}{}
				//arr = val.([]interface{})
				jsonByte, _ := json.Marshal(val)
				jsonStr := string(jsonByte)
				onTime := gojson.Json(jsonStr).Arrayindex(1)
				rt := gojson.Json(jsonStr).Getindex(2).Tostring()
				//rt := arr[1].(string)
				if rt == "NaN" {
					unTJ = true
					break
				}
				//rt := make(map[string]interface{})
				//rt[""]
				//append(rts,rt )
				key := fmt.Sprintf("%v", onTime)
				//rtMap[key] = arr[1].(string)
				index := strings.Index(rt, ".")
				if index != -1 {
					rt = string([]byte(rt)[:index])
				}
				// todo 将时间戳转换为时间字符串
				sec, _ := strconv.Atoi(key)
				secTime := time.Unix(int64(sec), 0)
				key = secTime.Format(models.Time_format)
				rtMap[key] = rt
			}
			if !unTJ {
				//break
				lastRes[uri] = rtMap
			}
		}
		// todo
		strs, _ := json.Marshal(lastRes)
		fmt.Printf("打印结果为: %s", string(strs))
	}

	//fmt.Println(res)
	// 将body总结
	c.SuccessJson(lastRes)
}

/**
分别计算出每个接口过去7天每个整点的平均值，
*/
func (c *ZYMonitorController) mvp1() {
	// 最终结构体
	dateMap := map[string]interface{}{}
	// 当前时间向前取7天
	last7DaysZeroTime := getLast7DaysZeroClock(getTodayZeroClock())
	for i := 0; i < len(last7DaysZeroTime); i++ {
		// 获取当天0时
		zeroTime := last7DaysZeroTime[i]
		step := "3600"
		startStr := strconv.Itoa(zeroTime)
		endStr := strconv.Itoa(zeroTime + 86399)
		date := time.Unix(int64(zeroTime), 0).Format("2006-01-02")

		url := baseUrl + "&start=" + startStr + "&end=" + endStr + "&step=" + step
		client := &http.Client{Timeout: 5 * time.Second}
		reqest, err := http.NewRequest("GET", url, nil)
		reqest.Header.Set("Cache-Control", "no-cache")
		resp, err := client.Do(reqest)
		if err != nil {
			logs.Error("发送get请求报错, err: ", err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logs.Error("发送get请求报错, err: ", err)
		}

		res := make(map[string]interface{})
		json.Unmarshal(body, &res)
		status := res["status"].(string)
		if status != "success" {
			fmt.Printf("请求结果不是success")
			fmt.Printf("请求结果为: %v", res)
		} else {
			fmt.Printf("打印请求结果为: %v", res)
			data := make(map[string]interface{})
			data = res["data"].(map[string]interface{})
			results := []interface{}{}
			results = data["result"].([]interface{})
			uriMap := map[string]interface{}{}
			for _, r := range results {
				result := make(map[string]interface{})
				result = r.(map[string]interface{})
				// 当有NaN时，该条数据不进行统计
				//unTJ := false
				metric := make(map[string]interface{})
				metric = result["metric"].(map[string]interface{})
				uri := metric["uri"].(string)
				fmt.Printf(uri + "\n")
				values := []interface{}{}
				values = result["values"].([]interface{})
				rtMap := make(map[string]string)
				for _, val := range values {
					jsonByte, _ := json.Marshal(val)
					jsonStr := string(jsonByte)
					rt := gojson.Json(jsonStr).Getindex(2).Tostring()
					//if rt == "NaN" {
					//	unTJ = true
					//	break
					//}
					// 该rt的时间戳
					onDate := gojson.Json(jsonStr).Arrayindex(1)
					key := fmt.Sprintf("%v", onDate)
					index := strings.Index(rt, ".")
					if index != -1 {
						// 响应时间向下取整
						rt = string([]byte(rt)[:index])
					}
					// todo 将时间戳转换为时间字符串
					sec, _ := strconv.Atoi(key)
					secTime := time.Unix(int64(sec), 0)
					key = secTime.Format(models.Time_format)
					rtMap[key] = rt
				}
				uriMap[uri] = rtMap
			}
			dateMap[date] = uriMap
		}
	}
	strs, _ := json.Marshal(dateMap)
	fmt.Printf("打印结果为: %s", string(strs))

	// 将body总结
	c.SuccessJson(dateMap)
}



func (c *ZYMonitorController) test() {
	//todayTimeZero := getTodayZeroClock()
	//getLast7DaysZeroClock(todayTimeZero)
	res := monitor.GetLast14DaysRtData("xmcs_gateway_acnt")
	c.SuccessJson(res)
}

func getTodayZeroClock() int {
	now := time.Now() //获取当前时间
	fmt.Printf("current time:%v\n", now)

	year := now.Year()     //年
	month := now.Month()   //月
	day := now.Day()       //日
	hour := now.Hour()     //小时
	minute := now.Minute() //分钟
	second := now.Second() //秒
	fmt.Printf("%d-%02d-%02d %02d:%02d:%02d\n", year, month, day, hour, minute, second)

	loc, _ := time.LoadLocation("Local")
	todayZero := time.Date(year, month, day, 0, 0, 0, 0, loc)
	shijianchuo := todayZero.Unix()
	//shijianchuo = shijianchuo - 86400
	fmt.Printf("打印时间戳：%d", shijianchuo)
	return int(shijianchuo)
}

/**
时间戳倒序
*/
func getLast7DaysZeroClock(todayZoreTime int) []int {
	result := []int{}
	for i := 0; i < 7; i++ {
		oneTime := todayZoreTime - 86400*(i+1)
		result = append(result, oneTime)
	}
	fmt.Printf("拿到的7个时间戳为：%v", result)
	return result
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
