package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"go_autoapi/db_proxy"
	"time"
)

var ORM orm.Ormer

type AdMockCase struct {
	Id         int       `json:"id"`
	CaseName   string    `json:"case_name"`
	CaseDesc   string    `json:"case_desc"`
	Wish       string    `json:"wish"`
	Uuid       string    `json:"uuid" orm:"uuid"`
	Location   int       `json:"location"`
	Source     int       `json:"source"`
	AppId      int       `json:"app_id"`
	DeviceType int       `json:"device_type"`
	Version    int       `json:"version"`
	Priority   int       `json:"priority"`
	Status     int       `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func init() {
	ORM = db_proxy.GetOrmObject()
	orm.RegisterModel(new(AdMockCase))
	////构造conn连接
	////用户名:密码@tcp(url地址)/数据库
	////conn := dbuser + ":" + dbpass + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?charset=utf8"
	//
	//orm.RegisterDataBase("default", "mysql", "francolin:123456@tcp(172.16.4.228:3306)/ad_mock?charset=utf8")
	////conn :="francolin:123456@172.16.4.228:3306/ad_mock?charset=utf8"
	////注册数据库连接
	////orm.RegisterDataBase("default", "mysql", conn)
	//fmt.Printf("数据库连接成功！%s\n")
}
func (a *AdMockCase) TableName() string {
	return "ad_case"
}

func (a *AdMockCase) QueryByUUid(uuid string) (AdMockCase, error) {
	amc := AdMockCase{Uuid: uuid}
	o := orm.NewOrm()
	qs := o.QueryTable("ad_case")
	err := qs.Filter("uuid", uuid).One(&amc)
	//err := o.Read(&amc)
	if err != nil {
		fmt.Println(amc.Id)
		//return amc,nil
	}
	fmt.Println(uuid, amc.CaseName)
	return amc, err
}
