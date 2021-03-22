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
	//
	//	//dbuser, _ := beego.AppConfig.String("dbuser")
	//	//dbpass, _ := beego.AppConfig.String("dbpass")
	//	//dbhost, _ := beego.AppConfig.String("dbhost")
	//	//dbport, _ := beego.AppConfig.String("dbport")
	//	//dbname, _ := beego.AppConfig.String("dbname")
	//	//orm.RegisterModel(new(models.AdMockCase))
	//	//orm.RegisterDriver("mysql", orm.DRMySQL)
	//	////构造conn连接
	//	////用户名:密码@tcp(url地址)/数据库
	//	////conn := dbuser + ":" + dbpass + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?charset=utf8"
	//	//
	//	//orm.RegisterDataBase("default", "mysql", "francolin:123456@tcp(172.16.4.228:3306)/ad_mock?charset=utf8")
	//	////conn :="francolin:123456@172.16.4.228:3306/ad_mock?charset=utf8"
	//	////注册数据库连接
	//	////orm.RegisterDataBase("default", "mysql", conn)
	//	//fmt.Printf("数据库连接成功！%s\n")
	//
}

func main() {
	orm.Debug = true
	beego.Run()
}
