package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/globalsign/mgo/bson"
	"go_autoapi/db_proxy"
)

type AllActiveApiMongo struct {
	Id           int64  `form:"id" json:"id" bson:"id"`
	BusinessCode int64  `form:"business_code" json:"business_code" bson:"business_code"`  // 平台链接
	BusinessName string `form:"business_name"  json:"business_name" bson:"business_name"` // 平台名称
	ApiName      string `form:"api_name" json:"api_name" bson:"api_name"`
	Use          int64  `form:"use" json:"use" bson:"use"` //1=废弃接口
}

func init() {
	db_proxy.InitMongoDB()
}
func (a *AllActiveApiMongo) TableName() string {
	return "all_active_api"
}

func (a *AllActiveApiMongo) Insert(acm AllActiveApiMongo) error {
	ms, db := db_proxy.Connect(db, "all_active_api")
	defer ms.Close()

	// id自增
	cnt, err := db.Count()
	if err != nil {
		logs.Error("Insert 错误: %v", err)
		return err
	}
	acm.Id = int64(cnt) + 1
	// 处理添加时间字段
	err = db.Insert(acm)
	if err != nil {
		logs.Error("Insert 错误: %v", err)
	}
	return err
}

//根据业务线查询所有活跃并需要统计的接口
func (a *AllActiveApiMongo) QueryByBusiness(businessCode int64) ([]AllActiveApiMongo, error) {
	ms, db := db_proxy.Connect(db, "all_active_api")
	defer ms.Close()

	query := bson.M{"business_code": businessCode, "use": 0}
	apiList := []AllActiveApiMongo{}
	err := db.Find(query).All(&apiList)
	if err != nil {
		if err.Error() == "not found" {
			err = nil
			//return nil, nil
			return nil, nil
		}
		logs.Error("QueryByBusiness 错误: %v", err)
	}
	return apiList, err
}
