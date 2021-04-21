package models

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
	"go_autoapi/db_proxy"
	"gopkg.in/mgo.v2/bson"
)

type AutoResult struct {
	Id     int64  `json:"id" bson:"_id"`
	RunId  string `json:"run_id" bson:"run_id"`
	CaseId int64  `json:"case_id" bson:"case_id"`
	Result int64  `json:"result" bson:"description"`
	Reason string `json:"method" bson:"method"`
	// omitempty 表示该字段为空时，不返回
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at"`
}

func init() {
	//orm.RegisterModel(new(AdMockCaseMongo))
	db_proxy.InitMongoDB()
	//ORM = db_proxy.GetOrmObject()
}
func (a *AutoResult) TableName() string {
	return "auto_result"
}

func (a *AutoResult) InsertResult(ar AutoResult) error {
	ms, db := db_proxy.Connect("auto_api", "auto_result")
	defer ms.Close()
	return db.Insert(ar)
}

func (a *AutoResult) GetResultByRunId(id int64) (ar []*AutoResult, err error) {
	fmt.Println(id)
	query := bson.M{"run_id": id}
	ms, db := db_proxy.Connect("auto_api", "auto_result")
	defer ms.Close()
	err = db.Find(query).All(&ar)
	fmt.Println(ar)
	if err != nil {
		logs.Error(1024, err)
	}
	return ar, err
}
