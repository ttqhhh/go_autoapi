package main

import (
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	//"github.com/zyx4843/gojson"
	"reflect"
)

const left_data = "{\"ct\":\"1625820849\",\"mid\":\"11668967\",\"type\":\"click\",\"stype\":\"video\",\"id\":\"1585481065\",\"oid\":\"1585481065\",\"frominfo\":\"selectbulletvote\",\"extdata\":{\"app\":\"zuiyou\",\"app_name\":\"zuiyou\",\"app_ver\":\"5.7.13.318\",\"channel\":\"appstore\",\"did\":\"3d7c3cc366fe893164d634d61e9a04f6\",\"dt\":1,\"model\":\"iPhone 12 Pro\",\"net_type\":1,\"pid\":235125605,\"ver\":\"5.7.13.318\",\"video_id\":1585481065}}"

func main(){
	actionList :=[] string {"detail_post"}
	fmt.Println(actionList)
	r := GetRealPoints("3d7c3cc366fe893164d634d61e9a04f6","click_video","1625801485.586624","NaN")
	//jsonLeft := `{"from":[{"sb":"sba"},{"sb":"sba"}],"to":"zs"}`
	l := make(map[string]interface{})
	err := json.Unmarshal([]byte(left_data), &l)
	if err != nil{
		return
	}
	var result1 string
	var result2 bool
	result1, result2 = JsonCompare(l,r,-1)
	fmt.Println(result2)
	fmt.Println(result1)
}

type JsonDiff struct {
	HasDiff    bool
	Result     string
	Path       string

	//DataCompareResult  []string
	//FrameCompareResult []string
}

func GetRealPoints(did,event, timeBegin, timeEnd string) map[string]interface{} {
	fmt.Println("准备拉取数据，action:" + event)
	//时间是空位NaN
	req, err := http.NewRequest("GET",
		"http://172.16.2.217:8090/search?user="+did+"&event="+event+"&time_begin="+timeBegin+"&time_end="+timeEnd,
		nil)
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
				if !reflect.DeepEqual(value, json2[key]) {
					diff.HasDiff = true
					diff.Result = diff.Result + "\n path:" + diff.Path + ";值不相等 -- 期待值：" + value.(string) + "  实际值：" +json2[key].(string)
				}
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
		//diff.Path =  diff.Path + "["+ strconv.Itoa(i) +"]"
		switch json1[i].(type) {
		case map[string]interface{}:
			if _, ok := json2[i].(map[string]interface{}); ok {
				jsonDiffDict(json1[i].(map[string]interface{}), json2[i].(map[string]interface{}), depth+1, diff)
			} else {
				diff.HasDiff = true
				diff.Result = diff.Result + "\n path:" + diff.Path + ";实际值类型非 map[string]interface{} " + marshal(json2[i])
			}
		case []interface{}:
			if _, ok2 := json2[i].([]interface{}); !ok2 {
				diff.HasDiff = true
				diff.Result = diff.Result + "\n path:" + diff.Path + ";实际值类型非 interface{} -- " + marshal(json2[i])
			} else {
				jsonDiffList(json1[i].([]interface{}), json2[i].([]interface{}), depth+1, diff)
			}
		default:
			if !reflect.DeepEqual(json1[i], json2[i]) {
				diff.HasDiff = true
				diff.Result = diff.Result + "\n path:" + diff.Path + ";值不相等 -- 期待值：" + json1[i].(string) + "  实际值：" +json2[i].(string)
			}
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