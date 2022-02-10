package models

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/globalsign/mgo/bson"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"go_autoapi/db_proxy"
	"time"
)

const (
	ad_h5_collection = "ad_h5"
)

type AdH5DataMongo struct {
	Id           int64     `json:"id" bson:"_id"`
	DataName     string    `json:"data_name"`
	DataUrl      string    `json:"data_url"`
	Status       int       `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Business     string    `json:"business"`
	BusinessName string    `json:"business_name"` //zuiyou,pipi,
}

func init() {
	//orm.RegisterModel(new(AdMockCaseMongo))
	db_proxy.InitMongoDB()
	//ORM = db_proxy.GetOrmObject()
}
func (a *AdH5DataMongo) TableName() string {
	return "ad_h5"
}

func (a *AdH5DataMongo) InsertAdCase(acm AdH5DataMongo) error {
	ms, db := db_proxy.Connect("auto_api", "ad_h5")
	defer ms.Close()
	err := db.Insert(acm)
	if err != nil {
		logs.Error("Insert 错误: %v", err)
	}
	return err
}

func (a *AdH5DataMongo) GetAllData() (acm []AdH5DataMongo, err error) {
	query := bson.M{"status": 0}
	ms, db := db_proxy.Connect("auto_api", "ad_h5")
	defer ms.Close()
	cursor := db.Find(query)
	if cursor == nil {
		fmt.Println("no data!")
		return
	}

	err = cursor.All(&acm)
	if err != nil {
		err = errors.Wrap(err, "err")
		fmt.Println(err)
		return
	}
	return
}

//更新数据
func (a *AdH5DataMongo) UpdateDataById(id int64, data_name string, data_url string, business string, acm AdH5DataMongo) (AdH5DataMongo, error) {
	fmt.Println(id)
	query := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{
		"dataname": data_name, "dataurl": data_url, "business": business,
	}}
	ms, db := db_proxy.Connect("auto_api", "ad_h5")
	defer ms.Close()

	err := db.Update(query, update)
	fmt.Println(acm)
	if err != nil {
		logs.Error("UpdateDataById更新h5数据失败", err)
	}
	return acm, err
}

// 删
func (a *AdH5DataMongo) DeleteDataById(id int64) error {
	ms, db := db_proxy.Connect("auto_api", "ad_h5")
	defer ms.Close()
	err := db.Remove(bson.M{"_id": id})
	if err != nil {
		logs.Error("Delete 错误: %v", err)
	}
	logs.Info("upsert函数返回的响应为：%v")
	return err
}

func (a *AdH5DataMongo) GetDataByBusiness(business string) (acm []AdH5DataMongo, err error) {
	query := bson.M{"status": 0, "business": business}

	ms, db := db_proxy.Connect("auto_api", "ad_h5")
	defer ms.Close()

	cursor := db.Find(query)
	if cursor == nil {
		fmt.Println("no data!")
		return
	}

	err = cursor.All(&acm)
	if err != nil {
		err = errors.Wrap(err, "err")
		fmt.Println(err)
		return
	}
	return
}
func (a *AdH5DataMongo) GetDataById(id int64) (acm []AdH5DataMongo, err error) {
	query := bson.M{"_id": id}
	ms, db := db_proxy.Connect("auto_api", "ad_h5")
	defer ms.Close()
	cursor := db.Find(query)
	if cursor == nil {
		fmt.Println("no data!")
		return
	}
	err = cursor.All(&acm)
	if err != nil {
		err = errors.Wrap(err, "err")
		fmt.Println(err)
		return
	}
	return
}
