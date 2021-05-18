package models

import (
	"github.com/astaxie/beego/logs"
	"go_autoapi/constants"
	"go_autoapi/db_proxy"
	"go_autoapi/utils"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const domainCollection = "domain"


type Domain struct {
	Id           	int64  `json:"id" bson:"_id"`
	Business 		int8 `json:"business" bson:"business" form:"business"`
	DomainName   	string `json:"domain_name" bson:"domain_name" form:"domain_name"`
	Author 			string `json:"author"  bson:"author" form:"author"`
	//0：正常，1：删除
	Status 			int    `json:"status"  bson:"status"`
	// omitempty 表示该字段为空时，不返回
	CreatedAt 		string `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt 		string `json:"updated_at,omitempty" bson:"updated_at"`
}

func init() {
	//orm.RegisterModel(new(AdMockCaseMongo))
	db_proxy.InitMongoDB()
	//ORM = db_proxy.GetOrmObject()
}
func (d *Domain) TableName() string {
	return domainCollection
}

// 增加
func (d *Domain) DomainInsert(ab Domain) error{
	ms, db := db_proxy.Connect(db, domainCollection)
	defer ms.Close()
	// 查询是否存在 存在则不进行插入操作
	query := bson.M{"business":ab.Business,"domain_name":ab.DomainName}
	err := db.Find(query).One(&ab)
	if err == nil{
		logs.Info("域名已经存在")
		return err
	}
	//接收
	now := time.Now().Format(constants.TimeFormat)
	ab.CreatedAt = now
	ab.UpdatedAt = now
	ab.Status = 0
	r := utils.GetRedis()
	id, err := r.Incr(constants.DO_MAIN_PRIMARY_KEY).Result()
	ab.Id = id
	if err != nil {
		logs.Error("auto_result新增数据时，从redis获取自增主键错误, err:", err)
	}
	if err = db.Insert(ab) ; err != nil{
		logs.Error("域名数据插入失败")
	}
	return err
}

// 查询（通过)

func (d *Domain) GetDomainByBusiness(business int8) (domains []*Domain, err error) {
	query := bson.M{"business": business, "status": 0}
	ms, db := db_proxy.Connect(db, domainCollection)
	defer ms.Close()
	err = db.Find(query).All(&domains)
	if err != nil {
		logs.Error("get business error", err)
	}
	return
}

