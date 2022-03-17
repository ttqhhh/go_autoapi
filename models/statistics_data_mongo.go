package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/globalsign/mgo/bson"
	"go_autoapi/db_proxy"
)

type StatisticsMongo struct {
	Id                         int64   `form:"id" json:"_id"`
	BusinessName               string  `form:"business_name" json:"business_name" bson:"business_name"`
	AllApiCount                float64 `form:"all_api_count" json:"all_api_count" bson:"all_api_count"`
	NewApiConut                float64 `form:"new_api_count" json:"new_api_count" bson:"new_api_count"`
	AllCaseCount               int     `form:"all_case_count" json:"all_case_count" bson:"all_case_count"`
	NewCaseConut               int     `form:"new_case_count" json:"new_case_count" bson:"new_case_count"`
	AllApi                     int     `form:"all_api" json:"all_api" bson:"all_api"` //完成度
	DegreeOfCompletion         string  `form:"degree_of_completion" json:"degree_of_completion" bson:"degree_of_completion"`
	LastWeekDegreeOfCompletion string  `form:"last_week_degree_of_completion" json:"last_week_degree_of_completion"` //上周完成度
	UnUseApi                   int     `form:"un_use_api" json:"un_use_api"`                                         //废弃接口数
}

func (mongo *StatisticsMongo) TableName() string {
	return "statistics_data"
}

func (mongo *StatisticsMongo) Insert(statisticsData StatisticsMongo) error {
	ms, db := db_proxy.Connect(db, "statistics_data")
	defer ms.Close()

	// id自增
	cnt, err := db.Count()
	if err != nil {
		logs.Error("Insert 错误: %v", err)
		return err
	}
	statisticsData.Id = int64(cnt) + 1
	err = db.Insert(statisticsData)
	if err != nil {
		logs.Error("Insert 错误: %v", err)
	}
	return err
}

func (mongo *StatisticsMongo) QueryAll() ([]StatisticsMongo, int, error) {
	ms, db := db_proxy.Connect(db, "statistics_data")
	defer ms.Close()

	// 查询分页数据
	query := bson.M{}

	dataList := []StatisticsMongo{}
	err := db.Find(query).Sort("-_id").Limit(6).All(&dataList)
	if err != nil {
		logs.Error("查询出错")
	}
	return dataList, 7, err
}
