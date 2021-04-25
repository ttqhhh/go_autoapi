package libs

import (
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
