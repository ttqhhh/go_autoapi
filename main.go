package main

import (
	//"fmt"
	"github.com/beego/beego/v2/client/orm"
	beego "github.com/beego/beego/v2/server/web"
	"go_autoapi/db_proxy"
	_ "go_autoapi/routers"
)

//
func init() {
	db_proxy.InitDB()
}

func main() {
	orm.Debug = true
	beego.Run()
}
