package main

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"github.com/prometheus/common/log"
	"github.com/xuri/excelize/v2"
	"go_autoapi/constants"
	"go_autoapi/models"
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
	if false { //重新统计并入库 但不会晴空之前的数据
		cookie := getLogin()
		cookiehaiwai := getLoginHaiWai()
		cookiehaiwaius := getLoginHaiWaiUS()
		cookieZD := getLoginZhongDong()

		getZyAllApiCount(cookie)
		getPPAllApiCount(cookie)
		gethaiwaiAllApiCount(cookiehaiwai)
		getzhongdongAllApiCount(cookieZD)
		getSyhAllApiCount(cookie)
		gethaiwaiUSAllApiCount(cookiehaiwaius)

	}

	if true { //将数据输出exal
		var numbers = []int64{
			constants.ZuiyYou, constants.PiPi, constants.HaiWai, constants.ZhongDong, constants.ShangYeHua, constants.HaiWaiUS,
		}
		for _, business_code := range numbers {
			mongo := models.AllActiveApiMongo{}
			apiList, err := mongo.QueryByBusiness(business_code)
			if err != nil {
				log.Error("根据业务线获取全部接口出错")
			}
			line := 1
			f := excelize.NewFile() // 设置单元格的值
			for _, oneApi := range apiList {
				//开始创建excal 写入数据
				// 这里设置表头
				f.SetCellValue("Sheet1", "A1", "id")
				f.SetCellValue("Sheet1", "B1", "业务线")
				f.SetCellValue("Sheet1", "C1", "接口名")
				f.SetCellValue("Sheet1", "D1", "是否使用")

				line++
				f.SetCellValue("Sheet1", fmt.Sprintf("A%d", line), oneApi.Id)
				f.SetCellValue("Sheet1", fmt.Sprintf("B%d", line), oneApi.BusinessName)
				f.SetCellValue("Sheet1", fmt.Sprintf("C%d", line), oneApi.ApiName)
				f.SetCellValue("Sheet1", fmt.Sprintf("D%d", line), oneApi.Use)

				// 保存文件

			}
			file := strconv.FormatInt(business_code, 10)
			filename := file + ".xlsx"
			if err := f.SaveAs("/Users/tangtianqing/Desktop/active_api/" + filename); err != nil {
				fmt.Println(err)
			}

		}

	}

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
			if ones[1] != "0" {
				count++
				acm := models.AllActiveApiMongo{}
				acm.BusinessName = "商业化"
				acm.BusinessCode = constants.ShangYeHua
				acm.ApiName = one.Meturc.Uri
				acm.Insert(acm)
				break
			}
		}
	}
	return float64(count)
}

func getZyAllApiCount(cookie *http.Cookie) float64 {
	count := 0
	zuiyouURLlsit := []string{
		"http://grafana.ixiaochuan.cn/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_gateway_rec_http_latency_count%5B1m%5D))by(uri)&start=1646064000&end=1648742400&step=7200",
		"http://grafana.ixiaochuan.cn/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_gateway_topic_http_latency_count%5B1m%5D))by(uri)&start=1646064000&end=1648742400&step=7200",
		"http://grafana.ixiaochuan.cn/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_gateway_post_http_latency_count%5B1m%5D))by(uri)&start=1646064000&end=1648742400&step=7200",
		"http://grafana.ixiaochuan.cn/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_gateway_rev_http_latency_count%5B1m%5D))by(uri)&start=1646064000&end=1648742400&step=7200",
		"http://grafana.ixiaochuan.cn/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_gateway_acnt_http_latency_count%5B1m%5D))by(uri)&start=1646064000&end=1648742400&step=7200",
		"http://grafana.ixiaochuan.cn/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_gateway_cfg_http_latency_count%5B1m%5D))by(uri)&start=1646064000&end=1648742400&step=7200",
		"http://grafana.ixiaochuan.cn/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_gateway_danmaku_http_latency_count%5B1m%5D))by(uri)&start=1646064000&end=1648742400&step=7200",
		"http://grafana.ixiaochuan.cn/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_gateway_misc_http_latency_count%5B1m%5D))by(uri)&start=1646064000&end=1648742400&step=7200",
		"http://grafana.ixiaochuan.cn/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_snssrv_gateway_native_http_latency_count%5B1m%5D))by(uri)&start=1646064000&end=1648742400&step=7200",
		"http://grafana.ixiaochuan.cn/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_snssrv_gateway-nearby_http_latency_count%5B1m%5D))by(uri)&start=1646064000&end=1648742400&step=7200",
		"http://grafana.ixiaochuan.cn/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_zy_gateway_teamchat_http_latency_count%5B1m%5D))by(uri)&start=1646064000&end=1648742400&step=7200",
		"http://grafana.ixiaochuan.cn/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_chatsrv_gateway_http_latency_count%5B1m%5D))by(uri)&start=1646064000&end=1648742400&step=7200",
	}
	for _, i := range zuiyouURLlsit {
		respData := doReq(i, cookie)
		toJson := JsonData{}
		err := json.Unmarshal([]byte(respData), &toJson)
		if err != nil {
			log.Error("转换类型出错", err)
		}
		for _, one := range toJson.Data.Result {
			for _, ones := range one.Values {
				if ones[1] != "0" {
					count++
					acm := models.AllActiveApiMongo{}
					acm.BusinessName = "最右"
					acm.BusinessCode = constants.ZuiyYou
					acm.ApiName = one.Meturc.Uri
					acm.Insert(acm)
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
				if ones[1] != "0" {
					count++
					acm := models.AllActiveApiMongo{}
					acm.BusinessName = "皮皮"
					acm.BusinessCode = constants.PiPi
					acm.ApiName = one.Meturc.Uri
					acm.Insert(acm)
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
				if ones[1] != "0" {
					count++
					acm := models.AllActiveApiMongo{}
					acm.BusinessName = "海外"
					acm.BusinessCode = constants.HaiWai
					acm.ApiName = one.Meturc.Uri
					acm.Insert(acm)
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
				if ones[1] != "0" {
					count++
					acm := models.AllActiveApiMongo{}
					acm.BusinessName = "海外US"
					acm.BusinessCode = constants.HaiWaiUS
					acm.ApiName = one.Meturc.Uri
					acm.Insert(acm)
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
		"http://grafana.mehiya.com/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_me_live_chat-gateway_http_latency_count%5B1m%5D))by(uri)&start=1644303600&end=1646895600&step=1800",
		"http://grafana.mehiya.com/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_me_live_gamestore_gateway_http_latency_count%5B1m%5D))by(uri)&start=1644303600&end=1646895600&step=1800",
		"http://grafana.mehiya.com/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_me_live_trade_gateway_http_latency_count%5B1m%5D))by(uri)&start=1644303600&end=1646895600&step=1800",
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
				if ones[1] != "0" {
					count++
					acm := models.AllActiveApiMongo{}
					acm.BusinessName = "中东"
					acm.BusinessCode = constants.ZhongDong
					acm.ApiName = one.Meturc.Uri
					acm.Insert(acm)
					break
				}
			}
		}

	}

	return float64(count)
}
