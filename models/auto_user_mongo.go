package models

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
	"go_autoapi/db_proxy"
	"gopkg.in/mgo.v2/bson"
)

type AutoUser struct {
	Id       int64  `json:"id" bson:"_id"`
	UserName string `json:"user_name" bson:"user_name"`
	Email    string `json:"email" bson:"email"`
	Mobile   string `json:"mobile" bson:"mobile"`
	// 0：最右，1：皮皮，2：海外，3：中东，4：妈妈
	Business int `json:"business" bson:"business"`
	//0：正常，1：删除
	Status int `json:"status"  bson:"status"`
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
	return "auto_user"
}

func (a *AutoUser) InsertCase(au AutoUser) error {
	ms, db := db_proxy.Connect("auto_api", "auto_user")
	defer ms.Close()
	return db.Insert(au)
}

func (a *AutoUser) GetUserInfoById(id int64) (AutoUser, error) {
	fmt.Println(id)
	query := bson.M{"_id": id}
	au := AutoUser{}
	ms, db := db_proxy.Connect("auto_api", "auto_user")
	defer ms.Close()
	err := db.Find(query).One(&au)
	fmt.Println(au)
	if err != nil {
		logs.Error(1024, err)
	}
	return au, err
}

func (a *AutoUser) GetUserByName(name string) (au AutoUser, err error) {
	query := bson.M{"user_name": name}
	ms, db := db_proxy.Connect("auto_api", "auto_user")
	defer ms.Close()
	err = db.Find(query).One(&au)
	fmt.Println(au)
	if err != nil {
		logs.Error(59, err)
	}
	return au, err
}

func (a *AutoUser) UpdateUserById(id int64, au AutoUser) (AutoUser, error) {
	fmt.Println(id)
	query := bson.M{"_id": id}
	ms, db := db_proxy.Connect("auto_api", "auto_user")
	defer ms.Close()
	err := db.Update(query, au)
	fmt.Println(au)
	if err != nil {
		logs.Error(1024, err)
	}
	return au, err
}

func (a *AutoUser) GetUserList(offset, page int) (au []*AutoUser, err error) {
	query := bson.M{"status": 0}
	ms, db := db_proxy.Connect("auto_api", "auto_user")
	defer ms.Close()
	err = db.Find(query).Select(bson.M{"id": 1, "user_name": 1}).Skip(page * offset).Limit(offset).All(&au)
	fmt.Println(au)
	if err != nil {
		logs.Error(1024, err)
	}
	return au, err
}
