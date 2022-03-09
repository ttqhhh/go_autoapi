package main

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"net/http"
)

const ZY_grafana_login_url = "http://grafana.ixiaochuan.cn/login"
const AD_gateway_path_url = "http://grafana.ixiaochuan.cn/api/datasources/proxy/1/api/v1/query_range?query=sum(rate(xms_gateway_ad_http_latency_count%5B1m%5D))by(uri)&start=1646064000&end=1648742400&step=7200"

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
	//for i := range list{  //遍历所有getway
	//
	//}
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
				break
			}
		}
	}
	fmt.Print("活跃接口数为：")
	fmt.Print(count)
	fmt.Print("\n")

}
