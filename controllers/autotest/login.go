package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/beego/beego/v2/server/web"
	"github.com/go-ldap/ldap/v3"
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

	//searchRequest := ldap.NewSearchRequest(searchDn, ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
	//	fmt.Sprintf("(&(objectClass=inetOrgPerson)(mail=%s))", u.UserName),
	//	[]string{"dn"},
	//	nil, )

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

	for _, entry := range sr.Entries {
		fmt.Printf("%s: %v\n", entry.DN, entry.GetAttributeValue("cn"))
	}
	userDN := sr.Entries[0].DN
	err = conn.Bind(userDN, u.Passwrod)
	if err != nil {
		logs.Error("password does not exist or too many entries returned")
		c.ErrorJson(-1, "登录失败", nil)
	}
	c.Ctx.SetSecureCookie(sercetKey, "userid", "liuweiqiang")
	c.SuccessJson("登录成功")
}
