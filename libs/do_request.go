package libs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"go_autoapi/db_proxy"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

func init() {
	fmt.Println("please init me")
	_ = db_proxy.InitClient()
}

//模拟请求方法
func HttpPost(postUrl string, headers map[string]string, jsonMap map[string]interface{}) (int, string, string) {
	client := &http.Client{}
	//转换成postBody
	bytesData, err := json.Marshal(jsonMap)
	if err != nil {
		fmt.Println(err.Error())
		return 0, "", ""
	}
	postBody := bytes.NewReader(bytesData)
	client = &http.Client{}
	//post请求
	req, _ := http.NewRequest("POST", postUrl, postBody)
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	resp, _ := client.Do(req)
	//logs.Error("requests err:", resp)
	//返回内容
	body, err := ioutil.ReadAll(resp.Body)
	logs.Error("requests err:", err)
	//解析返回的cookie
	var cookieStr string
	cookies := resp.Cookies()
	if cookies != nil {
		for _, c := range cookies {
			cookieStr += c.Name + "=" + c.Value + ";"
		}
	}
	return resp.StatusCode, string(body), cookieStr
}

func DoRequest(url string, uuid string, data map[string]interface{}, verify map[string]map[string]interface{}) {
	//密码
	r := db_proxy.GetRedisObject()
	statusCode, body, _ := HttpPost(url, nil, data)
	//body jsonStr转map
	var jmap map[string]interface{}
	if err := json.Unmarshal([]byte(body), &jmap); err != nil {
		fmt.Println("解析失败", err)
		return
	}
	// 判断某个字段的类型
	//fmt.Println("type:", reflect.TypeOf(jmap["code"]))
	//判断登录是否成功
	doVerify(statusCode, body, verify)
	r.Incr(uuid)

}

// 增加验证函数，比较响应和需要验证的内容
func doVerify(statusCode int, response string, verify map[string]map[string]interface{}) {
	var jmap map[string]interface{}
	if err := json.Unmarshal([]byte(response), &jmap); err != nil {
		fmt.Println("解析失败", err)
		return
	}
	data := jmap["data"].(map[string]interface{})
	if statusCode != 200 {
		fmt.Println("请求返回状态不是200，请求失败")
		return
	}
	for k, v := range verify {
		//fmt.Println(k, v, verify, reflect.TypeOf(verify))
		logs.Error("k,v is ", k, v, reflect.TypeOf(k))
		if k != "code" && data[k] == nil {
			logs.Error("the verify key is not exist in the response", k)
			return
		}
		for subK, subV := range v {
			if subK == "eq" {
				if k == "code" {
					logs.Error("code here ,OK no problem", jmap["code"], subV)
					if jmap["code"] != subV {
						logs.Error("code not equal", jmap["code"], subV, subK)
						return
					}
				} else if data[k] != subV {
					logs.Error("not equal", data[k], subV)
					return
				}
			} else if subK == "lt" {
				if subV.(float64) >= data[k].(float64) {
					logs.Error("not lt", data[k], subV)
					return
				}
			} else if subK == "gt" {
				if subV.(int64) <= data[k].(int64) {
					logs.Error("not gt", data[k], subV)
					return
				}
			} else if subK == "lte" {
				if subV.(int64) > data[k].(int64) {
					logs.Error("not lte", data[k], subV)
					return
				}
			} else if subK == "gte" {
				if subV.(int64) < data[k].(int64) {
					logs.Error("not gte", data[k], subV)
					return
				}
			} else if subK == "need" {
				logs.Error("need string", data[k], subV)
				if data[k] == nil {
					logs.Error("not need", data[k], subV)
					return
				}
			} else if subK == "in" {
				b := strings.ContainsAny(data[k].(string), subV.(string))
				if b == false {
					logs.Error("not in", data[k], subV)
					return
				}
			} else {
				logs.Error("do not support")
				return
			}
		}
	}
}
