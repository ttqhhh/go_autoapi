package models

import (
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	_ "github.com/go-sql-driver/mysql"
	"go_autoapi/db_proxy"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type TestCaseMongo struct {
	Id          int64     `json:"id" bson:"_id"`
	ApiName     string    `json:"api_name" bson:"api_name"`
	CaseName    string    `json:"case_name" bson:"case_name"`
	Description string    `json:"description" bson:"description"`
	Method      string    `json:"method" bson:"method"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	//zen
	AppName 	string	   `json:"app_name"`
	ServiceName	string	   `json:"service_name"`
	ApiUrl 		string 	   `json:"api_url"`
	TestEnv 	string 	   `json:"test_env"`
	Mock		string		`json:"mock"`
	RequestMethod string	`json:"request_method"`
	Parameter   string		`json:"parameter"`
	Checkpoint	string		`json:"check_point"`
	Level		string		`json:"level"`
}

//db:操作的数据库
//collection:操作的文档(表)
//query:查询条件
//selector:需要过滤的数据(projection)
//result:查询到的结果

// 获取指定server下的所有case

func (t *TestCaseMongo) GetCasesByQuery(query interface{})(TestCaseMongo, error){
	//query := TestCaseMongo{}
	var acm = TestCaseMongo{}
	ms, c := db_proxy.Connect("auto_api", "case")
	defer ms.Close()
	err := c.Find(query).All(&acm)
	if err != nil {
		logs.Error(1024, err)
	}
	return acm, err
}

// 获取全部case
func (t *TestCaseMongo) GetAllCases()(TestCaseMongo, error){
	var acm = TestCaseMongo{}
	ms, c := db_proxy.Connect("auto_api", "case")
	defer ms.Close()
	err := c.Find(nil).All(&acm)
	if err != nil {
		logs.Error(1024, err)
	}
	return acm, err
}


// 通过id获取指定case

func (t *TestCaseMongo) GetOneCase(id int64) (TestCaseMongo, error) {

	fmt.Println(id)
	query := bson.M{"_id": id}
	acm := TestCaseMongo{}
	ms, db := db_proxy.Connect("auto_api", "case")
	defer ms.Close()
	err := db.Find(query).One(&acm)
	fmt.Println(acm)
	if err != nil {
		logs.Error(1024, err)
	}
	return acm, err
}

// 添加一条case

func (t *TestCaseMongo) AddCase(acm TestCaseMongo) error{
	ms, db := db_proxy.Connect("auto_api", "case")
	defer  ms.Close()
	err := db.Insert(acm)
	if err != nil {
		logs.Error(1024, err)
	}
	return err
}

// 通过id修改case

func (t *TestCaseMongo) UpdateCase(id int64, acm TestCaseMongo) (TestCaseMongo, error) {
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