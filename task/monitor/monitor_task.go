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
	ZyPormtheusQueryUrl     = "http://172.16.3.127:1090/api/v1/query_range?query="
)

func MonitorTask() error {
	logs.Info("生产接口RT监控定时任务启动执行...")
	// 查询出来当前业务线下，所有的服务，拼凑出来不同的
	serviceCode := "xmcs_gateway_acnt"
	//url := ZyPormtheusQueryUrl+serviceCode+"_latency_quantile%7Bquantile%3D%22p99%22%7D"
	//getRtDetailByRange(serviceCode, "", "", 3600)
	last14DaysRtData, last14DaysZeroTimes, uris := GetLast14DaysRtData(serviceCode)

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

	url := ZyPormtheusQueryUrl + serviceCode + "_latency_quantile%7Bquantile%3D%22p99%22%7D"
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
func GetLast14DaysRtData(serviceCode string) (last14DayRtData map[string]interface{}, last14DaysZeroTimes map[int]int, uris []string) {
	// 最终结构体
	//last14DayRtDatas = map[string]interface{}{}
	// 该服务下的所有接口
	uris = []string{}
	// 当前时间向前取14天
	last14DaysZeroTimes = getLast14DaysZeroClock(getTodayZeroClock())
	//for i := 0; i < len(last14DaysZeroTimes); i++ {
	for k, v := range last14DaysZeroTimes {
		// 当天0时
		zeroTime := v
		stepStr := "3600"
		startStr := strconv.Itoa(zeroTime)
		endStr := strconv.Itoa(zeroTime + 86399)
		//date := time.Unix(int64(zeroTime), 0).Format("2006-01-02")

		//url := baseUrl + "&start=" + startStr + "&end=" + endStr + "&step=" + step
		url := ZyPormtheusQueryUrl + serviceCode + "_http_latency_quantile%7Bquantile%3D%22p99%22%7D"
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
			//uriMap := map[string]interface{}{}
			for _, r := range results {
				result := make(map[string]interface{})
				result = r.(map[string]interface{})
				// 当有NaN时，该条数据不进行统计
				//unTJ := false
				metric := make(map[string]interface{})
				metric = result["metric"].(map[string]interface{})
				uri := metric["uri"].(string)
				// 路径收集
				if ElementIndexInSlice(uris, uri) == -1 {
					uris = append(uris, uri)
				}

				fmt.Printf(uri + "\n")
				values := []interface{}{}
				values = result["values"].([]interface{})
				rtMap := make(map[string]int)
				for _, val := range values {
					jsonByte, _ := json.Marshal(val)
					jsonStr := string(jsonByte)
					rt := gojson.Json(jsonStr).Getindex(2).Tostring()
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
					rtInt := -1
					if rt != "NAN" && rt != "" {
						rtInt, _ = strconv.Atoi(rt)
					}
					rtMap[key] = rtInt
				}
				// todo 把rtMap入库
				rtDetailInDb(k, serviceCode, uri, rtMap)
			}
		}
	}
	return
}

func rtDetailInDb(someday int, serviceCode string, uri string, rtMap map[string]int) {
	mongo := models.RtDetailMongo{}
	mongo, err := mongo.GetByServiceAndUri(serviceCode, uri)
	if err != nil {
		return
	}
	rtDetailBytes, err := json.Marshal(rtMap)
	if err != nil {
		fmt.Printf("格式化rtMap为json字符串时失败 ...")
		return
	}
	rtDetailStr := string(rtDetailBytes)
	switch someday {
	case -1:
		mongo.Last1DayRt = rtDetailStr
	case -2:
		mongo.Last2DayRt = rtDetailStr
	case -3:
		mongo.Last3DayRt = rtDetailStr
	case -4:
		mongo.Last4DayRt = rtDetailStr
	case -5:
		mongo.Last5DayRt = rtDetailStr
	case -6:
		mongo.Last6DayRt = rtDetailStr
	case -7:
		mongo.Last7DayRt = rtDetailStr
	case -8:
		mongo.Last8DayRt = rtDetailStr
	case -9:
		mongo.Last9DayRt = rtDetailStr
	case -10:
		mongo.Last10DayRt = rtDetailStr
	case -11:
		mongo.Last11DayRt = rtDetailStr
	case -12:
		mongo.Last12DayRt = rtDetailStr
	case -13:
		mongo.Last13DayRt = rtDetailStr
	case -14:
		mongo.Last14DayRt = rtDetailStr
	}
	if mongo.Id == 0 {
		mongo.CreatedAt = time.Now().Format(models.Time_format)
		mongo.ServiceCode = serviceCode
		mongo.Uri = uri
	} else {
		mongo.UpdatedAt = time.Now().Format(models.Time_format)
		//
	}
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
func getLast14DaysZeroClock(todayZoreTime int) map[int]int {
	result := map[int]int{}
	for i := 0; i < 14; i++ {
		oneTime := todayZoreTime - 86400*(i+1)
		//result = append(result, oneTime)
		result[-i] = oneTime
	}
	fmt.Printf("拿到的14个时间戳为：%v", result)
	return result
}

func ElementIndexInSlice(arr []string, ele string) (index int) {
	for i, s := range arr {
		if s == ele {
			return i
		}
	}
	return -1
}
