package tuijian

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	constant "go_autoapi/constants"
	"go_autoapi/libs"
	"go_autoapi/models"
	"math/rand"
	"os"
	"path"
	"time"
)

const uploadDir = "~/upload"

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
	case "run":
		//c.run()

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
	//userId, _ := c.GetSecureCookie(constant.CookieSecretKey, "user_id")
	serviceName := c.GetString("service_name")
	//business, err := c.GetInt8("business", -1)
	//if err != nil {
	//	logs.Warn("/service/page接口 参数异常, err: %v", err)
	//	c.ErrorJson(-1, "参数异常", nil)
	//}
	//businessCodeList := GetUserBusinessesList(userId)
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
	//result := make(map[string]interface{})
	//result["total"] = total
	//result["data"] = services

	//c.SuccessJson(result)
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
		//获取上传的文件
		f, h, _ := c.GetFile("flow_file")
		fileName := h.Filename
		ext := path.Ext(fileName)
		//验证后缀名是否符合要求
		var AllowExtMap map[string]bool = map[string]bool{
			".gor": true,
		}
		if _, ok := AllowExtMap[ext]; !ok {
			c.Ctx.WriteString("流量文件后缀名不符合要求")
			return
		}
		//创建目录
		uploadDir := "/Users/sunzhiying/upload"
		err := os.MkdirAll(uploadDir, 777)
		if err != nil {
			c.Ctx.WriteString(fmt.Sprintf("%v", err))
			return
		}
		//构造文件名称
		rand.Seed(time.Now().UnixNano())
		randNum := fmt.Sprintf("%d", rand.Intn(9999)+1000)
		hashName := md5.Sum([]byte(time.Now().Format("2006_01_02_15_04_05_") + randNum))

		fileName = fmt.Sprintf("%x", hashName) + fileName
		//this.Ctx.WriteString(  fileName )
		//fpath := uploadDir + h.Filename
		fpath := uploadDir + "/" + fileName
		defer f.Close() //关闭上传的文件，不然的话会出现临时文件不能清除的情况

		//c.SaveToFile("flow_file", fpath)
		err = c.SaveToFile("flow_file", fpath)
		if err != nil {
			//c.Ctx.WriteString( fmt.Sprintf("%v",err) )
			c.ErrorJson(-1, "保存文件失败", nil)
		}

		flowreplay := &models.FlowReplayMongo{}
		err = c.ParseForm(flowreplay)
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
		flowreplay.FlowFile = fileName
		err = flowreplay.Insert(*flowreplay)
		if err != nil {
			c.ErrorJson(-1, "服务添加数据异常", nil)
		}
		//c.SuccessJson(nil)
		c.TplName = "replay.html"
		//c.ErrorJson(-1, "参数异常", nil)
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
	//获取上传的文件
	f, h, _ := c.GetFile("flow_file")
	if f == nil {
		flowreplay := &models.FlowReplayMongo{}
		err := c.ParseForm(flowreplay)
		if err != nil {
			logs.Warn("/flowreplay/update接口 参数异常, err: %v", err)
			c.ErrorJson(-1, "参数异常", nil)
		}
		logs.Info("请求参数：%v", flowreplay)

		//验证服务名 唯一性
		//serviceName := flowreplay.ServiceName
		//flowreplay.FlowFile = fileName
		//temp, err := flowreplay.QueryByName(serviceName)
		//if err != nil {
		//	logs.Error("流量回放编辑时, 验证serviceName唯一性时报错")
		//}
		//if temp != nil {
		//	c.ErrorJson(-1, "存在服务名相同的流量", nil)
		//}

		flowreplay.UpdateBy = userId
		err = flowreplay.Update(*flowreplay)
		if err != nil {
			c.ErrorJson(-1, "服务更新数据异常", nil)
		}
		//c.SuccessJson(nil)
		c.TplName = "replay.html"
	}
	if f != nil {
		fileName := h.Filename
		ext := path.Ext(fileName)
		//验证后缀名是否符合要求
		var AllowExtMap map[string]bool = map[string]bool{
			".gor": true,
		}
		if _, ok := AllowExtMap[ext]; !ok {
			c.Ctx.WriteString("流量文件后缀名不符合要求")
			return
		}
		//创建目录
		uploadDir := "/Users/sunzhiying/upload"
		//uploadDir := "~/upload/" + time.Now().Format("2006/01/02/")
		err := os.MkdirAll(uploadDir, 777)
		if err != nil {
			c.Ctx.WriteString(fmt.Sprintf("%v", err))
			return
		}
		//构造文件名称
		rand.Seed(time.Now().UnixNano())
		randNum := fmt.Sprintf("%d", rand.Intn(9999)+1000)
		hashName := md5.Sum([]byte(time.Now().Format("2006_01_02_15_04_05_") + randNum))

		fileName = fmt.Sprintf("%x", hashName) + fileName
		//this.Ctx.WriteString(  fileName )

		fpath := uploadDir + "/" + fileName
		defer f.Close() //关闭上传的文件，不然的话会出现临时文件不能清除的情况
		err = c.SaveToFile("flow_file", fpath)
		if err != nil {
			//c.Ctx.WriteString( fmt.Sprintf("%v",err) )
			c.ErrorJson(-1, "保存文件失败", nil)
		}

		flowreplay := &models.FlowReplayMongo{}
		err = c.ParseForm(flowreplay)
		if err != nil {
			logs.Warn("/flowreplay/update接口 参数异常, err: %v", err)
			c.ErrorJson(-1, "参数异常", nil)
		}
		logs.Info("请求参数：%v", flowreplay)

		//验证服务名 唯一性
		//serviceName := flowreplay.ServiceName
		flowreplay.FlowFile = fileName
		//temp, err := flowreplay.QueryByName(serviceName)
		if err != nil {
			logs.Error("流量回放编辑时, 验证serviceName唯一性时报错")
		}
		//if temp != nil {
		//	c.ErrorJson(-1, "存在服务名相同的流量", nil)
		//}

		flowreplay.UpdateBy = userId
		err = flowreplay.Update(*flowreplay)
		if err != nil {
			c.ErrorJson(-1, "服务更新数据异常", nil)
		}
		c.SuccessJson(nil)
		c.TplName = "replay.html"
	}
}

//func (c *FlowReplayController) run() {
//	id, err := c.GetInt64("id")
//	if err != nil {
//		logs.Warn("/flowreplay/getById接口 参数异常, err: %v", err)
//		c.ErrorJson(-1, "参数异常", nil)
//	}
//	logs.Info("请求参数: id=%v", id)
//	flowReplayMongo := models.FlowReplayMongo{}
//	flowReplay, err := flowReplayMongo.QueryById(id)
//	if err != nil {
//		c.ErrorJson(-1, "服务查询数据异常", nil)
//	}
//	fmt.Println(flowReplay)
//	flowfile :=flowReplay.FlowFile
//	fth := flowReplay.FlowTargetHost
//	rt := flowReplay.ReplayTimes
//	fmt.Println(flowfile,fth,rt)
//
//	c.SuccessJson(flowReplay)
//}
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
func (c *FlowReplayController) replay() {
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
	// todo 回放文件名称
	flowFileName := flowreplay.FlowFile
	
	// 回放频率
	replayTimes := param.ReplayTimes
	if replayTimes == 0 {
		replayTimes = flowreplay.ReplayTimes
	}
	// 回放目标机器
	targetHost := param.TargetHost
	if targetHost == "" {
		targetHost = flowreplay.FlowTargetHost
	}

	c.SuccessJson(nil)
}