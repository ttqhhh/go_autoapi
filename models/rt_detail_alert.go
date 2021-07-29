package models

import (
	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
	"go_autoapi/constants"
	"go_autoapi/db_proxy"
	"go_autoapi/utils"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	QUICK_INCREASE_RT_ALERT = iota
	SLOW_INCREASE_RT_ALERT
)

type RtDetailAlertMongo struct {
	Id          int64  `json:"id" bson:"_id"`
	Business	int 	`json:"business" bson:"business"`
	ServiceCode string `json:"service_code" bson:"service_code"`
	Uri         string `json:"uri" bson:"uri"`
	Type 		int `json:"type" bson:"type"`
	AvgRt       int `json:"avg_rt", bson:"avg_rt"`
	ThresholdRt int `json:"threshold_rt" bson:"threshold_rt"`
	Rt			int `json:"rt" bson:"rt"`
	CreatedAt string `json:"created_at" bson:"created_at"`
}

func init() {
	db_proxy.InitMongoDB()
}
func (a *RtDetailAlertMongo) TableName() string {
	return "rt_detail_alert"
}

func (a *RtDetailAlertMongo) Insert(rtDetailAlert RtDetailAlertMongo) error {
	ms, db := db_proxy.Connect("auto_api", "rt_detail")
	defer ms.Close()
	// 当没有Id主键时，进行Id赋值
	if rtDetailAlert.Id == 0 {
		r := utils.GetRedis()
		id, err := r.Incr(constants.RT_DETAIL_PRIMARY_KEY).Result()
		if err != nil {
			logs.Error("保存接口响应详情时，从redis获取唯一主键报错，err: ", err)
			return err
		}
		rtDetailAlert.Id = id
	}
	return db.Insert(rtDetailAlert)
}

func (a *RtDetailAlertMongo) GetOneWeekAlertInfo() ([]RtDetailAlertMongo, error) {
	query := bson.M{}
	timestamp := time.Now().Unix()
	queryCond := []interface{}{}
	// 查询过去7天该点的响应时间，不包括今天
	for i := 0; i < 7; i++ {
		one := timestamp - 60*60*24
		oneDate := time.Unix(one, 0).Format(Time_format)[:11]
		queryOr := bson.M{}
		queryOr["created_at"] = bson.M{"$regex": oneDate}
		queryCond = append(queryCond, queryOr)
	}
	query["$or"] = queryCond

	queryList := []RtDetailAlertMongo{}
	ms, db := db_proxy.Connect("auto_api", "rt_detail")
	defer ms.Close()
	err := db.Find(query).All(&queryList)
	if err != nil {
		if err.Error() == "not found" {
			return nil, nil
		}
		logs.Error("GetById未查询到RtDetail数据，err: ", err)
	}
	return queryList, err
}

func (a *RtDetailAlertMongo) SummaryLast2WeekAlert() ([]RtDetailAlertMongo, error) {
	query := bson.M{}
	timestamp := time.Now().Unix()
	queryCond := []interface{}{}
	// 查询过去14天该点的响应时间
	for i := 0; i < 14; i++ {
		one := timestamp - 60*60*24
		oneDate := time.Unix(one, 0).Format(Time_format)[:11]
		queryOr := bson.M{}
		queryOr["created_at"] = bson.M{"$regex": oneDate}
		queryCond = append(queryCond, queryOr)
	}
	query["$or"] = queryCond

	queryList := []RtDetailAlertMongo{}
	ms, db := db_proxy.Connect("auto_api", "rt_detail")
	defer ms.Close()
	err := db.Find(query).All(&queryList)
	if err != nil {
		if err.Error() == "not found" {
			return nil, nil
		}
		logs.Error("GetById未查询到RtDetail数据，err: ", err)
	}
	return queryList, err
}

func (a *RtDetailAlertMongo) GetByServiceAndUri(serviceCode string, uri string) (*RtDetailAlertMongo, error) {
	query := bson.M{"service_code": serviceCode, "uri": uri}
	rtDetailAlert := RtDetailAlertMongo{}
	ms, db := db_proxy.Connect("auto_api", "rt_detail")
	defer ms.Close()
	err := db.Find(query).One(&rtDetailAlert)
	if err != nil {
		if err.Error() == "not found" {
			return nil, nil
		}
		logs.Error("GetByServiceAndUri未查询到RtDetail数据，err: ", err)
	}
	return &rtDetailAlert, err
}