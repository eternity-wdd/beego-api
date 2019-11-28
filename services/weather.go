package services

import (
	"api_wx_klagri_com_cn_go/util"
	"fmt"
	"github.com/astaxie/beego/httplib"
	"github.com/shopspring/decimal"
	"math"
	// "reflect"
	"strconv"
	"time"
)

// 5*5格点数据, 含未来累计降雨量等信息
// lon 纬度
// lat 经度
// days 预报天数
// start 起始时间是否进行偏移，0为当天，-1为昨天
func FutureGrid(lon, lat float64, days, start int64) map[string]interface{} {
	// 根据请求的天数获取起止日期
	now := time.Now().Format("2006-01-02")
	loc, _ := time.LoadLocation("Local") //获取时区
	tmp, _ := time.ParseInLocation("2006-01-02", now, loc)
	nowStamp := tmp.Unix() //转化为时间戳 类型是int64

	startDate := time.Unix(nowStamp+(86400*start), 0).Format("2006010215")
	endDate := time.Unix(nowStamp+86400*start+86400*days, 0).Format("2006010215")
	gridUrl := "http://api.mlogcn.com/scmocservice/v2/scmoc/gridding/range/point?var=rh&var=tmp&var=tp&var=u10&var=v10&token=5ccc96af717842a5ad410a0ede8bfc6b&lon=" + strconv.FormatFloat(lon, 'f', 6, 64) + "&lat=" + strconv.FormatFloat(lat, 'f', 6, 64) + "&start=" + startDate + "&end=" + endDate + "&var=rh&var=tmp&var=v10&var=u10&var=tmp_min&var=rh_min&var=rh_max&var=tmp_max&var=tp&format=yyyyMMdd"

	//获取未来一周的降水预测
	grid, err := httplib.Get(gridUrl).String()
	if err != nil {
		fmt.Println("调用格点降雨错误")
	}

	gridArray := util.JSONToArray(grid)
	//统计未来N天，每天的累计降水和上下午风力风速
	futureGrid := make(map[string]interface{})
	futureGrid = formatGrid(gridArray)
	// fmt.Println(futureGrid)
	return futureGrid

}

// 整理5*5预报数据，每3小时一条数据转为每天一条数据
// 把径向风速与纬向风速合成为标量风速
func formatGrid(grid []interface{}) map[string]interface{} {

	// 最终返回数据
	data := make(map[string]interface{})

	// 各个要素的切片数组
	// rain, wet, t, tMax, tMin, windSpeed := make(map[string][]float64)
	rain := make(map[string][]float64)
	wet := make(map[string][]float64)
	t := make(map[string][]float64)
	tMax := make(map[string][]float64)
	tMin := make(map[string][]float64)
	windSpeed := make(map[string][]float64)
	// windLevel := make(map[string]string)

	for _, value := range grid {
		temp := value.(map[string]interface{})
		datatime := temp["datatime"].(string) // 使用时间作为key
		data[datatime] = make(map[string]interface{})
		// 每天的降雨量组成一个切片
		rain[datatime] = append(rain[datatime], temp["tp"].(float64))
		// 湿度
		wet[datatime] = append(wet[datatime], temp["rh"].(float64))
		// 温度
		t[datatime] = append(t[datatime], temp["tmp"].(float64))
		// 最高温
		tMax[datatime] = append(tMax[datatime], temp["tmp_max"].(float64))
		// 最低温
		tMin[datatime] = append(tMin[datatime], temp["tmp_min"].(float64))
		// 风速
		windSpeed[datatime] = append(windSpeed[datatime], math.Sqrt(math.Pow(temp["u10"].(float64), 2)+math.Pow(temp["v10"].(float64), 2)))

	}
	//strconv.FormatInt(util.WindSpeed(temp["wind_level"].(float64)), 10)+"级"
	decimal.DivisionPrecision = 3
	for key, _ := range data {
		mData := make(map[string]interface{})
		mData["rain"] = util.ArraySum(rain[key])
		mData["wet"] = decimal.NewFromFloat(util.ArraySum(wet[key])).Div(decimal.NewFromFloat(float64(len(wet[key]))))
		mData["tMax"] = tMax[key][util.ArrayMax(tMax[key]).(int)]
		mData["tMin"] = tMax[key][util.ArrayMin(tMax[key]).(int)]
		mData["t"] = decimal.NewFromFloat(util.ArraySum(wet[key])).Div(decimal.NewFromFloat(float64(len(wet[key]))))
		mData["windSpeed"] = decimal.NewFromFloat(util.ArraySum(windSpeed[key])).Div(decimal.NewFromFloat(float64(len(windSpeed[key]))))

		// decimal.NewFromFloat(mData["windSpeed"].(float64))
		windLevel, _ := decimal.Decimal.Float64(decimal.NewFromFloat(16.12))
		mData["windLevel"] = strconv.FormatInt(util.WindSpeed(windLevel), 10) + "级" //
		// fmt.Println(mData)
		data[key] = mData
	}
	return data
}
