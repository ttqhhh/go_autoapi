package models

import (
	"github.com/astaxie/beego/logs"
	"go_autoapi/constants"
	"go_autoapi/db_proxy"
	"go_autoapi/utils"
	"gopkg.in/mgo.v2/bson"
	"time"
)

// 本常量代码块用于枚举RunReportMongo中的IsPass状态
const (
	RUNNING = iota // 执行中
	SUCCESS        // 成功
	FAIL           // 失败
	ALL
)

const (
	run_record_collection = "run_record"
)

type RunReportMongo struct {
	Id             int64  `form:"id" json:"id" bson:"_id"`
	Name           string `form:"name" json:"name" bson:"name"`
	ServiceName    string `form:"service_name" json:"service_name" bson:"service_name"` // 服务名
	Business       int8   `form:"business" json:"business" bson:"business"`
	TotalCases     int    `form:"total_cases" json:"total_cases" bson:"total_case"`                // 本次执行操作运行的case总条数
	TotalFailCases int    `form:"total_fail_cases" json:"total_fail_cases" bson:"total_fail_case"` // 本次执行操作运行失败的case总条数
	IsPass         int8   `form:"is_pass" json:"is_pass"  bson:"is_pass"`                          // 0：执行中，1：成功，2：失败
	RunId          string `form:"run_id" json:"run_id" bson:"run_id"`
	CreateBy       string `form:"create_by" json:"create_by" bson:"create_by"`    // 添加人
	UpdateBy       string `form:"update_by" json:"update_by" bson:"update_by"`    // 修改人
	CreatedAt      string `form:"created_at" json:"created_at" bson:"created_at"` // omitempty 表示该字段为空时，不返回
	UpdatedAt      string `form:"updated_at" json:"updated_at" bson:"updated_at"`
}

func (mongo *RunReportMongo) TableName() string {
	return run_record_collection
}

// 增
func (mongo *RunReportMongo) Insert(service RunReportMongo) (int64, error) {
	ms, db := db_proxy.Connect(db, run_record_collection)
	defer ms.Close()

	r := utils.GetRedis()
	id, err := r.Incr(constants.RUN_RECORD_PRIMARY_KEY).Result()
	if err != nil {
		logs.Error("新增运行记录时报错, err:", err)
		return -1, err
	}
	service.Id = id
	// 处理添加时间字段
	service.CreatedAt = time.Now().Format(Time_format)
	err = db.Insert(service)
	if err != nil {
		logs.Error("Insert 错误: %v", err)
	}
	return id, err
}

// 删
//func (mongo *RunReportMongo) Delete(id int64) error {
//	ms, db := db_proxy.Connect(db, run_record_collection)
//	defer ms.Close()
//
//	// 处理更新时间字段
//	data := bson.M{
//		"$set": bson.M{
//			"status":     1,
//			"updated_at": time.Now().Format(time_format),
//		},
//	}
//	changeInfo, err := db.UpsertId(id, data)
//	if err != nil {
//		logs.Error("Delete 错误: %v", err)
//	}
//	logs.Info("upsert函数返回的响应为：%v", changeInfo)
//	return err
//}

//改
func (mongo *RunReportMongo) UpdateIsPass(id int64, isPass int8, totalFailCase int64, username string) error {
	ms, db := db_proxy.Connect(db, run_record_collection)
	defer ms.Close()

	runReport := RunReportMongo{}

	// 处理更新时间字段
	runReport.UpdatedAt = time.Now().Format(Time_format)
	data := bson.M{
		"$set": bson.M{
			"is_pass":         isPass,
			"total_fail_case": totalFailCase,
			"updated_at":      time.Now().Format(Time_format),
			"update_by":       username,
		},
	}
	changeInfo, err := db.UpsertId(id, data)
	if err != nil {
		logs.Error("Update 错误: %v", err)
	}
	logs.Info("upsert函数返回的响应为：%v", changeInfo)
	return err
}

// todo 根据服务名模糊查询待实现
// 分页查询
func (mongo *RunReportMongo) QueryByPage(businesses []int, serviceName string, pageNo int, pageSize int, runReportStatus int, isInspection int) ([]RunReportMongo, int64, error) {
	ms, db := db_proxy.Connect(db, run_record_collection)
	defer ms.Close()

	// 查询分页数据
	//query := bson.M{"status": 0}
	query := bson.M{}
	if len(businesses) > 0 {
		//queryCond := []interface{}{bson.D{"business"}}
		queryCond := []interface{}{}
		for _, v := range businesses {
			queryCond = append(queryCond, bson.M{"business": v})
		}
		query["$or"] = queryCond
	}
	if serviceName != "" {
		query["service_name"] = bson.M{"$regex": serviceName}
	}
	// 根据是否成功进行数据过滤掉
	if runReportStatus == FAIL {
		query["is_pass"] = FAIL
	} else if runReportStatus == SUCCESS {
		query["is_pass"] = SUCCESS
	}
	// 线上巡检case筛选
	if isInspection == INSPECTION {
		query["create_by"] = bson.M{"$regex":"线上巡检"}
	} else {
		//  限制非【线上巡检】数据
		query["create_by"] = bson.M{"$ne": "线上巡检"}
	}
	runReportList := []RunReportMongo{}
	skip := (pageNo - 1) * pageSize
	err := db.Find(query).Sort("-_id").Skip(skip).Limit(pageSize).All(&runReportList)
	if err != nil {
		logs.Error("QueryByPage 错误: %v", err)
	}
	// 查询总共条数
	runReportTotalList := []RunReportMongo{}
	err = db.Find(query).All(&runReportTotalList)
	if err != nil {
		logs.Error("QueryByPage 错误: %v", err)
	}
	return runReportList, int64(len(runReportTotalList)), err
}

// 查(根据业务线)
//func (mongo *RunReportMongo) QueryByBusiness(business int8) ([]RunReportMongo, error) {
//	ms, db := db_proxy.Connect(db, run_record_collection)
//	defer ms.Close()
//
//	query := bson.M{"status": 0}
//	if business != -1 {
//		query["business"] = business
//	}
//	serviceList := []ServiceMongo{}
//	err := db.Find(query).All(&serviceList)
//	if err != nil {
//		logs.Error("QueryByBusiness 错误: %v", err)
//	}
//	return serviceList, err
//}

// 查(根据id)
func (mongo *RunReportMongo) QueryById(id int64) (*RunReportMongo, error) {
	ms, db := db_proxy.Connect(db, run_record_collection)
	defer ms.Close()

	query := bson.M{"_id": id}
	runReport := RunReportMongo{}
	err := db.Find(query).One(&runReport)
	if err != nil {
		if err.Error() == "not found" {
			err = nil
			//return nil, nil
			return nil, nil
		}
		logs.Error("QueryById 错误: %v", err)
	}
	return &runReport, err
}
