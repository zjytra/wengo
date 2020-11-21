/*
创建时间: 2020/08/2020/8/29
作者: Administrator
功能介绍:

*/
package apploginsv

import (
	"wengo/app/netmsgsys"
	"wengo/appdata"
	"wengo/xlog"
)

// 每秒定时器
func (this *LogionServer) PerOneSTimer() {
	xlog.DebugLogNoInScene("PerOneSTimer")
}



// 分钟定时器
func (this *LogionServer) PerOneMinuteTimer() {
	xlog.WarningLogNoInScene("workpool.Running = %d workpool.Free = %d workpool.Cap = %d",
		appdata.WorkPool.Running(), appdata.WorkPool.Free(), appdata.WorkPool.Cap())
	netmsgsys.DelExpireData()
}

// 小时定时器
func (this *LogionServer) PerOneHourTimer() {

}