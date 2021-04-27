package models

import (
	"github.com/astaxie/beego/logs"
	"go_autoapi/db_proxy"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	db          = "auto_api"
	collection  = "service"
	time_format = "2006-01-02 15:04:05"
)

type ServiceMongo struct {
	Id          int64  `form:"id" json:"id" bson:"_id"`
	ServiceName string `form:"service_name" json:"service_name" bson:"service_name"` // 服务名
	Business    int8   `form:"business" json:"business" bson:"business"`             // 0：最右，1：皮皮，2：海外，3：中东，4：妈妈
	Status      int8   `form:"status" json:"status"  bson:"status"`                  // 0：正常，1：删除
	CreateBy    string `form:"create_by" json:"create_by" bson:"create_by"`          // 添加人
	UpdateBy    string `form:"update_by" json:"update_by" bson:"update_by"`          // 修改人
	//CreatedAt time.Time `form:"created_at" json:"created_at,omitempty" bson:"created_at"` // omitempty 表示该字段为空时，不返回
	//UpdatedAt time.Time `form:"updated_at" json:"updated_at" bson:"updated_at"`
	CreatedAt string `form:"created_at" json:"created_at" bson:"created_at"` // omitempty 表示该字段为空时，不返回
	UpdatedAt string `form:"updated_at" json:"updated_at" bson:"updated_at"`
}

func (mongo *ServiceMongo) TableName() string {
	return collection
}

// 增
func (mongo *ServiceMongo) Insert(service ServiceMongo) error {
	ms, db := db_proxy.Connect(db, collection)
	defer ms.Close()

	// id自增
	cnt, err := db.Count()
	if err != nil {
		logs.Error("Insert 错误: %v", err)
		return err
	}
	service.Id = int64(cnt) + 1
	// 处理添加时间字段
	service.CreatedAt = time.Now().Format(time_format)
	// 新增时，默认status为0
	service.Status = 0
	err = db.Insert(service)
	if err != nil {
		logs.Error("Insert 错误: %v", err)
	}
	return err
}

// 删
func (mongo *ServiceMongo) Delete(id int64) error {
	ms, db := db_proxy.Connect(db, collection)
	defer ms.Close()

	// 处理更新时间字段
	data := bson.M{
		"$set": bson.M{
			"status":     1,
			"updated_at": time.Now().Format(time_format),
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
func (mongo *ServiceMongo) Update(service ServiceMongo) error {
	ms, db := db_proxy.Connect(db, collection)
	defer ms.Close()

	// 处理更新时间字段
	service.UpdatedAt = time.Now().Format(time_format)
	data := bson.M{
		"$set": bson.M{
			"service_name": service.ServiceName,
			"business":     service.Business,
			"updated_at":   service.UpdatedAt,
			"update_by":    service.UpdateBy,
		},
	}
	changeInfo, err := db.UpsertId(service.Id, data)
	if err != nil {
		logs.Error("Update 错误: %v", err)
	}
	logs.Info("upsert函数返回的响应为：%v", changeInfo)
	return err
}

// todo 根据服务名模糊查询待实现
// 分页查询
func (mongo *ServiceMongo) QueryByPage(business int8, serviceName string, pageNo int, pageSize int) ([]ServiceMongo, int, error) {
	ms, db := db_proxy.Connect(db, collection)
	defer ms.Close()

	// 查询分页数据
	query := bson.M{"status": 0}
	if business != -1 {
		query["business"] = business
	}
	if serviceName != "" {
		query["service_name"] = bson.M{"$regex": serviceName}
	}
	serviceList := []ServiceMongo{}
	skip := (pageNo - 1) * pageSize
	err := db.Find(query).Skip(skip).Limit(pageSize).All(&serviceList)
	if err != nil {
		logs.Error("QueryByPage 错误: %v", err)
	}
	// 查询总共条数
	serviceTotalList := []ServiceMongo{}
	err = db.Find(query).All(&serviceTotalList)
	if err != nil {
		logs.Error("QueryByPage 错误: %v", err)
	}
	return serviceList, len(serviceTotalList), err
}

// 查(根据业务线)
func (mongo *ServiceMongo) QueryByBusiness(business int8) ([]ServiceMongo, error) {
	ms, db := db_proxy.Connect(db, collection)
	defer ms.Close()

	query := bson.M{"status": 0}
	if business != -1 {
		query["business"] = business
	}
	serviceList := []ServiceMongo{}
	err := db.Find(query).All(&serviceList)
	if err != nil {
		logs.Error("QueryByBusiness 错误: %v", err)
	}
	return serviceList, err
}

// 查(根据id)
func (mongo *ServiceMongo) QueryById(id int64) (*ServiceMongo, error) {
	ms, db := db_proxy.Connect(db, collection)
	defer ms.Close()

	query := bson.M{"_id": id, "status": 0}
	service := ServiceMongo{}
	err := db.Find(query).One(&service)
	if err != nil {
		if err.Error() == "not found" {
			err = nil
			//return nil, nil
			return nil, nil
		}
		logs.Error("QueryById 错误: %v", err)
	}
	return &service, err
}
