package libs

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	jsonpath "github.com/spyzhov/ajson"
	constant "go_autoapi/constants"
	"go_autoapi/db_proxy"
	"go_autoapi/models"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
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

func DoRequestWithNoneVerify(business int, url string, param string) (respStatus int, body []byte, err error) {
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
	// 当Case所属业务线为麻团时，增加debug请求头
	if business == constant.Matuan {
		headers["debug"] = "1"
	}

	client := &http.Client{}
	// todo 千万不要删，用于处理json格式化问题（删了后某些服务会报504问题）
	v := make(map[string]interface{})
	err = json.Unmarshal([]byte(strings.TrimSpace(param)), &v)
	if err != nil {
		logs.Error("发送冒烟请求前，解码json报错，err：", err)
		return
	}
	paramByte, err := json.Marshal(v)
	//logs.Info("打印json", string(paramByte))
	if err != nil {
		logs.Error("发送冒烟请求前，处理请求json报错， err:", err)
		return
	}
	postData := bytes.NewReader(paramByte)
	req, err := http.NewRequest("POST", url, postData)
	if err != nil {
		logs.Error("冒烟发送业务请求失败，err: ", err)
		return
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	response, err := client.Do(req)
	if err != nil {
		logs.Error("冒烟请求失败, err:", err)
		return
	}
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

func DoRequestV2(domain string, url string, uuid string, m string, checkPoint string, caseId int64, isInspection int, runBy string) (isPass bool) {
	isPass = true
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
	// 当Case所属业务线为麻团时，增加debug请求头
	var businessCode string
	if isInspection == models.INSPECTION {
		icm := models.InspectionCaseMongo{}
		icm = icm.GetOneCase(caseId)
		businessCode = icm.BusinessCode
	} else {
		testCaseMongo := models.TestCaseMongo{}
		testCaseMongo = testCaseMongo.GetOneCase(caseId)
		businessCode = testCaseMongo.BusinessCode
	}
	business, _ := strconv.Atoi(businessCode)
	if business == constant.Matuan {
		headers["debug"] = "1"
	}

	client := &http.Client{}
	postData := bytes.NewReader([]byte(m))
	// 对domain和url进行兼容性拼接
	if strings.HasSuffix(domain, "/") {
		domain = domain[:len(domain)-1]
	}
	if strings.HasPrefix(url, "/") {
		url = url[1:]
	}
	var path = ""
	if domain != "" {
		path = domain + "/" + url
	} else {
		path = url
	}
	if path == "" {
		logs.Error("请求路径为空，请检查Case配置是否有误")
		return
	}
	req, err := http.NewRequest("POST", path, postData)
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
	isPass = doVerifyV2(respStatus, uuid, string(body), verify, caseId, isInspection, runBy)
	return
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
func doVerifyV2(statusCode int, uuid string, response string, verify map[string]map[string]interface{}, caseId int64, isInspection int, runBy string) (isPass bool) {
	isPass = true
	reason := ""
	result := models.AUTO_RESULT_FAIL
	if statusCode != 200 {
		logs.Error("请求返回状态不是200，请求失败")
		reason = "状态码不是200"
		saveTestResult(uuid, caseId, isInspection, result, reason, runBy, response, statusCode)
		isPass = false
		return
	}
	// 提前检查jsonpath是否存在，不存在就报错
	for k := range verify {
		verifyO, err := jsonpath.JSONPath([]byte(response), k)
		if err != nil {
			logs.Error("doVerifyV2 jsonpath error，test failed", err)
			//saveTestResult(uuid, caseId, result, k+" jsonpath err", runBy, response)
			reason = "checkpoint表达式有误，请检查您的checkpoint (" + k + ")"
			saveTestResult(uuid, caseId, isInspection, result, reason, runBy, response, statusCode)
			isPass = false
			return
		}
		if len(verifyO) == 0 {
			logs.Error("the verify key is not exist in the response", k)
			reason = "json路径: 【" + k + "】, 未配置有效的校验规则"
			//saveTestResult(uuid, caseId, result, k+" the verify key not exist err", runBy, response)
			saveTestResult(uuid, caseId, isInspection, result, reason, runBy, response, statusCode)
			isPass = false
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
					//reason += ";" + fmt.Sprintf("not equal, key %s, actual value %v,expected %v", k, vv, subV)
					reason += ";" + fmt.Sprintf("不满足【相等】, json路径: 【%s】, 实际值: 【%v】, 期望值: 【%v】", k, vv, subV)
					continue
				}
			} else if subK == "need" {
				if subV != vv {
					logs.Error("not need, key %s, actual value %v,expected %v", k, vv, subV)
					//reason += ";" + fmt.Sprintf("not need, key %s, actual value %v,expected %v", k, vv, subV)
					reason += ";" + fmt.Sprintf("不满足【必须】, json路径: 【%s】, 实际值: 【%v】, 期望值: 【%v】", k, vv, subV)
					continue
				}
			} else if subK == "in" {
				if !strings.Contains(vv.(string), subV.(string)) {
					logs.Error("not in, key %s, actual value %v,expected %v", k, vv, subV)
					//reason += ";" + fmt.Sprintf("not in, key %s, actual value %v,expected %v", k, vv, subV)
					reason += ";" + fmt.Sprintf("不满足【包含】, json路径: 【%s】, 实际值: 【%v】, 期望值: 【包含于%v】", k, vv, subV)
					continue
				}
			} else if subK == "lt" {
				if !(vv.(float64) < subV.(float64)) {
					logs.Error("not lt, key %s, actual %v < expected %v", k, vv, subV)
					//reason += ";" + fmt.Sprintf("not lt, key %s, actual %v < expected %v", k, vv, subV)
					reason += ";" + fmt.Sprintf("不满足【小于】, json路径: 【%s】, 实际值: 【%v】, 期望值: 【<%v】", k, vv, subV)
					continue
				}
			} else if subK == "gt" {
				if !(vv.(float64) > subV.(float64)) {
					logs.Error("not gt, key %s, actual %v > expected %v", k, vv, subV)
					//reason += ";" + fmt.Sprintf("not gt, key %s, actual %v > expected %v", k, vv, subV)
					reason += ";" + fmt.Sprintf("不满足【大于】, json路径: 【%s】, 实际值: 【%v】, 期望值: 【>%v】", k, vv, subV)
					continue
				}
			} else if subK == "lte" {
				if !(vv.(float64) <= subV.(float64)) {
					logs.Error("not lte, key %s, actual %v <= expected %v", k, vv, subV)
					//reason += ";" + fmt.Sprintf("not lte, key %s, actual %v <= expected %v", k, vv, subV)
					reason += ";" + fmt.Sprintf("不满足【小于等于】, json路径: 【%s】, 实际值: 【%v】, 期望值: 【<=%v】", k, vv, subV)
					continue
				}
			} else if subK == "gte" {
				if !(vv.(float64) >= subV.(float64)) {
					logs.Error("not gte, key %s, actual %v >= expected %v", k, vv, subV)
					//reason += ";" + fmt.Sprintf("not gte, key %s, actual %v >= expected %v", k, vv, subV)
					reason += ";" + fmt.Sprintf("不满足【大于等于】, json路径: 【%s】, 实际值: 【%v】, 期望值: 【>=%v】", k, vv, subV)
					continue
				}
			} else {
				logs.Error("do not support, subK: ", subK)
				//reason += ";" + fmt.Sprintf("do not support this operator")
				reason += ";" + fmt.Sprintf("不支持的验证类型: %s", subK)
				continue
			}
		}
	}
	// 将该case执行结果聚合入库
	if reason == "" {
		result = models.AUTO_RESULT_SUCCESS
	} else {
		isNeedSub := strings.HasPrefix(reason, ";")
		if isNeedSub {
			resultDescRune := []rune(reason)
			reason = string(resultDescRune[1:])
		}
		isPass = false
	}
	saveTestResult(uuid, caseId, isInspection, result, reason, runBy, response, statusCode)
	return
}

func saveTestResult(uuid string, caseId int64, isInspection int, result int, reason string, author string, resp string, statusCode int) {
	err := models.InsertResult(uuid, caseId, isInspection, result, reason, author, resp, statusCode)
	if err != nil {
		logs.Error("save test result error,please check the db connection", err)
	}
	return
}
