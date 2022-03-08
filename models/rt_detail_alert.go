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
const (
	QUICK_INCREASE_RT_ALERT_REASON = "接口响应时间大幅度高出近两周平均值"
	SLOW_INCREASE_RT_ALERT_REASON  = "接口近期响应时间呈缓增趋势"
)

type RtDetailAlertMongo struct {
	Id             int64  `json:"id" bson:"_id"`
	Business       int    `json:"business" bson:"business"`
	ServiceCode    string `json:"service_code" bson:"service_code"`
	Uri            string `json:"uri" bson:"uri"`
	Type           int    `json:"type" bson:"type"`
	AvgRt          int    `json:"avg_rt", bson:"avg_rt"`
	ThresholdRt    int    `json:"threshold_rt" bson:"threshold_rt"`
	AvgThresholdRt int    `json:"avg_threshold_rt", bson:"avg_threshold_rt"` //历史平均响应时间
	Rt             int    `json:"rt" bson:"rt"`
	CreatedAt      string `json:"created_at" bson:"created_at"`
	Reason         string `json:"reason" bson:"reason"` //警报原因
}

func init() {
	db_proxy.InitMongoDB()
}
func (a *RtDetailAlertMongo) TableName() string {
	return "rt_detail_alert"
}

func (a *RtDetailAlertMongo) Insert(rtDetailAlert RtDetailAlertMongo) error {
	ms, db := db_proxy.Connect("auto_api", "rt_detail_alert")
	defer ms.Close()
	// 当没有Id主键时，进行Id赋值
	if rtDetailAlert.Id == 0 {
		r := utils.GetRedis()
		defer r.Close()
		id, err := r.Incr(constants.RT_DETAIL_PRIMARY_KEY).Result()
		if err != nil {
			logs.Error("保存接口响应详情时，从redis获取唯一主键报错，err: ", err)
			return err
		}
		rtDetailAlert.Id = id
	}
	return db.Insert(rtDetailAlert)
}

//通过id删除性能报告
func (a *RtDetailAlertMongo) DeleteById(id int64) error {
	ms, db := db_proxy.Connect("auto_api", "rt_detail_alert")
	defer ms.Close()
	date := bson.M{
		"_id": id,
	}
	err := db.Remove(date)
	if err != nil {
		logs.Info("根据id删除错误")
	}
	return err
}

//查询一周之前的性能监控报告
func (a *RtDetailAlertMongo) QueryResult() ([]RtDetailAlertMongo, error) {
	nowTimeStr := time.Now().AddDate(0, 0, -7).Format("2006-01-02 15:04:0")
	ms, db := db_proxy.Connect("auto_api", "rt_detail_alert")
	defer ms.Close()
	query := bson.M{
		"created_at": bson.M{
			"$lt": nowTimeStr,
		},
	}
	RtDetailAlertList := []RtDetailAlertMongo{}
	err := db.Find(query).All(&RtDetailAlertList)
	if err != nil {
		if err.Error() == "not found" {
			err = nil
			//return nil, nil
			return nil, nil
		}
		logs.Error("查询错误: %v", err)
	}
	return RtDetailAlertList, err
}

func (a *RtDetailAlertMongo) GetOneWeekAlertInfo(page int, limit int) ([]RtDetailAlertMongo, int, error) {
	query := bson.M{}

	timestamp := time.Now().Unix()
	queryCond := []interface{}{}
	//查询过去7天该点的响应时间，不包括今天
	for i := 0; i < 7; i++ {
		one := timestamp - 60*60*24
		oneDate := time.Unix(one, 0).Format(Time_format)[:11]
		queryOr := bson.M{}
		queryOr["created_at"] = bson.M{"$regex": oneDate}
		queryCond = append(queryCond, queryOr)
	}
	query["$or"] = queryCond

	queryList := []RtDetailAlertMongo{}
	ms, db := db_proxy.Connect("auto_api", "rt_detail_alert")
	defer ms.Close()
	//err := db.Find(query).All(&queryList)
	err := db.Find(query).Sort("-_id").Skip((page - 1) * limit).Limit(limit).All(&queryList)
	if err != nil {
		if err.Error() == "not found" {
			return nil, 0, nil
		}
		logs.Error("分页查询警报出错 ", err)
	}
	intol, err := db.Find(query).Count()
	if err != nil {
		logs.Info("取得全部警报出错")
	}
	return queryList, intol, err
}

//查一周内的警报记录
func (a *RtDetailAlertMongo) GetAllOneWeekAlertInfo(page int, limit int) ([]RtDetailAlertMongo, int, error) {
	nowTimeStr := time.Now().AddDate(0, 0, -7).Format("2006-01-02 15:04:0")
	query := bson.M{
		"created_at": bson.M{
			"$gt": nowTimeStr,
		},
	}
	queryList := []RtDetailAlertMongo{}
	ms, db := db_proxy.Connect("auto_api", "rt_detail_alert")
	defer ms.Close()
	//err := db.Find(query).All(&queryList)
	err := db.Find(query).Sort("-_id").Skip((page - 1) * limit).Limit(limit).All(&queryList)
	if err != nil {
		if err.Error() == "not found" {
			return nil, 0, nil
		}
		logs.Error("分页查询警报出错 ", err)
	}
	intol, err := db.Find(query).Count()
	if err != nil {
		logs.Info("查询一周内警报总数数量错误")
	}
	return queryList, intol, err
}

//通过id查询警报详情
func (a *RtDetailAlertMongo) GetOneAlertById(id int64) (RtDetailAlertMongo, error) {

	ms, db := db_proxy.Connect("auto_api", "rt_detail_alert")
	defer ms.Close()
	alert := RtDetailAlertMongo{}
	query := bson.M{"_id": id}

	err := db.Find(query).One(&alert)
	if err != nil {
		logs.Error("查询警报出错 ", err)
	}
	return alert, err
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
	ms, db := db_proxy.Connect("auto_api", "rt_detail_alert")
	defer ms.Close()
	err := db.Find(query).All(&queryList)
	if err != nil {
		if err.Error() == "not found" {
			return nil, nil
		}
		logs.Error("分页查询警报出错 ", err)
	}
	return queryList, err
}

func (a *RtDetailAlertMongo) GetByServiceAndUri(serviceCode string, uri string) (*RtDetailAlertMongo, error) {
	query := bson.M{"service_code": serviceCode, "uri": uri}
	rtDetailAlert := RtDetailAlertMongo{}
	ms, db := db_proxy.Connect("auto_api", "rt_detail_alert")
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
