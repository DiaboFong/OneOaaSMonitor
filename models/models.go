package models

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var (
	//public
	Orm orm.Ormer

	//db conf
	dbtype string
	dbuser string
	dbpass string
	dbhost string
	dbport string
	dbname string
	dsn    string

	//mysql
	maxIdle int
	maxConn int
)

//注册模型
func setModels() {
	orm.RegisterModel(new(User))
	orm.RegisterModel(new(SmsCode))
	orm.RegisterModel(new(License))

}

func setMysql() {
	dbtype = beego.AppConfig.String("dbtype")
	dbuser = beego.AppConfig.String("dbuser")
	dbpass = beego.AppConfig.String("dbpass")
	dbhost = beego.AppConfig.String("dbhost")
	dbport = beego.AppConfig.String("dbport")
	dbname = beego.AppConfig.String("dbname")
	switch dbtype {
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", dbuser, dbpass, dbhost, dbport, dbname) + "&loc=" + url.QueryEscape("Local")
		break
	case "postgres":
		dsn = fmt.Sprintf("postgres://%s:%s@:%s/%s?sslmode=%s&host=%s", url.QueryEscape(dbuser), url.QueryEscape(dbpass), dbport, dbname, "disable", dbhost)
		break
	case "sqlite3":
		dsn = fmt.Sprintf("./db/oneoaas_web.db")
	default:
		beego.Error("不支持此类型数据库")
	}

	maxIdle = 30
	maxConn = 50
}

func setBeego() {
	beego.BConfig.WebConfig.Session.SessionProviderConfig = dsn
}

func init() {
	setMysql()
	setModels()
	//setBeego()

	switch dbtype {
	case "mysql":
		orm.RegisterDriver("mysql", orm.DRMySQL)
	case "postgres":
		orm.RegisterDriver("postgres", orm.DRPostgres)
	case "sqlite3":
		orm.RegisterDriver("sqlite3", orm.DRSqlite)
		//orm.RegisterDataBase("default", "sqlite3", "./conf/slaver.db")
	default:
		beego.Info("不支持", dbtype, "数据库")
	}
	orm.RegisterDataBase("default", dbtype, dsn, maxIdle, maxConn)
	orm.Debug = true

	//orm.RegisterDataBase("default", dbtype, "oneoaas_event:oneoaas_event@tcp(127.0.0.1:3306)/oneoaas_event?charset=utf8&loc=Local", maxIdle, maxConn)
	//force=true 会先删除表,后重建表
	//verbose=true 显示执行信息
	orm.RunSyncdb("default", false, true)
	Orm = orm.NewOrm()
	Orm.Using("default") // 默认使用 default，你可以指定为其他数据库
}

func dbErrorParse(err string) (string, int64) {
	Parts := strings.Split(err, ":")
	errorMessage := Parts[1]
	Code := strings.Split(Parts[0], "Error ")
	errorCode, _ := strconv.ParseInt(Code[1], 10, 32)
	return errorMessage, errorCode
}
