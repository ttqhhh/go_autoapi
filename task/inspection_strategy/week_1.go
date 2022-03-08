package inspection_strategy

import (
	"github.com/astaxie/beego/logs"
	//controllers "go_autoapi/controllers/autotest"
	"go_autoapi/models"
)

func Delete1Week() error {
	logs.Info("【周】定期批量清理报告任务启动执行。。。。。。") /**
	删除报告库
	*/
	var rp = models.RunReportMongo{}
	runReportList, err := rp.Query()
	if nil != err {
		logs.Info("查询错误", err)
	}
	for _, one := range runReportList {
		rp.Delete(one.Id)
	}
	/**
	删除结果库(线上库中报告中和报告及结果对应的报告已经删除，所以通过时间删除)
	*/
	var rpResult = models.AutoResult{}
	autoResultList, err := rpResult.QueryResult()

	if err != nil {
		logs.Info("查询错误err:", err)
	}
	for _, result := range autoResultList {
		rpResult.DeleteResult(result.RunId)
	}
	return nil
}
