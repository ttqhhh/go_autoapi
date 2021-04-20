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
	logs.Error("requests err:", resp)
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
	if statusCode != 200 {
		fmt.Println("请求返回状态不是200，请求失败")
		return
	}
	//ret := verify["code"]["code"]
	//fmt.Println("ret is ", ret)
	//if ret != jmap["code"] {
	//	fmt.Println("接口返回状态码不正确", jmap["code"], verify["ret"]["code"])
	//}
	for k, v := range verify {
		//fmt.Println(k, v, verify, reflect.TypeOf(verify))
		for subK, subV := range v {
			data := jmap["data"].(map[string]interface{})
			fmt.Println("sub is ", subK, reflect.TypeOf(subV), reflect.TypeOf(data[k]))
			if subK == "eq" {
				logs.Error("eq here")
				if jmap["code"] != subV {
					logs.Error("not equal", jmap["code"], subV)
				}
			} else if subK == "lt" {
				if subV.(int64) >= data[k].(int64) {
					logs.Error("not lt", data[k], subV)
				}
			} else if subK == "gt" {
				if subV.(int64) <= data[k].(int64) {
					logs.Error("not gt", data[k], subV)
				}
			} else if subK == "lte" {
				if subV.(int64) > data[k].(int64) {
					logs.Error("not lte", data[k], subV)
				}
			} else if subK == "gte" {
				if subV.(int64) < data[k].(int64) {
					logs.Error("not gte", data[k], subV)
				}
			} else if subK == "need" {
				if data[k] == nil {
					logs.Error("not need", data[k], subV)
				}
			} else if subK == "in" {
				b := strings.ContainsAny(data[k].(string), subV.(string))
				if b {
					logs.Error("not in", data[k], subV)
				}
			} else {
				logs.Error("do not support")
			}
		}
	}
}
