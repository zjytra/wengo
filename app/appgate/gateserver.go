/*
创建时间: 2019/11/24
作者: zjy
功能介绍:
登录服
*/

package appgatesv

import (
	"wengo/appdata"
	"wengo/dispatch"
	"wengo/network"
	"wengo/xlog"
	"sync"
)



type GateServer struct {
	tcpserver    *network.TCPServer
	conns        sync.Map
	dispSys      *dispatch.DispatchSys
}


// 程序启动
func (this *GateServer)OnStart() {
	this.dispSys = dispatch.NewDispatchSys()
	if this.dispSys == nil {
		panic("GateServer OnStart  this.dispSys == nil  ")
	}
	this.dispSys.SetNetObserver(this)
	this.tcpserver = network.NewTcpServer(this.dispSys, appdata.NetConf,appdata.WorkPool)
	this.tcpserver.Start()
}

//初始化
func (this *GateServer)OnInit() bool{

	return true
}
// 程序运行
func (this *GateServer)OnUpdate() bool{
	
	return true
}
// 关闭
func (this *GateServer)OnRelease(){
	this.tcpserver.Close()
	this.dispSys.Release()
	network.Release()
}

func (this *GateServer)OnNetWorkConnect(conn network.Conner) error{
	xlog.DebugLogNoInScene("OnNetWorkConnect %v",conn.RemoteAddr())
	return  nil
}


func (this *GateServer)OnNetWorkClose(conn network.Conner) error{
	xlog.DebugLogNoInScene("OnNetWorkClose %v",conn.RemoteAddr())
	return  nil
}

func (this *GateServer)OnNetWorkRead(msgdata *network.MsgData) error{
	xlog.ErrorLog(appdata.GetSecenName(), "GateServer OnNetWorkRead",)
	return  hanlerRead(msgdata.Conn,msgdata.MainCmd,msgdata.SubCmd,msgdata.Msgdata)
}


func hanlerRead(conn network.Conner,maincmd,subcmd uint16,msgdata []byte) error{
	
	
	return  nil
}


