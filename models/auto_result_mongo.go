package models

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
	"go_autoapi/constants"
	"go_autoapi/db_proxy"
	"go_autoapi/utils"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	result_collection       = "auto_result"
	auto_result_primary_key = "auto_result_primary_key"
)

const (
	AUTO_RESULT_FAIL = iota
	AUTO_RESULT_SUCCESS
)

type AutoResult struct {
	Id           int64  `json:"id,omitempty" bson:"_id"`
	RunId        string `json:"run_id,omitempty" bson:"run_id"`
	CaseId       int64  `json:"case_id" bson:"case_id"`
	IsInspection int    `json:"is_inspection" bson:"is_inspection"`
	Result       int    `json:"result" bson:"result"` // 0：失败 1-成功
	Reason       string `json:"reason" bson:"reason"`
	Author       string `json:"author" bson:"author"`
	Response     string `json:"response,omitempty" bson:"response"`
	// 0705新增codeStatus
	StatusCode int `json:"status_code" bson:"status_code"`
	// omitempty 表示该字段为空时，不返回
	CreatedAt string `json:"created_at,omitempty" bson:"created_at"`
	UpdatedAt string `json:"updated_at,omitempty" bson:"updated_at"`
}

func init() {
	//orm.RegisterModel(new(AdMockCaseMongo))
	db_proxy.InitMongoDB()
	//ORM = db_proxy.GetOrmObject()
}
func (a *AutoResult) TableName() string {
	return "auto_result"
}

func InsertResult(uuid string, case_id int64, isInspection int, result int, reason string, author string, resp string, statusCode int) error {
	now := time.Now().Format(constants.TimeFormat)
	ar := AutoResult{}
	//id := GetId("result")
	r := utils.GetRedis()
	id, err := r.Incr(auto_result_primary_key).Result()
	if err != nil {
		logs.Error("auto_result新增数据时，从redis获取自增主键错误, err:", err)
	}
	logs.Info("插入操作时，获取到的自增主键为: ", id)
	ar.Id = id
	ar.RunId = uuid
	ar.CaseId = case_id
	ar.IsInspection = isInspection
	ar.Result = result
	ar.Reason = reason
	ar.Author = author
	ar.Response = resp
	ar.StatusCode = statusCode
	ar.CreatedAt = now
	ar.UpdatedAt = now
	ms, db := db_proxy.Connect(db, result_collection)
	defer ms.Close()
	return db.Insert(ar)
}

func GetResultByRunId(id string) (ar []*AutoResult, err error) {
	fmt.Println(id)
	query := bson.M{"run_id": id}
	ms, db := db_proxy.Connect(db, result_collection)
	defer ms.Close()
	err = db.Find(query).Select(bson.M{"status_code": 1, "case_id": 1, "is_inspection": 1, "reason": 1, "result": 1, "author": 1, "response": 1, "created_at": 1}).All(&ar)
	fmt.Println(ar)
	if err != nil {
		logs.Error(1024, err)
	}
	return ar, err
}

func (a *AutoResult) GetAllResult(page, limit int) (ar []*AutoResult, totalCount int64, err error) {
	ms, db := db_proxy.Connect(db, result_collection)
	defer ms.Close()

	var query interface{} = nil
	// 查询分页列表数据
	err = db.Find(query).Select(bson.M{"_id": 1, "run_id": 1, "case_id": 1, "result": 1, "reason": 1, "author": 1, "created_at": 1, "update_at": 1}).Skip((page - 1) * limit).Sort("-_id").Limit(limit).All(&ar)
	fmt.Println(ar)
	if err != nil {
		logs.Error("查询分页列表数据报错, err: ", err)
		return nil, 0, err
	}
	// 查询数据总条数用于分页
	total, err := db.Find(query).Count()
	fmt.Println(ar)
	if err != nil {
		logs.Error("查询分页列表数据报错, err: ", err)
		return nil, 0, err
	}
	totalCount = int64(total)
	return
}

func (a *AutoResult) GetFailCount(uuid string) (failCount int64, err error) {
	ms, db := db_proxy.Connect(db, result_collection)
	defer ms.Close()

	var query interface{} = bson.M{"run_id": uuid, "result": AUTO_RESULT_FAIL}
	var allFail = []AutoResult{}
	// 查询分页列表数据
	err = db.Find(query).All(&allFail)
	noReapetFail := RemoveRepeatedElement(allFail)
	if err != nil {
		logs.Error("获取全部数据出错，err:", err)
	}
	failCount = int64(len(noReapetFail))
	return

}

//删(直接删除库中数据)
func (a *AutoResult) DeleteResult(id string) error {
	ms, db := db_proxy.Connect(db, result_collection)
	defer ms.Close()
	// 处理更新时间字段
	data := bson.M{
		"run_id": id,
	}
	err := db.Remove(data)
	if err != nil {
		logs.Error("Delete 错误，err:", err)
	}
	return err
}

//查询一周之前的的报告结果
func (a *AutoResult) QueryResult() ([]AutoResult, error) {
	nowTimeStr := time.Now().AddDate(0, 0, -2).Format("2006-01-02 15:04:0")
	ms, db := db_proxy.Connect(db, result_collection)
	defer ms.Close()
	query := bson.M{
		"created_at": bson.M{
			"$lt": nowTimeStr,
		},
	}
	autoResultList := []AutoResult{}
	err := db.Find(query).All(&autoResultList)
	if err != nil {
		if err.Error() == "not found" {
			err = nil
			//return nil, nil
			return nil, nil
		}
		logs.Error("查询错误: %v", err)
	}
	return autoResultList, err
}

//去重
func RemoveRepeatedElement(arr []AutoResult) (newArr []AutoResult) {
	newArr = make([]AutoResult, 0)
	for i := 0; i < len(arr); i++ {
		repeat := false
		for j := i + 1; j < len(arr); j++ {
			if arr[i] == arr[j] {
				repeat = true
				break
			}
		}
		if !repeat {
			newArr = append(newArr, arr[i])
		}
	}
	return
}
