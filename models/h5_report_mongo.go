package models

import (
	"github.com/astaxie/beego/logs"
	"go_autoapi/db_proxy"
	"gopkg.in/mgo.v2/bson"
)

// 本常量代码块用于枚举RunReportMongo中的IsPass状态
//const (
//	RUNNING = iota // 执行中
//	SUCCESS        // 成功
//	FAIL           // 失败
//	ALL
//)

//const (
//	run_record_collection = "run_record"
//)

type H5RunReportMongo struct {
	Id           int64  `json:"id" bson:"_id"`
	Business     string `form:"business" json:"business" bson:"business"`
	BusinessName string `json:"business_name"` //zuiyou,pipi,
	DataName     string `json:"data_name"`
	DataUrl      string `json:"data_url"`
	ErrorCode    string `json:"error_code"`
	CreatedAt    string `form:"created_at" json:"created_at" bson:"created_at"` // omitempty 表示该字段为空时，不返回
	Status       string `json:"status"`
}

func init() {
	db_proxy.InitMongoDB()
}

const (
	h5_run_report = "h5_run_report"
)

func (mongo *H5RunReportMongo) TableName() string {
	return h5_run_report
}

//查询所有的报告
func (mongo *H5RunReportMongo) Query() ([]H5RunReportMongo, error) {
	ms, db := db_proxy.Connect("auto_api", "h5_run_report")
	defer ms.Close()
	//timeUnix := time.Now().Unix() //当前时间转换为时间戳
	//nowTime:=time.Unix(timeUnix,0).Format("2006-01-02 15:04:05") //时间戳转换为字符串
	//NOWTIME,_:=time.ParseInLocation("2006-01-02 15:04:05",nowTime,time.Local)
	query := bson.M{}
	runReportList := []H5RunReportMongo{}
	err := db.Find(query).All(&runReportList)
	if err != nil {
		if err.Error() == "not found" {
			err = nil
			//return nil, nil
			return nil, nil
		}
		logs.Error("查询错误: %v", err)
	}
	return runReportList, err
}

//增
func (mongo *H5RunReportMongo) Insert(acm H5RunReportMongo) error {
	ms, db := db_proxy.Connect("auto_api", h5_run_report)
	defer ms.Close()
	err := db.Insert(acm)
	if err != nil {
		logs.Error("Insert 错误: %v", err)
	}
	return err
}

//删(直接删除库中数据)
//func (mongo *H5RunReportMongo) Delete(id int64) error {
//	ms, db := db_proxy.Connect(db, h5_run_report)
//	defer ms.Close()
//	// 处理更新时间字段
//	data := bson.M{
//		"_id":id,
//	}
//	err := db.Remove(data)
//	if err != nil {
//		logs.Error("Delete 错误，err:", err)
//	}
//	return err
//}

//改
//func (mongo *H5RunReportMongo) UpdateIsPass(id int64, isPass int8, totalFailCase int64, username string) error {
//	ms, db := db_proxy.Connect(db, h5_run_report)
//	defer ms.Close()
//
//	runReport := H5RunReportMongo{}
//
//	// 处理更新时间字段
//	runReport.UpdatedAt = time.Now().Format(Time_format)
//	data := bson.M{
//		"$set": bson.M{
//			"is_pass":         isPass,
//			"total_fail_case": totalFailCase,
//			"updated_at":      time.Now().Format(Time_format),
//			"update_by":       username,
//		},
//	}
//	changeInfo, err := db.UpsertId(id, data)
//	if err != nil {
//		logs.Error("Update 错误: %v", err)
//	}
//	logs.Info("upsert函数返回的响应为：%v", changeInfo)
//	return err
//}

// todo 根据服务名模糊查询待实现
// 分页查询
//func (mongo *H5RunReportMongo) QueryByPage(businesses []int, serviceName string, pageNo int, pageSize int, runReportStatus int, isInspection int) ([]RunReportMongo, int64, error) {
//	ms, db := db_proxy.Connect(db, h5_run_report)
//	defer ms.Close()
//
//	// 查询分页数据
//	//query := bson.M{"status": 0}
//	query := bson.M{}
//	if len(businesses) > 0 {
//		//queryCond := []interface{}{bson.D{"business"}}
//		queryCond := []interface{}{}
//		for _, v := range businesses {
//			queryCond = append(queryCond, bson.M{"business": v})
//		}
//		query["$or"] = queryCond
//	}
//	if serviceName != "" {
//		query["service_name"] = bson.M{"$regex": serviceName}
//	}
//	// 根据是否成功进行数据过滤掉
//	if runReportStatus == FAIL {
//		query["is_pass"] = FAIL
//	} else if runReportStatus == SUCCESS {
//		query["is_pass"] = SUCCESS
//	}
//	// 线上巡检case筛选
//	if isInspection == INSPECTION {
//		query["create_by"] = bson.M{"$regex":"线上巡检"}
//	} else {
//		//  限制非【线上巡检】数据
//		query["create_by"] = bson.M{"$ne": "线上巡检"}
//	}
//	runReportList := []H5RunReportMongo{}
//	skip := (pageNo - 1) * pageSize
//	err := db.Find(query).Sort("-_id").Skip(skip).Limit(pageSize).All(&runReportList)
//	if err != nil {
//		logs.Error("QueryByPage 错误: %v", err)
//	}
//	// 查询总共条数
//	runReportTotalList := []H5RunReportMongo{}
//	err = db.Find(query).All(&runReportTotalList)
//	if err != nil {
//		logs.Error("QueryByPage 错误: %v", err)
//	}
//	return runReportList, int64(len(runReportTotalList)), err
//}

// 查(根据id)
//func (mongo *H5RunReportMongo) QueryById(id int64) (*H5RunReportMongo, error) {
//	ms, db := db_proxy.Connect(db, h5_run_report)
//	defer ms.Close()
//
//	query := bson.M{"_id": id}
//	runReport := H5RunReportMongo{}
//	err := db.Find(query).One(&runReport)
//	if err != nil {
//		if err.Error() == "not found" {
//			err = nil
//			//return nil, nil
//			return nil, nil
//		}
//		logs.Error("QueryById 错误: %v", err)
//	}
//	return &runReport, err
//}
