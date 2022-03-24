package models

import (
	"github.com/astaxie/beego/logs"
	"go_autoapi/db_proxy"
	"gopkg.in/mgo.v2/bson"
)

const (
	publishMsgCollection = "publish_msg"
)

type PublishMsg struct {
	Id           int64  `json:"id" bson:"_id"`
	BusinessId 	 int8 	`json:"business_id" bson:"business_id"`
	BusinessName string `json:"business_name" bson:"business_name"`
	Kind 		 string `json:"kind"  bson:"kind"`
	User 		 string `json:"user"  bson:"user"`
	Project 	 string `json:"project" bson:"project"`
	OnlineTime   string `json:"online_time" bson:"online_time"`
	Status 		 int    `json:"status,omitempty"  bson:"status" valid:"Range(0, 1)"`
	CreatedAt 	 string `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt 	 string `json:"updated_at,omitempty" bson:"updated_at"`
}

func init() {
	//orm.RegisterModel(new(AdMockCaseMongo))
	db_proxy.InitMongoDB()
	//ORM = db_proxy.GetOrmObject()
}
func (a *PublishMsg) TableName() string {
	return publishMsgCollection
}

//增加服务发表信息
func (a *PublishMsg) InsertPubMsg(pm PublishMsg) error {
	ms, dbs := db_proxy.Connect(db, publishMsgCollection)
	defer ms.Close()
	cnt, err := dbs.Count()
	if err != nil {
		logs.Error("Insert 错误: %v", err)
		return err
	}
	pm.Id = int64(cnt) + 1
	errs := dbs.Insert(pm)
	return errs
}

//删除服务发布信息
func (a *PublishMsg) DeletePubMsg(ab PublishMsg) (err error) {
	query := bson.M{"_id": ab.Id}
	ms, db := db_proxy.Connect(db, publishMsgCollection)
	defer ms.Close()
	err = db.Remove(query)
	//err = db.Update(query, bson.M{"status": 1, "updated_at": time.Now().Format(constants.TimeFormat)})
	return
}

//获取所有发布信息
func (a *PublishMsg) GetAllPubMsg(offset, page int) (ab []*PublishMsg, err error) {
	query := bson.M{"status": 0}
	ms, db := db_proxy.Connect(db, publishMsgCollection)
	defer ms.Close()
	err = db.Find(query).Select(bson.M{"_id": 1, "business_name": 1}).Skip(page * offset).Limit(offset).All(&ab)
	if err != nil {
		logs.Error("query pub_msg list error", err)
	}
	return ab, err
}

//根据业务线获取指定发布信息
func (a *PublishMsg) GetPubMsgByBusiness(businessId int8) (result []*PublishMsg, err error) {
	query := bson.M{"business_id": businessId, "status": 0}
	ms, db := db_proxy.Connect(db, publishMsgCollection)
	defer ms.Close()
	err = db.Find(query).All(&result)
	if err != nil {
		logs.Error("get pub_msg error", err)
	}
	return result, err
}

//获取所有的业务线
func (a *PublishMsg) GetAllBusiness() (business []*PublishMsg, err error) {
	query := bson.M{"status": 0}
	ms, db := db_proxy.Connect(db, publishMsgCollection)
	defer ms.Close()
	err = db.Find(query).Select(bson.M{"_id": 1, "business_name": 1}).All(&business)
	if err != nil {
		logs.Error("get all business error", err)
	}
	return business, err
}
