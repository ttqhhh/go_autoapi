package models

import (
	"github.com/astaxie/beego/logs"
	"go_autoapi/constants"
	"go_autoapi/db_proxy"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	business_collection = "auto_business"
)

type AutoBusiness struct {
	Id           int64  `json:"id" bson:"_id"`
	BusinessName string `json:"business_name" bson:"business_name" valid:"Required"`
	//0：正常，1：删除
	Status int    `json:"status,omitempty"  bson:"status" valid:"Range(0, 1)"`
	Author string `json:"author,omitempty"  bson:"author"`
	// omitempty 表示该字段为空时，不返回
	CreatedAt string `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt string `json:"updated_at,omitempty" bson:"updated_at"`
}

func init() {
	//orm.RegisterModel(new(AdMockCaseMongo))
	db_proxy.InitMongoDB()
	//ORM = db_proxy.GetOrmObject()
}
func (a *AutoBusiness) TableName() string {
	return business_collection
}

//添加业务线
func (a *AutoBusiness) InsertBusiness(ab AutoBusiness) error {
	ms, db := db_proxy.Connect(db, business_collection)
	defer ms.Close()
	return db.Insert(ab)
}

//删除某条业务线
func (a *AutoBusiness) DeleteBusiness(ab AutoBusiness) (err error) {
	query := bson.M{"_id": ab.Id}
	ms, db := db_proxy.Connect(db, business_collection)
	defer ms.Close()
	err = db.Update(query, bson.M{"status": 1, "updated_at": time.Now().Format(constants.TimeFormat)})
	return
}

//获取所有业务线
func (a *AutoBusiness) GetBusinessList(offset, page int) (ab []*AutoBusiness, err error) {
	query := bson.M{"status": 0}
	ms, db := db_proxy.Connect(db, business_collection)
	defer ms.Close()
	err = db.Find(query).Select(bson.M{"_id": 1, "business_name": 1}).Skip(page * offset).Limit(offset).All(&ab)
	if err != nil {
		logs.Error("query business list error", err)
	}
	return ab, err
}

//根据名字获取所有业务线
func (a *AutoBusiness) GetBusinessByName(businessName string) (business []*AutoBusiness, err error) {
	query := bson.M{"business_name": businessName, "status": 0}
	ms, db := db_proxy.Connect(db, business_collection)
	defer ms.Close()
	err = db.Find(query).Select(bson.M{"_id": 1}).One(&business)
	if err != nil {
		logs.Error("get business error", err)
	}
	return business, err
}

//获取所有的业务线
func (a *AutoBusiness) GetAllBusiness() (business []*AutoBusiness, err error) {
	query := bson.M{"status": 0}
	ms, db := db_proxy.Connect(db, business_collection)
	defer ms.Close()
	err = db.Find(query).Select(bson.M{"_id": 1, "business_name": 1}).All(&business)
	if err != nil {
		logs.Error("get all business error", err)
	}
	return business, err
}
