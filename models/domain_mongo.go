package models

import (
	"github.com/astaxie/beego/logs"
	"go_autoapi/db_proxy"
	"gopkg.in/mgo.v2/bson"
)

const domain_colection = "domain"


type Domain struct {
	Id           int64  `json:"id" bson:"_id"`
	Business string `json:"business" bson:"business" form:"business"`
	DomainName   string `json:"domain_name" bson:"domain_name" form:"domain_name"`
	Author string `json:"author"  bson:"author" form:"author"`
	//0：正常，1：删除
	Status int    `json:"status"  bson:"status"`
	// omitempty 表示该字段为空时，不返回
	CreatedAt string `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt string `json:"updated_at,omitempty" bson:"updated_at"`
}

func init() {
	//orm.RegisterModel(new(AdMockCaseMongo))
	db_proxy.InitMongoDB()
	//ORM = db_proxy.GetOrmObject()
}
func (d *Domain) TableName() string {
	return domain_colection
}

// 增加
func (d *Domain) DomainInsert(ab Domain) error{
	ms, db := db_proxy.Connect(db, domain_colection)
	defer ms.Close()
	// 查询是否存在 存在则不进行插入操作
	query := bson.M{"business":ab.Business,"domain_name":ab.DomainName}
	err := db.Find(query).One(&ab)
	if err != nil{
		logs.Error("查询失败，数据可能不存在，即将执行插入")
	}
	return db.Insert(ab)
}

// 查询（通过)

func (d *Domain) GetBusinessByName(business string) (domains []*Domain, err error) {
	query := bson.M{"business": business, "status": 0}
	ms, db := db_proxy.Connect(db, domain_colection)
	defer ms.Close()
	err = db.Find(query).All(&domains)
	if err != nil {
		logs.Error("get business error", err)
	}
	return domains, err
}

