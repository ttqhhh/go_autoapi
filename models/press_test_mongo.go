package models

import (
	"github.com/astaxie/beego/logs"
	"go_autoapi/db_proxy"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	presstest_collection = "press_test"
)

type PressTestMongo struct {
	Id          int64  `form:"id" json:"id" bson:"_id"`
	TestTimes   int64  `form:"test_times" json:"test_times" bson:"test_times"`       // 执行次数
	Concurrent  int8   `form:"concurrent" json:"concurrent" bson:"concurrent"`       // 并发数
	RequestMode string `form:"request_mode" json:"request_mode" bson:"request_mode"` //请求方式
	URL         string `form:"url" json:"url" bson:"url"`                            // 请求地址
	ServiceName string `form:"service_name" json:"service_name" bson:"service_name"` // 服务名
	ApiName     string `form:"api_name" json:"api_name"  bson:"api_name"`            // 接口
	Args        string `form:"args" json:"args" bson:"args"`                         // 参数
	CreateBy    string `form:"created_by" json:"created_by" bson:"created_by"`       // 创建人
	CreatedAt   string `form:"created_at" json:"created_at" bson:"created_at"`       // 创建时间
	Status      int8   `form:"status" json:"status"  bson:"status"`                  //状态 1=删除 0 =正常
}

func (mongo *PressTestMongo) TableName() string {
	return presstest_collection
}

// 增
func (mongo *PressTestMongo) Insert(presstest PressTestMongo) error {
	ms, db := db_proxy.Connect(db, presstest_collection)
	defer ms.Close()

	// id自增
	cnt, err := db.Count()
	if err != nil {
		logs.Error("Insert 错误: %v", err)
		return err
	}
	presstest.Id = int64(cnt) + 1
	// 处理添加时间字段
	presstest.CreatedAt = time.Now().Format(Time_format)
	// 新增时，默认status为0
	presstest.Status = 0
	err = db.Insert(presstest)
	if err != nil {
		logs.Error("Insert 错误: %v", err)
	}
	return err
}

// 查(根据id)
func (mongo *PressTestMongo) QueryById(id int64) (*PressTestMongo, error) {
	ms, db := db_proxy.Connect(db, presstest_collection)
	defer ms.Close()

	query := bson.M{"_id": id, "status": 0}
	presstest := PressTestMongo{}
	err := db.Find(query).One(&presstest)
	if err != nil {
		if err.Error() == "not found" {
			err = nil
			return nil, nil
		}
		logs.Error("QueryById 错误: %v", err)
	}
	return &presstest, err
}

// 改
func (mongo *PressTestMongo) Update(presstest PressTestMongo) error {
	ms, db := db_proxy.Connect(db, presstest_collection)
	defer ms.Close()

	// 处理更新时间字段

	data := bson.M{
		"$set": bson.M{
			"test_times":   presstest.TestTimes,
			"concurrent":   presstest.Concurrent,
			"request_mode": presstest.RequestMode,
			"url":          presstest.URL,
			"service_name": presstest.ServiceName,
			"api_name":     presstest.ApiName,
			"args":         presstest.Args,
		},
	}

	changeInfo, err := db.UpsertId(presstest.Id, data)
	if err != nil {
		logs.Error("Update 错误: %v", err)
	}
	logs.Info("upsert函数返回的响应为：%v", changeInfo)
	return err
}

// 删
func (mongo *PressTestMongo) Delete(id int64) error {
	ms, db := db_proxy.Connect(db, presstest_collection)
	defer ms.Close()

	// 处理更新时间字段
	data := bson.M{
		"$set": bson.M{
			"status": 1,
		},
	}
	changeInfo, err := db.UpsertId(id, data)
	if err != nil {
		logs.Error("Delete 错误: %v", err)
	}
	logs.Info("upsert函数返回的响应为：%v", changeInfo)
	return err
}

//
//// 分页查询
func (mongo *PressTestMongo) QueryByPage(serviceName string, pageNo int, pageSize int) ([]PressTestMongo, int, error) {
	ms, db := db_proxy.Connect(db, presstest_collection)
	defer ms.Close()

	// 查询分页数据
	query := bson.M{"status": 0}
	//if business != -1 {
	//	query["business"] = business
	//}
	if serviceName != "" {
		query["service_name"] = bson.M{"$regex": serviceName}
	}
	// 进行业务线筛选
	//if len(businesses) > 0 {
	//	//queryCond := []interface{}{bson.D{"business"}}
	//	queryCond := []interface{}{}
	//	for _, v := range businesses {
	//		queryCond = append(queryCond, bson.M{"business": v})
	//	}
	//	query["$or"] = queryCond
	//}

	serviceList := []PressTestMongo{}
	skip := (pageNo - 1) * pageSize
	err := db.Find(query).Sort("-_id").Skip(skip).Limit(pageSize).All(&serviceList)
	if err != nil {
		logs.Error("QueryByPage 错误: %v", err)
	}
	// 查询总共条数
	serviceTotalList := []PressTestMongo{}
	err = db.Find(query).All(&serviceTotalList)
	if err != nil {
		logs.Error("QueryByPage 错误: %v", err)
	}
	return serviceList, len(serviceTotalList), err
}
