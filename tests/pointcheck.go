package main

import (
	"bytes"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"net/http"
)

func main(){
	// 一天的时间是
	getLoginCookie()
	//GetBasePoints(2)
	//actionList :=[] string {"detail_post"}
	//fmt.Println(actionList)
	//r := GetRealPoints("3d7c3cc366fe893164d634d61e9a04f6","expose_search","1625801485.586624","NaN")
	//fmt.Println(r)
	//jsonRight := `{"from":[{"sb":"sba"},{"sb":"sba"}],"to":"zs"}`
	//r := make(map[string]interface{})
	//err := json.Unmarshal([]byte(jsonRight), &r)
	//if err != nil{
	//	return
	//}
	//jsonLeft := `{"from":[{"sb":"sba"},{"sbs":"sbas"}],"to":"zs"}`
	//l := make(map[string]interface{})
	//err = json.Unmarshal([]byte(jsonLeft), &l)
	//var result1 string
	//var result2 bool
	//result1, result2 = JsonCompare(l,r,-1)
	//fmt.Println(result2)
	//fmt.Println(result1)
}

func getLoginCookie(){
	postBody := `{"username": "wangzhen01", "password": "Iepohg5go4iawoo"}`
	postData := bytes.NewReader([]byte(postBody))
	req, err := http.NewRequest("POST", "http://et.ixiaochuan.cn/proxy/api/user",postData)
	if err !=nil{
		logs.Error(err)
	}
	client := &http.Client{}
	response, _ := client.Do(req)
	//return response.Cookies()
	//fmt.Println(response.Cookies())
	//fmt.Println(response.Cookies()[0])
	str := fmt.Sprintf("%v",  response.Cookies())
	fmt.Println(str)
}
