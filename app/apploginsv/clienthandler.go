/*
创建时间: 2020/6/8
作者: zjy
功能介绍:

*/

package apploginsv

import (
	"github.com/panjf2000/ants/v2"
	"github.com/zjytra/wengo/app/netmsgsys"
	"github.com/zjytra/wengo/cmdconst"
	"github.com/zjytra/wengo/cmdconst/cmdaccount"
	"github.com/zjytra/wengo/csvdata"
	"github.com/zjytra/wengo/dispatch"
	"github.com/zjytra/wengo/network"
)

type ClientHandler struct {
	netmsgsys *netmsgsys.NetMsgSys
	apool     *ants.Pool
	netCof    *csvdata.Networkconf
	tcpServer *network.TCPServer    //为客户端服务的
	svDispSys *dispatch.DispatchSys
}

func NewClientHandle(conf *csvdata.Networkconf, apool *ants.Pool,dispSys *dispatch.DispatchSys)*ClientHandler{
	if  conf == nil {
		panic("NewClientHandle  conf is nil")
		return nil
	}
	if  apool == nil {
		panic("NewClientHandle  apool is nil")
		return nil
	}
	if dispSys == nil {
		panic("NewClientHandle  dispSys is nil")
		return nil
	}
	
	this := new(ClientHandler)
	this.apool = apool
	this.netCof = conf
	this.svDispSys = dispSys
	if !this.OnInit(){
	    panic("创建客户端处理失败")
		return nil
	}
	return this
}

func (this *ClientHandler)OnInit() bool{
	this.svDispSys.SetNetObserver(this)
	//接收客户端的消息
	this.tcpServer = network.NewTcpServer(this.svDispSys, this.netCof,this.apool)
	if this.tcpServer  == nil {
		return false
	}
	this.tcpServer.Start()
	this.netmsgsys = netmsgsys.NewMsgHandler()
	this.LoginRegisterMsg()
	return true
}

//客户端连接
func (this *ClientHandler)OnNetWorkConnect(conn network.Conner) error{
	return  nil
}

//客户端关闭连接
func (this *ClientHandler)OnNetWorkClose(conn network.Conner) error{
	pLoginAccountMgr.Offline(conn.GetConnID())
	return  nil
}

//读取客戶端发来的消息
func (this *ClientHandler)OnNetWorkRead(msgdata *network.MsgData) error{
	return	this.netmsgsys.OnNetWorkMsgHandle(msgdata)
}

// 关闭
func (this *ClientHandler)OnRelease(){
	this.netmsgsys.Release()
	this.tcpServer.Close()
}

func (this *ClientHandler)LoginRegisterMsg(){
	this.netmsgsys.RegisterMsgHandle(cmdconst.Main_Account, cmdaccount.Sub_C_LS_RegisterAccount, RegisterAccountMsgHandler)
	this.netmsgsys.RegisterMsgHandle(cmdconst.Main_Account, cmdaccount.Sub_C_LS_LoginAccount, LoginAccountMsgHandler)
	
	this.netmsgsys.RegisterMsgHandle(cmdconst.Main_Account, cmdaccount.Sub_C_LS_AccountLeave, CL_LSReqLeave)
}


