package db_proxy

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"go_autoapi/models"
	"time"
)

var ormObject orm.Ormer

func InitDB() {
	dbuser := beego.AppConfig.String("dbuser")
	dbpass := beego.AppConfig.String("dbpass")
	dbhost := beego.AppConfig.String("dbhost")
	dbport := beego.AppConfig.String("dbport")
	dbname := beego.AppConfig.String("dbname")
	orm.RegisterModel(new(models.AdMockCase))
	//构造conn连接
	//用户名:密码@tcp(url地址)/数据库
	conn := dbuser + ":" + dbpass + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?charset=utf8"

	orm.RegisterDriver("mysql", orm.DRMySQL)
	//orm.RegisterDataBase("default", "mysql", "francolin:123456@tcp(172.16.4.228:3306)/ad_mock?charset=utf8")
	orm.RegisterDataBase("default", "mysql", conn)
	//orm.RunSyncdb("default", false, true)
	orm.DefaultTimeLoc = time.UTC
}

func GetOrmObject() orm.Ormer {
	return ormObject
}
