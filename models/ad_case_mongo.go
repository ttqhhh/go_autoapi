package models

import (
	"fmt"
	"github.com/globalsign/mgo/bson"
	_ "github.com/go-sql-driver/mysql"
	"go_autoapi/db_proxy"
	"time"
)

type AdMockCaseMongo struct {
	Id         int       `json:"id" bson:"_id"`
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
	//orm.RegisterModel(new(AdMockCaseMongo))
	db_proxy.InitMongoDB()
	//ORM = db_proxy.GetOrmObject()
}
func (a *AdMockCaseMongo) TableName() string {
	return "ad_case"
}

func (a *AdMockCaseMongo) InsertAdCase(amcm AdMockCaseMongo) error {
	ms, db := db_proxy.Connect("ad_mock", "ad_case")
	defer ms.Close()
	return db.Insert(amcm)
}

func (a *AdMockCaseMongo) GetAdCase(id int64) (AdMockCaseMongo, error) {
	query := bson.M{"_id": id}
	amcm := AdMockCaseMongo{}
	ms, db := db_proxy.Connect("ad_mock", "ad_case")
	defer ms.Close()
	err := db.Find(query).One(&amcm)
	if err != nil {
		fmt.Println("Error")
	}
	return amcm, err
}
