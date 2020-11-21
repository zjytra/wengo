/*
创建时间: 2020/7/5
作者: zjy
功能介绍:

*/

package dispatch

import (
	"errors"
	"wengo/network"
	"wengo/xlog"
)

// 设置客戶端网络消息处理者
func (this *DispatchSys) SetNetObserver(netobserver network.NetWorkObserver) error {
	if netobserver == nil {
		xlog.DebugLogNoInScene( "noticeNetWorkAccept  netObserver is nil")
		return errors.New(" netObserver is nil")
	}
	this.netObserver = netobserver
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// 客户端网络事件结束
// 客户端网络连接事件
func (this *DispatchSys) OnNetWorkConnect(conn network.Conner) error {
	if this.endFlag.IsTrue() {
		return  errors.New("OnNetWorkConnect已经关闭调度系统")
	}
	data := this.eventDatapool.GetDisPatchDataByPool(NetWorkAccept_Event, conn)
	if data == nil {
		return errors.New("OnNetWorkConnect data nil")
	}
	this.qet.AddEvent(data)
	return nil
}

// 客户端网络读取事件
func (this *DispatchSys) OnNetWorkRead(msgdata *network.MsgData) error {
	if this.endFlag.IsTrue() {
		return  errors.New("OnNetWorkRead已经关闭调度系统")
	}
	data := this.eventDatapool.GetDisPatchDataByPool(NetWorkRead_Event, msgdata)
	if data == nil {
		return errors.New("OnNetWorkRead data nil")
	}
	this.qet.AddEvent(data)
	return nil
}

// 客户端网络关闭事件
func (this *DispatchSys) OnNetWorkClose(conn network.Conner) error {
	if this.endFlag.IsTrue() {
		return  errors.New("OnNetWorkClose已经关闭调度系统")
	}
	data := this.eventDatapool.GetDisPatchDataByPool(NetWorkClose_Event, conn)
	if data == nil {
		return errors.New("OnNetWorkClose data nil")
	}
	this.qet.AddEvent(data)
	return nil
}
//客户端网络事件
////////////////////////////////////////////////////////////////////////////////

////////////////////////////////////////////////////////////////////////////////
// 客户端网络事件处理
// 客户端网络连接事件
func (this *DispatchSys) noticeNetWorkAccept(val interface{}) error {
	conn,ok:= val.(network.Conner)
	if !ok {
		return errors.New("noticeNetWorkAccept  Assert network.Conner type Erro")
	}
	return this.netObserver.OnNetWorkConnect(conn)
}

// 客户端网络读取事件
func (this *DispatchSys) noticeNetWorkRead(val interface{}) error {
	msgData,ok:= val.(*network.MsgData)
	if !ok {
		return errors.New("noticeNetWorkRead  Assert network.MsgData type Erro")
	}
	return this.netObserver.OnNetWorkRead(msgData)
}

// 客户端网络关闭事件
func (this *DispatchSys) noticeNetWorkClose(val interface{}) error {
	conn,ok:= val.(network.Conner)
	if !ok {
		return errors.New("noticeNetWorkClose  Assert network.Conner type Erro")
	}
	return this.netObserver.OnNetWorkClose(conn)
}
