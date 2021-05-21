package api

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/go-ldap/ldap/v3"
	constant "go_autoapi/constants"
	"go_autoapi/libs"
	"go_autoapi/models"
	"go_autoapi/utils"
	"gopkg.in/mgo.v2"
	"time"
)

type ApiController struct {
	libs.BaseController
}

func (c *ApiController) Get() {
	do := c.GetMethodName()
	switch do {
	default:
		logs.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持", nil)
	}
}

func (c *ApiController) Post() {
	do := c.GetMethodName()
	switch do {
	case "login":
		c.login()
	default:
		logs.Warn("action: %s, not implemented", do)
		c.ErrorJson(-1, "不支持", nil)
	}
}

type User struct {
	UserName string `json:"user_name"`
	Passwrod string `json:"password"`
}

func (c *ApiController) login() {
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
		autoUserId, err := r.Incr(constant.AUTO_USER_PRIMARY_KEY).Result()
		au := models.AutoUser{CreatedAt: timestamp, UpdatedAt: timestamp, Id: autoUserId, UserName: u.UserName, Email: u.UserName + "2014@xiaochuankeji.cn"}
		err = au.InsertUser(au)
		if err != nil {
			logs.Error("failed to store in db")
			c.ErrorJson(-1, "登录失败", nil)
		}
	}
	c.Ctx.SetSecureCookie(constant.CookieSecretKey, "user_id", u.UserName)
	c.Ctx.SetSecureCookie(constant.CookieSecretKey, "user_type", string(loginUser.Business))

	c.SuccessJsonWithMsg(nil, "登录成功")
}
