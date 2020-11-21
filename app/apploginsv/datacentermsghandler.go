/*
创建时间: 2020/5/24
作者: zjy
功能介绍:
服务器处理,这里主要是区别客户端与服务器消息,以便于分别处理
*/

package apploginsv

import (
	"github.com/panjf2000/ants/v2"
	"wengo/app/netmsgsys"
	"wengo/appdata"
	"wengo/cmdconst"
	"wengo/cmdconst/cmdaccount"
	"wengo/cmdconst/cmddatacenter"
	"wengo/csvdata"
	"wengo/dispatch"
	"wengo/network"
	"wengo/protobuf/pb/dc_proto"
	"wengo/timersys"
	"wengo/xlog"
	"time"
)

type DataCenterMsgHandler struct {
	svDispSys      *dispatch.DispatchSys   // 服务器调度对象 共用一个将服务器业务逻辑搞成单线程
	netmsgsys      *netmsgsys.NetMsgSys
	apool          *ants.Pool
	worldCof       *csvdata.Networkconf
	DataCenterConn *network.ServeiceClient //数据中心连接
	sndHeartTimer  uint32                  //给世界服发送心跳
}

func NewDataCenterHandle(conf *csvdata.Networkconf, apool *ants.Pool,dispSys *dispatch.DispatchSys)*DataCenterMsgHandler {
	if  conf == nil {
		panic("NewDataCenterHandle  conf is nil")
		return nil
	}
	if  apool == nil {
		panic("NewDataCenterHandle  apool is nil")
		return nil
	}
	if  dispSys == nil {
		panic("NewDataCenterHandle  dispSys is nil")
		return nil
	}
	this := new(DataCenterMsgHandler)
	this.apool = apool
	this.worldCof = conf
	this.svDispSys = dispSys
	if !this.OnInit() {
		xlog.ErrorLogNoInScene("DataCenterMsgHandler 初始化 失败")
		return  nil
	}
	return this
}

func (this *DataCenterMsgHandler)OnInit() bool{
	// this.dispSys = dispatch.NewDispatchSys(this.apool)
	this.svDispSys.SetServiceNet(this)
	//与世界服对接，获取网关信息
	this.DataCenterConn = network.NewServeiceClient(this.svDispSys,this.worldCof,this.apool)
	if this.DataCenterConn == nil {
		return  false
	}
	this.DataCenterConn.Start()
	this.netmsgsys = netmsgsys.NewMsgHandler()
	this.RegisterServerMsg()
	this.SetTimer()
	return true
}

func (this *DataCenterMsgHandler)OnServiceLink(conn network.Conner) error{
	return  nil
}

//世界服关闭连接
func (this *DataCenterMsgHandler)OnServiceClose(conn network.Conner) error{
	return  nil
}

//读取世界服发来的消息
func (this *DataCenterMsgHandler)OnServiceMsg(msgdata *network.MsgData) error{
	return	this.netmsgsys.OnNetWorkMsgHandle(msgdata)
}

// 关闭
func (this *DataCenterMsgHandler)OnRelease(){
	this.netmsgsys.Release()
	this.DataCenterConn.Close()
	timersys.StopTimer(this.sndHeartTimer)
}


// 发送心跳给世界服
func (this *DataCenterMsgHandler) SendHeartToWS() {
	
	//给中心服发送心跳
	sverSysInfo := &dc_proto.ServerSysInfo{}
	sverSysInfo.FromAppId = appdata.AppID
	sverSysInfo.UserConnnum = GetClientConnSize()
	erro := this.DataCenterConn.WritePBMsg(cmdconst.Main_DataCenter, cmddatacenter.Sub_Req_Server_HeartBeat,sverSysInfo)
	if erro != nil {
		xlog.ErrorLogNoInScene("SendHeartToWS %v",erro)
	}
}

// 设置定时器
func (this *DataCenterMsgHandler)SetTimer() {
	this.sndHeartTimer = timersys.NewWheelTimer(time.Second * 30,this.SendHeartToWS,this.svDispSys)
}

// 注册服务器 的消息
func (this *DataCenterMsgHandler)RegisterServerMsg(){
	this.netmsgsys.RegisterMsgHandle(cmdconst.Main_DataCenter, cmddatacenter.Sub_Repsn_Connet_DCSucceed, ConnectDCServerRepsn)
	this.netmsgsys.RegisterMsgHandle(cmdconst.Main_DataCenter, cmddatacenter.Sub_Repsn_RegisterServer, RegisterServerRepsn)
	this.netmsgsys.RegisterMsgHandle(cmdconst.Main_DataCenter, cmddatacenter.Sub_Repsn_Server_HeartBeat, RepsnDataCenterHeartBeat)
	this.netmsgsys.RegisterMsgHandle(cmdconst.Main_DataCenter, cmddatacenter.Sub_DC_Close_FarEnd_Conn, DCLSCloseFarEndConn)
	this.netmsgsys.RegisterMsgHandle(cmdconst.Main_DataCenter, cmddatacenter.Sub_DC_LS_GateWayInfo, DCLSGateWayInfo)
	//账号相关
	this.netmsgsys.RegisterMsgHandle(cmdconst.Main_Account, cmdaccount.Sub_DC_LS_RegisterAccount, DCLSRespnRegisterAccount)
	this.netmsgsys.RegisterMsgHandle(cmdconst.Main_Account, cmdaccount.Sub_DC_LS_LoginAccount, DCLSRespnLoginAccount)
}