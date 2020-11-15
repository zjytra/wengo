/*
创建时间: 2020/2/3
作者: zjy
功能介绍:

*/

package global

import (
	"fmt"
	"github.com/snowflake"
	"runtime"
)



var (
	SnowGid  *snowflake.NodeGID //雪花算法gid生成
) // 协程池


func InitGid(nodeId int16) bool  {
	var erro error
	SnowGid,erro = snowflake.NewNodeGID(nodeId,12)
	if erro != nil {
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
