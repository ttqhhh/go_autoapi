package main

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"github.com/prometheus/common/log"
	"net/http"
	"strconv"
)

const ZY_grafana_login_url = "http://grafana.ixiaochuan.cn/login"
const HW_grafana_login_url = "http://dashboard.icocofun.net/login"
const HWUS_grafana_login_url = "http://grafanaus.icocofun.net/login"
const AD_gateway_path_url = "http://grafana.ixiaochuan.cn/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_gateway_ad_http_latency_count%5B1m%5D))by(uri)&start=1646064000&end=1648742400&step=7200"
const ZD_grafana_login_url = "http://grafana.mehiya.com/login"

type JsonData struct {
	Status string `json:"status"`
	Data   Data   `json:"data"`
}

type Data struct {
	ResultType string   `json:"resultType"`
	Result     []Result `json:"result"`
}

type Result struct {
	Meturc Metruc          `json:"metric"`
	Values [][]interface{} `json:"values"`
}

type Metruc struct {
	Uri string `json:"uri"`
}

//type Values struct {
//	Values [][]int`json:"values"`
//}

func getLogin() *http.Cookie {
	req := httplib.Post(ZY_grafana_login_url)
	req.Param("user", "wangzhen01")
	req.Param("password", "Iepohg5go4iawoo")
	req.Param("email", "")
	resp, err := req.Response()
	if err != nil {
	}
	cookies := resp.Cookies()
	for _, cookie := range cookies {
		return cookie
	}
	return nil
}

func getLoginHaiWai() *http.Cookie {
	req := httplib.Post(HW_grafana_login_url)
	req.Param("user", "wangzhen01")
	req.Param("password", "Iepohg5go4iawoo")
	resp, err := req.Response()
	if err != nil {
	}
	cookies := resp.Cookies()
	for _, cookie := range cookies {
		return cookie
	}
	return nil
}

func getLoginHaiWaiUS() *http.Cookie {
	req := httplib.Post(HWUS_grafana_login_url)
	req.Param("user", "wangzhen01")
	req.Param("password", "Iepohg5go4iawoo")
	resp, err := req.Response()
	if err != nil {
	}
	cookies := resp.Cookies()
	for _, cookie := range cookies {
		return cookie
	}
	return nil
}

func getLoginZhongDong() *http.Cookie {
	req := httplib.Post(ZD_grafana_login_url)
	req.Param("user", "chengxiaojing")
	req.Param("password", "ls51HGb8y0MA")
	resp, err := req.Response()
	if err != nil {
	}
	cookies := resp.Cookies()
	for _, cookie := range cookies {
		return cookie
	}
	return nil
}

func doReq(url string, cookie *http.Cookie) string {
	req := httplib.Get(url)
	req.SetCookie(cookie)
	str, err := req.String()
	if err != nil {
		//
	}
	return str
}

func main() {
	cookie := getLogin()
	cookiehaiwai := getLoginHaiWai()
	cookiehaiwaius := getLoginHaiWaiUS()
	//商业化
	shangyehua := getSyhAllApiCount(cookie)
	////最右
	zuiyou := getZyAllApiCount(cookie)
	pipi := getPPAllApiCount(cookie)
	haiwai := gethaiwaiAllApiCount(cookiehaiwai)
	haiwaius := gethaiwaiUSAllApiCount(cookiehaiwaius)
	fmt.Print(shangyehua)
	fmt.Print("\n" + "商业化")
	fmt.Print(zuiyou)
	fmt.Print("\n" + "最右")
	fmt.Print(pipi)
	fmt.Print("\n" + "皮皮")
	fmt.Print(haiwai)
	fmt.Print("\n" + "海外")
	fmt.Print(haiwaius)
	fmt.Print("\n" + "海外US")
	cookieZD := getLoginZhongDong()
	zhcount := getzhongdongAllApiCount(cookieZD)
	fmt.Print(zhcount)
	fmt.Print("\n" + "中东")

}

func getSyhAllApiCount(cookie *http.Cookie) float64 {
	respData := doReq(AD_gateway_path_url, cookie)
	toJson := JsonData{}
	err := json.Unmarshal([]byte(respData), &toJson)
	if err != nil {

	}
	count := 0
	for _, one := range toJson.Data.Result {
		for _, ones := range one.Values {
			str_one := fmt.Sprintf("%v", ones[1])
			float_one, err := strconv.ParseFloat(str_one, 64)
			if err != nil {
				log.Error("类型转换出错", err)
			}
			if float_one > 1 {
				count++
				break
			}
		}
	}
	return float64(count)
}

func getZyAllApiCount(cookie *http.Cookie) float64 {
	count := 0
	ZYlist := []string{
		"gateway_rec", "gateway_topic", " gateway_post", "gateway_rev ", " gateway_acnt", "gateway_cfg", "gateway_danmaku",
		"gateway_misc", "snssrv_gateway_native", "snssrv_gateway-nearby", "zy_gateway_teamchat", "chatsrv_gateway",
	}
	for _, i := range ZYlist {
		ZY_gateway_path_url := "http://grafana.ixiaochuan.cn/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_" + i + "_http_latency_count%5B1m%5D))by(uri)&start=1646064000&end=1648742400&step=7200"
		respData := doReq(ZY_gateway_path_url, cookie)
		toJson := JsonData{}
		err := json.Unmarshal([]byte(respData), &toJson)
		if err != nil {

		}
		for _, one := range toJson.Data.Result {
			for _, ones := range one.Values {
				str_one := fmt.Sprintf("%v", ones[1])
				float_one, err := strconv.ParseFloat(str_one, 64)
				if err != nil {
					log.Error("类型转换出错", err)
				}
				if float_one > 1 {
					count++
					break
				}
			}
		}

	}

	return float64(count)
}

func getPPAllApiCount(cookie *http.Cookie) float64 {
	count := 0
	//PPlist :=[]string { //使用循环遍历拼接的方法 会出现bad request 400
	//	"pp-gateway-acnt","pp-gateway-internal"," pp-gateway-misc","pp-gateway-point "," pp-gateway-post","pp-gateway-rec ","pp-gateway-review",
	//	"pp-gateway-topic","pp-gateway-town",
	//}
	pipiURLlsit := []string{
		"http://grafana.ixiaochuan.cn/api/datasources/proxy/5/api/v1/query_range?query=sum%20by(uri)%20(irate(xcmetrics_httpsrv_qps%7Bjob%3D%22pp-gateway-misc%22%7D%5B1m%5D))&start=1646064000&end=1648742400&step=7200",
		"http://grafana.ixiaochuan.cn/api/datasources/proxy/5/api/v1/query_range?query=sum%20by(uri)%20(irate(xcmetrics_httpsrv_qps%7Bjob%3D%22pp-gateway-acnt%22%7D%5B1m%5D))&start=1646064000&end=1648742400&step=7200",
		"http://grafana.ixiaochuan.cn/api/datasources/proxy/5/api/v1/query_range?query=sum%20by(uri)%20(irate(xcmetrics_httpsrv_qps%7Bjob%3D%22pp-gateway-internal%22%7D%5B1m%5D))&start=1646064000&end=1648742400&step=7200",
		"http://grafana.ixiaochuan.cn/api/datasources/proxy/5/api/v1/query_range?query=sum%20by(uri)%20(irate(xcmetrics_httpsrv_qps%7Bjob%3D%22pp-gateway-point%22%7D%5B1m%5D))&start=1646064000&end=1648742400&step=7200",
		"http://grafana.ixiaochuan.cn/api/datasources/proxy/5/api/v1/query_range?query=sum%20by(uri)%20(irate(xcmetrics_httpsrv_qps%7Bjob%3D%22pp-gateway-post%22%7D%5B1m%5D))&start=1646064000&end=1648742400&step=7200",
		"http://grafana.ixiaochuan.cn/api/datasources/proxy/5/api/v1/query_range?query=sum%20by(uri)%20(irate(xcmetrics_httpsrv_qps%7Bjob%3D%22pp-gateway-rec%22%7D%5B1m%5D))&start=1646064000&end=1648742400&step=7200",
		"http://grafana.ixiaochuan.cn/api/datasources/proxy/5/api/v1/query_range?query=sum%20by(uri)%20(irate(xcmetrics_httpsrv_qps%7Bjob%3D%22pp-gateway-review%22%7D%5B1m%5D))&start=1646064000&end=1648742400&step=7200",
		"http://grafana.ixiaochuan.cn/api/datasources/proxy/5/api/v1/query_range?query=sum%20by(uri)%20(irate(xcmetrics_httpsrv_qps%7Bjob%3D%22pp-gateway-topic%22%7D%5B1m%5D))&start=1646064000&end=1648742400&step=7200",
		"http://grafana.ixiaochuan.cn/api/datasources/proxy/5/api/v1/query_range?query=sum%20by(uri)%20(irate(xcmetrics_httpsrv_qps%7Bjob%3D%22pp-gateway-town%22%7D%5B1m%5D))&start=1646064000&end=1648742400&step=7200",
	}
	for _, i := range pipiURLlsit {
		//print(i)
		respData := doReq(i, cookie)
		toJson := JsonData{}
		err := json.Unmarshal([]byte(respData), &toJson)
		if err != nil {

		}
		for _, one := range toJson.Data.Result {
			for _, ones := range one.Values {
				str_one := fmt.Sprintf("%v", ones[1])
				float_one, err := strconv.ParseFloat(str_one, 64)
				if err != nil {
					log.Error("类型转换出错", err)
				}
				if float_one > 1 {
					count++
					break
				}
			}
		}

	}

	return float64(count)
}

func gethaiwaiAllApiCount(cookie *http.Cookie) float64 {
	count := 0
	haiwaiURLlsit := []string{
		"http://dashboard.icocofun.net/api/datasources/proxy/25/api/v1/query_range?query=sum(rate(xms_omg_gateway_http_latency_count%5B1m%5D))by(uri)&start=1644288000&end=1646880000&step=1200&timeout=300s",
		"http://dashboard.icocofun.net/api/datasources/proxy/25/api/v1/query_range?query=sum(rate(xms_chatsrv_gateway_http_latency_count%5B1m%5D))by(uri)&start=1644288000&end=1646880000&step=1200&timeout=300s",
		"http://dashboard.icocofun.net/api/datasources/proxy/25/api/v1/query_range?query=sum(rate(xms_gateway_ad_http_latency_count%5B1m%5D))by(uri)&start=1644288000&end=1646880000&step=1200&timeout=300s",
		"http://dashboard.icocofun.net/api/datasources/proxy/25/api/v1/query_range?query=sum(rate(xms_omg_gateway_acnt_http_latency_count%5B1m%5D))by(uri)&start=1644288000&end=1646880000&step=1200&timeout=300s",
		"http://dashboard.icocofun.net/api/datasources/proxy/25/api/v1/query_range?query=sum(rate(xms_omg_gateway_index_http_latency_count%5B1m%5D))by(uri)&start=1644288000&end=1646880000&step=1200&timeout=300s",
		"http://dashboard.icocofun.net/api/datasources/proxy/25/api/v1/query_range?query=sum(rate(xms_omg_gateway_post_http_latency_count%5B1m%5D))by(uri)&start=1644288000&end=1646880000&step=1200&timeout=300s",
		"http://dashboard.icocofun.net/api/datasources/proxy/25/api/v1/query_range?query=sum(rate(xms_omg_gateway_review_http_latency_count%5B1m%5D))by(uri)&start=1644288000&end=1646880000&step=1200&timeout=300s",
		"http://dashboard.icocofun.net/api/datasources/proxy/25/api/v1/query_range?query=sum(rate(xms_omg_gateway_topic_http_latency_count%5B1m%5D))by(uri)&start=1644288000&end=1646880000&step=1200&timeout=300s",
	}
	for _, i := range haiwaiURLlsit {
		//print(i)
		respData := doReq(i, cookie)
		toJson := JsonData{}
		err := json.Unmarshal([]byte(respData), &toJson)
		if err != nil {

		}
		for _, one := range toJson.Data.Result {
			for _, ones := range one.Values {
				str_one := fmt.Sprintf("%v", ones[1])
				float_one, err := strconv.ParseFloat(str_one, 64)
				if err != nil {
					log.Error("类型转换出错", err)
				}
				if float_one > 1 {
					count++
					break
				}
			}
		}

	}

	return float64(count)
}

func gethaiwaiUSAllApiCount(cookie *http.Cookie) float64 {
	count := 0
	haiwaiUSURLlsit := []string{
		"http://grafanaus.icocofun.net/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_maga_gateway_http_latency_count%7Buri!%3D%22%2Fhealthcheck%22%7D%5B1m%5D))by(uri)&start=1646496000&end=1646668680&step=120",
		"http://grafanaus.icocofun.net/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_maga_chatsrv_gateway_http_latency_count%7Buri!%3D%22%2Fhealthcheck%22%7D%5B1m%5D))by(uri)&start=1646496000&end=1646668680&step=120",
		"http://grafanaus.icocofun.net/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_maga_gateway_account_http_latency_count%7Buri!%3D%22%2Fhealthcheck%22%7D%5B1m%5D))by(uri)&start=1646496000&end=1646668680&step=120",
	}
	for _, i := range haiwaiUSURLlsit {
		//print(i)
		respData := doReq(i, cookie)
		toJson := JsonData{}
		err := json.Unmarshal([]byte(respData), &toJson)
		if err != nil {

		}
		for _, one := range toJson.Data.Result {
			for _, ones := range one.Values {
				str_one := fmt.Sprintf("%v", ones[1])
				float_one, err := strconv.ParseFloat(str_one, 64)
				if err != nil {
					log.Error("类型转换出错", err)
				}
				if float_one > 1 {
					count++
					break
				}
			}
		}

	}

	return float64(count)
}

func getzhongdongAllApiCount(cookie *http.Cookie) float64 {
	count := 0
	ZhongDongURLlsit := []string{
		"http://grafana.mehiya.com/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_me_live_account_http_latency_count%5B1m%5D))by(uri)&start=1644303600&end=1646895600&step=1800",
		"http://grafana.mehiya.com/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_me_live_privilege_http_latency_count%5B1m%5D))by(uri)&start=1644303600&end=1646895600&step=1800",
		"http://grafana.mehiya.com/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_me_live_user_http_latency_count%5B1m%5D))by(uri)&start=1644303600&end=1646895600&step=1800",
		"http://grafana.mehiya.com/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_me_live_chat-gateway_http_latency_count%5B1m%5D))by(uri)&start=1644303600&end=1646895600&step=1800",
		"http://grafana.mehiya.com/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_me_live_turntable_http_latency_count%5B1m%5D))by(uri)&start=1644303600&end=1646895600&step=1800",
		"http://grafana.mehiya.com/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_me_live_family_http_latency_count%5B1m%5D))by(uri)&start=1644303600&end=1646895600&step=1800",
		"http://grafana.mehiya.com/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_me_live_gamestore_http_latency_count%5B1m%5D))by(uri)&start=1644303600&end=1646895600&step=1800",
		"http://grafana.mehiya.com/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_me_live_pk_http_latency_count%5B1m%5D))by(uri)&start=1644303600&end=1646895600&step=1800",
		"http://grafana.mehiya.com/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_me_live_misc_http_latency_count%5B1m%5D))by(uri)&start=1644303600&end=1646895600&step=1800",
		"http://grafana.mehiya.com/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_me_live_trade_gateway_http_latency_count%5B1m%5D))by(uri)&start=1644303600&end=1646895600&step=1800",
		"http://grafana.mehiya.com/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_me_live_relation_http_latency_count%5B1m%5D))by(uri)&start=1644303600&end=1646895600&step=1800",
		"http://grafana.mehiya.com/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_me_live_gateway_http_latency_count%5B1m%5D))by(uri)&start=1644303600&end=1646895600&step=1800",
	}
	for _, i := range ZhongDongURLlsit {
		//print(i)
		respData := doReq(i, cookie)
		toJson := JsonData{}
		err := json.Unmarshal([]byte(respData), &toJson)
		if err != nil {

		}
		for _, one := range toJson.Data.Result {
			for _, ones := range one.Values {
				str_one := fmt.Sprintf("%v", ones[1])
				float_one, err := strconv.ParseFloat(str_one, 64)
				if err != nil {
					log.Error("类型转换出错", err)
				}
				if float_one > 1 {
					count++
					break
				}
			}
		}

	}

	return float64(count)
}
