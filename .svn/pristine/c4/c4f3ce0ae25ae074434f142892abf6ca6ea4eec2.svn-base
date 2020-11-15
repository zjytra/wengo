/*
创建时间: 2020/4/28
作者: zjy
功能介绍:
反射相关处理
*/

package xutil

import (
	"fmt"
)

// 验证内置类型数组
func ValidArrIndex(arr interface{}, index int) bool {
	if arr == nil {
		return false
	}
	// 下标为负
	if index < 0 {
		return false
	}
	switch val := arr.(type) {
	case []int:
		return index < len(val)
	case []string:
		return index < len(val)
	case []float32:
		return index < len(val)
	default:
		fmt.Println(arr, "is an unknown type. ")
		return false
	}
	return true
}



