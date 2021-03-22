package db_proxy

import (
	"github.com/astaxie/beego/orm"
	"time"
)

var ormObject orm.Ormer

func InitDB() {
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "francolin:123456@tcp(172.16.4.228:3306)/ad_mock?charset=utf8")
	//orm.RunSyncdb("default", false, true)
	orm.DefaultTimeLoc = time.UTC
}

func GetOrmObject() orm.Ormer {
	return ormObject
}
