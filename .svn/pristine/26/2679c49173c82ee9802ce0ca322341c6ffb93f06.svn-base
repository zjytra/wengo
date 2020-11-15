/*
创建时间: 2020/08/2020/8/26
作者: Administrator
功能介绍:
数据库异步查询事件
*/
package dbsys

import (
	"github.com/wengo/dispatch"
	"github.com/wengo/xlog"
	"github.com/wengo/xutil"
	"github.com/wengo/xutil/strutil"
	"github.com/wengo/xutil/timeutil"
	"reflect"
)

//查询事件处理
func (this *MySqlDBStore) OnDBQuereyEvent(eventdata *dispatch.EventData) {
	if eventdata == nil {
		return
	}
	//自定义查询
	if eventdata.DipatchType == DBQueryCustomOneRowQuery_Event  {
		this.onCustomOneRowQuery(eventdata)
		return
	}
	//处理其他查询方式
	this.dbQuery(eventdata)
}

func (this *MySqlDBStore) dbQuery(eventdata *dispatch.EventData) {
	queryEvent, ok := eventdata.Val.(*AsyncDBQueryData)
	if !ok {
		xlog.ErrorLogNoInScene("OnDBQuereyEvent Assert *DBQueryRowsCallbackEventData")
		return
	}
	startT := timeutil.GetCurrentTimeMs() //计算当前时间
	row, erro := this.Query(queryEvent.Querystr)
	if erro != nil {
		//回收数据
		PDBParamPool.Recycle(queryEvent.BDParam)
		return
	}
	//处理的方式
	switch eventdata.DipatchType {
	case DBQueryRowsCB_Event:
		queryEvent.BDParam.DBRows = row
	case DBQueryRowToStructCb_Event:
		RowToStruct(row, queryEvent.BDParam.ReflectObj)
		erro = row.Close()
	case DBQueryRowsToStrSlicesCb_Event:
		queryEvent.BDParam.StrRows = RowsToStringsSlices(row)
		erro = row.Close()
	case DBQueryRowToStrSliceCb_Event:
		queryEvent.BDParam.StrRow = RowToStringSlice(row)
		erro = row.Close()
	case DBQueryRowsToStructSliceCb_Event:
		queryEvent.BDParam.Objs = RowsToStructSlice(row, queryEvent.BDParam.ReflectObj)
		erro = row.Close()
	default:
		xlog.ErrorLogNoInScene("OnDBQuereyEvent  查询类型%v未处理", eventdata.DipatchType)
		PDBParamPool.Recycle(queryEvent.BDParam)
		erro = row.Close()
		return
	}
	since := xutil.MaxInt64(0, timeutil.GetCurrentTimeMs()-startT)
	if since >= 200 {
		xlog.WarningLogNoInScene("onQueryRowsEvent sql =%v,执行时间%v ms", queryEvent.Querystr, since)
	}
	if queryEvent.BDParam.CbDispSys == nil || queryEvent.Cb == nil {
		xlog.WarningLogNoInScene("OnDBQuereyEvent 调度事件=nil 或者回调 = nil 就不向逻辑队列投递事件了")
		return
	}
	//进入调度队列
	onEventErro := queryEvent.BDParam.CbDispSys.PostDBQuery(&dispatch.DBQueryData{
		BDParam: queryEvent.BDParam,
		Cb:      queryEvent.Cb,
	})
	if onEventErro != nil {
		xlog.ErrorLogNoInScene("投递查询事件 %v", onEventErro)
	}
}

// 异步查询
//@param dbParam 查询回调调度器，这里为了方便选择哪个线程回调方法
//@param dbcb 查询回调,由dbParam.CbDispSys 调度
//@param query 查询字符串
func (this *MySqlDBStore) AsyncRowsQuery(dbParam *dispatch.DBEventParam, dbcb dispatch.OnDBQueryCB, query string) {
	this.asyncQuery(DBQueryRowsCB_Event,dbParam, dbcb, query)
}

// 异步查询返回结构体
//@param dbParam 查询回调调度器，这里为了方便选择哪个线程回调方法
//@param dbcb 查询回调,由dbParam.CbDispSys 调度
//@param query 查询字符串
func (this *MySqlDBStore) AsyncRowToStructQuery(dbParam *dispatch.DBEventParam, dbcb dispatch.OnDBQueryCB, query string) {
	this.asyncQuery(DBQueryRowToStructCb_Event,dbParam, dbcb, query)
}

// 异步查询返回结构体数组
//@param dbParam 查询回调调度器，这里为了方便选择哪个线程回调方法
//@param dbcb 查询回调,由dbParam.CbDispSys 调度
//@param query 查询字符串
func (this *MySqlDBStore) AsyncRowsToStructSliceQuery(dbParam *dispatch.DBEventParam, dbcb dispatch.OnDBQueryCB, query string) {
	this.asyncQuery(DBQueryRowsToStructSliceCb_Event,dbParam, dbcb, query)
}

//投递查询数据
func (this *MySqlDBStore) asyncQuery(dtype int, dbParam *dispatch.DBEventParam, dbcb dispatch.OnDBQueryCB, query string) {
	
	if strutil.StringIsNil(query) {
		xlog.ErrorLogNoInScene("AsyncQuery  query == nil")
		return
	}
	//由于逻辑线程使用了对象池，这里投递的时候就不用池子了避免提前被回收
	data := dispatch.NewEventData(dtype, &AsyncDBQueryData{
		BDParam:  dbParam,
		Cb:       dbcb,
		Querystr: query,
	})
	this.quereyEvent.AddEvent(data)
}



// 异步查询
//@param oneRow 自定义查询 每个参数作为一个
func (this *MySqlDBStore) AsyncCustomOneRowQuery(oneRow dispatch.CustomDBOperate) {
	if oneRow == nil {
		return
	}
	//由于逻辑线程使用了对象池，这里投递的时候就不用池子了避免提前被回收
	data := dispatch.NewEventData(DBQueryCustomOneRowQuery_Event, oneRow)
	this.quereyEvent.AddEvent(data)
}

//向逻辑线程返回原始的查询结果
func (this *MySqlDBStore) onCustomOneRowQuery(eventdata *dispatch.EventData) {
	startT := timeutil.GetCurrentTimeMs() //计算当前时间
	queryEvent, ok := eventdata.Val.(dispatch.CustomDBOperate)
	if !ok {
		xlog.ErrorLogNoInScene("OnDBQuereyEvent Assert *DBQueryRowsCallbackEventData")
		return
	}
	queryEvent.ExcouteQueryFun()
	since := xutil.MaxInt64(0, timeutil.GetCurrentTimeMs()-startT)
	if since >= 200 {
		reFlecttype := reflect.TypeOf(queryEvent)
		xlog.WarningLogNoInScene("自定义接口%v,执行时间%v ms", reFlecttype.String(), since)
	}
}

