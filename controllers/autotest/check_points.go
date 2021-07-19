package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const cookie = "99BFD42401E3660BFE97D2268BB1EC5A"

type pointData struct{
	Limit    int 		`form:"limit" json:"limit"`
	Business string		`form:"business" json:"business"`
	Did      string		`form:"did" json:"did"`
}

func (c*AutoTestController) showCheckPoints(){
	c.TplName = "check_points.html"
}

func (c*AutoTestController) checkPoints(){
	pd := pointData{}
	if err := c.ParseForm(&pd); err != nil { // 传入user指针
		c.Ctx.WriteString("出错了！")
	}
	resultMsg := GetBasePoints(pd.Limit, pd.Business, pd.Did)
	c.Data["result"] = resultMsg
	c.TplName = "check_points.html"
}

func getLoginCookie() string {
	postBody := `{"username": "wangzhen01", "password": "Iepohg5go4iawoo"}`
	postData := bytes.NewReader([]byte(postBody))
	req, err := http.NewRequest("POST", "http://et.ixiaochuan.cn/proxy/api/user",postData)
	if err !=nil{
		logs.Error(err)
	}
	client := &http.Client{}
	response, _ := client.Do(req)
	ck := fmt.Sprintf("%v",  response.Cookies())
	sep := ";"
	sep2 := "="
	result := strings.Split(strings.Split(ck,sep)[0],sep2)[1]
	// 再用这个cookie登录一次
	postBody2 := `{"username": "wangzhen01", "password": "Iepohg5go4iawoo"}`
	postData2 := bytes.NewReader([]byte(postBody2))
	req2, _ := http.NewRequest("POST", "http://et.ixiaochuan.cn/proxy/api/user",postData2)
	req2.Header.Add("Cookie","JSESSIONID="+result)
	req2.Header.Add("Content-Type","application/json")
	client2 := &http.Client{}
	response2, _ := client2.Do(req2)
	var reader io.ReadCloser
	reader = response2.Body
	body2, _:= ioutil.ReadAll(reader)
	fmt.Println(string(body2))
	return result
	//fmt.Println(response)

}

type JsonDiff struct {
	HasDiff    bool
	Result     string
	Path       string
}

// 通过埋点系统平台获取标准校验点,获取前一周时间范围的校验点（当前取最后一个校验点）

// todo 输入的参数有：limit(查询的个数)；

func GetBasePoints(limit int, business string, did string) []string{
	var resultMsg []string
	fmt.Println("第一次执行获取total总数")
	limits := strconv.Itoa(limit)
	// 当前是按照时间倒序查询，limit限制查询总数
	postBody := `{"offset": 0,"limit": `+ limits +`,"app_name": "` +business+ `","sort_field":"update_time","sort_flag":"desc"}`
	postData := bytes.NewReader([]byte(postBody))
	req, err := http.NewRequest("POST", "http://et.ixiaochuan.cn/proxy/api/event_list",postData)
	if err !=nil{
		logs.Error(err)
	}
	cookies := getLoginCookie()
	req.Header.Add("Cookie","JSESSIONID="+cookies)
	req.Header.Add("Content-Type","application/json")
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		logs.Error("请求失败, err:", err)
	}
	var reader io.ReadCloser
	reader = response.Body
	body, _:= ioutil.ReadAll(reader)
	v := make(map[string]interface{})
	_ = json.Unmarshal(body,&v)
	records := v["data"].(map[string]interface{})["records"].([]interface{})
	for _,val := range records {
		fmt.Println("循环获取event_detail")
		vals := val.(map[string]interface{})
		appName := vals["app_name"].(string)
		types := vals["type"].(string)
		stype := vals["stype"].(string)
		frominfo := vals["frominfo"].(string)
		postBody2 := `{"app_name": "`+appName+`","frominfo":"`+frominfo+`","is_approval": "false","type":"`+types+`","stype":"`+stype+`"}`
		postData2:= bytes.NewReader([]byte(postBody2))
		req2, _ := http.NewRequest("POST", "http://et.ixiaochuan.cn/proxy/api/event_detail", postData2)
		req2.Header.Add("Cookie","JSESSIONID="+cookies)
		req2.Header.Add("Content-Type","application/json")
		if err != nil {
			logs.Error("请求失败，err: ", err)
		}
		client2 := &http.Client{}
		response2, err := client2.Do(req2)
		if err != nil {
			logs.Error("请求失败, err:", err)
		}
		var reader2 io.ReadCloser
		reader2 = response2.Body
		body2, _:= ioutil.ReadAll(reader2)
		v2 := make(map[string]interface{})
		_ = json.Unmarshal(body2,&v2)
		// 主要获取拓展字段自定义字段
		newExtended := make(map[string]interface{})
		extendedCustom := v2["data"].(map[string]interface{})["extended_custom"].([]interface{})
		for _,valss := range extendedCustom{
			newExtended[valss.(map[string]interface{})["field_name"].(string)] = "none"
		}
		finallExtended := make(map[string]interface{})
		finallExtended["cur_page"] = "none"
		finallExtended["from_page"] = "none"
		finallExtended["extdata"] = newExtended
		fmt.Println(finallExtended)
		// 开始获取真实入库数据
		r := GetRealPoints(did,types+"_"+stype,"1625801485.586624","NaN",appName)
		if r == nil{
			resultMsg = append(resultMsg,"没有查询到数据，已跳过："+types+"_"+stype + "\n")
			logs.Error("没有查询到数据，已跳过："+types+"_"+stype)
		}else{
			var result1 string
			var result2 bool
			result1, result2 = JsonCompare(finallExtended,r,-1)
			if result2 == true {
				resultMsg = append(resultMsg,"检查到异常，事件:"+types+"_"+stype +"  ; " + result1)
				fmt.Println("检查到异常，事件:"+types+"_"+stype)
				fmt.Println(result1)
			}else{
				resultMsg = append(resultMsg,"检查结构通过，事件:"+types+"_"+stype +"\n")
				fmt.Println("检查通过，事件:"+types+"_"+stype)
			}
		}
	}
	return resultMsg
}

func GetRealPoints(did,event, timeBegin, timeEnd ,appName string) map[string]interface{} {
	fmt.Println("准备拉取数据，action:" + event)
	//时间是空位NaN
	urls := ""
	if appName != "omg" {
		urls = "http://172.16.2.217:8090/search?user="+did+"&event="+event+"&time_begin="+timeBegin+"&time_end="+timeEnd
	}else{
		urls = "http://10.12.44.53:9090//search?user="+did+"&event="+event+"&time_begin="+timeBegin+"&time_end="+timeEnd
	}
	req, err := http.NewRequest("GET", urls,nil)
	if err != nil {
		logs.Error("请求失败，err: ", err)
	}
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		logs.Error("请求失败, err:", err)
	}
	var reader io.ReadCloser
	reader = response.Body
	body, _:= ioutil.ReadAll(reader)
	//Sbody := string(body)
	var result []string
	_ = json.Unmarshal(body,&result)
	if len(result)==0 {
		logs.Error("当前did下没有发现行为埋点：",event)
		return nil
	}
	arr :=strings.Fields(result[0])
	realJson := arr[len(arr)-1]
	v := make(map[string]interface{})
	_ = json.Unmarshal([]byte(realJson),&v)
	ext := make(map[string]interface{})
	_ = json.Unmarshal([]byte(v["extdata"].(string)),&ext)
	v["extdata"] = ext
	return v
}

func JsonCompare(left, right map[string]interface{}, n int) (string, bool) {
	diff := &JsonDiff{HasDiff: false, Result: ""}
	jsonDiffDict(left, right, 1, diff)
	if diff.HasDiff {
		if n < 0 {
			return diff.Result, diff.HasDiff
		} else {
			return processContext(diff.Result, n), diff.HasDiff
		}
	}
	fmt.Println(diff.Path)
	return "", diff.HasDiff
}

func marshal(j interface{}) string {
	value, _ := json.Marshal(j)
	return string(value)
}

func jsonDiffDict(json1, json2 map[string]interface{}, depth int, diff *JsonDiff) {
	for key, value := range json1 {
		diff.Path = diff.Path + "[" + key + "]"
		if _, ok := json2[key]; ok {
			switch value.(type) {
			case map[string]interface{}:
				if _, ok2 := json2[key].(map[string]interface{}); !ok2 {
					diff.HasDiff = true
					diff.Result = diff.Result + "\n path:" + diff.Path + ";实际值类型非 map[string]interface{} " + marshal(json2[key])
				} else {
					jsonDiffDict(value.(map[string]interface{}), json2[key].(map[string]interface{}), depth+1, diff)
				}
			case []interface{}:
				if _, ok2 := json2[key].([]interface{}); !ok2 {
					diff.HasDiff = true
					diff.Result = diff.Result + "\n path:" + diff.Path + ";实际值类型非 interface{} -- " + marshal(json2[key])
				} else {
					jsonDiffList(value.([]interface{}), json2[key].([]interface{}), depth+1, diff)
				}
			default:
				//if !reflect.DeepEqual(value, json2[key]) {
				//	diff.HasDiff = true
				//	diff.Result = diff.Result + "\n path:" + diff.Path + ";值不相等 -- 期待值：" + value.(string) + "  实际值：" +json2[key].(string)
				//}
			}
		} else {
			diff.HasDiff = true
			diff.Result = diff.Result + "\n 键不存在：" + key
		}
	}
}

func jsonDiffList(json1, json2 []interface{}, depth int, diff *JsonDiff) {

	size := len(json1)
	if size > len(json2) {
		size = len(json2)
	}
	for i := 0; i < size; i++ {
		switch json1[i].(type) {
		case map[string]interface{}:
			if _, ok := json2[i].(map[string]interface{}); ok {
				jsonDiffDict(json1[i].(map[string]interface{}), json2[i].(map[string]interface{}), depth+1, diff)
			} else {
				diff.HasDiff = true
				diff.Path =  diff.Path + "["+ strconv.Itoa(i) +"]"
				diff.Result = diff.Result + "\n path:" + diff.Path + ";实际值类型非 map[string]interface{} " + marshal(json2[i])
			}
		case []interface{}:
			if _, ok2 := json2[i].([]interface{}); !ok2 {
				diff.HasDiff = true
				diff.Path =  diff.Path + "["+ strconv.Itoa(i) +"]"
				diff.Result = diff.Result + "\n path:" + diff.Path + ";实际值类型非 interface{} -- " + marshal(json2[i])
			} else {
				jsonDiffList(json1[i].([]interface{}), json2[i].([]interface{}), depth+1, diff)
			}
		default:
			//if !reflect.DeepEqual(json1[i], json2[i]) {
			//	diff.HasDiff = true
			//	diff.Path =  diff.Path + "["+ strconv.Itoa(i) +"]"
			//	diff.Result = diff.Result + "\n path:" + diff.Path + ";值不相等 -- 期待值：" + json1[i].(string) + "  实际值：" +json2[i].(string)
			//}
		}
	}
}

func processContext(diff string, n int) string {
	index1 := strings.Index(diff, "\n-")
	index2 := strings.Index(diff, "\n+")
	begin := 0
	end := 0
	if index1 >= 0 && index2 >= 0 {
		if index1 <= index2 {
			begin = index1
		} else {
			begin = index2
		}
	} else if index1 >= 0 {
		begin = index1
	} else if index2 >= 0 {
		begin = index2
	}
	index1 = strings.LastIndex(diff, "\n-")
	index2 = strings.LastIndex(diff, "\n+")
	if index1 >= 0 && index2 >= 0 {
		if index1 <= index2 {
			end = index2
		} else {
			end = index1
		}
	} else if index1 >= 0 {
		end = index1
	} else if index2 >= 0 {
		end = index2
	}
	pre := diff[0:begin]
	post := diff[end:]
	i := 0
	l := begin
	for i < n && l >= 0 {
		i++
		l = strings.LastIndex(pre[0:l], "\n")
	}
	r := 0
	j := 0
	for j <= n && r >= 0 {
		j++
		t := strings.Index(post[r:], "\n")
		if t >= 0 {
			r = r + t + 1
		}
	}
	if r < 0 {
		r = len(post)
	}
	return pre[l+1:] + diff[begin:end] + post[0:r+1]
}

func LoadJson(path string, dist interface{}) (err error) {
	var content []byte
	if content, err = ioutil.ReadFile(path); err == nil {
		err = json.Unmarshal(content, dist)
	}
	return err
}