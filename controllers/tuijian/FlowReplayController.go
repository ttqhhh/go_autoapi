package tuijian

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	constant "go_autoapi/constants"
	"go_autoapi/libs"
	"go_autoapi/models"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

const uploadDir = "/home/work/efficiency/upload/tuijian"

type FlowReplayController struct {
	libs.BaseController
}

func (c *FlowReplayController) Get() {
	do := c.GetMethodName()
	switch do {
	case "index":
		c.index()
	case "list":
		c.list()
	case "getById":
		c.getById()
	case "showAllFlowFiles":
		c.showAllFlowFiles()

	default:
		logs.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持", nil)
	}
}

func (c *FlowReplayController) Post() {
	do := c.GetMethodName()
	switch do {
	case "save":
		c.add()
	case "update":
		c.update()
	case "remove":
		c.remove()
	case "kill":
		c.Killreplay()
	case "replay":
		c.Replay()
	case "replaycycle":
		c.ReplayCycle()
	case "collect_flow_file":
		c.collectFlowFile()
	case "pressure_measurement":
		c.pressure_Measurement()

	default:
		logs.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持", nil)
	}
}

func (c *FlowReplayController) index() {
	c.TplName = "replay.html"
}

func (c *FlowReplayController) list() {
	// 只能看到自己有权限的服务
	serviceName := c.GetString("service_name")
	pageNo, err := c.GetInt("page", 1)
	if err != nil {
		logs.Warn("/service/page接口 参数异常, err: %v", err)
		c.ErrorJson(-1, "参数异常", nil)
	}
	pageSize, err := c.GetInt("limit", 10)
	if err != nil {
		logs.Warn("/service/page接口 参数异常, err: %v", err)
		c.ErrorJson(-1, "参数异常", nil)
	}
	logs.Info("请求参数: service_name=%v, business=%v, page_no=%v, page_size=%v", serviceName, pageNo, pageSize)
	flowReplayMongo := models.FlowReplayMongo{}
	services, total, err := flowReplayMongo.QueryByPage(serviceName, pageNo, pageSize)
	if err != nil {
		c.ErrorJson(-1, "服务查询数据异常", nil)
	}

	res := make(map[string]interface{})
	res["code"] = 0
	res["msg"] = "成功"
	res["count"] = total
	res["data"] = services

	c.Data["json"] = res
	c.ServeJSON() //对json进行序列化输出
	c.StopRun()
}

func (c *FlowReplayController) add() {
	userId, _ := c.GetSecureCookie(constant.CookieSecretKey, "user_id")
	id, err1 := c.GetInt64("id")
	logs.Info("请求参数: id=%v", id)
	if err1 != nil {
		flowreplay := &models.FlowReplayMongo{}
		err := c.ParseForm(flowreplay)
		if err != nil {
			logs.Warn("/flowreplay/add接口 参数异常, err: %v", err)
			c.ErrorJson(-1, "参数异常", nil)
		}
		logs.Info("请求参数：%v", flowreplay)

		//验证服务名 唯一性
		serviceName := flowreplay.ServiceName
		temp, err := flowreplay.QueryByName(serviceName)
		if err != nil {
			logs.Error("流量回放添加时, 验证serviceName唯一性时报错")
		}
		if temp != nil {
			c.ErrorJson(-1, "存在服务名相同的流量", nil)
		}

		flowreplay.CreateBy = userId
		err = flowreplay.Insert(*flowreplay)
		if err != nil {
			c.ErrorJson(-1, "服务添加数据异常", nil)
		}
		c.Redirect("/flowreplay/index", http.StatusFound)
	}
	if err1 == nil {
		c.update()
	}

}

func (c *FlowReplayController) getById() {
	id, err := c.GetInt64("id")
	if err != nil {
		logs.Warn("/flowreplay/getById接口 参数异常, err: %v", err)
		c.ErrorJson(-1, "参数异常", nil)
	}
	logs.Info("请求参数: id=%v", id)
	flowReplayMongo := models.FlowReplayMongo{}
	flowReplay, err := flowReplayMongo.QueryById(id)
	if err != nil {
		c.ErrorJson(-1, "服务查询数据异常", nil)
	}
	c.SuccessJson(flowReplay)
}

func (c *FlowReplayController) update() {
	userId, _ := c.GetSecureCookie(constant.CookieSecretKey, "user_id")
	flowreplay := &models.FlowReplayMongo{}
	err := c.ParseForm(flowreplay)
	if err != nil {
		logs.Warn("/flowreplay/update接口 参数异常, err: %v", err)
		c.ErrorJson(-1, "参数异常", nil)
	}
	logs.Info("请求参数：%v", flowreplay)

	//验证服务名 唯一性
	serviceName := flowreplay.ServiceName
	temp, err := flowreplay.QueryByName(serviceName)
	if err != nil {
		logs.Error("流量回放添加时, 验证serviceName唯一性时报错")
	}
	if temp != nil && temp.Id != flowreplay.Id {
		c.ErrorJson(-1, "存在服务名相同的流量", nil)
	}
	flowreplay.UpdateBy = userId
	err = flowreplay.Update(*flowreplay)
	if err != nil {
		c.ErrorJson(-1, "服务更新数据异常", nil)
	}
	c.Redirect("/flowreplay/index", http.StatusFound)

}

type RemoveParam struct {
	Id int64 `form:"id" json:"id"`
}

// 删除
func (c *FlowReplayController) remove() {
	param := &RemoveParam{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, param)
	if err != nil {
		logs.Warn("/flowreplay/remove接口 参数异常, err: %v", err)
		c.ErrorJson(-1, "参数异常", nil)
	}
	logs.Info("请求参数: id=%v", param)
	flowreplayMongo := models.FlowReplayMongo{}
	err = flowreplayMongo.Delete(param.Id)
	if err != nil {
		c.ErrorJson(-1, "服务删除数据异常", nil)
	}
	c.SuccessJson(nil)
}

// 流量回放
type ReplayParam struct {
	Id          int64   `form:"id" json:"id"`
	ReplayTimes float64 `form:"replay_times" json:"replay_times"`
	TargetHost  string  `form:"target_host" json:"target_host"`
}

func (c *FlowReplayController) Replay() {
	param := &ReplayParam{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, param)
	if err != nil {
		logs.Error("/flowreplay/replay接口 参数异常, err: %v", err)
		c.ErrorJson(-1, "参数异常", nil)
	}
	logs.Info("请求参数: id=%v", param)

	flowreplayMongo := models.FlowReplayMongo{}
	flowreplay, err := flowreplayMongo.QueryById(param.Id)
	if err != nil {
		logs.Error("执行流量回放时, 查询指定回放报错")
		c.ErrorJson(-1, "查询指定回放报错", nil)
	}
	//回放文件名称
	flowFileName := flowreplay.FlowFile
	//回放路径
	flowFileName = uploadDir + "/" + flowFileName
	//机器
	filePath := flowreplay.FlowTargetHost
	//回放频率
	replayTimes := flowreplay.ReplayTimes
	//回放并发数
	replayurl := flowreplay.ReplayUri

	//默认不循环
	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("gor --input-file '%s|%v' --output-http=%s --stats --output-http-stats --output-http-timeout 1s --output-http-workers %s &", flowFileName, replayTimes, filePath, replayurl))

	//execCommand := "./gor --input-file \"./rankingmm.gor|1000%\" --output-http=\"http://172.16.1.22:8766\" --stats --output-http-stats --output-http-timeout 1s  --output-http-workers 1000"
	//execCommand := "./gor --input-file \"./"+filePath+"|"+strconv.Itoa(int(replayTimes*100))+"%\" --output-http=\"http://"+targetHost+"\" --stats --output-http-stats --output-http-timeout 1s  --output-http-workers 1000"
	//execCommand := "   "
	//cmd := exec.Command("/bin/bash", "-c", "gor --input-file '/Users/sunzhiying/rankingmm.gor|1%' --output-http=http://172.16.1.22:8766 --stats --output-http-stats --output-http-timeout 1s --output-http-workers 1000")
	//cmd := exec.Command("/bin/sh", "-c", "gor --input-file '/Users/sunzhiying/rankingmm.gor|1%' --output-http=http://172.16.1.22:8766 --stats --output-http-stats --output-http-timeout 1s --output-http-workers 1000")
	//fmt.Sprintf("gor --input-file '%s|%v' --output-http=%s --stats --output-http-stats --output-http-timeout 1s --output-http-workers 1000", flowFileName, replayTimes, filePath)

	//循环
	//cmd1 := exec.Command("/bin/bash", "-c", fmt.Sprintf("gor --input-file '%s|%v' --output-http=%s --stats --output-http-stats --output-http-timeout 1s --input-file-loop --output-http-workers %s", flowFileName, replayTimes, filePath,replayurl))
	//不循环
	//cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("gor --input-file '%s|%v' --output-http=%s --stats --output-http-stats --output-http-timeout 1s --output-http-workers %s", flowFileName, replayTimes, filePath,replayurl))
	//结束进程
	//kill := exec.Command("/bin/bash", "-c", fmt.Sprintf("ps aux | grep '%s' | cut -c 18-22 | xargs kill -9", flowFileName))
	fmt.Println(cmd)
	go func() {
		//cmd := exec.Command("gor", "--input-file", "'/Users/xueyibing/Desktop/小川文件夹/rankingmm.gor|100%'", "--output-http=http://172.16.1.22:8766", "--stats", "--output-http-stats", "--output-http-timeout", "1s", "--output-http-workers", "1000")
		body, err := cmd.CombinedOutput()
		//杀掉进程
		//body, err := kill.CombinedOutput()
		if err != nil {
			fmt.Printf("打印错误: %s", err.Error())
			os.Exit(1)
		} else {
			fmt.Printf("shell执行结果为~~~~~~~~/n: %s", string(body))
		}
	}()
	c.SuccessJson(nil)
}
func (c *FlowReplayController) ReplayCycle() {
	param := &ReplayParam{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, param)
	if err != nil {
		logs.Error("/flowreplay/replay接口 参数异常, err: %v", err)
		c.ErrorJson(-1, "参数异常", nil)
	}
	logs.Info("请求参数: id=%v", param)

	flowreplayMongo := models.FlowReplayMongo{}
	flowreplay, err := flowreplayMongo.QueryById(param.Id)
	if err != nil {
		logs.Error("执行流量回放时, 查询指定回放报错")
		c.ErrorJson(-1, "查询指定回放报错", nil)
	}
	//回放文件名称
	flowFileName := flowreplay.FlowFile
	//回放路径
	flowFileName = uploadDir + "/" + flowFileName
	//机器
	filePath := flowreplay.FlowTargetHost
	//回放频率
	replayTimes := flowreplay.ReplayTimes
	//回放并发数
	replayurl := flowreplay.ReplayUri

	//默认循环
	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("gor --input-file '%s|%v' --output-http=%s --stats --output-http-stats --output-http-timeout 1s --input-file-loop --output-http-workers %s &", flowFileName, replayTimes, filePath, replayurl))

	fmt.Println(cmd)
	go func() {
		//cmd := exec.Command("gor", "--input-file", "'/Users/xueyibing/Desktop/小川文件夹/rankingmm.gor|100%'", "--output-http=http://172.16.1.22:8766", "--stats", "--output-http-stats", "--output-http-timeout", "1s", "--output-http-workers", "1000")
		body, err := cmd.CombinedOutput()
		//杀掉进程
		//body, err := kill.CombinedOutput()
		if err != nil {
			fmt.Printf("打印错误: %s", err.Error())
			os.Exit(1)
		} else {
			fmt.Printf("shell执行结果为~~~~~~~~/n: %s", string(body))
		}
	}()
	c.SuccessJson(nil)
}
func (c *FlowReplayController) Killreplay() {
	param := &ReplayParam{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, param)
	if err != nil {
		logs.Error("/flowreplay/replay接口 参数异常, err: %v", err)
		c.ErrorJson(-1, "参数异常", nil)
	}
	logs.Info("请求参数: id=%v", param)

	flowreplayMongo := models.FlowReplayMongo{}
	flowreplay, err := flowreplayMongo.QueryById(param.Id)
	if err != nil {
		logs.Error("执行流量回放时, 查询指定回放报错")
		c.ErrorJson(-1, "查询指定回放报错", nil)
	}
	//回放文件名称
	flowFileName := flowreplay.FlowFile
	//回放路径
	flowFileName = uploadDir + "/" + flowFileName
	//本地mac结束进程
	//kill := exec.Command("/bin/bash", "-c", fmt.Sprintf("ps aux | grep '%s' | cut -c 18-22 | xargs kill -9", flowFileName))
	//服务linux结束进程
	kill := exec.Command("/bin/bash", "-c", fmt.Sprintf("ps aux | grep '%s' | cut -c 9-15 | xargs kill -9", flowFileName))

	fmt.Println(kill)
	go func() {
		//cmd := exec.Command("gor", "--input-file", "'/Users/xueyibing/Desktop/小川文件夹/rankingmm.gor|100%'", "--output-http=http://172.16.1.22:8766", "--stats", "--output-http-stats", "--output-http-timeout", "1s", "--output-http-workers", "1000")
		//杀掉进程
		body, err := kill.CombinedOutput()
		if err != nil {
			fmt.Printf("打印错误: %s", err.Error())
			//os.Exit(1)
		} else {
			fmt.Printf("shell执行结果为~~~~~~~~/n: %s", string(body))
		}
	}()
	//c.Redirect("/flowreplay/index", http.StatusFound)
	c.SuccessJson(nil)
}

func (c *FlowReplayController) collectFlowFile() {
	//获取上传的文件
	f, h, _ := c.GetFile("file")
	fileName := h.Filename
	//创建目录
	err := os.MkdirAll(uploadDir, 777)
	if err != nil {
		logs.Error("创建文件目录报错, err: ", err)
		c.ErrorJson(-1, "保存文件失败", nil)
	}

	fpath := uploadDir + "/" + fileName
	defer f.Close() //关闭上传的文件，不然的话会出现临时文件不能清除的情况

	err = c.SaveToFile("file", fpath)
	if err != nil {
		logs.Error("保存文件报错,  err: ", err)
		c.ErrorJson(-1, "保存文件失败", nil)
	}
	c.SuccessJson(nil)
}

func (c *FlowReplayController) showAllFlowFiles() {
	fileNames := []string{}
	files, err := ioutil.ReadDir(uploadDir)
	if err != nil {
		logs.Error("获取文件夹下文件列表时报错, err: ", err)
		c.ErrorJson(-1, "获取文件夹下文件列表时报错", nil)
	}
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}
	c.SuccessJson(fileNames)
}

func (c *FlowReplayController) pressure_Measurement() {
	times, err := c.GetInt64("test_times")
	if err != nil {
		logs.Info(err, "取得压测次数报错")
	}
	concurrent, err := c.GetInt64("concurrent", 10)
	if err != nil {
		logs.Info(err, "取得并发数错误")
	}
	RequestMode := c.GetString("request_mode")
	URL := c.GetString("url")
	ServiceName := c.GetString("service_name")
	Apiname := c.GetString("api_name")
	Args := c.GetString("args")
	cmd := exec.Command("/bin/bash", "-c", fmt.Sprintf("hey -n %v -c %v -t 3000 -m %s -T \"application/x-www-form-urlencoded\" %s/%s/%s -d '%s' ", times, concurrent, RequestMode, URL, ServiceName, Apiname, Args))
	fmt.Println(cmd)

	out, err := cmd.CombinedOutput()
	if err != nil {
		logs.Info("请求接口错误 err：", err)
	}
	fmt.Println(string(out))
	stringOut := string(out)
	//解析数据并封装
	stringList := strings.Split(stringOut, "\n")
	fmt.Printf(stringList[0])
	data := make(map[string]interface{})
	data["Total"] = stringList[2]
	data["Slowest"] = stringList[3]
	data["Fastest"] = stringList[4]
	data["Average"] = stringList[5]
	data["Requests/sec"] = stringList[6]
	data["histogram1"] = stringList[12] //响应直方图
	data["histogram2"] = stringList[13]
	data["histogram3"] = stringList[14]
	data["histogram4"] = stringList[15]
	data["histogram5"] = stringList[16]
	data["histogram6"] = stringList[17]
	data["histogram7"] = stringList[18]
	data["histogram8"] = stringList[19]
	data["histogram9"] = stringList[20]
	data["histogram10"] = stringList[21]
	data["histogram11"] = stringList[22]
	data["distribution1"] = stringList[26] //http请求时延分布
	data["distribution2"] = stringList[27]
	data["distribution3"] = stringList[28]
	data["distribution4"] = stringList[29]
	data["distribution5"] = stringList[30]
	data["distribution6"] = stringList[31]
	data["distribution7"] = stringList[32]
	data["status"] = stringList[42]

	c.Data["resp"] = data
	c.TplName = "show_testdata.html"
}

//Summary:
//Total:	0.0554 secs
//Slowest:	0.0094 secs
//Fastest:	0.0041 secs
//Average:	0.0055 secs
//Requests/sec:	180.4073
//
//Total data:	22830 bytes
//Size/request:	2283 bytes
//
//Response time histogram:
//0.004 [1]	|■■■■■■■■■■■■■
//0.005 [3]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
//0.005 [1]	|■■■■■■■■■■■■■
//0.006 [2]	|■■■■■■■■■■■■■■■■■■■■■■■■■■■
//0.006 [1]	|■■■■■■■■■■■■■
//0.007 [0]	|
//0.007 [1]	|■■■■■■■■■■■■■
//0.008 [0]	|
//0.008 [0]	|
//0.009 [0]	|
//0.009 [1]	|■■■■■■■■■■■■■
//
//
//Latency distribution:
//10% in 0.0042 secs
//25% in 0.0045 secs
//50% in 0.0052 secs
//75% in 0.0070 secs
//90% in 0.0094 secs
//0% in 0.0000 secs
//0% in 0.0000 secs
//
//Details (average, fastest, slowest):
//DNS+dialup:	0.0001 secs, 0.0000 secs, 0.0005 secs
//DNS-lookup:	0.0000 secs, 0.0000 secs, 0.0000 secs
//req write:	0.0000 secs, 0.0000 secs, 0.0001 secs
//resp wait:	0.0027 secs, 0.0021 secs, 0.0035 secs
//resp read:	0.0001 secs, 0.0001 secs, 0.0002 secs
//
//Status code distribution:
//[200]	10 responses
