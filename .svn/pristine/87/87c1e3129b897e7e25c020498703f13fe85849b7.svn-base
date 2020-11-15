/*
创建时间: 2020/7/5
作者: zjy
功能介绍:

*/

package dispatch

import (
	"errors"
	"github.com/wengo/network"
	"github.com/wengo/xlog"
)


// 设置服务器网络消息处理者
func (this *DispatchSys) SetServiceNet (serviceNet network.ServiceNetEvent) error {
	if serviceNet == nil {
		xlog.DebugLogNoInScene( "noticeNetWorkAccept  netObserver is nil")
		return errors.New(" serviceNet is nil")
	}
	this.serviceNet = serviceNet
	return nil
}

////////////////////////////////////////////////////////////////////////////////
//服务器网络相关事件
//服务器网络连接事件
func (this *DispatchSys) OnServiceLink(conn network.Conner) error {
	if this.endFlag.IsTrue() {
		return  errors.New("OnServiceLink已经关闭调度系统")
	}
	data := this.eventDatapool.GetDisPatchDataByPool(ServiceLink_Event, conn)
	if data == nil {
		return errors.New("OnServiceLink data nil")
	}
	this.qet.AddEvent(data)
	return nil
}

// 服务器网络读取事件
func (this *DispatchSys) OnServiceMsg(msgdata *network.MsgData) error {
	if this.endFlag.IsTrue() {
		return  errors.New("OnServiceMsg已经关闭调度系统")
	}
	data := this.eventDatapool.GetDisPatchDataByPool(ServiceMsg_Event, msgdata)
	if data == nil {
		return errors.New("OnServiceMsg data nil")
	}
	this.qet.AddEvent(data)
	return nil
}

// 服务器网络关闭事件
func (this *DispatchSys) OnServiceClose(conn network.Conner) error {
	if this.endFlag.IsTrue() {
		return  errors.New("OnServiceClose已经关闭调度系统")
	}
	data := this.eventDatapool.GetDisPatchDataByPool(ServiceClose_Event, conn)
	if data == nil {
		return errors.New("OnServiceClose data nil")
	}
	this.qet.AddEvent(data)
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// 服务器网络事件处理
// 服务器网络连接事件
func (this *DispatchSys) noticeServiceLink(val interface{}) error {
	conn,ok:= val.(network.Conner)
	if !ok {
		return errors.New("noticeServiceLink  Assert network.Conner type Erro")
	}
	return this.serviceNet.OnServiceLink(conn)
}

// 服务器网络读取事件
func (this *DispatchSys) noticeServiceMsg(val interface{}) error {
	msgData,ok:= val.(*network.MsgData)
	if !ok {
		return errors.New("noticeServiceMsg  Assert network.MsgData type Erro")
	}
	return this.serviceNet.OnServiceMsg(msgData)
}

// 服务器网络关闭闭事件
func (this *DispatchSys) noticeServiceClose(val interface{}) error {
	conn,ok:= val.(network.Conner)
	if !ok {
		return errors.New("noticeServiceClose  Assert network.Conner type Erro")
	}
	return this.serviceNet.OnServiceClose(conn)
}
//服务器网络相关事件结束
////////////////////////////////////////////////////////////////////////////////
