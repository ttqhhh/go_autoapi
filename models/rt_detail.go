package models

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
	"go_autoapi/db_proxy"
	"gopkg.in/mgo.v2/bson"
)

type RtDetailMongo struct {
	Id          int64  `json:"id" bson:"_id"`
	Service     string `json:"service" bson:"service"`
	Uri         string `json:"uri" bson:"uri"`
	AvgRt       int    `json:"description" bson:"description"`
	ThresholdRt int    `json:"threshold_rt" bson:"threshold_rt"`
	last1DayRt  string `json:"last_1_day_rt" bson:"last_1_day_rt"`
	last2DayRt  string `json:"last_2_day_rt" bson:"last_2_day_rt"`
	last3DayRt  string `json:"last_3_day_rt" bson:"last_3_day_rt"`
	last4DayRt  string `json:"last_4_day_rt" bson:"last_4_day_rt"`
	last5DayRt  string `json:"last_5_day_rt" bson:"last_5_day_rt"`
	last6DayRt  string `json:"last_6_day_rt" bson:"last_6_day_rt"`
	last7DayRt  string `json:"last_7_day_rt" bson:"last_7_day_rt"`
	last8DayRt  string `json:"last_8_day_rt" bson:"last_8_day_rt"`
	last9DayRt  string `json:"last_9_day_rt" bson:"last_9_day_rt"`
	last10DayRt string `json:"last_10_day_rt" bson:"last_10_day_rt"`
	last11DayRt string `json:"last_11_day_rt" bson:"last_11_day_rt"`
	last12DayRt string `json:"last_12_day_rt" bson:"last_12_day_rt"`
	last13DayRt string `json:"last_13_day_rt" bson:"last_13_day_rt"`
	last14DayRt string `json:"last_14_day_rt" bson:"last_14_day_rt"`
	last15DayRt string `json:"last_15_day_rt" bson:"last_15_day_rt"`
	// omitempty 表示该字段为空时，不返回
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at"`
}

func init() {
	//orm.RegisterModel(new(AdMockCaseMongo))
	db_proxy.InitMongoDB()
	//ORM = db_proxy.GetOrmObject()
}
func (a *RtDetailMongo) TableName() string {
	return "rt_detail"
}

func (a *RtDetailMongo) Insert(rtDetail RtDetailMongo) error {
	ms, db := db_proxy.Connect("auto_api", "rt_detail")
	defer ms.Close()
	return db.Insert(rtDetail)
}

func (a *RtDetailMongo) GetById(id int64) (RtDetailMongo, error) {

	fmt.Println(id)
	query := bson.M{"_id": id}
	rtDetail := RtDetailMongo{}
	ms, db := db_proxy.Connect("auto_api", "rt_detail")
	defer ms.Close()
	err := db.Find(query).One(&rtDetail)
	fmt.Println(rtDetail)
	if err != nil {
		logs.Error("GetCaseById获取AutoCase失败", err)
	}
	return rtDetail, err
}

func (a *RtDetailMongo) UpdateById(id int64, rtDetail RtDetailMongo) (RtDetailMongo, error) {
	fmt.Println(id)
	query := bson.M{"_id": id}
	ms, db := db_proxy.Connect("auto_api", "rt_detail")
	defer ms.Close()
	err := db.Update(query, rtDetail)
	fmt.Println(rtDetail)
	if err != nil {
		logs.Error("UpdateCaseById更新AutoCase失败", err)
	}
	return rtDetail, err
}
