package controllers

func (c *CaseManageController) showCases(){
	//c.Data["Website"] = "beego.me"
	//c.Data["Email"] = "astaxie@gmail.com"

	c.TplName = "case_manager.html"

}

func (c *CaseManageController) getAllCase(){

	//读取数据库全部case
	//all_case = models.TestCaseMongo{}
	c.Data["json"] = "hhhh"
	c.ServeJSON()
}