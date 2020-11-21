/*
创建时间: 2020/2/3
作者: zjy
功能介绍:

*/

package global

import (
	"fmt"
	"github.com/sony/sonyflake"
	"runtime"
	"time"
)



var (
	Sonyflk *sonyflake.Sonyflake //雪花算法gid生成
)

func InitGid(t time.Time) bool  {
	var st sonyflake.Settings
	st.StartTime = t
	Sonyflk = sonyflake.NewSonyflake(st)
	if Sonyflk == nil {
		return false
	}
	return true
}



// 拉起宕机标准输出
func GrecoverToStd() {
	if rec := recover(); rec != nil {
		buf := make([]byte, 4096)
		l := runtime.Stack(buf, false)
		fmt.Printf("%v\n%s \n", rec, buf[:l])
	}
}
