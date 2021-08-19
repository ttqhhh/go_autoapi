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
	case "excute_at_first_time":
		c.excuteAtFirstTime()
	case "excute_one_time":
		c.excuteOneTime()
	case "index":
		c.index()
	case "alert_visual":
		c.alertVisual()
	case "this_week_alert":
		c.ThisWeekAlert()
	case "list_2_week_trend":
		c.Last2WeekTrend()

	default:
		logs.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持", nil)
	}
}

func (c *ZYMonitorController) Post() {
	do := c.GetMethodName()
	switch do {
	case "set_rt_threshold":
		c.setRtThreshold()
	case "query_alert_by_id": //暂时不使用
		c.queryAlertById()
	default:
		logs.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持", nil)
	}
}

func (c *ZYMonitorController) excuteOneTime() {
	monitor.MonitorTask()
	c.SuccessJson(nil)
}

//const url = "http://172.16.3.127:1090/api/v1/query_range?query=xmcs_gateway_acnt_http_latency_quantile%7Bquantile%3D%22p99%22%7D&start=1625558640&end=1625559540&step=15"
//const baseUrl = "http://172.16.3.127:1090/api/v1/query_range?query=xmcs_gateway_acnt_http_latency_quantile%7Bquantile%3D%22p99%22%7D"

/**
功能上线时，跑出来所有接口的阈值
todo only在上线时调用一次!!!
*/
func (c *ZYMonitorController) excuteAtFirstTime() {
	// 起一个协程，串行的去执行该逻辑
	go func() {
		for _, serv := range monitor.Zuiyou_Servs {
			fmt.Printf("开始构建数据的服务是: %s", serv)
			GenerateLast14DaysRtData(serv)
		}
		fmt.Printf("首次执行时，数据构建结束~~~ ~~~")
	}()
	c.SuccessJson(nil)
}

func (c *ZYMonitorController) index() {
	c.TplName = "alert_list.html"
}

func (c *ZYMonitorController) alertVisual() {
	alerts := c.GetString("alerts")
	data := map[interface{}]interface{}{}
	data["alerts"] = alerts

	c.Data = data
	c.TplName = "alert_visual.html"
}

func (c *ZYMonitorController) ThisWeekAlert() {
	//serviceCode := c.GetString("service_code")
	//uri := c.GetString("uri")
	//oclock := c.GetString("oclock")
	//
	//todayZoreTimestamp := GetTodayZeroClock()
	//last14DateMap := GetLast14DaysDate(todayZoreTimestamp)
	mongo := models.RtDetailAlertMongo{}
	page, _ := strconv.Atoi(c.GetString("page"))
	limit, _ := strconv.Atoi(c.GetString("limit"))
	alertInfos, count, err := mongo.GetAllOneWeekAlertInfo(page, limit)
	if err != nil {
		//logs.Error("获取本周报警数据报错, err: ", err)
		c.ErrorJson(-1, "获取本周报警数据报错", nil)
	}
	//alertMap := map[string]models.RtDetailAlertMongo{}
	alertList := []interface{}{}
	for _, alert := range alertInfos {
		oneAlert := map[string]interface{}{}
		oneAlert["service_code"] = alert.ServiceCode
		oneAlert["uri"] = alert.Uri
		oneAlert["create_at"] = alert.CreatedAt
		oneAlert["id"] = alert.Id
		oneAlert["threshold_rt"] = alert.ThresholdRt
		oneAlert["avg_rt"] = alert.AvgRt
		oneAlert["type"] = alert.Type
		oneAlert["business"] = alert.Business
		oneAlert["rt"] = alert.Rt
		oneAlert["avg_threshold_rt"] = alert.AvgThresholdRt //历史平均响应时间
		oneAlert["reason"] = alert.Reason
		alertList = append(alertList, oneAlert)
	}
	c.FormSuccessJson(int64(count), alertList)
}
func (c *ZYMonitorController) queryAlertById() {
	idstr := c.GetString("id")
	fmt.Printf(idstr)
	id, err := strconv.Atoi(idstr)
	if err != nil {
		logs.Error("请求参数类型转换报错， err:", err)
		c.ErrorJson(-1, "请求参数转换异常", nil)
	}
	mongo := models.RtDetailAlertMongo{}
	alert, err := mongo.GetOneAlertById(int64(id))
	if err != nil {
		logs.Info("通过id查询警报错误")
	}
	c.SuccessJson(alert)

}
func (c *ZYMonitorController) Last2WeekTrend() {
	serviceCode := c.GetString("service_code")
	uri := c.GetString("uri")
	oclock := c.GetString("oclock")

	todayZoreTimestamp := GetTodayZeroClock()
	last14DateMap := GetLast14DaysDate(todayZoreTimestamp)

	rtDetail := &models.RtDetailMongo{}
	rtDetail, err := rtDetail.GetByServiceAndUri(serviceCode, uri)
	if err != nil {
		logs.Error("查询接口响应数据详情报错, err: ", err)
		c.ErrorJson(-1, "查询接口响应数据详情报错", nil)
	}

	times := []string{}
	rts := []int{}
	times = append(times, last14DateMap[-14]+" "+oclock)
	rts = append(rts, getRtByOclock(rtDetail.Last14DayRt, oclock))

	times = append(times, last14DateMap[-13]+" "+oclock)
	rts = append(rts, getRtByOclock(rtDetail.Last13DayRt, oclock))

	times = append(times, last14DateMap[-12]+" "+oclock)
	rts = append(rts, getRtByOclock(rtDetail.Last13DayRt, oclock))

	times = append(times, last14DateMap[-11]+" "+oclock)
	rts = append(rts, getRtByOclock(rtDetail.Last13DayRt, oclock))

	times = append(times, last14DateMap[-10]+" "+oclock)
	rts = append(rts, getRtByOclock(rtDetail.Last13DayRt, oclock))

	times = append(times, last14DateMap[-9]+" "+oclock)
	rts = append(rts, getRtByOclock(rtDetail.Last13DayRt, oclock))

	times = append(times, last14DateMap[-8]+" "+oclock)
	rts = append(rts, getRtByOclock(rtDetail.Last13DayRt, oclock))

	times = append(times, last14DateMap[-7]+" "+oclock)
	rts = append(rts, getRtByOclock(rtDetail.Last13DayRt, oclock))

	times = append(times, last14DateMap[-6]+" "+oclock)
	rts = append(rts, getRtByOclock(rtDetail.Last13DayRt, oclock))

	times = append(times, last14DateMap[-5]+" "+oclock)
	rts = append(rts, getRtByOclock(rtDetail.Last13DayRt, oclock))

	times = append(times, last14DateMap[-4]+" "+oclock)
	rts = append(rts, getRtByOclock(rtDetail.Last13DayRt, oclock))

	times = append(times, last14DateMap[-3]+" "+oclock)
	rts = append(rts, getRtByOclock(rtDetail.Last13DayRt, oclock))

	times = append(times, last14DateMap[-2]+" "+oclock)
	rts = append(rts, getRtByOclock(rtDetail.Last13DayRt, oclock))

	times = append(times, last14DateMap[-1]+" "+oclock)
	rts = append(rts, getRtByOclock(rtDetail.Last13DayRt, oclock))

	todayDate := time.Now().Format(models.Time_format)[:10]
	times = append(times, todayDate+" "+oclock)
	rts = append(rts, getRtByOclock(rtDetail.Last0DayRt, oclock))

	result := map[string]interface{}{}
	result["times"] = times
	result["rts"] = rts
	c.SuccessJson(&result)
}

func getRtByOclock(rtStr string, oclock string) (rt int) {
	rtMap := map[string]int{}
	json.Unmarshal([]byte(rtStr), &rtMap)
	rt = rtMap[oclock]
	return
}

//func (c *ZYMonitorController) listAlert() {
//	mongo := models.RtDetailAlertMongo{}
//	list, err := mongo.SummaryLast2WeekAlert()
//	if err != nil {
//	    logs.Error("查询报警数据报错")
//	}
//	alertMap := map[string][]models.RtDetailAlertMongo{}
//	for _, alert := range list {
//		alertType := alert.Type
//		if alertType == models.SLOW_INCREASE_RT_ALERT {
//			continue
//		}
//		serviceCode := alert.ServiceCode
//		uri := alert.Uri
//		time := alert.CreatedAt
//	}
//}

type SetRtThresholdParam struct {
	ServiceCode string `json:"service_code"`
	Uri         string `json:"uri"`
	RtThreshold int    `json:"rt_threshold"`
}

/**
设置响应时间阈值
*/
func (c *ZYMonitorController) setRtThreshold() {
	param := SetRtThresholdParam{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &param); err != nil {
		c.ErrorJson(-1, "请求错误", nil)
	}
	c.SuccessJson(nil)
}

/**
分别计算出每个接口过去7天每个整点的平均值，并入库In2DB
*/
func GenerateLast14DaysRtData(serviceCode string) {
	// 当前时间向前取14天,昨天算第一天
	// todo 验证下返回的数据在遍历时顺序是否正确
	last14DaysZeroTimes := GetLast14DaysZeroClock(GetTodayZeroClock())
	for lastN, zeroTime := range last14DaysZeroTimes {
		// 当天0时
		stepStr := "3600"
		startStr := strconv.Itoa(zeroTime)
		endStr := strconv.Itoa(zeroTime + 60*60*24 - 1)

		url := monitor.ZyPormtheusQueryRangeUrl + serviceCode + "_http_latency_quantile%7Bquantile%3D%22p99%22%7D&start=" + startStr + "&end=" + endStr + "&step=" + stepStr
		client := &http.Client{Timeout: 10 * time.Second}
		reqest, err := http.NewRequest("GET", url, nil)
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
			fmt.Printf("请求结果不是success, url为: %s\n, 响应体为: %v\n", url, res)
			continue
		} else {
			fmt.Printf("请求成功，打印请求结果为: %v\n", res)
			data := make(map[string]interface{})
			data = res["data"].(map[string]interface{})
			results := []interface{}{}
			results = data["result"].([]interface{})
			for _, r := range results {
				result := make(map[string]interface{})
				result = r.(map[string]interface{})
				metric := make(map[string]interface{})
				metric = result["metric"].(map[string]interface{})
				if metric["uri"] == nil {
					continue
				}
				uri := metric["uri"].(string)
				// 当uri包含如下字符时，不对该uri进行统计
				if strings.Contains(uri, "sign") || strings.Contains(uri, "%") || strings.Contains(uri, "=") || strings.Contains(uri, "?") || strings.Contains(uri, "+") || strings.Contains(uri, ":") || strings.Contains(uri, "(") || strings.Contains(uri, ")") || strings.Contains(uri, "\\") || strings.Contains(uri, "\"") || strings.Contains(uri, "'") || strings.Contains(uri, ".") || strings.Contains(uri, " ") || strings.Contains(uri, "|") || strings.Contains(uri, "@") || strings.Contains(uri, ";") || strings.Contains(uri, ":") || strings.Contains(uri, ",") || strings.Contains(uri, "!") || strings.Contains(uri, "<") || strings.Contains(uri, ">") || strings.Contains(uri, "--") || strings.Contains(uri, "//") {
					continue
				}
				values := []interface{}{}
				values = result["values"].([]interface{})
				rtMap := make(map[string]int)
				for _, val := range values {
					jsonByte, _ := json.Marshal(val)
					jsonStr := string(jsonByte)
					rt := gojson.Json(jsonStr).Getindex(2).Tostring()
					// 响应时间向下取整
					index := strings.Index(rt, ".")
					if index != -1 {
						rt = string([]byte(rt)[:index])
					}
					// todo 将时间戳转换为时间字符串
					// 该rt的时间戳
					onDate := gojson.Json(jsonStr).Arrayindex(1)
					key := fmt.Sprintf("%v", onDate)
					sec, _ := strconv.Atoi(key)
					secTime := time.Unix(int64(sec), 0)
					key = secTime.Format(models.Time_format)
					key = key[11:]
					rtInt := -1
					if rt != "NAN" && rt != "" {
						rtInt, _ = strconv.Atoi(rt)
					}
					rtMap[key] = rtInt
				}
				// todo 把rtMap入库
				RtDetailInDb(lastN, serviceCode, uri, rtMap)
			}
		}
	}
}

/**
时间戳倒序
*/
func GetLast14DaysZeroClock(todayZoreTime int) map[int]int {
	result := map[int]int{}
	for i := 0; i < 14; i++ {
		oneTime := todayZoreTime - 86400*(i+1)
		result[-i-1] = oneTime
	}
	fmt.Printf("拿到的14个时间戳为：%v\n", result)
	return result
}

func GetLast14DaysDate(todayZoreTime int) map[int]string {
	result := map[int]string{}
	for i := 0; i < 14; i++ {
		oneTime := todayZoreTime - 86400*(i+1)
		oneDate := time.Unix(int64(oneTime), 0).Format(models.Time_format)[:10]
		result[-i-1] = oneDate
	}
	fmt.Printf("拿到的14个日期为：%v\n", result)
	return result
}

func GetTodayZeroClock() int {
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
	fmt.Printf("打印时间戳：%d", shijianchuo)
	return int(shijianchuo)
}

func RtDetailInDb(lastNday int, serviceCode string, uri string, rtMap map[string]int) {
	mongo := &models.RtDetailMongo{}
	mongo, err := mongo.GetByServiceAndUri(serviceCode, uri)
	if err != nil {
		return
	}
	if mongo == nil {
		mongo = &models.RtDetailMongo{}
	}
	rtDetailBytes, err := json.Marshal(rtMap)
	if err != nil {
		fmt.Printf("格式化rtMap为json字符串时失败 ...")
		return
	}
	rtDetailStr := string(rtDetailBytes)
	switch lastNday {
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
		// todo--done 新增
		mongo.CreatedAt = time.Now().Format(models.Time_format)
		mongo.ServiceCode = serviceCode
		mongo.Uri = uri
		mongo.Insert(*mongo)
	} else {
		// todo--done 更新
		mongo.UpdatedAt = time.Now().Format(models.Time_format)
		// 符合条件时，计算出avgRt的值
		avgRtMap := avgRt(*mongo)
		avgRtJsonStr, _ := json.Marshal(avgRtMap)
		mongo.AvgRt = string(avgRtJsonStr)
		mongo.UpdateById(mongo.Id, *mongo)
	}
}

/**
计算出该接口过去14天里每个整点的rt平均值，结果值向下取整
*/
func avgRt(mongo models.RtDetailMongo) (avgRtMap map[string]int) {
	avgRtMap = map[string]int{}
	last1Dayrt := mongo.Last1DayRt
	last2Dayrt := mongo.Last2DayRt
	last3Dayrt := mongo.Last3DayRt
	last4Dayrt := mongo.Last4DayRt
	last5Dayrt := mongo.Last5DayRt
	last6Dayrt := mongo.Last6DayRt
	last7Dayrt := mongo.Last7DayRt
	last8Dayrt := mongo.Last8DayRt
	last9Dayrt := mongo.Last9DayRt
	last10Dayrt := mongo.Last10DayRt
	last11Dayrt := mongo.Last11DayRt
	last12Dayrt := mongo.Last12DayRt
	last13Dayrt := mongo.Last13DayRt
	last14Dayrt := mongo.Last14DayRt
	// 循环 计算出 过去14天每个时间点的接口响应平均值，低于7天不计算平均值(向下取整)
	//res := map[string]int{}
	for i := 0; i < 24; i++ {
		avgRt := -1
		count := 0
		totalRt := 0

		sumRt(i, &count, &totalRt, last1Dayrt)
		sumRt(i, &count, &totalRt, last2Dayrt)
		sumRt(i, &count, &totalRt, last3Dayrt)
		sumRt(i, &count, &totalRt, last4Dayrt)
		sumRt(i, &count, &totalRt, last5Dayrt)
		sumRt(i, &count, &totalRt, last6Dayrt)
		sumRt(i, &count, &totalRt, last7Dayrt)
		sumRt(i, &count, &totalRt, last8Dayrt)
		sumRt(i, &count, &totalRt, last9Dayrt)
		sumRt(i, &count, &totalRt, last10Dayrt)
		sumRt(i, &count, &totalRt, last11Dayrt)
		sumRt(i, &count, &totalRt, last12Dayrt)
		sumRt(i, &count, &totalRt, last13Dayrt)
		sumRt(i, &count, &totalRt, last14Dayrt)

		// 根据totalRt / count 计算平均值
		key := ""
		if i < 10 {
			//key = "0"+string(i) + ":00:00"
			key = fmt.Sprintf("0%v:00:00", i)
		} else {
			//key = string(i) + ":00:00"
			key = fmt.Sprintf("%v:00:00", i)
		}
		// 只有当统计的天数大于7天时，才计算平均响应时间
		if count >= 7 {
			avgRt = totalRt / count
		}
		avgRtMap[key] = avgRt
	}
	return
}

func sumRt(oclock int, count *int, totalRt *int, rtOneDay string) {
	// 当存在当天的24小时数据时，才会进if逻辑
	if rtOneDay != "" {
		jsonObject := map[string]int{}
		json.Unmarshal([]byte(rtOneDay), &jsonObject)
		key := ""
		if oclock < 10 {
			//key = "0"+string(oclock) + ":00:00"
			key = fmt.Sprintf("0%v:00:00", oclock)
		} else {
			key = fmt.Sprintf("%v:00:00", oclock)
		}
		rt := jsonObject[key]
		// 只有当响应时间字符串有有效值（非0、非空、非NaN）时，才会去取值计算
		if rt != 0 && rt != -1 {
			//rt, _ := strconv.Atoi(RtStr)
			*totalRt += rt
			*count++
		}
	}
}

func ElementIndexInSlice(arr []string, ele string) (index int) {
	for i, s := range arr {
		if s == ele {
			return i
		}
	}
	return -1
}
