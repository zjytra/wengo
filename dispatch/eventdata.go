/*
创建时间: 2020/2/29
作者: zjy
功能介绍:

*/

package dispatch

import (
	"github.com/wengo/appdata"
	"github.com/wengo/xlog"
	"sync"
)

type EventData struct {
	DipatchType int
	Val         interface{}
}

// 事件数据
func EventDataPoolNewFun() interface{} {
	return new(EventData)
}

// 事件数据
func NewEventData(dtype int, val interface{}) *EventData {
	data := new(EventData)
	data.SetEventData(dtype, val)
	return data
}

// 事件数据
func (this *EventData) SetEventData(dtype int, val interface{}) {
	this.DipatchType = dtype
	this.Val = val
}

type EventDataPool struct {
	datapool sync.Pool
}

func NewEventDataPool() *EventDataPool {
	pool := new(EventDataPool)
	pool.datapool.New = EventDataPoolNewFun
	return pool
}

func (this *EventDataPool) GetDisPatchDataByPool(dtype int, val interface{}) *EventData {
	dpd := this.datapool.Get()
	data, ok := dpd.(*EventData)
	if !ok {
		xlog.WarningLog(appdata.GetSecenName(), " GetDisPatchDataByPool is nil")
		return NewEventData(dtype, val)
	}
	data.SetEventData(dtype, val)
	return data
}

func (this *EventDataPool) Put(val *EventData) {
	if val == nil {
		return
	}
	this.datapool.Put(val)
}



