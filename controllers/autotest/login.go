package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/go-ldap/ldap/v3"
	constant "go_autoapi/constants"
	"go_autoapi/models"
	"go_autoapi/utils"
	"gopkg.in/mgo.v2"
	"strconv"
	"time"
)

//url = 172.16.1.233:10389
//baseDn = ou = People, dc = xiaochuankeji, dc= cn
//fiter = uid
//adminDn = uid= admin, ou = system
//adminPass = FSbxiULPVB8kyUUk

type User struct {
	UserName string `json:"user_name"`
	Passwrod string `json:"password"`
}

func (c *AutoTestController) toLogin() {
	c.TplName = "login.html"
}

// Login 登录
func (c *AutoTestController) login() {
	u := User{}
	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &u); err != nil {
		c.ErrorJson(-1, "请求错误", nil)
	}
	addr, _ := web.AppConfig.String("addr")
	searchDn, _ := web.AppConfig.String("searchDn")
	conn, err := ldap.Dial("tcp", addr)
	if err != nil {
		logs.Error("getbasedn error:%v\n", err)
		return
	}
	defer conn.Close()

	searchRequest := ldap.NewSearchRequest(
		searchDn, // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=organizationalPerson)(uid=%s))", u.UserName), // The filter to apply
		[]string{"dn", "cn"}, // A list attributes to retrieve
		nil,
	)
	sr, err := conn.Search(searchRequest)
	if err != nil {
		fmt.Println(err)
	}
	if len(sr.Entries) < 1 {
		logs.Error("user does not exist")
		c.ErrorJson(-1, "用户名或密码错误，登录失败", nil)
	}
	userDN := sr.Entries[0].DN
	err = conn.Bind(userDN, u.Passwrod)
	if err != nil {
		logs.Error("password does not exist or too many entries returned")
		c.ErrorJson(-1, "用户名或密码错误，登录失败", nil)
	}
	now := time.Now()
	timestamp := now.Format(constant.TimeFormat)
	au := models.AutoUser{}
	loginUser, err := au.GetUserInfoByName(u.UserName)
	if err == mgo.ErrNotFound {
		r := utils.GetRedis()
		defer r.Close()
		autoUserId, err := r.Incr(constant.AUTO_USER_PRIMARY_KEY).Result()
		au := models.AutoUser{CreatedAt: timestamp, UpdatedAt: timestamp, Id: autoUserId, UserName: u.UserName, Email: u.UserName + "2014@xiaochuankeji.cn"}
		err = au.InsertUser(au)
		if err != nil {
			logs.Error("failed to store in db")
			c.ErrorJson(-1, "登录失败", nil)
		}
	}
	c.Ctx.SetSecureCookie(constant.CookieSecretKey, "user_id", u.UserName, 31536000)
	c.Ctx.SetSecureCookie(constant.CookieSecretKey, "user_type", string(loginUser.Business), 31536000)
	var ul []*models.AutoUser
	for i := 0; i < 10; i++ {
		ul = append(ul, &au)
	}
	//c.SuccessJsonWithMsg(ul, "OK")
	// 默认跳转到第一个有权限的业务线case页面
	businesses := GetBusinesses(u.UserName)
	business := businesses[0]
	code := business["code"]
	codeInt := code.(int)
	redirectUrl := "/case/show_cases?business=" + strconv.Itoa(codeInt)
	//if !isTester(u.UserName) {
	//	redirectUrl = "/flowreplay/index"
	//}
	//redirectUrl := "/case/show_cases?business=ZuiyYou"
	c.SuccessJson(redirectUrl)
}

func (c *AutoTestController) logout() {
	c.SuccessJsonWithMsg(map[string]string{"location": "/login"}, "OK")
}

//func (c *AutoTestController) loginUserDetail() {
//	userId, _ := c.GetSecureCookie(constant.CookieSecretKey, "user_id")
//	res := make(map[string]interface{})
//	res["username"] = userId
//	res["isTester"] = isTester(userId)
//	c.SuccessJson(res)
//}
//
//func isTester(username string) bool {
//	alltester := append(TesterUserList, ShixiTesterUserList...)
//	alltester = append(alltester, WaibaoTesterUserList...)
//	isTester := false
//	for _, tester := range alltester {
//		if username == tester {
//			isTester = true
//		}
//	}
//	return isTester
//}
//
//var TesterUserList = []string{
//	// 4个超级管理员
//	"fengmanlong",
//	//"xueyibing",
//	"wangzhen01",
//	"wangjun",
//	"sunzhiying",
//	"zhangjuan",
//	"liuyinan",
//	"sunmingyao",
//	"zhangdanbing",
//	"xufei",
//	"suhuimin",
//	"wanglanjin",
//	"yangjingfang",
//	"chengxiaojing",
//	"zhaoxiaodong",
//}
//
//var ShixiTesterUserList = []string{
//	"tangtianqing",
//}
//
//var WaibaoTesterUserList = []string{
//	"zhaomengxin",
//	"huoyongjian",
//}
