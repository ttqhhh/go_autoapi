package web_report

import (
	"database/sql"
	"fmt"
	"github.com/astaxie/beego/logs"
	"go_autoapi/libs"
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
	case "queryAll":
		c.AllWebReport()
	case "query":
		c.Query()
	case "queryId":
		c.Queryid()
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
	case "queryAll":
		c.Query()
	case "query":
		c.AllWebReport()

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

func (c *WebreportController) Insert() {
	n := DataSt{}
	if err := c.ParseForm(&n); err != nil { //传入user指针
		c.Ctx.WriteString("出错了！")
	}
	if len(n.Describe) > 0 && len(n.Name) > 0 {
		// 插入数据操作
		insert(n.Name, n.Describe, n.Xmyl, n.Jszb, n.Fx, n.Sm, n.Zb, n.Recipient)
		log.Println("插入成功")
		n.Describe = strings.Replace(n.Describe, "\n", "</br>", -1)
		n.Xmyl = strings.Replace(n.Xmyl, "\n", "</br>", -1)
		n.Jszb = strings.Replace(n.Jszb, "\n", "</br>", -1)
		n.Fx = strings.Replace(n.Fx, "\n", "</br>", -1)
		n.Zb = strings.Replace(n.Zb, "\n", "</br>", -1)
		n.Sm = strings.Replace(n.Sm, "\n", "</br>", -1)
		n.Recipient = strings.Replace(n.Recipient, "\n", "</br>", -1)
		// 发送邮件
		SendEmail(n)
		//c.Ctx.WriteString(n.Name + ",插入成功！")
		c.SuccessJson(nil)
	} else {
		log.Println("插入失败")
	}

}

func SendEmail(n DataSt) {
	user := "noapply2014@xiaochuankeji.cn"
	password := "X6uQ4R1u"
	host := "smtp.exmail.qq.com:25"
	to := n.Recipient
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
	db, err := sql.Open("mysql", "root:@tcp(172.16.2.86:3306)/test")
	//db, err := sql.Open("mysql", "root:12345678@tcp(127.0.0.1:3306)/test")
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
