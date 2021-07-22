package monitor

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/blinkbean/dingtalk"
	"github.com/widuu/gojson"
	"go_autoapi/models"
	"go_autoapi/task/inspection_strategy"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

/**
	最右所有网关服务：chatsrv-gateway、gateway-acnt、gateway-ad、gateway-apmserver、gateway-applist、gateway-applog、
gateway-danmaku、gateway-diagnosis、gateway-earn、gateway-earn-acnt、gateway-feedback、gateway-market、gateway-misc、
gateway-mixture、gateway-post、gateway-rec、gateway-rev、gateway-shop、gateway-stat、gateway-topic、gateway-urlresolver、
gateway-vasapi、gateway-wxapp、gray-gateway、mall-gateway、media-gateway、media-gateway-op、media-gateway-upload、
miniemoji-gateway、rtcsrv-gateway、snssrv-gateway-native、snssrv-gateway-nearby、zuiyou-trade-gateway、zy-gateway-teamchat、
zy-live-gateway

*/
const (
	// 每小时执行一次
	MONITOR_DING_SEND_IS_OPEN = false
	MONITOR_TASK_EXPRESSION   = "0 0 * * * *"
	ZyPormtheusQueryUrl       = "http://172.16.3.127:1090/api/v1/query?query=xmcs_"
	ZyPormtheusQueryRangeUrl  = "http://172.16.3.127:1090/api/v1/query_range?query=xmcs_"
)

/**
1、上线时，执行一次，抛出来基线数据
2、每个整点跑1次定时任务：①和平均值or阈值进行比较 ②和历史rt值进行比较，验证正相关曲线
*/

// 最右的网关服务集合
var Zuiyou_Servs = []string{
	"chatsrv_gateway",
	"gateway_acnt",
	//"gateway_ad",
	//"gateway_apmserver",
	"gateway_applist",
	"gateway_applog",
	"gateway_danmaku",
	//"gateway_diagnosis", todo
	//"gateway_earn",
	//"gateway_earn_acnt",
	"gateway_feedback",
	//"gateway_market",
	//"gateway_misc",
	//"gateway_mixture",
	"gateway_post",
	"gateway_rec",
	"gateway_rev",
	"gateway_shop",
	//"gateway_stat", // todo
	"gateway_topic",
	"gateway_urlresolver",
	"gateway_vasapi",
	"gateway_wxapp",
	"gray_gateway",
	"mall_gateway",
	"media_gateway",
	"media_gateway_op",
	"media_gateway_upload",
	"miniemoji_gateway",
	"rtcsrv_gateway",
	"snssrv_gateway_native",
	"snssrv_gateway_nearby",
	"zuiyou_trade_gateway",
	"zy_gateway_teamchat",
	"zy_live_gateway",
}

func MonitorTask() error {
	logs.Info("生产接口RT监控定时任务启动执行...")
	// 查询出来当前业务线下，所有的服务，拼凑出来不同的
	taskTime := time.Now().Format(models.Time_format)
	// Time_format = "2006-01-02 15:04:05"
	taskTime = taskTime[:14]
	taskTime = taskTime + "00:00"
	taskTimestamp, _ := time.ParseInLocation(models.Time_format, taskTime, time.Local)
	for _, serv := range Zuiyou_Servs {
		//OneHourExcute(serv, taskTimestamp.Unix())
		HalfHourExcute(serv, taskTimestamp.Unix())
	}
	return nil
}

//func OneHourExcute(serviceCode string, timestamp int64) {
//	// 该服务下的所有接口
//	url := ZyPormtheusQueryUrl + serviceCode + "_http_latency_quantile%7Bquantile%3D%22p99%22%7D"
//
//	client := &http.Client{Timeout: 5 * time.Second}
//	reqest, err := http.NewRequest("GET", url, nil)
//	resp, err := client.Do(reqest)
//	if err != nil {
//		logs.Error("发送get请求报错, err: ", err)
//	}
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		logs.Error("发送get请求报错, err: ", err)
//	}
//
//	res := make(map[string]interface{})
//	json.Unmarshal(body, &res)
//	status := res["status"].(string)
//	if status != "success" {
//		fmt.Printf("请求结果不是success")
//		fmt.Printf("请求结果为: %v", res)
//	} else {
//		resStr, _ := json.Marshal(res)
//		fmt.Printf("打印请求结果为: %s", resStr)
//		data := make(map[string]interface{})
//		data = res["data"].(map[string]interface{})
//		results := []interface{}{}
//		results = data["result"].([]interface{})
//		for _, r := range results {
//			result := make(map[string]interface{})
//			result = r.(map[string]interface{})
//			metric := make(map[string]interface{})
//			metric = result["metric"].(map[string]interface{})
//			uri := metric["uri"].(string)
//			// 当uri包含如下字符时，不对该uri进行统计
//			if strings.Contains(uri, "sign") || strings.Contains(uri, "%") || strings.Contains(uri, "=") || strings.Contains(uri, "?") || strings.Contains(uri, "+") || strings.Contains(uri, ":") || strings.Contains(uri, "(") || strings.Contains(uri, ")") || strings.Contains(uri, "\\") || strings.Contains(uri, "\"") || strings.Contains(uri, "'") || strings.Contains(uri, ".") || strings.Contains(uri, " ") || strings.Contains(uri, "|") || strings.Contains(uri, "@") || strings.Contains(uri, ";") || strings.Contains(uri, ":") || strings.Contains(uri, ",") || strings.Contains(uri, "!") || strings.Contains(uri, "<") || strings.Contains(uri, ">") || strings.Contains(uri, "--") || strings.Contains(uri, "//") {
//				continue
//			}
//			values := []interface{}{}
//			values = result["value"].([]interface{})
//			rt := values[1].(string)
//			// 该rt的时间戳
//			index := strings.Index(rt, ".")
//			if index != -1 {
//				// 响应时间向下取整
//				rt = string([]byte(rt)[:index])
//			}
//			// todo 将时间戳转换为时间字符串
//			secTime := time.Unix(timestamp, 0)
//			key := secTime.Format(models.Time_format)
//			rtInt := -1
//			if rt != "NAN" && rt != "" {
//				rtInt, _ = strconv.Atoi(rt)
//			}
//			// todo 将rt和库中的阈值（平均值）进行对比
//			mongo := &models.RtDetailMongo{}
//			mongo, err := mongo.GetByServiceAndUri(serviceCode, uri)
//			if err != nil {
//				return
//			}
//			// 存在阈值时，取阈值；不存在阈值时，取平均值
//			thresholdRtStr := mongo.ThresholdRt
//			if thresholdRtStr == "0" || thresholdRtStr == "" || thresholdRtStr == "-1" {
//				oclock := key[11:]
//				avgRt := mongo.AvgRt
//				avgRtMap := map[string]int{}
//				json.Unmarshal([]byte(avgRt), &avgRtMap)
//				thresholdRtStr = strconv.Itoa(avgRtMap[oclock])
//			}
//			// 没有有效的阈值时，不进行验证
//			if thresholdRtStr != "0" && thresholdRtStr != "" && thresholdRtStr != "-1" {
//				thresholdRt, _ := strconv.Atoi(thresholdRtStr)
//				isQuickIncrease := IsQuickIncrease(rtInt, thresholdRt)
//				if isQuickIncrease && MONITOR_DING_SEND_IS_OPEN {
//					content := fmt.Sprintf("【性能监控-线上巡检Alert】: 接口响应时间大幅度高出近两周平均值，请及时关注！\n【业务线】: 最右\n【网关服务】: %s\n【URI】: %s\n【当前响应时间】: %v\n【历史平均响应时间】: %v\n", serviceCode, uri, rtInt, thresholdRt)
//					//inspection_strategy.DingSend(content)
//					fmt.Printf(content)
//				}
//			}
//			// todo 把rtMap入库
//			RtDetailInDbNow(timestamp, serviceCode, uri, rtInt)
//			// todo 调用缓增函数，判断该数据是否为缓增数据
//			res := []int{}
//			timeStr := time.Unix(timestamp, 0).Format(models.Time_format)
//			timeStr = timeStr[11:]
//			res = append(res, getMouShijianRt(mongo.Last0DayRt, timeStr))
//			res = append(res, getMouShijianRt(mongo.Last1DayRt, timeStr))
//			res = append(res, getMouShijianRt(mongo.Last2DayRt, timeStr))
//			res = append(res, getMouShijianRt(mongo.Last3DayRt, timeStr))
//			res = append(res, getMouShijianRt(mongo.Last4DayRt, timeStr))
//			res = append(res, getMouShijianRt(mongo.Last5DayRt, timeStr))
//			res = append(res, getMouShijianRt(mongo.Last6DayRt, timeStr))
//			// todo 需要确保res为由早到晚的时间顺序
//			IsSlowIncrease := IsSlowIncrease(res)
//			if IsSlowIncrease {
//				// todo 发送钉钉出去
//				if MONITOR_DING_SEND_IS_OPEN {
//					content := fmt.Sprintf("【性能监控-线上巡检Alert】: 接口近期响应时间呈缓增趋势，请关注！\n【业务线】: 最右\n【服务】: %s\n【URI】: %s\n", serviceCode, uri)
//					//inspection_strategy.DingSend(content)
//					fmt.Printf(content)
//				}
//			}
//		}
//	}
//	return
//}

/**
每小时跑一次
*/
func HalfHourExcute(serviceCode string, timestamp int64) {
	// 该服务下的所有接口
	//url := ZyPormtheusQueryUrl + serviceCode + "_http_latency_quantile%7Bquantile%3D%22p99%22%7D"
	startStr := fmt.Sprintf("%v", timestamp-60*60)
	endStr := fmt.Sprintf("%v", timestamp)
	stepStr := fmt.Sprintf("%v", 60*10)
	url := ZyPormtheusQueryRangeUrl + serviceCode + "_http_latency_quantile%7Bquantile%3D%22p99%22%7D&start=" + startStr + "&end=" + endStr + "&step=" + stepStr

	client := &http.Client{Timeout: 5 * time.Second}
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
		fmt.Printf("请求结果不是success")
		fmt.Printf("请求结果为: %v", res)
	} else {
		resStr, _ := json.Marshal(res)
		fmt.Printf("打印请求结果为: %s", string(resStr))
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
			rtArr := []int{}
			for _, val := range values {
				jsonByte, _ := json.Marshal(val)
				jsonStr := string(jsonByte)
				rt := gojson.Json(jsonStr).Getindex(2).Tostring()
				// 响应时间向下取整
				index := strings.Index(rt, ".")
				if index != -1 {
					rt = string([]byte(rt)[:index])
				}
				rtInt := -1
				if rt != "NAN" && rt != "" {
					rtInt, _ = strconv.Atoi(rt)
				}
				rtArr = append(rtArr, rtInt)
			}
			count := 0
			total := 0
			for _, val := range rtArr {
				if val != 0 && val != -1 {
					count++
					total += val
				}
			}
			rtInt := -1
			if count != 0 {
				rtInt = total / count
			} else {
				// 当没有取到有效值时，进行return
				return
			}
			// todo 将rt和库中的阈值（平均值）进行对比
			mongo := &models.RtDetailMongo{}
			mongo, err := mongo.GetByServiceAndUri(serviceCode, uri)
			if err != nil {
				return
			}
			// 存在阈值时，取阈值；不存在阈值时，取平均值
			thresholdRtStr := mongo.ThresholdRt
			if thresholdRtStr == "0" || thresholdRtStr == "" || thresholdRtStr == "-1" {
				oclock := time.Unix(timestamp, 0).Format(models.Time_format)[11:]
				//oclock := key[11:]
				avgRt := mongo.AvgRt
				avgRtMap := map[string]int{}
				json.Unmarshal([]byte(avgRt), &avgRtMap)
				thresholdRtStr = strconv.Itoa(avgRtMap[oclock])
			}
			// 没有有效的阈值时，不进行验证
			if thresholdRtStr != "0" && thresholdRtStr != "" && thresholdRtStr != "-1" {
				thresholdRt, _ := strconv.Atoi(thresholdRtStr)
				isQuickIncrease := IsQuickIncrease(rtInt, thresholdRt)
				if isQuickIncrease && MONITOR_DING_SEND_IS_OPEN {
					content := fmt.Sprintf("【性能监控-线上巡检Alert】: 接口响应时间大幅度高出近两周平均值，请及时关注！\n【业务线】: 最右\n【网关服务】: %s\n【URI】: %s\n【当前响应时间】: %v\n【历史平均响应时间】: %v\n", serviceCode, uri, rtInt, thresholdRt)
					DingSend(content)
					fmt.Printf(content)
				}
			}
			// todo 把rtMap入库
			RtDetailInDbNow(timestamp, serviceCode, uri, rtInt)
			// todo 调用缓增函数，判断该数据是否为缓增数据
			res := []int{}
			timeStr := time.Unix(timestamp, 0).Format(models.Time_format)
			timeStr = timeStr[11:]
			res = append(res, getMouShijianRt(mongo.Last0DayRt, timeStr))
			res = append(res, getMouShijianRt(mongo.Last1DayRt, timeStr))
			res = append(res, getMouShijianRt(mongo.Last2DayRt, timeStr))
			res = append(res, getMouShijianRt(mongo.Last3DayRt, timeStr))
			res = append(res, getMouShijianRt(mongo.Last4DayRt, timeStr))
			res = append(res, getMouShijianRt(mongo.Last5DayRt, timeStr))
			res = append(res, getMouShijianRt(mongo.Last6DayRt, timeStr))
			// todo 需要确保res为由早到晚的时间顺序
			IsSlowIncrease := IsSlowIncrease(res)
			if IsSlowIncrease {
				// todo 发送钉钉出去
				if MONITOR_DING_SEND_IS_OPEN {
					content := fmt.Sprintf("【性能监控-线上巡检Alert】: 接口近期响应时间呈缓增趋势，请关注！\n【业务线】: 最右\n【服务】: %s\n【URI】: %s\n", serviceCode, uri)
					DingSend(content)
					fmt.Printf(content)
				}
			}
		}
	}
	return
}

/**
当取不到相关数据时，rt=-1
*/
func getMouShijianRt(rtMapStr string, shijiandian string) (rt int) {
	rt = -1
	if rtMapStr != "" {
		rtMap := map[string]int{}
		json.Unmarshal([]byte(rtMapStr), &rtMap)
		value, ok := rtMap[shijiandian]
		if ok {
			rt = value
		}
	}
	return rt
}

func RtDetailInDbNow(timestamp int64, serviceCode string, uri string, rt int) {
	mongo := &models.RtDetailMongo{}
	mongo, err := mongo.GetByServiceAndUri(serviceCode, uri)
	if err != nil {
		return
	}
	secTime := time.Unix(timestamp, 0)
	key := secTime.Format(models.Time_format)
	key = key[11:]

	rtDetailMap := map[string]int{}
	rtDetailStr := mongo.Last0DayRt
	if rtDetailStr == "" {
		rtDetailMap[key] = rt
	} else {
		json.Unmarshal([]byte(rtDetailStr), rtDetailMap)
		rtDetailMap[key] = rt
	}
	rtDetailBytes, err := json.Marshal(rtDetailMap)
	mongo.Last0DayRt = string(rtDetailBytes)
	mongo.UpdateById(mongo.Id, *mongo)
}

/**
当周7天里，有4天呈现rt增长趋势时，理解为缓增接口
*/
//func IsSlowIncrease(datas []int) (isSlowIncrease bool) {
//	//res := []bool{}
//	huanzengtianshu := 0
//	for i := 0; i < len(datas); i++ {
//		isSlowIncrease = true
//		for j := i+1; j > 0; j-- {
//			if datas[j] < datas[j-1] {
//				isSlowIncrease = false
//				break
//			}
//		}
//		//res = append(res, isSlowIncrease)
//		if isSlowIncrease {
//			huanzengtianshu++
//		}
//	}
//	if huanzengtianshu > 3 {
//		// todo 该接口的rt时间为缓增趋势 ，发送钉钉进行通知
//		return true
//	}
//	return false
//}

/**
暴增告警策略
*/
func IsQuickIncrease(nowRt int, avgRt int) (isSlowIncrease bool) {
	isQuickIncrease := false
	if avgRt <= 100 {
		// todo-done 平均值在100-时,当rt涨幅高于50ms时，报警
		if float64(nowRt) > float64(avgRt)*1.5 && nowRt > avgRt+50 {
			isQuickIncrease = true
		}
	} else if avgRt > 100 && avgRt <= 200 {
		// todo-done 平均值在100+时,当rt涨幅高于50%时，报警
		if float64(nowRt) > float64(avgRt)*1.5 {
			isQuickIncrease = true
		}
	} else if avgRt > 200 && avgRt <= 300 {
		if float64(nowRt) > float64(avgRt)*1.6 {
			isQuickIncrease = true
		}
	} else if avgRt > 300 {
		if float64(nowRt) > float64(avgRt)*1.7 {
			isQuickIncrease = true
		}
	}
	return isQuickIncrease
}

/**
缓增告警策略
*/
func IsSlowIncrease(datas []int) (isSlowIncrease bool) {
	for i := 0; i < len(datas)-1; i++ {
		if datas[i+1] > datas[i] {
			continue
		} else {
			return false
		}
	}
	return true
}

func DingSend(content string) {
	var dingToken = []string{inspection_strategy.XIAO_NENG_QUN_TOKEN}
	cli := dingtalk.InitDingTalk(dingToken, "")
	cli.SendTextMessage(content)
}
