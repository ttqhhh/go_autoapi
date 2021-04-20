package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/go-ldap/ldap/v3"
	constant "go_autoapi/constants"
	"go_autoapi/models"
	"gopkg.in/mgo.v2"
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
	userDN := sr.Entries[0].DN
	err = conn.Bind(userDN, u.Passwrod)
	if err != nil {
		logs.Error("password does not exist or too many entries returned")
		c.ErrorJson(-1, "登录失败", nil)
	}
	now := time.Now()
	timestamp := now.Format(constant.TimeFormat)
	au := models.AutoUser{CreatedAt: timestamp, UpdatedAt: timestamp, Id: models.GetId("user_id"), UserName: u.UserName, Email: u.UserName + "2014@xiaochuankeji.cn"}
	loginUser, err := au.GetUserByName(u.UserName)
	if err == mgo.ErrNotFound {
		err = au.InsertUser(au)
		if err != nil {
			logs.Error("failed to store in db")
			c.ErrorJson(-1, "登录失败", nil)
		}
	}
	c.Ctx.SetSecureCookie(constant.CookieSecretKey, "user_id", u.UserName)
	c.Ctx.SetSecureCookie(constant.CookieSecretKey, "user_type", string(loginUser.Business))
	c.SuccessJson(au)
}

func (c *AutoTestController) logout() {
	c.SuccessJson(map[string]string{"location": "/login"})
}
