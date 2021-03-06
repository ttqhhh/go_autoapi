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
	// todo 暂时把格式化处理的相关逻辑挪到了DoRequest中
	//v := make(map[string]interface{})
	//err = json.Unmarshal([]byte(strings.TrimSpace(param)), &v)
	//if err != nil {
	//	logs.Error("发送冒烟请求前，解码json报错，err：", err)
	//	return
	//}
	//paramByte, err := json.Marshal(v)
	//if err != nil {
	//	logs.Error("发送冒烟请求前，处理请求json报错， err:", err)
	//	return
	//}
	// todo 通过字符串处理的方式，进行了json压缩，以便发生数据时服务器不报参数异常
	handleJson := HandleJson(param)
	paramByte := []byte(handleJson)
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
func DoRequest(domain string, url string, uuid string, param string, checkPoint string, caseId int64, isInspection int, runBy string) (isPass bool, resp string) {
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
	} else if isInspection == models.NOT_INSPECTION {
		testCaseMongo := models.TestCaseMongo{}
		testCaseMongo = testCaseMongo.GetOneCase(caseId)
		businessCode = testCaseMongo.BusinessCode
	} else {
		setCaseMongo := models.SetCaseMongo{}
		setCaseMongo, err := setCaseMongo.GetSetCaseById(caseId)
		if err != nil {
			logs.Error("获取单条case出错")
		}
		businessCode = setCaseMongo.BusinessCode
	}
	business, _ := strconv.Atoi(businessCode)
	if business == constant.Matuan {
		headers["debug"] = "1"
	}

	client := &http.Client{}
	// todo 通过字符串处理的方式，进行了json压缩，以便发生数据时服务器不报参数异常
	handleJson := HandleJson(param)
	paramByte := []byte(handleJson)
	postData := bytes.NewReader(paramByte)
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
	response, err := client.Do(req)
	if err != nil {
		logs.Error("DoRequest发起请求调用时出错, err: ", err)
		reason := "该接口不通, 请求超时..."
		result := models.AUTO_RESULT_FAIL
		resp = ""
		statusCode := 0
		SaveTestResult(uuid, caseId, isInspection, result, reason, runBy, resp, statusCode)
		isPass = false
		return
	}

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
		/**
			<option value="eq" selected>等于</option>	number/string
		    <option value="neq">不等于</option>			number/string
		    <option value="in">包含于</option>			string
		    <option value="exist">存在此字段</option>
		    <option value="lt">小于</option>			number
		    <option value="gt">大于</option>			number
		    <option value="lte">小于等于</option>		number
		    <option value="gte">大于等于</option>		number
			<option value="isTrue">为真</option>
			<option value="isFalse">为假</option>
		*/
		for checkType, checkValue := range checkRule {
			var vv interface{}
			// 根据类型转换jsonpath获取的数组首位类型
			if checkType != "exist" { // 当checkType不为exist时
				valueType := valueInResp[0].Type()
				if valueType != jsonpath.Numeric && valueType != jsonpath.String && valueType != jsonpath.Bool { // ①确保valueInResp的value值是否为number/string/bool类型
					logs.Error("the verify key is not last level in the response", path)
					reason += ";" + fmt.Sprintf("按照配置的json路径取到的值非基本类型（number/string/bool）, 请重新配置。json路径：【%s】", path)
					continue
				} else { // ② 确保checkType和checkValue匹配，然后取出valueInResp的value值
					if checkType == "isTrue" || checkType == "isFalse" { // 当checkType为exist/isTrue/isFalse时，不进行checkValue数据类型校验和响应值数据类型转换
						if valueType != jsonpath.Bool {
							logs.Error("the verify value type isn't bool", path)
							reason += ";" + fmt.Sprintf("按照配置的json路径取到的值类型不是布尔类型, 请重新配置。json路径：【%s】", path)
							continue
						}
						vv = valueInResp[0].MustBool()
					} else { // 当checkType不为exist、isTrue和isFalse时
						switch checkValue.(type) {
						case string:
							if checkType == "lt" || checkType == "gt" || checkType == "lte" || checkType == "gte" {
								logs.Error("the check type should not string type'", path)
								reason += ";" + fmt.Sprintf("校验类型为:【%s】时, 校验值类型不能为string。json路径：【%s】", checkType, path)
								continue
							}
							if valueType != jsonpath.String {
								logs.Error("the verify key type and value type isn't same", path)
								reason += ";" + fmt.Sprintf("按照配置的json路径取到的值类型与校验点中配置的值类型不符, 请重新配置。json路径：【%s】", path)
								continue
							}
							vv = valueInResp[0].MustString()
						case float64:
							if checkType == "in" {
								logs.Error("the check type should not number type'", path)
								reason += ";" + fmt.Sprintf("校验类型为:【%s】时, 校验值类型不能为number。json路径：【%s】", checkType, path)
								continue
							}
							if valueType != jsonpath.Numeric {
								logs.Error("the verify key type and value type isn't same", path)
								reason += ";" + fmt.Sprintf("按照配置的json路径取到的值类型与校验点中配置的值类型不符, 请重新配置。json路径：【%s】", path)
								continue
							}
							vv = valueInResp[0].MustNumeric()
						}
					}
				}
			}
			if checkType == "eq" {
				if checkValue != vv {
					logs.Error("not equal, key %s, actual value %checkRule,expected %checkRule", path, vv, checkValue)
					reason += ";" + fmt.Sprintf("不满足【相等】, json路径: 【%s】, 实际值: 【%v】, 期望值: 【%v】", path, vv, checkValue)
					continue
				}
			} else if checkType == "neq" {
				if checkValue == vv {
					logs.Error("equal, key %s, actual value %checkRule,expected %checkRule", path, vv, checkValue)
					reason += ";" + fmt.Sprintf("不满足【不相等】, json路径: 【%s】, 实际值: 【%v】, 期望值: 【%v】", path, vv, checkValue)
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
			} else if checkType == "isTrue" {
				if vv.(bool) != true {
					logs.Error("not isTrue, key %s, actual %checkRule, expected True", path, vv)
					reason += ";" + fmt.Sprintf("不满足【为真】, json路径: 【%s】, 实际值: 【%v】, 期望值: 【true】", path, vv)
					continue
				}
			} else if checkType == "isFalse" {
				if vv.(bool) != false {
					logs.Error("not isFalse, key %s, actual %checkRule, expected False", path, vv)
					reason += ";" + fmt.Sprintf("不满足【为假】, json路径: 【%s】, 实际值: 【%v】, 期望值: 【false】", path, vv)
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
