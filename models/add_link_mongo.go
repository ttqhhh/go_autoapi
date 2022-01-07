package models

import (
	"github.com/astaxie/beego/logs"
	_ "github.com/beego/beego/v2/server/web"
	"go_autoapi/db_proxy"
	"gopkg.in/mgo.v2/bson"
)

const (
	_db         = "auto_api" // 数据库
	_collection = "add_link" // 表
)

type AddLinkMongo struct {
	//Id   int64  `form:"id" json:"id" bson:"id"`
	Link    string `form:"link" json:"link" bson:"link"`          // 平台链接
	Name    string `form:"name"  json:"name" bson:"name"`         // 平台名称
	Imglink string `form:"imglink" json:"imglink" bson:"imglink"` // 图片链接
}

func (mongo *AddLinkMongo) TableName() string {
	return _collection
}

// 增
func (mongo *AddLinkMongo) Insert(AddLink AddLinkMongo) error {
	ml, db := db_proxy.Connect(_db, _collection)
	defer ml.Close()

	_, err := db.Count()
	if err != nil {
		logs.Error("Insert 错误: %v", err)
		return err
	}

	err = db.Insert(AddLink)
	if err != nil {
		logs.Error("Insert 错误: %v", err)
	}
	return err
}

// 查
func (mongo *AddLinkMongo) Lists() ([]AddLinkMongo, error) {
	ml, db := db_proxy.Connect(_db, _collection)
	defer ml.Close()

	var links []AddLinkMongo

	err := db.Find(nil).All(&links)

	if err != nil {
		if err.Error() == "not found" {
			err = nil
			return nil, nil
		}
		logs.Error("find 错误: %v", err)
	}

	return links, err
}

// 删
func (mongo *AddLinkMongo) Delete(name string) error {
	ml, db := db_proxy.Connect(_db, _collection)
	defer ml.Close()

	query := bson.M{"name": name}
	err := db.Remove(query)
	if err != nil {
		logs.Error("err: ", err)
	}
	return err
}
