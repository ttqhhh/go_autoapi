package models

import (
	"errors"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	"github.com/globalsign/mgo"
	_ "github.com/go-sql-driver/mysql"
	"go_autoapi/db_proxy"
	"gopkg.in/mgo.v2/bson"
)

type CaseSetMongo struct {
	Id           int64  `form:"id" json:"id" bson:"_id"`
	CaseSetName  string `form:"case_set_name" json:"case_set_name" bson:"case_set_name"`
	BusinessName string `form:"business_name" json:"business_name" bson:"business_name"`
	BusinessCode string `form:"business_code" json:"business_code" bson:"business_code"`
	Parameter    string `form:"parameter" json:"parameter" bson:"parameter"` // 用于存放测试集合配置的公共参数
	Description  string `form:"description" json:"description" bson:"description"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	//zen
	Author string `form:"author" json:"author" bson:"author"`
	Status int64  `json:"status" bson:"status"`
}

func (t *CaseSetMongo) GetCaseSetByPage(page, limit int, business_code string, caseSetName string) (result []CaseSetMongo, totalCount int64, err error) {
	ms, c := db_proxy.Connect("auto_api", "case_set")
	defer ms.Close()
	var query = bson.M{"business_code": business_code, "status": status}
	if caseSetName != "" {
		query["case_set_name"] = bson.M{"$regex": caseSetName}
	}
	// 获取CaseSet列表
	err = c.Find(query).Sort("-_id").Skip((page - 1) * limit).Limit(limit).All(&result)
	if err != nil {
		logs.Error("查询分页列表数据报错, err: ", err)
		return nil, 0, err
	}
	// 获取全部CaseSet数量
	total, err := c.Find(query).Count()
	if err != nil {
		logs.Error("数据库查询指定业务线下case数量报错, err: ", err)
		return nil, 0, err
	}
	totalCount = int64(total)
	return
}

// 通过id获取指定case
func (t *CaseSetMongo) CaseSetById(id int64) (CaseSetMongo, error) {
	query := bson.M{"_id": id, "status": status}
	acm := CaseSetMongo{}
	ms, db := db_proxy.Connect("auto_api", "case_set")
	defer ms.Close()
	err := db.Find(query).One(&acm)
	fmt.Println(acm)
	if err != nil {
		if err == mgo.ErrNotFound {
			return acm, nil
		}
		logs.Error("根据Id查询CaseSet报错, err: ", err)
	}
	return acm, err
}

// 添加一条case
func (t *CaseSetMongo) AddCaseSet(acm CaseSetMongo) error {
	ms, db := db_proxy.Connect("auto_api", "case_set")
	defer ms.Close()

	query := bson.M{"case_set_name": acm.CaseSetName, "status": status}
	count, err := db.Find(query).Count()

	if err != nil {
		logs.Error("数据库查询SetCase报错, err: ", err)
		return err
	}
	if count > 0 {
		return errors.New("同名用例集已经存在，请更换其他名字")
	}

	err = db.Insert(acm)
	if err != nil {
		logs.Error("数据库插入测试用例集报错, err:", err)
	}
	return err
}

// 通过id修改case（全更新）

func (t *CaseSetMongo) UpdateCaseSet(id int64, acm CaseSetMongo) (CaseSetMongo, error) {
	query := bson.M{"_id": id}
	ms, db := db_proxy.Connect("auto_api", "case_set")
	defer ms.Close()
	acm.Status = 0
	err := db.Update(query, acm)

	if err != nil {
		logs.Error("数据库更新CaseSet报错，err: ", err)
	}
	return acm, err
}

// 修改status
func (t *CaseSetMongo) DelCaseSet(id int64) error {
	query := bson.M{"_id": id}
	ms, db := db_proxy.Connect("auto_api", "case_set")
	defer ms.Close()

	err := db.Update(query, bson.M{"$set": bson.M{"status": del_}})
	if err != nil {
		logs.Error("删除数据库测试用例集报错, err: ", err)
		err = errors.New("删除数据库测试用例集报错")
		return err
	}
	return nil
}
