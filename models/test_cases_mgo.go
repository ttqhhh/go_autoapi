package models

import (
	_ "github.com/go-sql-driver/mysql"
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

//func (t *TestCaseMongo) getALLCases(){
//	query := TestCaseMongo{}
//	//ms, db := db_proxy.FindAll("auto_api", "case", query)
//	if err := db.Find(); err !=nil{
//		logs.Error(1024, err)
//	}
//}


