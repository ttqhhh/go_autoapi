package libs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/bitly/go-simplejson"
	jsonpath "github.com/spyzhov/ajson"
	"go_autoapi/db_proxy"
	"go_autoapi/models"
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

func DoRequest(url string, uuid string, data map[string]interface{}, verify map[string]map[string]interface{}, caseId int64) {
	//密码
	r := db_proxy.GetRedisObject()
	statusCode, body, _ := HttpPost(url, nil, data)
	//body jsonStr转map
	var jmap map[string]interface{}
	if err := json.Unmarshal([]byte(body), &jmap); err != nil {
		fmt.Println("解析失败", err)
		return
	}
	// 此处采用go-simplejson来做个示例，用于以后扩展检查使用
	js, err := simplejson.NewJson([]byte(body))
	if err != nil {
		return
	}
	email, err := js.Get("data").Get("email").String()
	fmt.Println(js.Get("code"), email)

	// 判断某个字段的类型
	//fmt.Println("type:", reflect.TypeOf(jmap["code"]))
	//判断登录是否成功
	doVerifyV2(statusCode, uuid, body, verify, caseId)
	r.Incr(uuid)

}

// 采用jsonpath 对结果进行验证
func doVerifyV2(statusCode int, uuid string, response string, verify map[string]map[string]interface{}, caseId int64) {

	if statusCode != 200 {
		logs.Error("请求返回状态不是200，请求失败")
		err := models.InsertResult(uuid, caseId, "请求返回状态不是200")
		if err != nil {
			logs.Error("verify status code dump the result error", err)
		}
		return
	}
	// 提前检查jsonpath是否存在，不存在就报错
	for k := range verify {
		verifyO, err := jsonpath.JSONPath([]byte(response), k)
		if err != nil {
			logs.Error("doVerifyV2 jsonpath error，test failed", err)
		}
		if len(verifyO) == 0 {
			logs.Error("the verify key is not exist in the response", k)
			return
		}
	}
	for k, v := range verify {
		//fmt.Println(k, v, verify, reflect.TypeOf(verify))
		logs.Error("k,v is ", k, v, reflect.TypeOf(k))
		verifyO, _ := jsonpath.JSONPath([]byte(response), k)
		for subK, subV := range v {
			var vv interface{}
			// 根据类型转换jsonpath获取的数组首位类型
			switch subV.(type) {
			case string:
				vv = verifyO[0].MustString()
			case float64:
				vv = verifyO[0].MustNumeric()

			}
			if subK == "eq" {
				if subV != vv {
					logs.Error("not equal, key %s, actual value %v,expected %v", k, vv, subV)
					return
				}
			} else if subK == "need" {
				if subV != vv {
					logs.Error("not need, key %s, actual value %v,expected %v", k, vv, subV)
					return
				}
			} else if subK == "in" {
				if !strings.ContainsAny(vv.(string), subV.(string)) {
					logs.Error("not in, key %s, actual value %v,expected %v", k, vv, subV)
					return
				}
			} else if subK == "lt" {
				if !strings.ContainsAny(vv.(string), subV.(string)) {
					logs.Error("not lt, key %s, actual value %v,expected %v", k, vv, subV)
					return
				}
			} else if subK == "gt" {
				if !strings.ContainsAny(vv.(string), subV.(string)) {
					logs.Error("not gt, key %s, actual value %v,expected %v", k, vv, subV)
					return
				}
			} else if subK == "lte" {
				if !strings.ContainsAny(vv.(string), subV.(string)) {
					logs.Error("not lte, key %s, actual value %v,expected %v", k, vv, subV)
					return
				}
			} else if subK == "gte" {
				if !strings.ContainsAny(vv.(string), subV.(string)) {
					logs.Error("not gte, key %s, actual value %v,expected %v", k, vv, subV)
					return
				}
			} else {
				logs.Error("do not support")
				return
			}
		}
	}
}
