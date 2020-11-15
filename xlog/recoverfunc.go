/*
创建时间: 2020/2/17
作者: zjy
功能介绍:

*/

package xlog

import (
	"fmt"
	"runtime"
)
var(
	LenStackBuf = 4096
)

func RecoverToStd()  {
	if rec := recover(); rec != nil {
		buf := make([]byte, LenStackBuf)
		l := runtime.Stack(buf, false)
		fmt.Printf("%v\n%s \n", rec, buf[:l])
	}
}


//拉起宕机日志输出
func RecoverToLog(hanler func()) {
	if rec := recover(); rec != nil {
		buf := make([]byte, LenStackBuf)
		l := runtime.Stack(buf, false)
		ErrorLog("","%v\n%s", rec, buf[:l])
		hanler()
	}
}
