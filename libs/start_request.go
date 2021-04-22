package libs

import (
	"bytes"
	"fmt"
	"github.com/astaxie/beego/logs"
	"go_autoapi/db_proxy"
	"io"
	"net/http"
	"os"
)

func init() {
	fmt.Println("please init me")
	_ = db_proxy.InitClient()
}

func StartRequest(url string, m string, checkPoint string) {
	headers := map[string]string{
		"ZYP": "mid=248447243",
		"X-Xc-Agent": "av=5.7.1.001,dt=0",
		"User-Agent": "okhttp/3.12.2 Zuiyou/5.7.1.001 (Android/29)",
		"Request-Type": "text/json",
		"Content-Type": "application/json; charset=utf-8",
		"Content-Length": "",
		"Host": "api.izuiyou.com",
		"Accept-Encoding": "gzip",
		"Connection": "keep-alive"}
	//redis? 不知道干嘛的
	client := &http.Client{}
	postData := bytes.NewReader([]byte(m))
	req, err := http.NewRequest("POST", url, postData)
	if err != nil{
		logs.Error("请求失败")
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	response,_ := client.Do(req)
	result := response.Body
	stdout := os.Stdout
	_, err = io.Copy(stdout, response.Body)
	fmt.Println(result)

}

