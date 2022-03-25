package models

import (
	"github.com/astaxie/beego/logs"
	"github.com/globalsign/mgo/bson"
	"go_autoapi/db_proxy"
)

const NO_USING = 1
const USING = 0

type AllActiveApiMongo struct {
	Id           int64  `form:"id" json:"id" bson:"id"`
	BusinessCode int64  `form:"business_code" json:"business_code" bson:"business_code"`  // 平台链接
	BusinessName string `form:"business_name"  json:"business_name" bson:"business_name"` // 平台名称
	ApiName      string `form:"api_name" json:"api_name" bson:"api_name"`
	Use          int64  `form:"use" json:"use" bson:"use"`                   //1=废弃接口
	Calculate    int64  `form:"calculate" json:"calculate" bson:"calculate"` //是否被用于计算接口覆盖率 0 = 没有被统计 1= 已经被统计
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

	query := bson.M{"business_code": businessCode, "use": USING}
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

//
//根据业务线查询所有活跃接口的数量
func (a *AllActiveApiMongo) QueryAllCountByBusinessCount(businessCode int64) int {
	ms, db := db_proxy.Connect(db, "all_active_api")
	defer ms.Close()

	query := bson.M{"business_code": businessCode, "use": USING}
	count, err := db.Find(query).Count()
	if err != nil {
		logs.Error("查询出错err:", err)

	}
	return count
}

//判断接口是否存在与数据库
func (a *AllActiveApiMongo) NewApiIsInDatabase(api_name string, business int64) (AllActiveApiMongo, bool) {
	ms, db := db_proxy.Connect(db, "all_active_api")
	defer ms.Close()

	query := bson.M{"business_code": business, "api_name": api_name}
	api := AllActiveApiMongo{}
	err := db.Find(query).One(&api)
	if err != nil {
		if err.Error() == "not found" {
			err = nil
			//return nil, nil
			return api, false
		}
		logs.Error("QueryById 错误: %v", err)

	}
	return api, true
}

//获得所有废弃接口数量
func (a *AllActiveApiMongo) GetAllUnUseApiCount(business int64) int {
	ms, db := db_proxy.Connect(db, "all_active_api")
	defer ms.Close()

	query := bson.M{"business_code": business, "use": NO_USING}
	count, err := db.Find(query).Count()
	if err != nil {
		logs.Error("根据业务线获取废弃接口数量出错，err:", err)
	}

	return count
}

//更新数据的状态
func (a *AllActiveApiMongo) ChangeApiCalculate(id int64, acm AllActiveApiMongo) error {
	ms, db := db_proxy.Connect(db, "all_active_api")
	defer ms.Close()

	data := bson.M{
		"$set": bson.M{
			"calculate": 1,
		},
	}
	query := bson.M{
		"id": id,
	}
	err := db.Update(query, data)
	if err != nil {
		logs.Error("updata出错，err", err)
	}
	return err
}
