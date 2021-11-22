package task

import (
	"github.com/astaxie/beego/toolbox"
	"go_autoapi/controllers/inspection"
	"go_autoapi/task/inspection_strategy"
	"go_autoapi/task/monitor"
)

// 是否开启线上巡检任务(测试环境关闭)
const IS_OPEN_INSPECTION_TASK = true

// 是否开启线上巡检任务(测试环境关闭)
const IS_OPEN_RT_MONITOR_TASK = true

//是否开启定期删除任务
const IS_OPEN_REGULAR_DELETE_TASK = true

func init() {
	// new一个线上巡检的定时任务定时任务
	//onlineInspectionTask := toolbox.NewTask("Inspection_Task", ONLINE_INSPECTION_EXPRESSION, OnlineInspection)
	//toolbox.AddTask("Inspection_Task", onlineInspectionTask)
	// todo 调试使用
	//err := onlineInspectionTask.Run()
	//if err != nil {
	//	logs.Error("运行系统定时任务报错， err: ", err)
	//	return
	//}
	if IS_OPEN_INSPECTION_TASK {
		Inspection1M := toolbox.NewTask("Inspection_Task_5min", inspection.ONE_MIN_EXPRESSION, inspection_strategy.Strategy1Min)
		toolbox.AddTask("Inspection_Task_1min", Inspection1M)
		Inspection5M := toolbox.NewTask("Inspection_Task_5min", inspection.FIVE_MIN_EXPRESSION, inspection_strategy.Strategy5Min)
		toolbox.AddTask("Inspection_Task_5min", Inspection5M)
		Inspection10M := toolbox.NewTask("Inspection_Task_10min", inspection.TEN_MIN_EXPRESSION, inspection_strategy.Strategy10Min)
		toolbox.AddTask("Inspection_Task_10min", Inspection10M)
		Inspection1Q := toolbox.NewTask("Inspection_Task_15min", inspection.ONE_QUARTER_EXPRESSION, inspection_strategy.Strategy15Min)
		toolbox.AddTask("Inspection_Task_15min", Inspection1Q)
		InspectionHalfH := toolbox.NewTask("Inspection_Task_30min", inspection.HALF_HOUR_EXPRESSION, inspection_strategy.Strategy30Min)
		toolbox.AddTask("Inspection_Task_30min", InspectionHalfH)
		Inspection1H := toolbox.NewTask("Inspection_Task_1hour", inspection.ONE_HOUR_EXPRESSION, inspection_strategy.Strategy1Hour)
		toolbox.AddTask("Inspection_Task_1hour", Inspection1H)
		Inspection6H := toolbox.NewTask("Inspection_Task_6hour", inspection.SIX_HOUR_EXPRESSION, inspection_strategy.Strategy6Hour)
		toolbox.AddTask("Inspection_Task_6hour", Inspection6H)
		InspectionHalfD := toolbox.NewTask("Inspection_Task_12hour", inspection.HALF_DAY_EXPRESSION, inspection_strategy.Strategy12Hour)
		toolbox.AddTask("Inspection_Task_12hour", InspectionHalfD)
		Inspection1D := toolbox.NewTask("Inspection_Task_1day", inspection.ONE_DAY_EXPRESSION, inspection_strategy.Strategy1Day)
		toolbox.AddTask("Inspection_Task_1day", Inspection1D)
	}
	if IS_OPEN_RT_MONITOR_TASK {
		RtMonitorTask := toolbox.NewTask("Rt_Monitor_Task", monitor.MONITOR_TASK_EXPRESSION, monitor.MonitorTask)
		toolbox.AddTask("Rt_Monitor_Task", RtMonitorTask)
	}
	if IS_OPEN_REGULAR_DELETE_TASK { //新增，每周一执行一次的定时删除报告任务
		//Delete1Week := toolbox.NewTask("Delete_Date_1Week", inspection.ONE_WEEK_EXPRESSION, inspection_strategy.Delete1Week)
		//toolbox.AddTask("Delete_Date_1Week", Delete1Week)
		//性能监控的定期删除  一周一次，清理时间小于当前时间的7天前记录
		Delete1WeekAlert := toolbox.NewTask("Delete_Date_1Week_Alert", inspection.ONE_WEEK_EXPRESSION, inspection_strategy.DeleteOneWeek)
		toolbox.AddTask("Delete_Date_1Week_Alert", Delete1WeekAlert)

	}
	toolbox.StartTask()
}
