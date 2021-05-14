package task

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/toolbox"
)

func init() {
	// new一个线上巡检的定时任务定时任务
	onlineInspectionTask := toolbox.NewTask("Inspection_Task", ONLINE_INSPECTION_EXPRESSION, OnlineInspection)
	err := onlineInspectionTask.Run()
	if err != nil {
		logs.Error("运行系统定时任务报错， err: ", err)
		return
	}
	toolbox.AddTask("Inspection_Task", onlineInspectionTask)
	toolbox.StartTask()
}
