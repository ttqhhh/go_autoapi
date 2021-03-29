package main

import (
	beego "github.com/beego/beego/v2/server/web"
	_ "go_autoapi/routers"
)

func main() {
	beego.Run()
}
