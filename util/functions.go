package util

import (
	"encoding/json"
	// "fmt"
	// "github.com/shopspring/decimal"
	"reflect"
)

// json转Map ()
func JSONToMap(str string) map[string]interface{} {

	var tempMap = make(map[string]interface{})

	err := json.Unmarshal([]byte(str), &tempMap)

	if err != nil {
		panic(err)
	}

	return tempMap
}

// 将json转成一个map组成的数组，
func JSONToArray(str string) []interface{} {
	// 声明一个数组
	var temArray []interface{}
	err := json.Unmarshal([]byte(str), &temArray)
	if err != nil {
		panic(err)
	}
	// fmt.Println(123)
	return temArray

}

// 数组求和，一维
func ArraySum(nums []float64) float64 {
	sum := 0.00
	for _, num := range nums {
		sum += num
	}
	return sum
}

// 传入切片或数组，返回最大值的索引
func ArrayMax(nums interface{}) interface{} {
	numsValue := reflect.ValueOf(nums)
	// max := interface{}
	// reflect.TypeOf(haystack).Kind()
	tmp := numsValue.Index(0).Interface()
	key := 0
	numType := reflect.TypeOf(tmp).Kind()
	switch numType {
	case reflect.Float64, reflect.Float32, reflect.Int64, reflect.Int, reflect.Uint64, reflect.Uint32:
		// var maxVal interface{}
		length := numsValue.Len()
		if length == 1 {
			return tmp
		}
		for i := 1; i < length; i++ {
			if numsValue.Index(i).Interface().(float64) > tmp.(float64) {
				tmp = numsValue.Index(i).Interface()
				key = i
			}
		}
	}
	return key
}

// 传入切片或数组，返回最小值的索引
func ArrayMin(nums interface{}) interface{} {
	numsValue := reflect.ValueOf(nums)
	// max := interface{}
	// reflect.TypeOf(haystack).Kind()
	tmp := numsValue.Index(0).Interface()
	key := 0
	numType := reflect.TypeOf(tmp).Kind()
	switch numType {
	case reflect.Float64, reflect.Float32, reflect.Int64, reflect.Int, reflect.Uint64, reflect.Uint32:
		// var maxVal interface{}
		length := numsValue.Len()
		if length == 1 {
			return tmp
		}
		for i := 1; i < length; i++ {
			if numsValue.Index(i).Interface().(float64) < tmp.(float64) {
				tmp = numsValue.Index(i).Interface()
				key = i
			}
		}
	}
	return key
}

// 判断一个数组（切片）， map中是否含有某个值, 参照PHP的in_array
// needle 要查找的值， haystack 一个已知的数组
func Contain(needle interface{}, haystack interface{}) bool {
	haystackValue := reflect.ValueOf(haystack)
	switch reflect.TypeOf(haystack).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < haystackValue.Len(); i++ {
			if haystackValue.Index(i).Interface() == needle {
				return true
			}
		}
	case reflect.Map:
		if haystackValue.MapIndex(reflect.ValueOf(needle)).IsValid() {
			return true
		}
	}
	return false
}

// 两个map合并，A + B = C, 如有重复的值，则用 B的值 覆盖 A的值
func ArrayMerge(mapA, mapB map[string]interface{}) map[string]interface{} {
	for itemB, valueB := range mapB {
		mapA[itemB] = valueB
	}

	return mapA
}
