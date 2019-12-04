package controllers

import (
	"api_wx_klagri_com_cn_go/models"
	"api_wx_klagri_com_cn_go/util"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
	"time"
)

type baseController struct {
	beego.Controller
	log            *logs.BeeLogger
	o              orm.Ormer
	controllerName string
	actionName     string
	apiLogModel    models.ApiLog
	startTime      int64
}

// 控制器执行前会运行此方法
func (p *baseController) Prepare() {
	// 获取控制器名称与方法名称
	controllerName, actionName := p.GetControllerAndAction()
	p.controllerName = strings.ToLower(controllerName[0 : len(controllerName)-10])
	p.actionName = strings.ToLower(actionName)

	// 获取BeeLogger
	p.log = logs.NewLogger(10000)
	// 控制台打印日志
	p.log.SetLogger("console")
	// 输出文件打印 输出打印
	// p.l.SetLogger(logs.AdapterFile, `{"filename":"project.log","level":7,"maxlines":0,"maxsize":0,"daily":true,"maxdays":10,"color":true}`)

	// TODO 记录API日志到数据库
	var (
		query interface{}
	)
	// 以json格式记录参数, query
	// type query = interface{}
	if p.Ctx.Input.Method() == "GET" {
		params := p.Input()
		newParams := make(map[string]interface{})
		for key, val := range params {
			if len(val) < 2 {
				newParams[key] = val[0]
			} else {
				newParams[key] = val
			}
		}
		query, _ = json.Marshal(newParams) // 获取query参数

	} else {
		query = p.Ctx.Input.RequestBody // 获取Body中的参数

	}

	p.apiLogModel.Params = string(query.([]byte))              // 将[]byte转字符串
	p.apiLogModel.Name = p.Ctx.Input.URL()                     // 获取访问路径： /v1/weather/future
	p.apiLogModel.Method = p.Ctx.Input.Method()                // 获取访问方法， GET/POST
	p.apiLogModel.Server = p.Ctx.Input.Domain()                // 获取访问本服务的域名，当多个服务的日志记录到同一张表时区分
	p.apiLogModel.Client = p.Ctx.Input.Refer()                 // 访问接口的来源地址
	tem, _ := json.Marshal(p.Ctx.Input.Context.Request.Header) // 请求头
	p.apiLogModel.Header = string(tem)
	p.apiLogModel.CreateTime = time.Now() // 请求时间
	p.apiLogModel.StartTime, _ = strconv.Atoi(strconv.FormatInt(time.Now().Unix(), 10))
	p.startTime = time.Now().UnixNano() / 1e6
	fmt.Println(1111111111111)
	p.apiLogModel.Platform = "测试" //	接口所属平台，类似 p.apiLogModel.Server

}

// Title   API日志
// Description 记录日志到MySQL数据库
// Param   result
// return  bool
func (p *baseController) ApiLog(result map[string]interface{}) {
	res, _ := json.Marshal(result)
	p.apiLogModel.Response = string(res) // 响应结果
	lifeTime := time.Now().UnixNano()/1e6 - p.startTime
	p.apiLogModel.Life, _ = strconv.Atoi(strconv.FormatInt(lifeTime, 10)) //生命周期
	id, err := models.AddApiLog(&p.apiLogModel)
	if err == nil {
		fmt.Println("插入日志成功")
		var logId = id
		fmt.Println(logId)
	}
}

// @Title   成功返回json
// @Description 对返回数据进行统一的格式化之后返回，并终止生命周期
// @Param	response	返回数据，任意类型, 可以是 map["code":200, "msg":"成功", "data":任意类型]， 也可以直接是任意类型的数据，但将将会使用默认的code以及msg
// @response map["code": 200, "msg":"成功", "data":[]]
func (p *baseController) SuccessJson(response map[string]interface{}) {
	var result = make(map[string]interface{})
	result["code"] = 200
	result["msg"] = "成功"
	_, ok := response["data"].(interface{})
	if ok {
		// response 中可以重写 code 与 msg
		result = util.ArrayMerge(result, response)
	} else {
		result["data"] = response
	}
	// 如果返回数据为nil，则json后的data为[]
	if response == nil {
		result["data"] = [...]interface{}{}
	}
	p.Data["json"] = result
	// 返回json
	p.ServeJSON()

	// TODO 记录日志
	p.ApiLog(result)

	// 停止声明周期
	p.StopRun()
}

// @Title   成功返回json
// @Description 对返回数据进行统一的格式化之后返回，并终止生命周期
// @Param	response	返回数据，任意类型, 可以是 map["code":200, "msg":"成功", "data":任意类型]， 也可以直接是任意类型的数据，但将将会使用默认的code以及msg
// @response map["code": 200, "msg":"成功", "data":[]]
func (p *baseController) ErrorJson(response map[string]interface{}) {
	param, err := json.Marshal(response)
	if err != nil {
		p.log.Debug("有错误")
	}
	p.log.Debug(string(param))

	var result = make(map[string]interface{})
	result["code"] = -100
	result["msg"] = "查询失败"
	_, ok := response["data"].(interface{})
	if ok {
		// response 中可以重写 code 与 msg
		result = util.ArrayMerge(result, response)
	} else {
		result["data"] = response
	}
	// 如果返回数据为nil，则json后的data为[]
	if response == nil {
		result["data"] = [...]interface{}{}
	}
	p.Data["json"] = result
	// 返回json
	p.ServeJSON()

	// TODO 记录日志
	p.ApiLog(result)

	// 停止声明周期
	p.StopRun()
}

func ValidateResult() {

}
