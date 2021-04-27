package models

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"go_autoapi/constants"
	"go_autoapi/db_proxy"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	user_collection = "auto_user"
)

type AutoUser struct {
	Id       int64  `json:"id" bson:"_id"`
	UserName string `json:"user_name" bson:"user_name" valid:"Required"`
	Email    string `json:"email" bson:"email"`
	Mobile   string `json:"mobile" bson:"mobile"`
	// 0：最右，1：皮皮，2：海外，3：中东，4：妈妈，5：商业化
	Business int `json:"business" bson:"business" valid:"Range(0, 5)"`
	//0：正常，1：删除
	Status int `json:"status,omitempty"  bson:"status" valid:"Range(0, 1)"`
	// omitempty 表示该字段为空时，不返回
	CreatedAt string `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt string `json:"updated_at,omitempty" bson:"updated_at"`
}

func init() {
	//orm.RegisterModel(new(AdMockCaseMongo))
	db_proxy.InitMongoDB()
	//ORM = db_proxy.GetOrmObject()
}
func (a *AutoUser) TableName() string {
	return user_collection
}

//创建新用户
func (a *AutoUser) InsertUser(au AutoUser) error {
	ms, db := db_proxy.Connect(db, user_collection)
	defer ms.Close()
	return db.Insert(au)
}

//根据用户id获取用户信息
func (a *AutoUser) GetUserInfoById(id int64) (AutoUser, error) {
	query := bson.M{"_id": id, "status": 0}
	au := AutoUser{}
	ms, db := db_proxy.Connect(db, user_collection)
	defer ms.Close()
	err := db.Find(query).One(&au)
	if err != nil {
		logs.Error(1024, err)
	}
	return au, err
}

//根据用户名字获取用户信息
func (a *AutoUser) GetUserInfoByName(name string) (au AutoUser, err error) {
	query := bson.M{"user_name": name}
	ms, db := db_proxy.Connect(db, user_collection)
	defer ms.Close()
	err = db.Find(query).One(&au)
	if err != nil {
		logs.Error(err)
	}
	return au, err
}

//根据id更新用户相关手机信息
func (a *AutoUser) UpdateUserById(id int64, mobile string, business int) (err error) {
	query := bson.M{"_id": id}
	ms, db := db_proxy.Connect(db, user_collection)
	defer ms.Close()
	err = db.Update(query, bson.M{"$set": bson.M{"mobile": mobile, "business": business, "updated_at": time.Now().Format(constants.TimeFormat)}})
	if err != nil {
		logs.Error(1024, err)
	}
	return err
}

func (a *AutoUser) DeleteUserById(id int64) (err error) {
	fmt.Println(id)
	query := bson.M{"_id": id}
	ms, db := db_proxy.Connect(db, user_collection)
	defer ms.Close()
	err = db.Update(query, bson.M{"$set": bson.M{"status": 1, "updated_at": time.Now().Format(constants.TimeFormat)}})
	if err != nil {
		logs.Error("delete user failed", err)
	}
	return err
}

// 获取用户列表
func (a *AutoUser) GetUserList(offset, page int) (au []*AutoUser, err error) {
	query := bson.M{"status": 0}
	ms, db := db_proxy.Connect(db, user_collection)
	defer ms.Close()

	skip := (page - 1) * offset
	err = db.Find(query).Select(bson.M{"id": 1, "user_name": 1, "email": 1, "mobile": 1, "business": 1}).Skip(skip).Limit(offset).All(&au)
	if err != nil {
		logs.Error(1024, err)
	}
	return au, err
}

// 获取用户列表
func (a *AutoUser) GetActivateUserCount() (total int, err error) {
	query := bson.M{"status": 0}
	ms, db := db_proxy.Connect(db, user_collection)
	defer ms.Close()

	total, err = db.Find(query).Count()
	if err != nil {
		logs.Error(1024, err)
	}
	return total, err
}
