package models

import (
	"github.com/astaxie/beego/logs"
	"go_autoapi/db_proxy"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	db_email          = "auto_api"
	collection_email  = "email"
	Time_format_email = ""
)

type EmailMongo struct {
	Id             int64  `form:"id" json:"id" bson:"_id"`
	EmailName      string `form:"email_name" json:"email_name" bson:"email_name"`                // 邮箱组名称
	Emailrecipient string `form:"email_recipient" json:"email_recipient" bson:"email_recipient"` // 邮箱组成员
	Business       int8   `form:"business" json:"business" bson:"business"`                      // 0：最右，1：皮皮，2：海外，3：中东，4：妈妈
	Status         int8   `form:"status" json:"status"  bson:"status"`                           // 0：正常，1：删除
	CreateBy       string `form:"create_by" json:"create_by" bson:"create_by"`                   // 添加人
	UpdateBy       string `form:"update_by" json:"update_by" bson:"update_by"`                   // 修改人
	//CreatedAt time.Time `form:"created_at" json:"created_at,omitempty" bson:"created_at"` // omitempty 表示该字段为空时，不返回
	//UpdatedAt time.Time `form:"updated_at" json:"updated_at" bson:"updated_at"`
	CreatedAt string `form:"created_at" json:"created_at" bson:"created_at"` // omitempty 表示该字段为空时，不返回
	UpdatedAt string `form:"updated_at" json:"updated_at" bson:"updated_at"`
}

func (mongo *EmailMongo) TableName() string {
	return collection_email
}

// 增
func (mongo *EmailMongo) Insert(email EmailMongo) error {
	ms, db := db_proxy.Connect(db_email, collection_email)
	defer ms.Close()

	// id自增
	cnt, err := db.Count()
	if err != nil {
		logs.Error("Insert 错误: %v", err)
		return err
	}
	email.Id = int64(cnt) + 1
	// 处理添加时间字段
	email.CreatedAt = time.Now().Format(Time_format_email)
	// 新增时，默认status为0
	email.Status = 0
	err = db.Insert(email)
	if err != nil {
		logs.Error("Insert 错误: %v", err)
	}
	return err
}

// 删
func (mongo *EmailMongo) Delete(id int64) error {
	ms, db := db_proxy.Connect(db_email, collection_email)
	defer ms.Close()

	// 处理更新时间字段
	data := bson.M{
		"$set": bson.M{
			"status":     1,
			"updated_at": time.Now().Format(Time_format_email),
		},
	}
	changeInfo, err := db.UpsertId(id, data)
	if err != nil {
		logs.Error("Delete 错误: %v", err)
	}
	logs.Info("upsert函数返回的响应为：%v", changeInfo)
	return err
}

// 改
func (mongo *EmailMongo) Update(email EmailMongo) error {
	ms, db := db_proxy.Connect(db_email, collection_email)
	defer ms.Close()

	// 处理更新时间字段
	email.UpdatedAt = time.Now().Format(Time_format_email)
	data := bson.M{
		"$set": bson.M{
			"email_name":      email.EmailName,
			"email_recipient": email.Emailrecipient,
			"business":        email.Business,
			"updated_at":      email.UpdatedAt,
			"update_by":       email.UpdateBy,
		},
	}
	changeInfo, err := db.UpsertId(email.Id, data)
	if err != nil {
		logs.Error("Update 错误: %v", err)
	}
	logs.Info("upsert函数返回的响应为：%v", changeInfo)
	return err
}

//查询所有邮件组
func (mongo *EmailMongo) QueryAll(businesses []int, emailName string, pageNo int, pageSize int) ([]EmailMongo, int64, []EmailMongo) {
	ms, db := db_proxy.Connect(db_email, collection_email)
	defer ms.Close()

	// 查询分页数据
	//query := bson.M{"status": 0}
	query := bson.M{}
	if len(businesses) > 0 {
		//queryCond := []interface{}{bson.D{"business"}}
		queryCond := []interface{}{}
		for _, v := range businesses {
			queryCond = append(queryCond, bson.M{"business": v})
		}
		query["$or"] = queryCond
	}
	if emailName != "" {
		query["email_name"] = bson.M{"$regex": emailName}
	}
	emailList := []EmailMongo{}
	skip := (pageNo - 1) * pageSize
	err := db.Find(query).Sort("-_id").Skip(skip).Limit(pageSize).All(&emailList)
	if err != nil {
		logs.Error("QueryByPage 错误: %v", err)
	}
	// 查询总共条数
	emailTotalList := []EmailMongo{}
	err = db.Find(query).All(&emailTotalList)
	if err != nil {
		logs.Error("QueryByPage 错误: %v", err)
	}
	return emailList, int64(len(emailTotalList)), nil
}

// 查(根据id)
func (mongo *EmailMongo) QueryById(id int64) (*EmailMongo, error) {
	ms, db := db_proxy.Connect(db_email, collection_email)
	defer ms.Close()

	query := bson.M{"_id": id, "status": 0}
	email := EmailMongo{}
	err := db.Find(query).One(&email)
	if err != nil {
		if err.Error() == "not found" {
			err = nil
			//return nil, nil
			return nil, nil
		}
		logs.Error("QueryById 错误: %v", err)
	}
	return &email, err
}

func (mongo *EmailMongo) Query() (*EmailMongo, error) {
	ms, db := db_proxy.Connect(db_email, collection_email)
	defer ms.Close()

	email := EmailMongo{}
	err := db.Find(nil).All(&email)
	if err != nil {
		if err.Error() == "not found" {
			err = nil
			//return nil, nil
			return nil, nil
		}
		logs.Error("查询错误: %v", err)
	}
	return &email, err
}
