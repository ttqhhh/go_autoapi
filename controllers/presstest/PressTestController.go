package presstest

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	constant "go_autoapi/constants"
	"go_autoapi/libs"
	"go_autoapi/models"
	"net/http"
	"os/exec"
	"strings"
)

type PressTestController struct {
	libs.BaseController
}

func (c *PressTestController) Get() {
	do := c.GetMethodName()
	switch do {
	case "index":
		c.index()
	case "getById":
		c.getById()
	case "getList":
		c.list()
	case "toResult":
		c.toResult()

	default:
		logs.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持", nil)
	}
}

func (c *PressTestController) Post() {
	do := c.GetMethodName()
	switch do {
	case "save":
		c.add()
	case "update":
		c.update()
	case "remove":
		c.remove()
	case "pressure_Measurement":
		c.pressureMeasurement()

	default:
		logs.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持", nil)
	}
}

func (c *PressTestController) index() {
	c.TplName = "press_test.html"
}
func (c *PressTestController) toResult() {
	c.TplName = "show_test_data.html"
}

func (c *PressTestController) list() {
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
	pressTestMongo := models.PressTestMongo{}
	services, total, err := pressTestMongo.QueryByPage(serviceName, pageNo, pageSize)
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

func (c *PressTestController) add() {
	userId, _ := c.GetSecureCookie(constant.CookieSecretKey, "user_id")
	id, err1 := c.GetInt64("id")
	logs.Info("请求参数: id=%v", id)
	if err1 != nil {
		presstest := &models.PressTestMongo{}
		err := c.ParseForm(presstest)
		if err != nil {
			logs.Warn("/presstest/add接口 参数异常, err: %v", err)
			c.ErrorJson(-1, "参数异常", nil)
		}
		logs.Info("请求参数：%v", presstest)
		presstest.CreateBy = userId
		err = presstest.Insert(*presstest)
		if err != nil {
			c.ErrorJson(-1, "服务添加数据异常", nil)
		}
		c.Redirect("/presstest/index", http.StatusFound)
	}

}

func (c *PressTestController) getById() {
	id, err := c.GetInt64("id")
	if err != nil {
		logs.Warn("/presstest/getById接口 参数异常, err: %v", err)
		c.ErrorJson(-1, "参数异常", nil)
	}
	logs.Info("请求参数: id=%v", id)
	pressTestMongo := models.PressTestMongo{}
	pressTest, err := pressTestMongo.QueryById(id)
	if err != nil {
		c.ErrorJson(-1, "服务查询数据异常", nil)
	}
	c.SuccessJson(pressTest)
}

func (c *PressTestController) update() {
	presstest := &models.PressTestMongo{}
	err := c.ParseForm(presstest)
	if err != nil {
		logs.Warn("/presstest/update接口 参数异常, err: %v", err)
		c.ErrorJson(-1, "参数异常", nil)
	}
	logs.Info("请求参数：%v", presstest)
	err = presstest.Update(*presstest)
	if err != nil {
		c.ErrorJson(-1, "服务更新数据异常", nil)
	}
	c.Redirect("/presstest/index", http.StatusFound)

}

type RemoveParam struct {
	Id int64 `form:"id" json:"id"`
}

// 删除
func (c *PressTestController) remove() {
	param := &RemoveParam{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, param)
	if err != nil {
		logs.Warn("/presstest/remove接口 参数异常, err: %v", err)
		c.ErrorJson(-1, "参数异常", nil)
	}
	logs.Info("请求参数: id=%v", param)
	pressTestMongo := models.PressTestMongo{}
	err = pressTestMongo.Delete(param.Id)
	if err != nil {
		c.ErrorJson(-1, "服务删除数据异常", nil)
	}
	c.SuccessJson(nil)
}

func (c *PressTestController) pressureMeasurement() {
	id, err := c.GetInt64("id", 1)
	if err != nil {
		logs.Warn("/presstest/pressure接口 参数异常, err: %v", err)
		c.ErrorJson(-1, "参数异常", nil)
	}
	pressTestModels := models.PressTestMongo{}
	pressTestModel, err := pressTestModels.QueryById(id)
	if err != nil {
		logs.Error("通过id取得数据出错", err)
	}
	times := pressTestModel.TestTimes
	concurrent := pressTestModel.Concurrent
	RequestMode := pressTestModel.RequestMode
	URL := pressTestModel.URL
	ServiceName := pressTestModel.ServiceName
	Apiname := pressTestModel.ApiName
	Args := pressTestModel.Args
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
	c.TplName = "show_test_data.html"
	fmt.Printf("跳转页面")
}
