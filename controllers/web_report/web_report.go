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
		err = results.Scan(&Project.Id, &Project.Name, &Project.Describe, &Project.Xmyl, &Project.Jszb, &Project.Sm, &Project.Fx, &Project.Zb)
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
		err = results.Scan(&Project.Id, &Project.Name, &Project.Describe, &Project.Xmyl, &Project.Jszb, &Project.Sm, &Project.Fx, &Project.Zb)
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
func insert(name, desc, xmyl, jszb, fx, sm, zb string) error {
	db := GetLink()
	defer db.Close()
	sql := "INSERT INTO `project`( `name`,`describe`,`xmyl`,`jszb`,`fx`,`sm`,`zb`) VALUES ('" + name + "', '" + desc + "','" + xmyl + "','" + jszb + "','" + fx + "','" + sm + "','" + zb + "');"
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
		insert(n.Name, n.Describe, n.Xmyl, n.Jszb, n.Fx, n.Sm, n.Zb)
		log.Println("插入成功")
		// 发送邮件
		SendEmail(n)
		//c.Ctx.WriteString(n.Name + ",插入成功！")
		c.SuccessJson(nil)
	} else {
		log.Println("插入失败")
	}

}

func SendEmail(n DataSt) {
	user := "sunzhiying2014@xiaochuankeji.cn"
	password := "Szy0204."
	host := "smtp.exmail.qq.com:25"
	to := n.Recipient
	//to := "fengmanlong2014@xiaochuankeji.cn;xueyibing2014@xiaochuankeji.cn;sunzhiying2014@xiaochuankeji.cn"

	subject := "辛苦查收测试报告"

	body := `
		<html>
		<head>
			<style type="text/css">
        	table {
            	width: 50%;border-top: 0px solid #000;border-left: 0px solid #000;border-spacing: 0;
       		 }
        	table td {
            	border-bottom: 0px solid #000; border-right: 0px solid #000;
       		 }
   		 </style>
		</head>
		<body>
			<h3>
			测试报告
			</h3>
			<table border="1" >
    			<tr>
       			 	<td colspan="2" style="border-right:#000000 solid 0px;text-align: center">` + n.Name + `—测试报告</td>
				</tr>
				<tr>
					 <td >质量说明</td>
        			<td>` + n.Describe + `</td>
				</tr>
				<tr>
        			<td >项目遗留问题</td>
        		<td>` + n.Xmyl + `</td>
    			</tr>
    			<tr>
        			<td >技术指标</td>
        		<td>` + n.Jszb + `</td>
    			</tr>
    			<tr>
        		<td >发布风险及灰度计划</td>
        			<td>` + n.Fx + `</td>
    			</tr>
    			<tr>
        			<td >具体质量指标</td>
        			<td>` + n.Zb + `</td>
    			</tr>
    			<tr>
        			<td>其他说明</td>
        			<td>` + n.Sm + `</td>
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

	msg := []byte("To: " + to + "\r\nFrom: " + "平台发送" + "\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, send_to, msg)
	return err
}

//数据库配置
func GetLink() *sql.DB {
	// sql.Open的第一个参数是driver名称，第二个参数是driver连接数据库的信息，各个driver可能不同。
	// DB不是连接，并且只有当需要使用时才会创建连接，如果想立即验证连接，需要用Ping()方法
	db, err := sql.Open("mysql", "root:12345678@tcp(127.0.0.1:3306)/test")
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
