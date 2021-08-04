package inspection_strategy

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/blinkbean/dingtalk"
	"strconv"
	"time"

	//controllers "go_autoapi/controllers/autotest"
	"go_autoapi/models"
)

func Delete1Week() error {
	logs.Info("【周】定期批量清理报告任务启动执行。。。。。。")
	var once = 0  //报告计数器
	var once1 =0  //报告结果计数器
	/**
	删除报告库
	 */
	var rp = models.RunReportMongo{}
	runReportList, err := rp.Query()
	if nil !=err{
			logs.Info("查询错误",err)
		}
		count:=len(runReportList)  //返回的list个数
		str1 := strconv.Itoa(count)
		fmt.Printf("共查询到"+str1+"个报告,开始批量删除超过一周的报告\n\n")

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
	/**
	删除结果库(线上库中报告中和报告及结果对应的报告已经删除，所以通过时间删除)
	*/
		var rpResult = models.AutoResult{}
		autoResultList, err := rpResult.QueryResult()
			count1:=len(autoResultList)
			str2 := strconv.Itoa(count1)
			fmt.Printf("查询到报告的结果为"+str2+"，开始删除超过一周的报告结果\n\n")

	if err!= nil{
			logs.Info("查询错误err:",err)
		}
		for _,result:= range autoResultList{
			resultCreateStringTime :=result.CreatedAt
			const Layout = "2006-01-02 15:04:05"//时间常量
			resultCreateTime ,_ := time.ParseInLocation(Layout,resultCreateStringTime,time.Local) //转换为date类型
			if resultCreateTime.Before(time.Now().AddDate(0,0,-8)){ //如果报告时间超过一周（为了防止报告结果查不到，留下多一天的报告结果）
				rpResult.DeleteResult(result.RunId)
			once1++
			}
		}
		timess := strconv.Itoa(once1)
		logs.Info("批量删除成功,共删除"+timess+"条报告结果")
		DingSendDeleteReport("【线上巡检：定时删除】\n每周一8：00的定时删除功能启动：已经删除一周前的测试报告和巡检报告，请有关同学及时关注并检查！\n清除报告数量：\n\n共检查到"+str1+"个报告\n已删除"+times+"个报告")
		return nil
}
func DingSendDeleteReport(content string) {
	var dingToken = []string{XIAO_NENG_QUN_TOKEN}
	cli := dingtalk.InitDingTalk(dingToken, "")
	cli.SendTextMessage(content)
}
