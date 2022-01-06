package h5listen

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"go_autoapi/models"
	"strconv"
)

func (c *H5ListenController) updateH5DateByID() {
	acm := models.AdH5DataMongo{}
	id := c.GetString("id")
	data_name := c.GetString("data_name")
	data_url := c.GetString("data_url")
	business := c.GetString("business")
	Id, err := strconv.ParseInt(id, 10, 64)
	fmt.Println(id, data_name, data_url, business)

	acm, err = acm.UpdateDataById(Id, data_name, data_url, business, acm)
	if err != nil {
		logs.Error("更新Case报错，err: ", err)
		c.ErrorJson(-1, "请求错误", nil)
	}
	c.Ctx.Redirect(302, "/h5listen/show_h5")
}

func (c *H5ListenController) DelH5DateByID() {
	caseID := c.GetString("id")
	ac := models.AdH5DataMongo{}
	caseIDInt, err := strconv.ParseInt(caseID, 10, 64)
	if err != nil {
		logs.Error("在删除用例的时候类型转换失败")
	}
	//ac.DelCase(caseIDInt)
	ac.DeleteDataById(caseIDInt)
	logs.Info("删除成功")
}
