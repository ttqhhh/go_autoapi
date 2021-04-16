package libs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"go_autoapi/db_proxy"
	"io/ioutil"
	"net/http"
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

func DoRequest(url string, uuid string, data map[string]interface{}, verify interface{}) {
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
	if statusCode != 200 {
		fmt.Println("登录失败", jmap["message"])
		return
	}
	r.Incr(uuid)

}
