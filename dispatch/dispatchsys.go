/*
创建时间: 2020/2/3
作者: zjy
功能介绍:
事件系统
*/

package dispatch

import (
	"errors"
	"fmt"
	"wengo/appdata"
	"wengo/model"
	"wengo/network"
	"wengo/xlog"
	"wengo/xutil"
	"wengo/xutil/timeutil"
)

//type DispatchNoticeFun func(interface{}) error // 对应的解析函数

const (
	Event_NONE                 = 0
	Timer_Event                = 1  // 定时器事件
	NetWorkAccept_Event        = 2  // 客户端网络连接事件
	NetWorkRead_Event          = 3  // 客户端网络读取事件
	NetWorkClose_Event         = 4  // 客户端网络关闭
	ServiceLink_Event          = 5  // 服务器网络连接事件
	ServiceMsg_Event           = 6  // 服务器网络读取事件
	ServiceClose_Event         = 7  // 服务器网络关闭事件
	DBQuerey_Event             = 8  // 数据库查询事件参数是原始数据
	DBWrite_Event              = 9  // 数据库写事件事件
	RedisResult_Event          = 10 // redis结果事件
	DBCustomQuereyOneRow_Event = 11 // 数据库自定义查询参数是返回单个结构体数据
	DisPatch_max               = 12
)

//事件字符串名称
var EventStrArr []string

// 处理消息的函数对象
type HandleMsg func(conn network.Conner, msgdata []byte) error

type DispatchSys struct {
	qet               *QueueEvent
	netObserver       network.NetWorkObserver
	serviceNet        network.ServiceNetEvent
	endFlag           *model.AtomicBool
	eventDatapool     *EventDataPool
}

// go自动调用 初始化管理变量
func init() {
	// DispSys = NewDispatchSys()
}

func NewDispatchSys() *DispatchSys {
	disp := new(DispatchSys)
	disp.qet = NewQueueEvent()
	disp.qet.AddEventDealer(disp.OnQueueEvent) //添加处理函数 一个处理函数就是单线程处理
	disp.eventDatapool = NewEventDataPool()
	disp.init()
	return disp
}

func (this *DispatchSys) init() {
	this.endFlag = model.NewAtomicBool()
	this.endFlag.SetFalse()
	this.InitEventName()
}
//映射事件名称
func (this *DispatchSys) InitEventName() {
	EventStrArr = make([]string, DisPatch_max)
	EventStrArr[Event_NONE] = "Event_NONE"
	EventStrArr[Timer_Event] = "Timer_Event"
	EventStrArr[NetWorkAccept_Event] = "NetWorkAccept_Even"
	EventStrArr[NetWorkRead_Event] = "NetWorkRead_Event"
	EventStrArr[NetWorkClose_Event] = "NetWorkClose_Event"
	EventStrArr[ServiceLink_Event] = "ServiceLink_Event"
	EventStrArr[ServiceMsg_Event] = "ServiceMsg_Event"
	EventStrArr[ServiceClose_Event] = "ServiceClose_Event"
	EventStrArr[DBQuerey_Event] = "DBQuerey_Event"
	EventStrArr[DBWrite_Event] = "DBWrite_Event"
	EventStrArr[RedisResult_Event] = "RedisResult_Event"
	EventStrArr[DBCustomQuereyOneRow_Event] ="DBCustomQuereyOneRow_Event"
}

// 投递定时器事件
func (this *DispatchSys) PostTimerEvent(cb func()) error {
	if this.endFlag.IsTrue() {
		return errors.New("PostTimerEvent已经关闭调度系统")
	}
	data := this.eventDatapool.GetDisPatchDataByPool(Timer_Event, cb)
	if data == nil {
		return errors.New("PostTimerEvent data nil")
	}
	this.qet.AddEvent(data)
	return nil
}

//队列事件回调
func (this *DispatchSys) OnQueueEvent(data *EventData) {
	
	if data.DipatchType < Event_NONE || data.DipatchType >= DisPatch_max {
		xlog.ErrorLogNoInScene("DipatchType = %d 未找到处理函数", data.DipatchType)
		return
	}
	//xlog.DebugLogNoInScene("调度事件 %v",EventStrArr[data.DipatchType])
	
	var erro error
	startT := timeutil.GetCurrentTimeMs()		//计算当前时间
	// 查找对应的方法处理数据
	switch data.DipatchType {
	case Timer_Event:
		erro = this.onEventTimer(data.Val)
	case NetWorkAccept_Event:
		erro = this.noticeNetWorkAccept(data.Val)
	case NetWorkRead_Event:
		erro = this.noticeNetWorkRead(data.Val)
	case NetWorkClose_Event:
		erro = this.noticeNetWorkClose(data.Val)
	case ServiceLink_Event:
		erro = this.noticeServiceLink(data.Val)
	case ServiceMsg_Event:
		erro = this.noticeServiceMsg(data.Val)
	case ServiceClose_Event:
		erro = this.noticeServiceClose(data.Val)
	case DBQuerey_Event:
		erro = this.noticeDBQuery(data.Val)
	case DBWrite_Event:
		erro = this.noticeDBWrite(data.Val)
	case DBCustomQuereyOneRow_Event:
	    erro = this.noticeCustomDBOperate(data.Val)
	default:
		erro = errors.New(fmt.Sprint("调度事件%v未处理",data.DipatchType))
	}
	if erro != nil {
		xlog.DebugLogNoInScene(" OnQueueEvent DipatchType = %d, err %v", data.DipatchType, erro)
	}
	//自定义查询不回收数据
	if data.DipatchType != DBCustomQuereyOneRow_Event &&  data.DipatchType != DBQuerey_Event && data.DipatchType != DBWrite_Event {
		this.eventDatapool.Put(data) // 放回到池子
	}
	since := xutil.MaxInt64(0,timeutil.GetCurrentTimeMs() - startT)
	if since > 20 { //大于20毫秒
		xlog.ErrorLog(appdata.GetSecenName(), "调度类型%v耗时 %v",EventStrArr[data.DipatchType], since)
	}
	
	
}

// 关闭系统
func (this *DispatchSys) Release() {
	this.endFlag.SetTrue()
	this.qet.Release()
	fmt.Println("DispatchSys Release")
}
