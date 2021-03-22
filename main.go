package main

import (
	beego "github.com/beego/beego/v2/server/web"
	"go_autoapi/db_proxy"
	_ "go_autoapi/routers"
)

//
func init() {
	db_proxy.InitDB()
}

func main() {
	beego.Run()
}
