package models

import (
	"github.com/beego/beego/v2/core/logs"
	"github.com/globalsign/mgo"
	_ "github.com/go-sql-driver/mysql"
	"go_autoapi/db_proxy"
	"gopkg.in/mgo.v2/bson"
)

type SetCaseMongo struct {
	Id          int64  `form:"id" json:"id" bson:"_id"`
	CaseSetId   int64  `form:"case_set_id" json:"case_set_id" bson:"case_set_id"`
	CaseName    string `form:"case_name" json:"case_name" bson:"case_name"`
	Description string `form:"description" json:"description" bson:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	//zen
	Author        string `form:"author" json:"author" bson:"author"`
	Domain        string `form:"domain" json:"domain" bson:"domain"`
	BusinessName  string `form:"business_name" json:"business_name" bson:"business_name"`
	BusinessCode  string `form:"business_code" json:"business_code" bson:"business_code"`
	ServiceId     int64  `form:"service_id" json:"service_id" bson:"service_id"`
	ServiceName   string `form:"service_name" json:"service_name" bson:"service_name"`
	ApiUrl        string `form:"api_url" json:"api_url" bson:"api_url"`
	RequestMethod string `form:"request_method" json:"request_method" bson:"request_method"`
	Parameter     string `form:"parameter" json:"parameter" bson:"parameter"`
	ExtractResp   string `form:"extract_resp" json:"extract_resp" bson:"extract_resp"` // 	用于从响应中提取值的配置
	Checkpoint    string `form:"check_point" json:"check_point" bson:"check_point"`
	SmokeResponse string `form:"smoke_response" json:"smoke_response,omitempty" bson:"smoke_response"`
	Order         int    `form:"order" json:"order" bson:"order"` // 顺序，用于执行和页面展示
	BeforeWait    int    `form:"before_wait" json:"before_wait" bson:"before_wait"` // 前置等待时间
	AfterWait     int    `form:"after_wait" json:"after_wait" bson:"after_wait"`	//	后置等待时间
	Status        int64  `json:"status" bson:"status"`
}

// 通过id获取指定case

func (t *SetCaseMongo) GetSetCaseById(id int64) (SetCaseMongo, error) {
	query := bson.M{"_id": id, "status": status}
	acm := SetCaseMongo{}
	ms, db := db_proxy.Connect("auto_api", "set_case")
	defer ms.Close()
	err := db.Find(query).One(&acm)

	if err != nil {
		if err != mgo.ErrNotFound {
			return acm, nil
		}
		logs.Error("数据库通过指定Id查询SetCase失败, err: ", err)
	}
	return acm, err
}

// 添加一条case
func (t *SetCaseMongo) AddSetCase(acm SetCaseMongo) error {
	ms, db := db_proxy.Connect("auto_api", "set_case")
	defer ms.Close()
	err := db.Insert(acm)
	if err != nil {
		logs.Error("数据库插入测试用例报错, err: ", err)
	}
	return err
}

// 通过id修改case（全更新）
func (t *SetCaseMongo) UpdateSetCase(id int64, acm SetCaseMongo) (SetCaseMongo, error) {
	query := bson.M{"_id": id}
	ms, db := db_proxy.Connect("auto_api", "set_case")
	defer ms.Close()

	acm.Status = 0
	err := db.Update(query, acm)

	if err != nil {
		logs.Error("数据库更新SetCase时报错, err: ", err)
	}
	return acm, err
}

// 通过id修改case（全更新）
func (t *SetCaseMongo) UpdateSetCaseOrder(id int64, order int) error {
	query := bson.M{"_id": id}
	ms, db := db_proxy.Connect("auto_api", "set_case")
	defer ms.Close()

	err := db.Update(query, bson.M{"$set": bson.M{"order": order}})
	if err != nil {
		logs.Error("数据库设置SetCase的Order字段报错, err: ", err)
	}
	return err
}

// 修改status
func (t *SetCaseMongo) DelSetCase(id int64) error {
	query := bson.M{"_id": id}
	ms, db := db_proxy.Connect("auto_api", "set_case")
	defer ms.Close()

	err := db.Update(query, bson.M{"$set": bson.M{"status": del_}})
	if err != nil {
		logs.Error("数据库删除SetCase报错, err: ", err)
	}
	return err
}

// todo 注意：该方法是根据order字段排序
// 获取指定服务集合下所有Case,
func (t *SetCaseMongo) GetSetCaseListByCaseSetId(caseSetId int64) (result []*SetCaseMongo, err error) {
	ms, c := db_proxy.Connect("auto_api", "set_case")
	defer ms.Close()

	query := bson.M{"case_set_id": caseSetId, "status": status}
	// 获取测试用例集下全部case列表, 按id升序
	//err = c.Find(query).Sort("_id").All(&result)
	err = c.Find(query).Sort("order").All(&result)
	if err != nil {
		logs.Error("查询指定服务集合下所有Case数据报错, err: ", err)
		return nil, err
	}
	return
}
