package models

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
	"go_autoapi/constants"
	"go_autoapi/db_proxy"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	result_collection = "auto_result"
)

type AutoResult struct {
	Id       int64  `json:"id,omitempty" bson:"_id"`
	RunId    string `json:"run_id,omitempty" bson:"run_id"`
	CaseId   int64  `json:"case_id" bson:"case_id"`
	Result   int64  `json:"result" bson:"result"`
	Reason   string `json:"reason" bson:"reason"`
	Author   string `json:"author" bson:"author"`
	Response string `json:"response,omitempty" bson:"response"`
	// omitempty 表示该字段为空时，不返回
	CreatedAt string `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt string `json:"updated_at,omitempty" bson:"updated_at"`
}

func init() {
	//orm.RegisterModel(new(AdMockCaseMongo))
	db_proxy.InitMongoDB()
	//ORM = db_proxy.GetOrmObject()
}
func (a *AutoResult) TableName() string {
	return "auto_result"
}

func InsertResult(uuid string, case_id int64, reason string, author string, resp string) error {
	now := time.Now().Format(constants.TimeFormat)
	ar := AutoResult{}
	ar.Id = GetId("result")
	ar.RunId = uuid
	ar.CaseId = case_id
	ar.Result = 1
	ar.Reason = reason
	ar.Author = author
	ar.Response = resp
	ar.CreatedAt = now
	ar.UpdatedAt = now
	ms, db := db_proxy.Connect(db, result_collection)
	defer ms.Close()
	return db.Insert(ar)
}

func GetResultByRunId(id string) (ar []*AutoResult, err error) {
	fmt.Println(id)
	query := bson.M{"run_id": id}
	ms, db := db_proxy.Connect(db, result_collection)
	defer ms.Close()
	err = db.Find(query).Select(bson.M{"case_id": 1, "reason": 1, "author": 1, "response": 1, "created_at": 1}).All(&ar)
	fmt.Println(ar)
	if err != nil {
		logs.Error(1024, err)
	}
	return ar, err
}

func (a *AutoResult) GetAllResult(page, limit int) (ar []*AutoResult, totalCount int64, err error) {
	ms, db := db_proxy.Connect(db, result_collection)
	defer ms.Close()

	var query interface{} = nil
	// 查询分页列表数据
	err = db.Find(query).Select(bson.M{"_id": 1, "run_id": 1, "case_id": 1, "result": 1, "reason": 1, "author": 1, "created_at": 1, "update_at": 1}).Skip((page - 1) * limit).Sort("-_id").Limit(limit).All(&ar)
	fmt.Println(ar)
	if err != nil {
		logs.Error("查询分页列表数据报错, err: ", err)
		return nil, 0, err
	}
	// 查询数据总条数用于分页
	total, err := db.Find(query).Count()
	fmt.Println(ar)
	if err != nil {
		logs.Error("查询分页列表数据报错, err: ", err)
		return nil, 0, err
	}
	totalCount = int64(total)
	return
}
