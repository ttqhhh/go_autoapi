package models

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
	"go_autoapi/db_proxy"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type AutoUser struct {
	Id       int64  `json:"id" bson:"_id"`
	UserName string `json:"user_name" bson:"user_name"`
	Email    string `json:"email" bson:"email"`
	Mobile   string `json:"mobile" bson:"mobile"`
	Business string `json:"business" bson:"business"`
	// omitempty 表示该字段为空时，不返回
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at"`
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
