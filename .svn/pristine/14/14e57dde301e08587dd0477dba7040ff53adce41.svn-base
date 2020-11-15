/*
创建时间: 2020/7/18
作者: zjy
功能介绍:

*/

package dispatch

import (
	"errors"
)

// 数据库查询返回原始数据
func (this *DispatchSys) PostDBQuery(querydata *DBQueryData) error {
	if this.endFlag.IsTrue() {
		return  errors.New("OnDBQuerey已经关闭调度系统")
	}
	data := NewEventData(DBQuerey_Event, querydata)
	if data == nil {
		return errors.New("PostDBQuery data nil")
	}
	this.qet.AddEvent(data)
	return nil
}

// 数据库写事件返回
func (this *DispatchSys) PostDBWrite(excuteData *DBExcuteData) error {
	if this.endFlag.IsTrue() {
		return  errors.New("OnDBWrite已经关闭调度系统")
	}
	data := NewEventData(DBWrite_Event,excuteData)
	if data == nil {
		return errors.New("PostDBWrite data nil")
	}
	this.qet.AddEvent(data)
	return nil
}

// 数据库回调事件
func (this *DispatchSys) noticeDBQuery(val interface{}) error {
	quereyData,ok:= val.(*DBQueryData)
	if !ok {
		return errors.New("noticeDBQuery DBQueryDataRowsCallback type Erro")
	}
	if quereyData == nil {
		return errors.New("noticeDBQuery quereyData is nil")
	}
	if quereyData.Cb == nil{
		return  errors.New("noticeDBQuery quereyData.Cb is nil")
	}
	return quereyData.Cb (quereyData.BDParam)
}

func (this *DispatchSys) noticeDBWrite(val interface{}) error {
	excuteData,ok:= val.(DBExcuteData)
	if !ok {
		return errors.New("noticeDBWrite DBExcuteData type Erro")
	}
	if excuteData.Cb == nil{
		return  errors.New("noticeDBWrite excuteData.Cb is nil")
	}
	return excuteData.Cb (excuteData.BDParam,excuteData.Result)
}

// 数据库查询返回自定义接口
func (this *DispatchSys) PostCustomDBOperateOneRow(custom CustomDBOperate) error {
	if this.endFlag.IsTrue() {
		return  errors.New("PostCustomDBOperateOneRow已经关闭调度系统")
	}
	//自定查询不在池子里面取，应为对象是在逻辑线程分配
	data := NewEventData(DBCustomQuereyOneRow_Event, custom)
	if data == nil {
		return errors.New("PostCustomDBOperateOneRow data nil")
	}
	this.qet.AddEvent(data)
	return nil
}


// 数据库回调事件
func (this *DispatchSys) noticeCustomDBOperate(val interface{}) error {
	quereyData,ok:= val.(CustomDBOperate)
	if !ok {
		return errors.New("noticeCustomDBOperate CustomDBOperate type Erro")
	}
	if quereyData == nil{
		return  errors.New("noticeCustomDBOperate quereyData is nil")
	}
	return quereyData.OnQueryCB()
}


// redis回调事件
func (this *DispatchSys) OnRedisResultEvnet(val interface{}) error {
	
	return  nil
}