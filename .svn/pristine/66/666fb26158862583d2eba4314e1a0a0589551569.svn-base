/*
创建时间: 2020/08/2020/8/29
作者: Administrator
功能介绍:

*/
package datacenter

import (
	"github.com/wengo/appdata"
	"github.com/wengo/xlog"
)

// 每秒定时器
func (this *DataCenter) PerOneSTimer() {
	xlog.DebugLogNoInScene("PerOneSTimer")
}



// 分钟定时器
func (this *DataCenter) PerOneMinuteTimer() {
	xlog.WarningLogNoInScene("workpool.Running = %d workpool.Free = %d workpool.Cap = %d",
		appdata.WorkPool.Running(), appdata.WorkPool.Free(), appdata.WorkPool.Cap())
}

// 小时定时器
func (this *DataCenter) PerOneHourTimer() {
	xlog.DebugLogNoInScene("PerOneHourTimer")
	PaccountMgr.DelAccountOnTimer()
}