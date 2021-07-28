package inspection_strategy

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"strconv"
	"time"

	//controllers "go_autoapi/controllers/autotest"
	"go_autoapi/models"
)

func Delete1Month() error {
	logs.Info("【月】定期批量清理报告任务启动执行。。。。。。")
	var once = 0  //计数器
	var rp = models.RunReportMongo{}
	runReportList, err := rp.Query()
	if nil !=err{
			logs.Info("查询错误",err)
		}
		count:=len(runReportList)  //返回的list个数
		str1 := strconv.Itoa(count)
		fmt.Printf("共查询到"+str1+"个报告,开始批量删除超过一个月的报告\n\n")

		for _,report := range runReportList {   //遍历list
			reportCreateStringTime:=report.CreatedAt  //取出每个报告的创建时间
			const Layout = "2006-01-02 15:04:05"//时间常量
			reportCreateTime,_ := time.ParseInLocation(Layout,reportCreateStringTime,time.Local) //转换为date类型

			if reportCreateTime.Before(time.Now().AddDate(0,0,-7)){ //如果报告时间超过一周

				rp.Delete(report.Id)
				once=once+1  //执行一次，次数+1
			}
		}
		times := strconv.Itoa(once)
		logs.Info("批量删除成功,共删除"+times+"条数据")
		return nil
}
