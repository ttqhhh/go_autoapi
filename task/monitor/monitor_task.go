package monitor

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/widuu/gojson"
	"go_autoapi/models"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	// 每小时执行一次
	MONITOR_TASK_EXPRESSION = "0 0 * * * *"
    ZyPormtheusQueryUrl = "http://172.16.3.127:1090/api/v1/query_range?query=xmcs_"
)

func MonitorTask() error {
	logs.Info("生产接口RT监控定时任务启动执行...")
	// 查询出来当前业务线下，所有的服务，拼凑出来不同的
	serviceCode := ""
	//url := ZyPormtheusQueryUrl+serviceCode+"_latency_quantile%7Bquantile%3D%22p99%22%7D"
	getRtDetailByRange(serviceCode, "", "", 3600)
	return nil
}

/**
获取接口每个整点的响应时间详情
通过起始时间、终止时间、步长（建议3600，每一个整点）
*/
func getRtDetailByRange(serviceCode string, startTime string, endTime string, step int64) map[string]interface{} {
	loc, _ := time.LoadLocation("Local")
	st, _ := time.ParseInLocation(models.Time_format, startTime, loc)
	et, _ := time.ParseInLocation(models.Time_format, endTime, loc)

	startTimeStamp := st.Unix()
	endTimeStamp := et.Unix()

	lastRes := getRtDetailByRangeBase(serviceCode, startTimeStamp, endTimeStamp, step)

	return lastRes
}

/**
获取接口每个整点的响应时间详情
通过起始时间、终止时间、步长（建议3600，每一个整点）
*/
func getRtDetailByRangeBase(serviceCode string, startTimeStamp int64, endTimeStamp int64, step int64) map[string]interface{} {
	//loc, _ := time.LoadLocation("Local")
	//st, _ := time.ParseInLocation(models.Time_format, start, loc)
	//et, _ := time.ParseInLocation(models.Time_format, end, loc)

	//startTime := st.Unix()
	//endTime := et.Unix()

	startStr := strconv.Itoa(int(startTimeStamp))
	endStr := strconv.Itoa(int(endTimeStamp))
	stepStr := strconv.Itoa(int(step))

	url := ZyPormtheusQueryUrl+serviceCode+"_latency_quantile%7Bquantile%3D%22p99%22%7D"
	url = url + "&start=" + startStr + "&end=" + endStr + "&step=" + stepStr

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
				jsonByte, _ := json.Marshal(val)
				jsonStr := string(jsonByte)
				onTime := gojson.Json(jsonStr).Arrayindex(1)
				rt := gojson.Json(jsonStr).Getindex(2).Tostring()
				//rt := arr[1].(string)
				if rt == "NaN" {
					unTJ = true
					break
				}
				key := fmt.Sprintf("%v", onTime)
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
				lastRes[uri] = rtMap
			}
		}
		// todo
		strs, _ := json.Marshal(lastRes)
		fmt.Printf("打印结果为: %s", string(strs))
	}

	return lastRes
}

/**
分别计算出每个接口过去7天每个整点的平均值，
*/
func getLast7DaysRtData(serviceCode string) map[string]interface{} {
	// 最终结构体
	dateMap := map[string]interface{}{}
	// 当前时间向前取7天
	last7DaysZeroTime := getLast7DaysZeroClock(getTodayZeroClock())
	for i := 0; i < len(last7DaysZeroTime); i++ {
		// 获取当天0时
		zeroTime := last7DaysZeroTime[i]
		stepStr := "3600"
		startStr := strconv.Itoa(zeroTime)
		endStr := strconv.Itoa(zeroTime + 86399)
		date := time.Unix(int64(zeroTime), 0).Format("2006-01-02")

		//url := baseUrl + "&start=" + startStr + "&end=" + endStr + "&step=" + step
		url := ZyPormtheusQueryUrl+serviceCode+"_latency_quantile%7Bquantile%3D%22p99%22%7D"
		url = url + "&start=" + startStr + "&end=" + endStr + "&step=" + stepStr

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
	return dateMap
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