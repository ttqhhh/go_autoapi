package db_proxy

import (
	"github.com/astaxie/beego/orm"
	"github.com/beego/beego/v2/server/web"
	"time"
)

var ormObject orm.Ormer

func InitDB() {
	dbuser, _ := web.AppConfig.String("dbuser")
	dbpass, _ := web.AppConfig.String("dbpass")
	dbhost, _ := web.AppConfig.String("dbhost")
	dbport, _ := web.AppConfig.String("dbport")
	dbname, _ := web.AppConfig.String("dbname")
	////orm.RegisterModel(new(models.AdMockCase))
	////构造conn连接
	////用户名:密码@tcp(url地址)/数据库
	conn := dbuser + ":" + dbpass + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?charset=utf8"

	orm.RegisterDriver("mysql", orm.DRMySQL)
	//orm.RegisterDataBase("default", "mysql", "francolin:123456@tcp(172.16.4.228:3306)/ad_mock?charset=utf8")
	orm.RegisterDataBase("default", "mysql", conn)
	orm.RunSyncdb("default", false, true)
	orm.DefaultTimeLoc = time.UTC
	ormObject = orm.NewOrm()
}

func GetOrmObject() orm.Ormer {
	return ormObject
}
