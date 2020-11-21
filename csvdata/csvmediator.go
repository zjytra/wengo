/*
创建时间: 2020/2/11
作者: zjy
功能介绍:

*/

package csvdata

import (
	"fmt"
	"wengo/xutil/strutil"
)

var (
	csvPath string
)

func SetCsvPath(csvpath string) {
	if strutil.StringIsNil(csvpath) {
		fmt.Println("csvpath is nil")
	}
	csvPath = csvpath
}
