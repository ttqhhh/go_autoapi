package inspection_strategy

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"strconv"
	//controllers "go_autoapi/controllers/autotest"
	"go_autoapi/models"
)

func DeleteOneWeek() error {
	logs.Info("【周】定期批量清理性能监控任务启动执行。。。。。。")
	var once1 = 0 //报告结果计数器
	// 性能监控报告的定期删除
	var rpAlertResult = models.RtDetailAlertMongo{}
	RtDetailAlertList, err := rpAlertResult.QueryResult()
	count1 := len(RtDetailAlertList)
	str2 := strconv.Itoa(count1)
	fmt.Printf("查询到报告的结果为" + str2 + "，开始删除超过一周的报告结果\n\n")
	if err != nil {
		logs.Info("查询错误err:", err)
	}
	for _, result := range RtDetailAlertList {
		rpAlertResult.DeleteById(result.Id)
		once1++
	}
	timess := strconv.Itoa(once1)
	logs.Info("批量删除成功,共删除" + timess + "条报告结果")
	return nil
}
