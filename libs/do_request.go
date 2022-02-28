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
	"strconv"
	"strings"
)

func init() {
	fmt.Println("please init me")
	_ = db_proxy.InitClient()
}

// 获取冒烟响应函数
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

// 发送请求函数
func DoRequest(domain string, url string, uuid string, m string, checkPoint string, caseId int64, isInspection int, runBy string) (isPass bool, resp string) {
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
	resp = string(body)
	isPass = doVerify(respStatus, uuid, resp, verify, caseId, isInspection, runBy)
	return
}

// 结果验证函数
func doVerify(statusCode int, uuid string, response string, verify map[string]map[string]interface{}, caseId int64, isInspection int, runBy string) (isPass bool) {
	isPass = true
	reason := ""
	result := models.AUTO_RESULT_FAIL
	if statusCode != 200 {
		logs.Error("请求返回状态不是200，请求失败")
		reason = "状态码不是200"
		SaveTestResult(uuid, caseId, isInspection, result, reason, runBy, response, statusCode)
		isPass = false
		return
	}

	for path, checkRule := range verify {
		valueInResp, err := jsonpath.JSONPath([]byte(response), path)
		// 提前检查jsonpath是否存在，不存在就报错
		if err != nil {
			logs.Error("doVerify jsonpath error，test failed", err)
			reason = "checkpoint表达式有误 OR 不满足【存在】, json路径：【" + path + "】"
			SaveTestResult(uuid, caseId, isInspection, result, reason, runBy, response, statusCode)
			isPass = false
			return
		}
		if len(valueInResp) == 0 {
			logs.Error("the verify key is not exist in the response", path)
			reason += ";" + fmt.Sprintf("不满足【存在】, json路径：【%s】", path)
			continue
		}
		for checkType, checkValue := range checkRule {
			var vv interface{}
			// 根据类型转换jsonpath获取的数组首位类型
			if checkType != "exist" {
				switch checkValue.(type) {
				case string:
					vv = valueInResp[0].MustString()
				case float64:
					vv = valueInResp[0].MustNumeric()
				}
			}
			if checkType == "eq" {
				if checkValue != vv {
					logs.Error("not equal, key %s, actual value %checkRule,expected %checkRule", path, vv, checkValue)
					reason += ";" + fmt.Sprintf("不满足【相等】, json路径: 【%s】, 实际值: 【%v】, 期望值: 【%v】", path, vv, checkValue)
					continue
				}
			} else if checkType == "exist" {
				continue
			} else if checkType == "in" {
				if !strings.Contains(vv.(string), checkValue.(string)) {
					logs.Error("not in, key %s, actual value %checkRule,expected %checkRule", path, vv, checkValue)
					reason += ";" + fmt.Sprintf("不满足【包含】, json路径: 【%s】, 实际值: 【%v】, 期望值: 【包含于%v】", path, vv, checkValue)
					continue
				}
			} else if checkType == "lt" {
				if !(vv.(float64) < checkValue.(float64)) {
					logs.Error("not lt, key %s, actual %checkRule < expected %checkRule", path, vv, checkValue)
					reason += ";" + fmt.Sprintf("不满足【小于】, json路径: 【%s】, 实际值: 【%v】, 期望值: 【<%v】", path, vv, checkValue)
					continue
				}
			} else if checkType == "gt" {
				if !(vv.(float64) > checkValue.(float64)) {
					logs.Error("not gt, key %s, actual %checkRule > expected %checkRule", path, vv, checkValue)
					reason += ";" + fmt.Sprintf("不满足【大于】, json路径: 【%s】, 实际值: 【%v】, 期望值: 【>%v】", path, vv, checkValue)
					continue
				}
			} else if checkType == "lte" {
				if !(vv.(float64) <= checkValue.(float64)) {
					logs.Error("not lte, key %s, actual %checkRule <= expected %checkRule", path, vv, checkValue)
					reason += ";" + fmt.Sprintf("不满足【小于等于】, json路径: 【%s】, 实际值: 【%v】, 期望值: 【<=%v】", path, vv, checkValue)
					continue
				}
			} else if checkType == "gte" {
				if !(vv.(float64) >= checkValue.(float64)) {
					logs.Error("not gte, key %s, actual %checkRule >= expected %checkRule", path, vv, checkValue)
					reason += ";" + fmt.Sprintf("不满足【大于等于】, json路径: 【%s】, 实际值: 【%v】, 期望值: 【>=%v】", path, vv, checkValue)
					continue
				}
			} else {
				logs.Error("do not support, checkType: ", checkType)
				reason += ";" + fmt.Sprintf("不支持的验证类型: %s", checkType)
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
		SaveTestResult(uuid, caseId, isInspection, result, reason, runBy, response, statusCode)
	}
	return
}

/**
statusCode 为0时，表示场景测试中，前后校验逻辑 or 处理逻辑出错。
*/
func SaveTestResult(uuid string, caseId int64, isInspection int, result int, reason string, author string, resp string, statusCode int) {
	err := models.InsertResult(uuid, caseId, isInspection, result, reason, author, resp, statusCode)
	if err != nil {
		logs.Error("save test result error,please check the db connection", err)
	}
	return
}
