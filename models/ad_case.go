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
	orm.RegisterModel(new(AdMockCase))
	db_proxy.InitDB()
	ORM = db_proxy.GetOrmObject()
}
func (a *AdMockCase) TableName() string {
	return "ad_case"
}

func (a *AdMockCase) QueryByUUid(uuid string) (AdMockCase, error) {
	amc := AdMockCase{Uuid: uuid}
	qs := ORM.QueryTable("ad_case")
	err := qs.Filter("uuid", uuid).One(&amc)
	//err := o.Read(&amc)
	if err != nil {
		fmt.Println(amc.Id)
		//return amc,nil
	}
	fmt.Println(uuid, amc.CaseName)
	return amc, err
}
