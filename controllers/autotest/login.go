package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/go-ldap/ldap"
)

//url = 172.16.1.233:10389
//baseDn = ou = People, dc = xiaochuankeji, dc= cn
//fiter = uid
//adminDn = uid= admin, ou = system
//adminPass = FSbxiULPVB8kyUUk

type User struct {
	UserName string
	Passwrod string
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
	searchRequest := ldap.NewSearchRequest(searchDn, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=inetOrgPerson)(mail=%s))", userName),
		[]string{"dn"},
		nil, )
	sr, err := conn.Search(searchRequest)
	if err != nil {
		logs.Error("request ldap error:%v\n", err)
		return
	}
	if len(sr.Entries) != 1 {
		logs.Error("User does not exist or too many entries returned")
	}
	userDN := sr.Entries[0].DN
	err = conn.Bind(userDN, password)
	if err != nil {
		logs.Error("password does not exist or too many entries returned")
	}
	logs.Error("success")
}

//func (c *AutoTestController) login() {
//
//	now := time.Now()
//	id := models.GetId("case")
//	acm := models.AutoCaseMongo{Id: id, CreatedAt: now, UpdatedAt: now}
//	if err := json.Unmarshal(c.Ctx.Input.RequestBody, &acm); err != nil {
//		c.ErrorJson(-1, "请求错误", nil)
//	}
//	err := acm.InsertCase(acm)
//	if err != nil {
//		fmt.Println(err)
//		c.ErrorJson(-1, "请求错误", nil)
//	}
//	c.SuccessJson("添加成功")
//}
