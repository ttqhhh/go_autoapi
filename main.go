package main

import (
	"github.com/beego/beego/v2/core/logs"
	beego "github.com/beego/beego/v2/server/web"
	_ "go_autoapi/routers"
)

func init() {
	_ = logs.SetLogger(logs.AdapterFile, `{"filename":"/Users/xiaoqiang/data/log/project.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10,"color":true}`)
	logs.EnableFuncCallDepth(true)
	logs.SetLogFuncCallDepth(3)
}

func main() {
	beego.Run()
}
