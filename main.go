package main

import (
	"api_wx_klagri_com_cn_go/models"
	_ "api_wx_klagri_com_cn_go/routers"
	"api_wx_klagri_com_cn_go/util"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

func init() {
	models.Init()
	// 想模板中注册函数（首字母大写）
	beego.AddFuncMap("ToUpper", util.FirstUpper)
	beego.BConfig.WebConfig.Session.SessionOn = true
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
		orm.Debug = true
	}

	beego.Run()
}
