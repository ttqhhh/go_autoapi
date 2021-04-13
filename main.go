package main

import (
	"github.com/astaxie/beego/logs"
	beego "github.com/beego/beego/v2/server/web"
	_ "go_autoapi/routers"
)

func main() {
	logs.EnableFuncCallDepth(true)
	logs.Info("debug")
	beego.Run()
}
