package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	// "strings"
)

// type baseModel struct {
// 	o orm.Ormer
// }

func Init() {
	dbhost := beego.AppConfig.String("dbhost")
	dbport := beego.AppConfig.String("dbport")
	dbuser := beego.AppConfig.String("dbuser")
	dbpassword := beego.AppConfig.String("dbpassword")
	dbname := beego.AppConfig.String("dbname")
	if dbport == "" {
		dbport = "3306"
	}
	dsn := dbuser + ":" + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?charset=utf8&loc=Asia%2FShanghai"
	orm.RegisterDataBase("default", "mysql", dsn)
}

// 返回带前缀的表名
// func TableName(str string) string {
// 	prefix := beego.AppConfig.String("dbprefix")
// 	length := strings.Count(str, "") - 1
// 	if str[0:length-1] != prefix {
// 		return prefix + str
// 	}
// 	return str
// }
