package libs

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	jsonpath "github.com/spyzhov/ajson"
	"go_autoapi/db_proxy"
	"go_autoapi/models"
	"io"
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
func HttpPost(postUrl string, headers map[string]string, jsonMap string, method string) (int, string, string) {
	client := &http.Client{}
	//转换成postBody
	//bytesData, err := json.Marshal(jsonMap)
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return 0, "", ""
	//}
	postBody := bytes.NewReader([]byte(jsonMap))
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
	fmt.Printf("body is %v", string(body))
	return resp.StatusCode, string(body), cookieStr
}

func DoRequestWithNoneVerify(url string, param string) (respStatus int, body []byte, err error) {
	headers := map[string]string{
		"ZYP":             "mid=248447243",
		"X-Xc-Agent":      "av=5.7.1.001,dt=0",
		"User-Agent":      "okhttp/3.12.2 Zuiyou/5.7.1.001 (Android/29)",
		"Request-Type":    "text/json",
		"Content-Type":    "application/json; charset=utf-8",
		"Content-Length":  "",
		"Host":            "api.izuiyou.com",
		"Accept-Encoding": "gzip",
		"Connection":      "keep-alive",
		"Accept-Charset":  "utf-8"}
	//redis? 不知道干嘛的
	client := &http.Client{}
	postData := bytes.NewReader([]byte(param))
	req, err := http.NewRequest("POST", url, postData)
	if err != nil {
		logs.Error("请求失败")
		return
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	response, _ := client.Do(req)
	respStatus = response.StatusCode
	var reader io.ReadCloser
	if response.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(response.Body)
		if err != nil {
			logs.Error("gzip解析响应body失败, err:", err)
			return
		}
	} else {
		reader = response.Body
	}
	body, err = ioutil.ReadAll(reader)
	if err != nil {
		logs.Error("reader获取响应body失败, err: ", err)
		return
	}
	return
}

func DoRequestV2(url string, uuid string, m string, checkPoint string, caseId int64, runBy string) {
	headers := map[string]string{
		"ZYP":             "mid=248447243",
		"X-Xc-Agent":      "av=5.7.1.001,dt=0",
		"User-Agent":      "okhttp/3.12.2 Zuiyou/5.7.1.001 (Android/29)",
		"Request-Type":    "text/json",
		"Content-Type":    "application/json; charset=utf-8",
		"Content-Length":  "",
		"Host":            "api.izuiyou.com",
		"Accept-Encoding": "gzip",
		"Connection":      "keep-alive",
		"Accept-Charset":  "utf-8"}
	//redis? 不知道干嘛的
	client := &http.Client{}
	postData := bytes.NewReader([]byte(m))
	req, err := http.NewRequest("POST", url, postData)
	if err != nil {
		logs.Error("构建请求失败, err:", err)
		return
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	response, _ := client.Do(req)
	respStatus := response.StatusCode
	var reader io.ReadCloser
	if response.Header.Get("Content-Encoding") == "gzip" {
		reader, err = gzip.NewReader(response.Body)
		if err != nil {
			return
		}
	} else {
		reader = response.Body
	}
	body, err := ioutil.ReadAll(reader)
	var verify map[string]map[string]interface{}
	if err := json.Unmarshal([]byte(checkPoint), &verify); err != nil {
		logs.Error("checkpoint解析失败", err)
		return
	}
	doVerifyV2(respStatus, uuid, string(body), verify, caseId, runBy)
}

//func DoRequest(url string, method string, uuid string, data string, verify string, caseId int64) {
//	//密码
//	r := db_proxy.GetRedisObject()
//	statusCode, body, _ := HttpPost(url, nil, data, method)
//	//body jsonStr转map
//	var jmap map[string]interface{}
//	if err := json.Unmarshal([]byte(body), &jmap); err != nil {
//		fmt.Println("解析失败", err)
//		return
//	}
//	// 此处采用go-simplejson来做个示例，用于以后扩展检查使用
//	js, err := simplejson.NewJson([]byte(body))
//	if err != nil {
//		return
//	}
//	email, err := js.Get("data").Get("email").String()
//	fmt.Println(js.Get("code"), email)
//
//	// 判断某个字段的类型
//	//fmt.Println("type:", reflect.TypeOf(jmap["code"]))
//	//判断登录是否成功
//	doVerifyV2(statusCode, uuid, body, verify, caseId)
//	r.Incr(uuid)
//
//}

// 采用jsonpath 对结果进行验证
func doVerifyV2(statusCode int, uuid string, response string, verify map[string]map[string]interface{}, caseId int64, runBy string) {

	if statusCode != 200 {
		logs.Error("请求返回状态不是200，请求失败")
		saveTestResult(uuid, caseId, "状态码不是200", runBy, response)
		return
	}
	// 提前检查jsonpath是否存在，不存在就报错
	for k := range verify {
		verifyO, err := jsonpath.JSONPath([]byte(response), k)
		if err != nil {
			logs.Error("doVerifyV2 jsonpath error，test failed", err)
			saveTestResult(uuid, caseId, k+" jsonpath err", runBy, response)
		}
		if len(verifyO) == 0 {
			logs.Error("the verify key is not exist in the response", k)
			saveTestResult(uuid, caseId, k+" the verify key not exist err", runBy, response)
			return
		}
	}

	isPass := true
	resultDesc := ""
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
					isPass = false
					logs.Error("not equal, key %s, actual value %v,expected %v", k, vv, subV)
					resultDesc += ";" + fmt.Sprintf("not equal, key %s, actual value %v,expected %v", k, vv, subV)
					continue
				}
			} else if subK == "need" {
				if subV != vv {
					isPass = false
					logs.Error("not need, key %s, actual value %v,expected %v", k, vv, subV)
					resultDesc += ";" + fmt.Sprintf("not need, key %s, actual value %v,expected %v", k, vv, subV)
					continue
				}
			} else if subK == "in" {
				if !strings.Contains(vv.(string), subV.(string)) {
					isPass = false
					logs.Error("not in, key %s, actual value %v,expected %v", k, vv, subV)
					resultDesc += ";" + fmt.Sprintf("not in, key %s, actual value %v,expected %v", k, vv, subV)
					continue
				}
			} else if subK == "lt" {
				if !(vv.(float64) < subV.(float64)) {
					isPass = false
					logs.Error("not lt, key %s, actual %v < expected %v", k, vv, subV)
					resultDesc += ";" + fmt.Sprintf("not lt, key %s, actual %v < expected %v", k, vv, subV)
					continue
				}
			} else if subK == "gt" {
				if !(vv.(float64) > subV.(float64)) {
					isPass = false
					logs.Error("not gt, key %s, actual %v > expected %v", k, vv, subV)
					resultDesc += ";" + fmt.Sprintf("not gt, key %s, actual %v > expected %v", k, vv, subV)
					continue
				}
			} else if subK == "lte" {
				if !(vv.(float64) <= subV.(float64)) {
					isPass = false
					logs.Error("not lte, key %s, actual %v <= expected %v", k, vv, subV)
					resultDesc += ";" + fmt.Sprintf("not lte, key %s, actual %v <= expected %v", k, vv, subV)
					continue
				}
			} else if subK == "gte" {
				if !(vv.(float64) >= subV.(float64)) {
					isPass = false
					logs.Error("not gte, key %s, actual %v >= expected %v", k, vv, subV)
					resultDesc += ";" + fmt.Sprintf("not gte, key %s, actual %v >= expected %v", k, vv, subV)
					continue
				}
			} else {
				logs.Error("do not support")
				isPass = false
				resultDesc += ";" + fmt.Sprintf("do not support this operator")
				continue
			}
		}
	}
	// 将该case执行结果聚合入库
	if !isPass {
		isNeedSub := strings.HasPrefix(resultDesc, ";")
		if isNeedSub {
			resultDescRune := []rune(resultDesc)
			resultDesc = string(resultDescRune[1:])
		}
		saveTestResult(uuid, caseId, resultDesc, runBy, response)
	}
}

func saveTestResult(uuid string, caseId int64, reason string, author string, resp string) {
	err := models.InsertResult(uuid, caseId, reason, author, resp)
	if err != nil {
		logs.Error("save test result error,please check the db connection", err)
	}
	return
}
