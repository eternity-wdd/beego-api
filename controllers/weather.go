package controllers

import (
	"api_wx_klagri_com_cn_go/models"
	"api_wx_klagri_com_cn_go/services"
	"api_wx_klagri_com_cn_go/util"
	// "encoding/json"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"github.com/eternity-wdd/solarlunar"
	"github.com/qichengzx/coordtransform"
	// "reflect"
	// "regexp"
	// "encoding/json"
	"strconv"
	"strings"
	"time"
)

// Operations about Weather
type WeatherController struct {
	baseController
}

var snow = [...]string{"13", "14", "15", "16", "17", "26", "27", "28", "302"}
var response = make(map[string]interface{})

// @APIVersion 1.0.1
// @Title   POST参数测试
// @Description get the current location weather through latitude and longitude
// @Param	longitude,latitude		body  	json	false		"json参数测试"
// @Success 200 {"code": 200, "msg":"成功", "data":[]}
// @Failure 403 param is wrong
// @router / [Post]
func (c *WeatherController) Post() {

}

// @APIVersion 1.0.1
// @Title   获取及时天气预报
// @Description get the current location weather through latitude and longitude
// @Param	longitude		query  	string	false		"经度"
// @Param	latitude		query  	string	false		"纬度"
// @Success 200 {"code": 200, "msg":"成功", "data":[]}
// @Failure 403 param is wrong
// @router / [get]
func (c *WeatherController) Get() {
	data := make(map[string]interface{})
	// 计算当天的农历日期
	solarDate := time.Now().Format("2006-01-02")
	data["calendar"] = solarlunar.SolarToMonthDay(solarDate)

	// 接收经纬度，并转为 float64
	longitude, _ := strconv.ParseFloat(c.Ctx.Input.Query("longitude"), 64)
	latitude, _ := strconv.ParseFloat(c.Ctx.Input.Query("latitude"), 64)
	if longitude == 0 || latitude == 0 {
		response["code"] = 403
		response["msg"] = "参数错误"
		c.Data["json"] = response
		c.ServeJSON()
		c.StopRun()
		return
	}
	// 将GJC02坐标系转为WGS84
	longitude, latitude = coordTransform.GCJ02toWGS84(longitude, latitude)

	// 通过经纬度获取短临降水(两小时)
	rainUrl := "http://47.97.212.221:8082/?type=2&lon=" + strconv.FormatFloat(longitude, 'f', 6, 64) + "&lat=" + strconv.FormatFloat(latitude, 'f', 6, 64) + "&token=5ccc96af717842a5ad410a0ede8bfc6b"
	rain, err := httplib.Get(rainUrl).String()
	if err != nil {
		fmt.Println("调用短临降水接口失败")
	}

	// 测试数据
	// rain = "{\"reqTime\":\"201911141554\",\"startTime\":\"201911141555\",\"endTime\":\"201911141755\",\"series\":[0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.07,0.14,0.32,0.51,0.72,0.92,1.11,1.28,1.41,1.5,1.53,1.5,1.41,1.28,1.11,0.92,0.72,0.51,0.32,0.14,0.07,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0,0.0],\"amendNow\":\"多云\",\"amendNowwen\":\"Cloudy\",\"amendNowwcode\":\"01\",\"msg\":\"1小时后会有雪，出门记得带伞\"}"
	// fmt.Println(rain)
	// Json转map，  map[string]interface{}
	rainMap := util.JSONToMap(rain)
	// 处理接口调用失败的情况
	errorCode, ok := rainMap["errorCode"]
	if ok {
		response["code"] = errorCode
		response["msg"] = rainMap["errorInfo"]
		c.Data["json"] = response
		c.ServeJSON()
	}
	// 降雨量之和
	sum := 0.00
	// 遍历map中的series数组，降雨量相加
	for _, num := range rainMap["series"].([]interface{}) {
		sum += num.(float64) // 这里使用断言
	}
	data["rainfall"] = strconv.FormatFloat(sum, 'f', 2, 64) + " mm"
	data["weatherMsg"] = rainMap["msg"]
	// 根据经纬度获取该地的实时天气信息
	weatherUrl := "http://47.97.212.221:8082/?type=3&lon=" + strconv.FormatFloat(longitude, 'f', 6, 64) + "&lat=" + strconv.FormatFloat(latitude, 'f', 6, 64) + "&distance=100000&var=evp&var=p&var=ws_2mi_avg&var=t&var=rh&var=t&var=wp&var=wd_2mi_avg&var=ws_2mi_avg&var=t_min_24h&var=t_max_24h&var=staName&var=staCode&token=5ccc96af717842a5ad410a0ede8bfc6b"
	weather, err := httplib.Get(weatherUrl).String()
	if err != nil {
		fmt.Println("调用实时天气接口失败")
	}
	// json字符串转map
	weatherMap := util.JSONToMap(weather)
	// 处理当前信息未更新的情况，使用上一小时的数据
	errorCode, ok = weatherMap["errorCode"]
	if ok && errorCode == "DT01" {
		datatime := time.Now().Unix() - 3600
		x := time.Unix(datatime, 0)
		datatimeStr := x.Format("2006010215")
		weatherUrl += "&datatime=" + datatimeStr
		weather, err = httplib.Get(weatherUrl).String()
		if err != nil {
			fmt.Println("调用实时天气接口失败")
		}
		weatherMap = util.JSONToMap(weather)
		// 处理接口调用失败或者仍无数据时的情况
		errorCode, ok = weatherMap["errorCode"]
		if ok {
			c.Data["code"] = errorCode
			c.Data["msg"] = weatherMap["errorInfo"]
			c.ServeJSON()
		}
	}
	// 获取风向
	data["windDirection"] = util.WindDirection(weatherMap["wd_2mi_avg"].(float64))
	// 获取风力等级
	data["windLevel"] = strconv.FormatInt(util.WindSpeed(weatherMap["ws_2mi_avg"].(float64)), 10) + "级"
	// 温度
	data["t"] = strconv.FormatFloat(weatherMap["t"].(float64), 'f', 2, 64) + "℃"
	data["tMin"] = strconv.FormatFloat(weatherMap["t_min_24h"].(float64), 'f', 2, 64) + "℃"
	data["tMax"] = strconv.FormatFloat(weatherMap["t_max_24h"].(float64), 'f', 2, 64) + "℃"
	data["airPress"] = strconv.FormatFloat(weatherMap["p"].(float64), 'f', 2, 64) + "hpa"
	data["wet"] = strconv.FormatFloat(weatherMap["rh"].(float64), 'f', 2, 64) + "%"
	data["windSpeed"] = weatherMap["ws_2mi_avg"]
	data["evp"] = weatherMap["evp"]
	data["staCode"] = weatherMap["staCode"]

	// 通过经纬获取地名，腾讯地图
	txMapUrl := "https://apis.map.qq.com/ws/geocoder/v1/?location=" + strconv.FormatFloat(latitude, 'f', 2, 64) + "," + strconv.FormatFloat(longitude, 'f', 2, 64) + "&key=QLIBZ-QUVC6-QQES6-ECI37-BFIJ6-7HBBI&output=json" // 如果提示签名失败，去腾讯地图自己申请一个key即可
	cityInfo, err := httplib.Get(txMapUrl).String()
	if err != nil {
		fmt.Println("调用腾讯地图失败")
	}
	cityInfoMap := util.JSONToMap(cityInfo)
	cityName, ok := cityInfoMap["result"].(map[string]interface{})["address_component"].(map[string]interface{})["district"]

	if !ok || cityName == "" {
		cityName, ok = cityInfoMap["result"].(map[string]interface{})["address_component"].(map[string]interface{})["city"]
		if !ok || cityName == "" {
			cityName, ok = cityInfoMap["result"].(map[string]interface{})["address_component"].(map[string]interface{})["province"]
		}
		if !ok || cityName == "" {
			cityName, ok = cityInfoMap["result"].(map[string]interface{})["address_component"].(map[string]interface{})["nation"]
		}
	}
	data["cityName"] = cityName.(string)
	data["street"] = cityName.(string) + cityInfoMap["result"].(map[string]interface{})["address_component"].(map[string]interface{})["street_number"].(string)
	data["weawther"] = rainMap["amendNow"]

	// 根据气象编码拼接背景图片地址
	if rainMap["amendNowwcode"] == "00" {
		data["backPic"] = "/images/index/sunny.png"
	} else if util.Contain(rainMap["amendNowwcode"], snow) {
		data["backPic"] = "/images/index/snow.png"
	} else {
		data["backPic"] = "/images/index/cloudy.png"
	}

	response["code"] = 200
	response["msg"] = "查询成功"
	response["data"] = data
	// 接口成功统一返回
	c.SuccessJson(response)
}

// @APIVersion 1.0.0
// @Title   获取未来几天的天气情况
// @Description 通过经纬度和城市名称获取未来几天的天气情况
// @Param	area			query  	string  true		"行政区名称"
// @Param	longitude		query  	string	true		"经度"
// @Param	latitude		query  	string	true		"纬度"
// @Success 200 {"code": 200, "msg":"成功", "data":[]}
// @Failure 403 param is wrong
// @router /future [get]
func (c *WeatherController) Future() {
	area := c.Ctx.Input.Query("area")
	longitude, _ := strconv.ParseFloat(c.Ctx.Input.Query("longitude"), 64)
	latitude, _ := strconv.ParseFloat(c.Ctx.Input.Query("latitude"), 64)
	if area == "" || longitude == 0.00 || latitude == 0.00 {
		c.ErrorJson(map[string]interface{}{"code": 200, "msg": "参数错误", "data": map[string]interface{}{}})
	}
	// gcj02转wgs84
	longitude, latitude = coordTransform.GCJ02toWGS84(longitude, latitude)

	// 去除字符串中的省市区
	area = strings.TrimRight(area, "省")
	area = strings.TrimRight(area, "市")
	area = strings.TrimRight(area, "区")

	// 通过行政区名称，查询areaCode
	// areaCode := models.AreaCode{}
	areaInfo := models.QueryByArea(area)
	code, _ := areaInfo["code"]
	if code != 200 {
		c.ErrorJson(areaInfo)
	}
	// 从多次map中取值时，需要断言value是否是一个map
	data := areaInfo["data"].(map[string]interface{})
	areaCode := data["areaCode"]

	// 通过areaCode获取未来一周的天气预报
	start := time.Now().Unix()
	end := start + 86400*6

	startDate := time.Unix(start, 0).Format("20060102")
	endDate := time.Unix(end, 0).Format("20060102")
	futureUrl := "http://47.97.212.221:8082/?type=1&start=" + startDate + "&end=" + endDate + "&areaCode=" + areaCode.(string) + "&var=wp&token=5ccc96af717842a5ad410a0ede8bfc6b"
	future, err := httplib.Get(futureUrl).String()
	if err != nil {
		fmt.Println("请求未来N天的天气预报失败")
	}
	// fmt.Println(future)
	futureMap := util.JSONToArray(future)

	// 获取未来一周的降水预测
	rain := make(map[string]interface{})
	rain = services.FutureGrid(longitude, latitude, 7, 0)

	result := make(map[string]interface{})
	result["rain"] = rain
	result["weather"] = futureMap
	c.SuccessJson(result)

}
