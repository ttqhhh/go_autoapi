package web_report

import (
	"database/sql"
	"fmt"
	"github.com/astaxie/beego/logs"
	constant "go_autoapi/constants"
	controllers "go_autoapi/controllers/autotest"
	"go_autoapi/libs"
	"go_autoapi/models"
	"log"
	"net/smtp"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type WebreportController struct {
	libs.BaseController
}

type DataSt struct {
	Id        string `form:"-"`
	Name      string `form:"name"`
	Describe  string `form:"describe"`
	Xmyl      string `form:"xmyl"`
	Jszb      string `form:"jszb"`
	Fx        string `form:"fx"`
	Zb        string `form:"zb"`
	Sm        string `form:"sm"`
	Recipient string `form:"recipient"`
}

func (c *WebreportController) Get() {
	do := c.GetMethodName()
	switch do {
	case "show_web_report":
		c.ShowWebReport()
	case "allwebreport":
		c.AllWebReport()
	case "query":
		c.Query()
	case "queryId":
		c.Queryid()
	case "show_email":
		c.ShowEmail()
	case "allemail":
		c.AllEmail()
	case "queryemail":
		c.QueryEmail()
	case "getById":
		c.getById()

	default:
		logs.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持呀", nil)
	}
}

func (c *WebreportController) Post() {
	do := c.GetMethodName()
	switch do {
	case "submit":
		c.Insert()
	case "queryReport":
		c.Query()
	case "query":
		c.AllWebReport()
	case "insert_email":
		c.SaveEmail()
	case "queryEmail":
		c.QueryEmail()

	default:
		logs.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持", nil)
	}
}

func (c *WebreportController) ShowWebReport() {
	c.TplName = "web_report.html"
}

func (c *WebreportController) AllWebReport() {
	c.TplName = "web_report_history.html"
}

func (c *WebreportController) AllEmail() {
	c.TplName = "web_email_history.html"
}

func (c *WebreportController) ShowEmail() {
	c.TplName = "web_report_email.html"
}

//添加邮箱组
func (c *WebreportController) SaveEmail() {
	userId, _ := c.GetSecureCookie(constant.CookieSecretKey, "user_id")
	email := &models.EmailMongo{}
	err := c.ParseForm(email)
	if err != nil {
		logs.Warn("/service/save接口 参数异常, err: %v", err)
		c.ErrorJson(-1, "参数异常", nil)
	}
	logs.Info("请求参数：%v", email)

	if string(email.Id) == "" || email.Id == -1 || email.Id == 0 {
		email.CreateBy = userId
		err = email.Insert(*email)
		if err != nil {
			c.ErrorJson(-1, "服务添加数据异常", nil)
		}
	} else {
		email.UpdateBy = userId
		err = email.Update(*email)
		if err != nil {
			c.ErrorJson(-1, "服务更新数据异常", nil)
		}
	}
	c.SuccessJson(nil)
}

//根据id查询邮箱组
func (c *WebreportController) getById() {
	id, err := c.GetInt64("id")
	if err != nil {
		logs.Warn("/service/getById接口 参数异常, err: %v", err)
		c.ErrorJson(-1, "参数异常", nil)
	}
	logs.Info("请求参数: id=%v", id)
	emailMongo := models.EmailMongo{}
	email, err := emailMongo.QueryById(id)
	if err != nil {
		c.ErrorJson(-1, "服务查询数据异常", nil)
	}
	c.Data["data"] = email
	//c.Data["a"] = &email
	c.TplName = "web_email_detail.html"

	//c.SuccessJson(email)
}

//编辑邮件组
func (c *WebreportController) edit() {
	id, err := c.GetInt64("id")
	if err != nil {

		c.ErrorJson(-1, "参数异常", nil)
	}
	logs.Info("请求参数: id=%v", id)
	emailMongo := models.EmailMongo{}
	email, err := emailMongo.QueryById(id)
	if err != nil {
		c.ErrorJson(-1, "服务查询数据异常", nil)
	}
	c.Data["data"] = email
	//c.Data["a"] = &email
	c.TplName = "web_email_detail.html"

	//c.SuccessJson(email)
}

//查询所有邮箱组信息
func (c *WebreportController) QueryEmail() {
	userId, _ := c.GetSecureCookie(constant.CookieSecretKey, "user_id")
	businessList := []int{}
	businessMap := controllers.GetBusinesses(userId)
	for _, business := range businessMap {
		for k, v := range business {
			if k == "code" {
				businessList = append(businessList, v.(int))
			}
		}
	}
	var rp = models.EmailMongo{}
	page, _ := strconv.Atoi(c.GetString("page"))
	limit, _ := strconv.Atoi(c.GetString("limit"))
	result, count, err := rp.QueryAll(nil, "", page, limit)
	if err != nil {
		c.FormErrorJson(-1, "获取报告列表数据失败")
	}
	c.FormSuccessJson(count, result)
}

//根据id查询报告详细信息
func (c *WebreportController) Queryid() {
	db := GetLink()
	defer db.Close()
	id, _ := c.GetInt64("id")
	println(id)
	results, err := db.Query("select * from project where id = ?", id)
	var ary []DataSt
	for results.Next() {
		var Project DataSt
		err = results.Scan(&Project.Id, &Project.Name, &Project.Describe, &Project.Xmyl, &Project.Jszb, &Project.Sm, &Project.Fx, &Project.Zb, &Project.Recipient)
		//err = results.Scan(&Project)
		if err != nil {
			panic(err.Error())
		}
		ary = append(ary, Project)
	}
	//c.FormSuccessJson(1, ary)
	c.TplName = "web_report_detail.html"
	res := ary[0]
	c.Data["data"] = res
}

//分页查询所有报告
func (c *WebreportController) Query() {
	db := GetLink()
	defer db.Close()
	page, _ := strconv.Atoi(c.GetString("page"))
	limit, _ := strconv.Atoi(c.GetString("limit"))

	// 定义sql查询条件
	conditionSql := "from project order by id desc"

	// 拼接分页查询sql
	sql := fmt.Sprintf("select * "+conditionSql+" limit %d offset %d", limit, (page-1)*limit)
	results, err := db.Query(sql)
	//println(results)
	fmt.Println(results)
	if err != nil {
		panic(err.Error())
	}

	//var ret Ret
	var ary []DataSt
	for results.Next() {
		var Project DataSt
		err = results.Scan(&Project.Id, &Project.Name, &Project.Describe, &Project.Xmyl, &Project.Jszb, &Project.Sm, &Project.Fx, &Project.Zb, &Project.Recipient)
		//err = results.Scan(&Project)
		if err != nil {
			panic(err.Error())
		}
		ary = append(ary, Project)
	}

	// 拼接查询总条数sql
	var ct int64
	count, err := db.Query("select count(*) " + conditionSql)
	for count.Next() {
		err = count.Scan(&ct)
		if err != nil {
			panic(err.Error())
		}
	}

	c.FormSuccessJson(ct, ary)
}

//插入语句
func insert(name, desc, xmyl, jszb, fx, sm, zb, recipient string) error {
	db := GetLink()
	defer db.Close()
	sql := "INSERT INTO `project`( `name`,`describe`,`xmyl`,`jszb`,`fx`,`sm`,`zb`,`recipient`) VALUES ('" + name + "', '" + desc + "','" + xmyl + "','" + jszb + "','" + fx + "','" + sm + "','" + zb + "','" + recipient + "');"
	res, err := db.Exec(sql)
	if err != nil {
		fmt.Println(err.Error())
	}
	if i, _ := res.RowsAffected(); i == 0 {
		//return errors.New("insert error.")
		println("insert error")
	}
	return nil
}

//遍历切片是否包含某个值
func in(target string, str_array []string) bool {

	for _, element := range str_array {

		if target == element {

			return true
		}
	}
	return false
}

func (c *WebreportController) check() {
	userId, _ := c.GetSecureCookie(constant.CookieSecretKey, "user_id")
	us := "2014@xiaochuankeji.cn"
	n := DataSt{}
	if err := c.ParseForm(&n); err != nil { //传入user指针
		c.Ctx.WriteString("出错了！")
	}
	parts := strings.Split(n.Recipient, ";")
	fmt.Println(parts)
	userId = userId + us
	result := in(userId, parts)
	fmt.Println(result)
	if result == true {
		return
	}
	if result == false {
		return
	}

}

//添加测试报告
func (c *WebreportController) Insert() {
	c.check()
	userId, _ := c.GetSecureCookie(constant.CookieSecretKey, "user_id")
	us := "2014@xiaochuankeji.cn"
	n := DataSt{}
	if err := c.ParseForm(&n); err != nil { //传入user指针
		c.Ctx.WriteString("出错了！")
	}
	parts := strings.Split(n.Recipient, ";")
	fmt.Println(parts)
	userId = userId + us
	result := in(userId, parts)
	fmt.Println(result)
	if result == true {
		if len(n.Describe) > 0 && len(n.Name) > 0 {
			// 插入数据操作
			insert(n.Name, n.Describe, n.Xmyl, n.Jszb, n.Fx, n.Sm, n.Zb, n.Recipient)
			log.Println("插入成功")
			//插入成功后把\n换成<br>方便前端发送邮件解析
			n.Describe = strings.Replace(n.Describe, "\n", "</br>", -1)
			n.Xmyl = strings.Replace(n.Xmyl, "\n", "</br>", -1)
			n.Jszb = strings.Replace(n.Jszb, "\n", "</br>", -1)
			n.Fx = strings.Replace(n.Fx, "\n", "</br>", -1)
			n.Zb = strings.Replace(n.Zb, "\n", "</br>", -1)
			n.Sm = strings.Replace(n.Sm, "\n", "</br>", -1)
			n.Recipient = strings.Replace(n.Recipient, "\n", "</br>", -1)
			// 发送邮件
			SendEmail(n, userId)
			fmt.Println(userId)
			//c.Ctx.WriteString(n.Name + ",插入成功！")
			c.SuccessJson(nil)
		} else {
			log.Println("插入失败")
		}

	}
	if result == false {
		fmt.Println("收件人无自己~")
		n.Recipient = userId + ";" + n.Recipient
		if len(n.Describe) > 0 && len(n.Name) > 0 {
			// 插入数据操作
			insert(n.Name, n.Describe, n.Xmyl, n.Jszb, n.Fx, n.Sm, n.Zb, n.Recipient)
			log.Println("插入成功")
			//插入成功后把\n换成<br>方便前端发送邮件解析
			n.Describe = strings.Replace(n.Describe, "\n", "</br>", -1)
			n.Xmyl = strings.Replace(n.Xmyl, "\n", "</br>", -1)
			n.Jszb = strings.Replace(n.Jszb, "\n", "</br>", -1)
			n.Fx = strings.Replace(n.Fx, "\n", "</br>", -1)
			n.Zb = strings.Replace(n.Zb, "\n", "</br>", -1)
			n.Sm = strings.Replace(n.Sm, "\n", "</br>", -1)
			n.Recipient = strings.Replace(n.Recipient, "\n", "</br>", -1)
			// 发送邮件
			SendEmail(n, userId)
			fmt.Println(userId)
			//c.Ctx.WriteString(n.Name + ",插入成功！")
			c.SuccessJson(nil)
		} else {
			log.Println("插入失败")
		}
	}
}

func SendEmail(n DataSt, userid string) {

	user := "noapply2014@xiaochuankeji.cn"
	password := "X6uQ4R1u"
	host := "smtp.exmail.qq.com:25"
	to := n.Recipient
	//切分to，根据；遍历to是否有userid
	//to.

	//to := "fengmanlong2014@xiaochuankeji.cn;xueyibing2014@xiaochuankeji.cn;sunzhiying2014@xiaochuankeji.cn"

	subject := "辛苦查收测试报告"

	body := `
		<html>
		<body>
<style type="text/css">
    table.gridtable {
        font-family: verdana,arial,sans-serif;
        width: 70%;
        margin: 0 auto;
        font-size:11px;
        color:#333333;
        border-width: 1px;
        border-color: #666666;
        border-collapse: collapse;
        border-top: #CCCCCC solid 5px;
        border-left: #CCCCCC solid 5px;
        border-right: #CCCCCC solid 5px;
        border-bottom: #CCCCCC solid 5px;

    }
    table.gridtable th {
        width: 120px;
        border-width: 1px;
        padding: 8px;
        border-style: solid;
        border-color: #666666;
        background-color: #dedede;
        border-top: #CCCCCC solid 5px;
        border-left: #CCCCCC solid 5px;
        /*border-right: #F3F3F3 solid 5px;*/
        border-bottom: #CCCCCC solid 5px;
    }
    table.gridtable td {
        border-width: 1px;
        padding: 8px;
        border-style: solid;
        border-color: #666666;
        background-color: #ffffff;
        /*border-top: #CCCCCC solid 5px;*/
        border-left: #CCCCCC solid 5px;
        /*border-right: #CCCCCC solid 5px;*/
        /*border-bottom: #F3F3F3 solid 5px;*/
    }
</style>
<table class="gridtable" border="1px">
    <tr>
        <td colspan="2" style="border-right:#000000 solid 0px;text-align: center"><font color="red" size="5">` + n.Name + `—测试报告</font></td>
    </tr>
    <tr>
        <th >质量说明</th>
        <td>` + n.Describe + `</td>
    </tr>
    <tr>
        <th >项目遗留问题</th>
        <td>` + n.Xmyl + `</td>
    </tr>
    <tr>
        <th >技术指标</th>
        <td>` + n.Jszb + `</td>
    </tr>
    <tr>
        <th >发布风险及灰度计划</th>
        <td>` + n.Fx + `</td>
    </tr>
    <tr>
        <th >具体质量指标</th>
        <td>` + n.Zb + `</td>
    </tr>
    <tr>
        <th>其他说明</th>
        <td style="border-bottom: #CCCCCC solid 1px;">` + n.Sm + `</td>
    </tr>
	<tr>
        <td colspan="2" style="border-right:#000000 solid 0px;text-align: center;color:#FF0000">注：本测试报告由测试平台自动发送，不尽事宜联系测试负责人。</td>
    </tr>
</table>
</body>
		</html>
		`
	fmt.Println("send email")
	err := SendToMail(user, password, host, to, subject, body, "html")
	if err != nil {
		fmt.Println("Send mail error!")
		fmt.Println(err)
	} else {
		fmt.Println("Send mail success!")
	}
}

//发送邮件
func SendToMail(user, password, host, to, subject, body, mailtype string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])
	var content_type string
	if mailtype == "html" {
		content_type = "Content-Type: text/" + mailtype + "; charset=UTF-8"
	} else {
		content_type = "Content-Type: text/plain" + "; charset=UTF-8"
	}

	msg := []byte("To: " + to + "\r\nFrom: " + user + "\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, send_to, msg)
	return err
}

//数据库配置
func GetLink() *sql.DB {
	// sql.Open的第一个参数是driver名称，第二个参数是driver连接数据库的信息，各个driver可能不同。
	// DB不是连接，并且只有当需要使用时才会创建连接，如果想立即验证连接，需要用Ping()方法
	// 172.16.2.86
	db, err := sql.Open("mysql", "root:12345678@tcp(127.0.0.1:3306)/test")
	//db, err := sql.Open("mysql", "root:12345678@tcp(127.0.0.1:3306)/go_autoapi")
	if err != nil {
		fmt.Println(err)
	}
	// Ping验证与数据库的连接仍然存在，必要时建立连接。
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	}
	return db
}
