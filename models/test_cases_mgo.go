package models

import (
	"errors"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	_ "github.com/go-sql-driver/mysql"
	"go_autoapi/db_proxy"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	status = 0
	del_   = 1
)

const (
	NOT_INSPECTION = iota // 非线上巡检接口
	INSPECTION            // 线上巡检接口
)
const (
	SHANG_YE_HUA_TEST = "test" //测试环境域名
)

type TestCaseMongo struct {
	Id           int64  `form:"id" json:"id" bson:"_id"`
	ApiName      string `form:"api_name" json:"api_name" bson:"api_name"`
	CaseName     string `form:"case_name" json:"case_name" bson:"case_name"`
	IsInspection int8   `form:"is_inspection" json:"is_inspection" bson:"is_inspection"`
	Description  string `form:"description" json:"description" bson:"description"`
	Method       string `form:"method" json:"method" bson:"method"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	//zen
	Author string `form:"author" json:"author" bson:"author"`
	//AppName       string `form:"app_name" json:"app_name" bson:"app_name"`
	Domain        string `form:"domain" json:"domain" bson:"domain"`
	BusinessName  string `form:"business_name" json:"business_name" bson:"business_name"`
	BusinessCode  string `form:"business_code" json:"business_code" bson:"business_code"`
	ServiceId     int64  `form:"service_id" json:"service_id" bson:"service_id"`
	ServiceName   string `form:"service_name" json:"service_name" bson:"service_name"`
	ApiUrl        string `form:"api_url" json:"api_url" bson:"api_url"`
	TestEnv       string `form:"test_env" json:"test_env" bson:"test_env"`
	Mock          string `form:"mock" json:"mock" bson:"mock"`
	RequestMethod string `form:"request_method" json:"request_method" bson:"request_method"`
	Parameter     string `form:"parameter" json:"parameter" bson:"parameter"`
	Checkpoint    string `form:"check_point" json:"check_point" bson:"check_point"`
	SmokeResponse string `form:"smoke_response" json:"smoke_response,omitempty" bson:"smoke_response"`
	Level         string `form:"level" json:"level" bson:"level"`
	Status        int64  `json:"status" bson:"status"`
}

//db:操作的数据库
//collection:操作的文档(表)
//query:查询条件
//selector:需要过滤的数据(projection)
//result:查询到的结果

// 获取指定server下的所有case

func (t *TestCaseMongo) GetCasesByQuery(query interface{}) (TestCaseMongo, error) {
	//query := TestCaseMongo{}
	var acm = TestCaseMongo{}
	ms, c := db_proxy.Connect("auto_api", "case")
	defer ms.Close()
	err := c.Find(query).All(&acm)
	if err != nil {
		logs.Error(1024, err)
	}
	return acm, err
}

//通过id list 获取用例

//func (t *TestCaseMongo) GetCasesByIds(ids []string) []TestCaseMongo {
//	var caseList []TestCaseMongo
//	for _, i := range ids {
//		id64, err := strconv.ParseInt(i, 10, 64)
//		if err != nil {
//			logs.Error("类型转换失败")
//		}
//		acm := t.GetOneCase(id64)
//		caseList = append(caseList, acm)
//	}
//	return caseList
//}

func (t *TestCaseMongo) GetCasesByConfusedUrl(page, limit int, business string, url string, service string) (result []TestCaseMongo, totalCount int64, err error) {
	ms, c := db_proxy.Connect("auto_api", "case")
	defer ms.Close()
	var query = bson.M{}
	if service == "" && url != "" {
		query = bson.M{"status": status,
			"business_code": business,
			"api_url":       bson.M{"$regex": bson.RegEx{Pattern: url, Options: "im"}}}
	} else if service != "" && url == "" {
		query = bson.M{"status": status, "business_code": business, "service_name": service}
	} else if service != "" && url != "" {
		query = bson.M{"status": status,
			"business_code": business,
			"service_name":  service,
			"api_url":       bson.M{"$regex": bson.RegEx{Pattern: url, Options: "im"}}}
	} else {
		query = bson.M{"status": status, "business_code": business}
	}
	// 获取指定业务线下全部case列表
	err = c.Find(query).Sort("-_id").Skip((page - 1) * limit).Limit(limit).All(&result)
	if err != nil {
		logs.Error("查询分页列表数据报错, err: ", err)
		return nil, 0, err
	}
	// 获取指定业务线下全部case数量
	total, err := c.Find(query).Count()
	if err != nil {
		logs.Error("数据库查询指定业务线下case数量报错, err: ", err)
		return nil, 0, err
	}
	totalCount = int64(total)
	return
}

// 获取指定业务线下的指定页面case
func (t *TestCaseMongo) GetAllCases(page, limit int, business string) (result []TestCaseMongo, totalCount int64, err error) {
	//acm := TestCaseMongo{}
	//result := make([]TestCaseMongo, 0, 10)
	ms, c := db_proxy.Connect("auto_api", "case")
	defer ms.Close()
	query := bson.M{"status": status, "business_code": business}
	// 获取指定业务线下全部case列表
	err = c.Find(query).Sort("-_id").Skip((page - 1) * limit).Limit(limit).All(&result)
	//err := c.Find(bson.M{"api_name":"api_name"}).One(&acm)
	if err != nil {
		logs.Error("查询分页列表数据报错, err: ", err)
		return nil, 0, err
	}
	// 获取指定业务线下全部case数量
	total, err := c.Find(query).Count()
	if err != nil {
		logs.Error("数据库查询指定业务线下case数量报错, err: ", err)
		return nil, 0, err
	}
	totalCount = int64(total)
	return
}

// 通过id获取指定case

func (t *TestCaseMongo) GetOneCase(id int64) TestCaseMongo {

	fmt.Println(id)
	query := bson.M{"_id": id, "status": status}
	acm := TestCaseMongo{}
	ms, db := db_proxy.Connect("auto_api", "case")
	defer ms.Close()
	err := db.Find(query).One(&acm)
	fmt.Println(acm)
	if err != nil {
		logs.Info("查询case失败")
		logs.Error(1024, err)
	}
	return acm
}

// 添加一条case

func (t *TestCaseMongo) AddCase(acm TestCaseMongo) error {
	ms, db := db_proxy.Connect("auto_api", "case")
	defer ms.Close()
	query := bson.M{"api_url": acm.ApiUrl, "domain": acm.Domain, "parameter": acm.Parameter, "status": 0}
	err := db.Find(query).One(&acm)
	if err == nil {
		return errors.New("用例已经存在")
	}
	err = db.Insert(acm)
	if err != nil {
		logs.Error(1024, err)
	}
	return err
}

// 通过id修改case（全更新）

func (t *TestCaseMongo) UpdateCase(id int64, acm TestCaseMongo) (TestCaseMongo, error) {
	query := bson.M{"_id": id}
	ms, db := db_proxy.Connect("auto_api", "case")
	defer ms.Close()
	acm.Status = 0
	err := db.Update(query, acm)
	fmt.Println(acm)
	if err != nil {
		logs.Error(1024, err)
	}
	return acm, err
}

//
func (t *TestCaseMongo) SetInspection(id int64, is_inspection int8) error {
	ms, db := db_proxy.Connect("auto_api", "case")
	defer ms.Close()

	data := bson.M{
		"$set": bson.M{
			"is_inspection": is_inspection,
			"updated_at":    time.Now().Format(Time_format),
		},
	}
	changeInfo, err := db.UpsertId(id, data)
	if err != nil {
		logs.Error("设置巡检状态错误: err: ", err)
	}
	logs.Info("upsert函数返回的响应为：%v", changeInfo)
	return err
}

// 修改status

func (t *TestCaseMongo) DelCase(id int64) {
	query := bson.M{"_id": id}
	ms, db := db_proxy.Connect("auto_api", "case")
	defer ms.Close()
	//err := db.Find(query).One(&acm)
	err := db.Update(query, bson.M{"$set": bson.M{"status": del_}})
	if err != nil {
		logs.Error("删除case失败，更给状态为1失败")
		logs.Error(err)
	}
}

// 获取指定业务线下所有Case
func (t *TestCaseMongo) GetAllCasesByBusiness(business string) (result []*TestCaseMongo, err error) {
	ms, c := db_proxy.Connect("auto_api", "case")
	defer ms.Close()
	query := bson.M{"status": status, "business_code": business}
	// 获取指定业务线下全部case列表
	err = c.Find(query).All(&result)
	if err != nil {
		logs.Error("查询指定业务线下所有Case数据报错, err: ", err)
		return nil, err
	}
	return result, err
}

// 获取指定服务集合下所有Case
func (t *TestCaseMongo) GetAllCasesByServiceList(serviceIds []int64) (result []*TestCaseMongo, err error) {
	ms, c := db_proxy.Connect("auto_api", "case")
	defer ms.Close()

	query := bson.M{}
	if len(serviceIds) > 0 {
		//queryCond := []interface{}{bson.D{"business"}}
		queryCond := []interface{}{}
		for _, serviceId := range serviceIds {
			queryCond = append(queryCond, bson.M{"service_id": serviceId})
		}
		query["$or"] = queryCond
	}
	query["status"] = status
	//query := bson.M{"status": status, "service_id": business}
	// 获取指定业务线下全部case列表
	err = c.Find(query).All(&result)
	if err != nil {
		logs.Error("查询指定服务集合下所有Case数据报错, err: ", err)
		return nil, err
	}
	return
}

// 获取指定服务集合下所有Case
func (t *TestCaseMongo) GetAllInspectionCasesByService(serviceId int64) (result []*TestCaseMongo, err error) {
	ms, c := db_proxy.Connect("auto_api", "case")
	defer ms.Close()

	query := bson.M{"status": status, "service_id": serviceId, "is_inspection": INSPECTION}
	// 获取指定业务线下全部case列表
	err = c.Find(query).All(&result)
	if err != nil {
		logs.Error("查询指定服务下所有巡检Case数据报错, err: ", err)
		return nil, err
	}
	return
}
