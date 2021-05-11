package libs

import (
	"github.com/astaxie/beego/logs"
	"go_autoapi/db_proxy"
	"go_autoapi/models"
	"gopkg.in/mgo.v2/bson"
)

const (
	db              = "auto_api"
	case_collection = "case"
)

func init() {
	//orm.RegisterModel(new(AdMockCaseMongo))
	db_proxy.InitMongoDB()
	//ORM = db_proxy.GetOrmObject()
}

func GetCasesByIds(ids []int64) (acms []*models.TestCaseMongo, err error) {
	ms, db := db_proxy.Connect(db, case_collection)
	defer ms.Close()
	query := bson.M{"_id": bson.M{"$in": ids}, "status": 0}
	err = db.Find(query).All(&acms)
	return
}

func GetCasesByServices(service []string)(ids []int64){
	var acm []models.TestCaseMongo
	ms, db := db_proxy.Connect(db, case_collection)
	defer ms.Close()
	query := bson.M{"service_name":bson.M{"$in": service}, "status":0}
	if err := db.Find(query).All(&acm); err != nil{
		logs.Error("通过service查询case失败", err)
	}
	for _, i := range acm{
		ids = append(ids, i.Id)
	}
	return
}
