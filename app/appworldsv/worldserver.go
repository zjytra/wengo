/*
创建时间: 2019/11/24
作者: zjy
功能介绍:
登录服
*/

package appworldsv

import (
	"errors"
	"sync"
	"wengo/app/appdata"
	"wengo/app/netmsgsys"
	"wengo/csvdata"
	"wengo/dispatch"
	"wengo/network"
	"wengo/xlog"
)



type WorldServer struct {
	tcpserver *network.TCPServer
	conns     sync.Map
	dispSys   *dispatch.DispatchSys
	netmsgsys *netmsgsys.NetMsgSys
}


// 程序启动
func (this *WorldServer)OnStart() {
	this.OnInit()
}

//初始化
func (this *WorldServer)OnInit() bool{
	this.dispSys = dispatch.NewDispatchSys()
	this.dispSys.SetNetObserver(this)
	// 处理其他服务器的连接
	this.tcpserver = network.NewTcpServer(this.dispSys, appdata.NetConf, appdata.WorkPool)
	this.netmsgsys = netmsgsys.NewMsgHandler()
	this.tcpserver.Start()
	this.RegisterMsg()
	csvdata.LoadCommonCsvData()
	NewData()
	return true
}
// 程序运行
func (this *WorldServer)OnUpdate() bool{
	
	return true
}
// 关闭
func (this *WorldServer)OnRelease(){
	this.tcpserver.Close()
	this.dispSys.Release()
	network.Release()
	ClearAllServerData()
}

func (this *WorldServer)OnNetWorkConnect(conn network.Conner) error{
	//给其他服务器发送链接服务器成功的消息
	//restcode := &common_proto.RestInt32CodeMsg{
	//	ResCode:model.Link_Server_Succeed,
	//}
	////通知其他服务器器连接成功
	//erro := netmsgsys.SendMsg(conn,cmdconst.Main_WorldSv, cmddatacenter.Sub_Req_Connet_WorldSucceed, restcode)
	return  nil
}


func (this *WorldServer)OnNetWorkClose(conn network.Conner) error{
	//这里要移除掉线的服务器
	if RemoveServerInfo(conn.GetConnID()){
		return nil
	}
	return  errors.New("未找到相关服务器")
}

func (this *WorldServer)OnNetWorkRead(msgdata *network.MsgData) error{
	xlog.ErrorLog(appdata.GetSecenName(), "WorldServer OnNetWorkRead",)
	return  this.netmsgsys.OnNetWorkMsgHandle(msgdata)
}


//注册消息
func (this *WorldServer)RegisterMsg(){

}



