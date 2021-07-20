package models

import (
	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
	"go_autoapi/constants"
	"go_autoapi/db_proxy"
	"go_autoapi/utils"
	"gopkg.in/mgo.v2/bson"
)

type RtDetailMongo struct {
	Id          int64  `json:"id" bson:"_id"`
	ServiceCode string `json:"service_code" bson:"service_code"`
	Uri         string `json:"uri" bson:"uri"`
	AvgRt       string `json:"avg_rt", bson:"avg_rt"`
	ThresholdRt string `json:"threshold_rt" bson:"threshold_rt"`
	Last0DayRt  string `json:"last_0_day_rt" bson:"last_0_day_rt"`
	Last1DayRt  string `json:"last_1_day_rt" bson:"last_1_day_rt"`
	Last2DayRt  string `json:"last_2_day_rt" bson:"last_2_day_rt"`
	Last3DayRt  string `json:"last_3_day_rt" bson:"last_3_day_rt"`
	Last4DayRt  string `json:"last_4_day_rt" bson:"last_4_day_rt"`
	Last5DayRt  string `json:"last_5_day_rt" bson:"last_5_day_rt"`
	Last6DayRt  string `json:"last_6_day_rt" bson:"last_6_day_rt"`
	Last7DayRt  string `json:"last_7_day_rt" bson:"last_7_day_rt"`
	Last8DayRt  string `json:"last_8_day_rt" bson:"last_8_day_rt"`
	Last9DayRt  string `json:"last_9_day_rt" bson:"last_9_day_rt"`
	Last10DayRt string `json:"last_10_day_rt" bson:"last_10_day_rt"`
	Last11DayRt string `json:"last_11_day_rt" bson:"last_11_day_rt"`
	Last12DayRt string `json:"last_12_day_rt" bson:"last_12_day_rt"`
	Last13DayRt string `json:"last_13_day_rt" bson:"last_13_day_rt"`
	Last14DayRt string `json:"last_14_day_rt" bson:"last_14_day_rt"`
	// omitempty 表示该字段为空时，不返回
	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at"`
}

func init() {
	db_proxy.InitMongoDB()
}
func (a *RtDetailMongo) TableName() string {
	return "rt_detail"
}

func (a *RtDetailMongo) Insert(rtDetail RtDetailMongo) error {
	ms, db := db_proxy.Connect("auto_api", "rt_detail")
	defer ms.Close()
	// 当没有Id主键时，进行Id赋值
	if rtDetail.Id == 0 {
		r := utils.GetRedis()
		id, err := r.Incr(constants.RT_DETAIL_PRIMARY_KEY).Result()
		if err != nil {
			logs.Error("保存接口响应详情时，从redis获取唯一主键报错，err: ", err)
			return err
		}
		rtDetail.Id = id
	}
	return db.Insert(rtDetail)
}

func (a *RtDetailMongo) GetById(id int64) (*RtDetailMongo, error) {
	query := bson.M{"_id": id}
	rtDetail := RtDetailMongo{}
	ms, db := db_proxy.Connect("auto_api", "rt_detail")
	defer ms.Close()
	err := db.Find(query).One(&rtDetail)
	if err != nil {
		if err.Error() == "not found" {
			return nil, nil
		}
		logs.Error("GetById未查询到RtDetail数据，err: ", err)
	}
	return &rtDetail, err
}

func (a *RtDetailMongo) GetByServiceAndUri(serviceCode string, uri string) (*RtDetailMongo, error) {
	query := bson.M{"service_code": serviceCode, "uri": uri}
	rtDetail := RtDetailMongo{}
	ms, db := db_proxy.Connect("auto_api", "rt_detail")
	defer ms.Close()
	err := db.Find(query).One(&rtDetail)
	if err != nil {
		if err.Error() == "not found" {
			return nil, nil
		}
		logs.Error("GetByServiceAndUri未查询到RtDetail数据，err: ", err)
	}
	return &rtDetail, err
}

func (a *RtDetailMongo) UpdateById(id int64, rtDetail RtDetailMongo) (RtDetailMongo, error) {
	query := bson.M{"_id": id}
	ms, db := db_proxy.Connect("auto_api", "rt_detail")
	defer ms.Close()
	err := db.Update(query, rtDetail)
	if err != nil {
		logs.Error("UpdateCaseById更新AutoCase失败", err)
	}
	return rtDetail, err
}
