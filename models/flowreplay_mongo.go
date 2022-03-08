package models

import (
	"github.com/astaxie/beego/logs"
	"go_autoapi/constants"
	"go_autoapi/db_proxy"
	"go_autoapi/utils"
	"gopkg.in/mgo.v2/bson"
	"time"
)

const (
	flowreplay_collection = "flow_replay"
)

type FlowReplayMongo struct {
	Id             int64  `form:"id" json:"id" bson:"_id"`
	ServiceName    string `form:"service_name" json:"service_name" bson:"service_name"`             // 服务名，需唯一
	FlowFile       string `form:"flow_file" json:"flow_file" bson:"flow_file"`                      // 流量文件
	FlowTargetHost string `form:"flow_target_host" json:"flow_target_host" bson:"flow_target_host"` //回放目标机器
	ReplayTimes    string `form:"replay_times" json:"replay_times" bson:"replay_times"`             // 回放倍率
	ReplayUri      string `form:"replay_uri" json:"replay_uri" bson:"replay_uri"`                   // 回放地址变为并发数
	Status         int8   `form:"status" json:"status"  bson:"status"`                              // 0：正常，1：删除
	CreateBy       string `form:"create_by" json:"create_by" bson:"create_by"`                      // 添加人
	UpdateBy       string `form:"update_by" json:"update_by" bson:"update_by"`                      // 修改人
	CreatedAt      string `form:"created_at" json:"created_at" bson:"created_at"`                   // omitempty 表示该字段为空时，不返回
	UpdatedAt      string `form:"updated_at" json:"updated_at" bson:"updated_at"`
	Cycle          string `form:"cycle" json:"cycle" bson:"cycle"` //是否循环
}

func (mongo *FlowReplayMongo) TableName() string {
	return flowreplay_collection
}

// 增
func (mongo *FlowReplayMongo) Insert(flowreplay FlowReplayMongo) error {
	ms, db := db_proxy.Connect(db, flowreplay_collection)
	defer ms.Close()

	// id自增
	r := utils.GetRedis()
	defer r.Close()
	id, err := r.Incr(constants.Flow_Replay_PRIMARY_KEY).Result()
	if err != nil {
		logs.Error("保存Case时，获取从redis获取唯一主键报错，err: ", err)
		return err
	}
	flowreplay.Id = id
	// 处理添加时间字段
	flowreplay.CreatedAt = time.Now().Format(Time_format)
	// 新增时，默认status为0
	flowreplay.Status = 0
	err = db.Insert(flowreplay)
	if err != nil {
		logs.Error("Insert 错误: %v", err)
	}
	return err
}

// 查(根据id)
func (mongo *FlowReplayMongo) QueryById(id int64) (*FlowReplayMongo, error) {
	ms, db := db_proxy.Connect(db, flowreplay_collection)
	defer ms.Close()

	query := bson.M{"_id": id, "status": 0}
	flowreplay := FlowReplayMongo{}
	err := db.Find(query).One(&flowreplay)
	if err != nil {
		if err.Error() == "not found" {
			err = nil
			return nil, nil
		}
		logs.Error("QueryById 错误: %v", err)
	}
	return &flowreplay, err
}

// 查(根据id)
func (mongo *FlowReplayMongo) QueryByName(serviceName string) (*FlowReplayMongo, error) {
	ms, db := db_proxy.Connect(db, flowreplay_collection)
	defer ms.Close()

	query := bson.M{"service_name": serviceName, "status": 0}
	flowreplay := FlowReplayMongo{}
	err := db.Find(query).One(&flowreplay)
	if err != nil {
		if err.Error() == "not found" {
			err = nil
			return nil, nil
		}
		logs.Error("QueryById 错误: %v", err)
	}
	return &flowreplay, err
}

// 改
func (mongo *FlowReplayMongo) Update(flowReplay FlowReplayMongo) error {
	ms, db := db_proxy.Connect(db, flowreplay_collection)
	defer ms.Close()

	// 处理更新时间字段
	flowReplay.UpdatedAt = time.Now().Format(Time_format)
	data := bson.M{}
	if flowReplay.FlowFile == "" {
		data = bson.M{
			"$set": bson.M{
				"service_name":     flowReplay.ServiceName,
				"flow_target_host": flowReplay.FlowTargetHost,
				"replay_times":     flowReplay.ReplayTimes,
				"replay_uri":       flowReplay.ReplayUri,
				"updated_at":       flowReplay.UpdatedAt,
				"update_by":        flowReplay.UpdateBy,
			},
		}
	} else {
		data = bson.M{
			"$set": bson.M{
				"service_name":     flowReplay.ServiceName,
				"flow_file":        flowReplay.FlowFile,
				"flow_target_host": flowReplay.FlowTargetHost,
				"replay_times":     flowReplay.ReplayTimes,
				"replay_uri":       flowReplay.ReplayUri,
				"updated_at":       flowReplay.UpdatedAt,
				"update_by":        flowReplay.UpdateBy,
			},
		}
	}

	changeInfo, err := db.UpsertId(flowReplay.Id, data)
	if err != nil {
		logs.Error("Update 错误: %v", err)
	}
	logs.Info("upsert函数返回的响应为：%v", changeInfo)
	return err
}

// 删
func (mongo *FlowReplayMongo) Delete(id int64) error {
	ms, db := db_proxy.Connect(db, flowreplay_collection)
	defer ms.Close()

	// 处理更新时间字段
	data := bson.M{
		"$set": bson.M{
			"status":     1,
			"updated_at": time.Now().Format(Time_format),
		},
	}
	changeInfo, err := db.UpsertId(id, data)
	if err != nil {
		logs.Error("Delete 错误: %v", err)
	}
	logs.Info("upsert函数返回的响应为：%v", changeInfo)
	return err
}

// 分页查询
func (mongo *FlowReplayMongo) QueryByPage(serviceName string, pageNo int, pageSize int) ([]FlowReplayMongo, int, error) {
	ms, db := db_proxy.Connect(db, flowreplay_collection)
	defer ms.Close()

	// 查询分页数据
	query := bson.M{"status": 0}
	//if business != -1 {
	//	query["business"] = business
	//}
	if serviceName != "" {
		query["service_name"] = bson.M{"$regex": serviceName}
	}
	// 进行业务线筛选
	//if len(businesses) > 0 {
	//	//queryCond := []interface{}{bson.D{"business"}}
	//	queryCond := []interface{}{}
	//	for _, v := range businesses {
	//		queryCond = append(queryCond, bson.M{"business": v})
	//	}
	//	query["$or"] = queryCond
	//}

	serviceList := []FlowReplayMongo{}
	skip := (pageNo - 1) * pageSize
	err := db.Find(query).Sort("-_id").Skip(skip).Limit(pageSize).All(&serviceList)
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
