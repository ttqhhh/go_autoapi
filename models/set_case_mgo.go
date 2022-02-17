package models

import (
	"errors"
	"fmt"
	"github.com/beego/beego/v2/core/logs"
	_ "github.com/go-sql-driver/mysql"
	"go_autoapi/db_proxy"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

type SetCaseMongo struct {
	Id          int64  `form:"id" json:"id" bson:"_id"`
	SetId       int64  `form:"set_id" json:"set_id" bson:"set_id"`
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
	Status        int64  `json:"status" bson:"status"`
}

//db:操作的数据库
//collection:操作的文档(表)
//query:查询条件
//selector:需要过滤的数据(projection)
//result:查询到的结果

// 获取指定server下的所有case

func (t *SetCaseMongo) GetSetCaseByQuery(query interface{}) (SetCaseMongo, error) {
	//query := SetCaseMongo{}
	var acm = SetCaseMongo{}
	ms, c := db_proxy.Connect("auto_api", "case")
	defer ms.Close()
	err := c.Find(query).All(&acm)
	if err != nil {
		logs.Error(1024, err)
	}
	return acm, err
}

//通过id list 获取用例

//func (t *SetCaseMongo) GetCasesByIds(ids []string) []SetCaseMongo {
//	var caseList []SetCaseMongo
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

func (t *SetCaseMongo) GetSetCaseByConfusedUrl(page, limit int, business string, url string, service string) (result []SetCaseMongo, totalCount int64, err error) {
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
func (t *SetCaseMongo) GetAllSetCase(page, limit int, business string) (result []SetCaseMongo, totalCount int64, err error) {
	//acm := SetCaseMongo{}
	//result := make([]SetCaseMongo, 0, 10)
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

func (t *SetCaseMongo) GetOneSetCase(id int64) SetCaseMongo {

	fmt.Println(id)
	query := bson.M{"_id": id, "status": status}
	acm := SetCaseMongo{}
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

func (t *SetCaseMongo) AddSetCase(acm SetCaseMongo) error {
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

func (t *SetCaseMongo) UpdateSetCase(id int64, acm SetCaseMongo) (SetCaseMongo, error) {
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
//func (t *SetCaseMongo) SetInspection(id int64, is_inspection int8) error {
//	ms, db := db_proxy.Connect("auto_api", "case")
//	defer ms.Close()
//
//	data := bson.M{
//		"$set": bson.M{
//			"is_inspection": is_inspection,
//			"updated_at":    time.Now().Format(Time_format),
//		},
//	}
//	changeInfo, err := db.UpsertId(id, data)
//	if err != nil {
//		logs.Error("设置巡检状态错误: err: ", err)
//	}
//	logs.Info("upsert函数返回的响应为：%v", changeInfo)
//	return err
//}

// 修改status

func (t *SetCaseMongo) DelSetCase(id int64) {
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
func (t *SetCaseMongo) GetAllSetCaseByBusiness(business string, kind int) (result []*SetCaseMongo, err error) {
	var testList []*SetCaseMongo
	var onlineList []*SetCaseMongo
	ms, c := db_proxy.Connect("auto_api", "case")
	defer ms.Close()
	query := bson.M{"status": status, "business_code": business}
	// 获取指定业务线下全部case列表
	err = c.Find(query).All(&result)
	if err != nil {
		logs.Error("查询指定业务线下所有Case数据报错, err: ", err)
		return nil, err
	}
	if kind == 1 { //测试环境 通过域名筛选
		for _, one := range result {
			if strings.Contains(one.Domain, SHANG_YE_HUA_TEST) {
				testList = append(testList, one)
			}

		}
		return testList, err
	}
	if kind == 2 { //线上环境
		for _, one := range result {
			if strings.Contains(one.Domain, SHANG_YE_HUA_TEST) {
				//什么都不做
			} else {
				onlineList = append(onlineList, one)
			}

		}
		return onlineList, err
	}
	return result, err
}

// 获取指定服务集合下所有Case
func (t *SetCaseMongo) GetAllSetCaseByServiceList(serviceIds []int64) (result []*SetCaseMongo, err error) {
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
//func (t *SetCaseMongo) GetAllInspectionCasesByService(serviceId int64) (result []*SetCaseMongo, err error) {
//	ms, c := db_proxy.Connect("auto_api", "case")
//	defer ms.Close()
//
//	query := bson.M{"status": status, "service_id": serviceId, "is_inspection": INSPECTION}
//	// 获取指定业务线下全部case列表
//	err = c.Find(query).All(&result)
//	if err != nil {
//		logs.Error("查询指定服务下所有巡检Case数据报错, err: ", err)
//		return nil, err
//	}
//	return
//}

func (t *SetCaseMongo) GetSetCaseByCondition(business_code string, service_code string, case_name string) (acms []*SetCaseMongo, err error) {
	ms, c := db_proxy.Connect("auto_api", "case")
	defer ms.Close()

	var query = bson.M{"status": status}

	if business_code != "" {
		query["business_code"] = business_code
	}
	if service_code != "" {
		query["service"] = service_code
	}
	if case_name != "" {
		query["case_name"] = bson.M{"$regex": bson.RegEx{Pattern: case_name, Options: "im"}}
	}

	// 获取指定业务线下全部case列表
	err = c.Find(query).Select(bson.M{"_id": 1, "case_name": 1, "api_url": 1}).Sort("-_id").All(&acms)
	if err != nil {
		logs.Error("数据库按指定条件查询用例数据报错, err: ", err)
		return nil, err
	}

	if err != nil {
		logs.Error("数据库按指定条件查询用例数据报错, err: ", err)
		return nil, err
	}

	return
}
