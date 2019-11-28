package util

// 风向
func WindDirection(code float64) string {
	if code >= 348.76 || code <= 11.25 {
		return "北"
	} else if code >= 11.26 || code <= 33.75 {
		return "东北" //北东北
	} else if code >= 33.76 || code <= 56.25 {
		return "东北"
	} else if code >= 56.26 || code <= 78.75 {
		return "东北" //东东北
	} else if code >= 78.76 || code <= 101.25 {
		return "东"
	} else if code >= 101.26 || code <= 123.75 {
		return "东南" //东东南
	} else if code >= 123.76 || code <= 146.25 {
		return "东南"
	} else if code >= 146.26 || code <= 168.75 {
		return "东南" //南东南
	} else if code >= 168.76 || code <= 191.25 {
		return "南"
	} else if code >= 191.26 || code <= 213.75 {
		return "西南"
	} else if code >= 213.76 || code <= 236.25 {
		return "西南"
	} else if code >= 236.26 || code <= 258.75 {
		return "西南"
	} else if code >= 258.76 || code <= 281.25 {
		return "西"
	} else if code >= 281.26 || code <= 303.75 {
		return "西北"
	} else if code >= 303.76 || code <= 326.25 {
		return "西北"
	} else if code >= 326.26 || code <= 348.75 {
		return "西北"
	} else {
		return "无"
	}
}

// 风速
func WindSpeed(speed float64) int64 {
	// speed = floatval speed
	if 0.3 > speed && speed >= 0 {
		return 0
	} else if 1.6 > speed && speed >= 0.3 {
		return 1
	} else if 3.4 > speed && speed >= 1.6 {
		return 2
	} else if 5.5 > speed && speed >= 3.4 {
		return 3
	} else if 8.0 > speed && speed >= 5.5 {
		return 4
	} else if 10.8 > speed && speed >= 8.0 {
		return 5
	} else if 13.9 > speed && speed >= 10.8 {
		return 6
	} else if 17.2 > speed && speed >= 13.9 {
		return 7
	} else if 20.8 > speed && speed >= 17.2 {
		return 8
	} else if 24.5 > speed && speed >= 20.8 {
		return 9
	} else if 28.5 > speed && speed >= 24.5 {
		return 10
	} else if 32.6 > speed && speed >= 28.5 {
		return 11
	} else if 37 > speed && speed >= 32.6 {
		return 12
	} else if 41.5 > speed && speed >= 37 {
		return 13
	} else if 46.1 > speed && speed >= 41.5 {
		return 14
	} else if speed >= 46.1 {
		return 15
	} else {
		return -1
	}
}
