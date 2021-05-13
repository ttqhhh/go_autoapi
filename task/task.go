package task

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/toolbox"
	controllers "go_autoapi/controllers/autotest"
	"go_autoapi/models"
)

func init() {
	// new一个定时任务, 每天早上9点执行
	tk := toolbox.NewTask("Inspection_Task", "0 0 9 * * ? *",
		func() error {
			logs.Info("启动定时任务：Product Inspection")
			// 获取所以业务线，并进行遍历所有的业务线
			businesses := controllers.GetAllBusinesses()
			for _, business := range businesses {
				// 遍历业务线下的所有所有服务
				serviceMongo := models.ServiceMongo{}
				businessId := int8(business["code"].(int))
				serviceMongos, err := serviceMongo.QueryByBusiness(businessId)
				if err != nil {
					logs.Error("执行线上巡检定时任务时，查询指定业务线下的服务时报错， err: ", err)
					return err
				}
				// 遍历服务下边所有的巡检Case
				for _, service := range serviceMongos {
					controllers.PerformInspection(businessId, service.Id)
				}
			}
			return nil
		})
	err := tk.Run()
	if err != nil {
		logs.Error("运行系统定时任务报错， err: ", err)
		return
	}
	toolbox.AddTask("Inspection_Task", tk)
	toolbox.StartTask()
}
