/*
创建时间: 2020/08/2020/8/26
作者: Administrator
功能介绍:

*/
package dbsys

import (
	"wengo/dispatch"
	"wengo/xlog"
	"wengo/xutil"
	"wengo/xutil/strutil"
	"wengo/xutil/timeutil"
)

// 异步写
//dbcb 回调方法
func (this *MySqlDBStore) AsynExtute(query string,dbParam  *dispatch.DBEventParam,dbcb dispatch.DBExcuteCallback) {
	if strutil.StringIsNil(query) {
		return
	}
	//由于逻辑线程使用了对象池，这里投递的时候就不用池子了避免提前被回收
	data := dispatch.NewEventData(DBExcuteEvent, &DBExcuteEventData{
		BDParam : dbParam,
		Cb:         dbcb,
		Excutestr:  query,
	})
	this.writeEvent.AddEvent(data)
}

//写事件处理
func (this *MySqlDBStore) OnDBWriteEvent(eventdata *dispatch.EventData) {
	queryEvent,ok:= eventdata.Val.(*DBExcuteEventData)
	if !ok {
		xlog.ErrorLogNoInScene("OnDBWriteEvent Assert *DBQueryRowsCallbackEventData")
		return
	}
	startT := timeutil.GetCurrentTimeMs()		//计算当前时间
	result,erro := this.Excute(queryEvent.Excutestr)
	since := xutil.MaxInt64(0,timeutil.GetCurrentTimeMs() - startT)
	if since >= 200 {
		xlog.WarningLogNoInScene("sql ：%v,执行时间%v ms",queryEvent.Excutestr,since)
	}
	if erro != nil {
		return
	}
	if queryEvent.BDParam == nil || queryEvent.BDParam.CbDispSys == nil || queryEvent.Cb == nil {
		return
	}
	//进入调度队列
	onEventErro := queryEvent.BDParam.CbDispSys .PostDBWrite(&dispatch.DBExcuteData{
		BDParam: queryEvent.BDParam,
		Cb:queryEvent.Cb,
		Result:result,
	})
	if onEventErro != nil {
		xlog.ErrorLogNoInScene("投递查询事件 %v",onEventErro)
	}
}
