package models

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
	"go_autoapi/db_proxy"
	"gopkg.in/mgo.v2/bson"
)

type AutoCaseMongo struct {
	Id          int64  `json:"id" bson:"_id"`
	ApiName     string `json:"api_name" bson:"api_name"`
	CaseName    string `json:"case_name" bson:"case_name"`
	Description string `json:"description" bson:"description"`
	Method      string `json:"method" bson:"method"`
	// omitempty 表示该字段为空时，不返回
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at"`
}

type Ids struct {
	Name string `json:"name" bson:"name"`
	Id   int64  `json:"id" bson:"id"`
}

func init() {
	//orm.RegisterModel(new(AdMockCaseMongo))
	db_proxy.InitMongoDB()
	//ORM = db_proxy.GetOrmObject()
}
func (a *AutoCaseMongo) TableName() string {
	return "auto_case"
}

func (a *Ids) TableName() string {
	return "ids"
}

func (a *Ids) GetId(name string) int64 {
	ids := Ids{}
	query := bson.M{"name": name}
	ms, db := db_proxy.Connect("auto_api", "ids")
	defer ms.Close()
	err := db.Update(query, bson.M{"$inc": bson.M{"id": 1}})
	fmt.Println(err)
	_ = db.Find(query).One(&ids)
	return ids.Id
}

func GetId(name string) (id int64) {
	ids := Ids{}
	id = ids.GetId(name)
	return
}

func (a *Ids) GetCollectionLength(name string) int64 {
	ids := Ids{}
	query := bson.M{"name": name}
	ms, db := db_proxy.Connect("auto_api", "ids")
	defer ms.Close()
	if err := db.Find(query).One(&ids); err != nil {
		logs.Error("查询失败")
	}
	return ids.Id
}

func (a *AutoCaseMongo) InsertCase(acm AutoCaseMongo) error {
	ms, db := db_proxy.Connect("auto_api", "case")
	defer ms.Close()
	return db.Insert(acm)
}

func (a *AutoCaseMongo) GetCaseById(id int64) (AutoCaseMongo, error) {

	fmt.Println(id)
	query := bson.M{"_id": id}
	acm := AutoCaseMongo{}
	ms, db := db_proxy.Connect("auto_api", "case")
	defer ms.Close()
	err := db.Find(query).One(&acm)
	fmt.Println(acm)
	if err != nil {
		logs.Error(1024, err)
	}
	return acm, err
}

func (a *AutoCaseMongo) UpdateCaseById(id int64, acm AutoCaseMongo) (AutoCaseMongo, error) {
	fmt.Println(id)
	query := bson.M{"_id": id}
	ms, db := db_proxy.Connect("auto_api", "case")
	defer ms.Close()
	err := db.Update(query, acm)
	fmt.Println(acm)
	if err != nil {
		logs.Error(1024, err)
	}
	return acm, err
}
