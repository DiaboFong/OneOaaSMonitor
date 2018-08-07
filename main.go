package main

import (
	"github.com/astaxie/beego"
	"oneoaas.com/oneoaas_web/controllers"
	_ "oneoaas.com/oneoaas_web/routers"
	//"github.com/astaxie/beego/plugins/auth"
	"github.com/astaxie/beego/session"
	_"github.com/astaxie/beego/session/mysql"
	"fmt"
)

func init() {
	beego.BConfig.WebConfig.Session.SessionProvider = "mysql"
	sessionConfig := new(session.ManagerConfig)
	providerCfg := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
		beego.AppConfig.String("dbuser"),
		beego.AppConfig.String("dbpass"),
		beego.AppConfig.String("dbhost"),
		beego.AppConfig.String("dbport"),
		beego.AppConfig.String("dbname"))

	sessionConfig.ProviderConfig = providerCfg

	//cookkie 默认一天
	sessionConfig.CookieLifeTime=int(1*60*60)

	sessionConfig.Gclifetime=3600

	sessionConfig.CookieName="oneoaasessionId"

	sessionConfig.EnableSetCookie=true

	beego.GlobalSessions,_= session.NewManager("mysql",sessionConfig )
	go beego.GlobalSessions.GC()
}

func main() {
	//beego.InsertFilter("/", beego.BeforeRouter, filter.FilterIndex)
	// 自定义错误页面
	beego.ErrorController(&controllers.ErrorController{})
	//beego.InsertFilter("/users", beego.BeforeRouter,auth.Basic("oneoaas","IRS5NnnLLZS19qrM1#opdu*GWFS1tq"))
	//beego.InsertFilter("/vendor", beego.BeforeRouter,auth.Basic("oneoaas","oneoaas"))


	// 下载路径
	beego.Run()
}
