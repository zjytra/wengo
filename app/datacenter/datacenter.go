/*
创建时间: 2019/11/24
作者: zjy
功能介绍:
登录服
*/

package datacenter

import (
	"errors"
	"wengo/app/datacenter/dcmodel"
	"wengo/app/netmsgsys"
	"wengo/appdata"
	"wengo/cmdconst"
	"wengo/cmdconst/cmdaccount"
	"wengo/cmdconst/cmddatacenter"
	"wengo/dbsys"
	"wengo/dispatch"
	"wengo/network"
	"wengo/protobuf/pb/common_proto"
	"wengo/timersys"
	"sync"
	"time"
)



type DataCenter struct {
	tcpserver  *network.TCPServer
	conns      sync.Map
	dispSys    *dispatch.DispatchSys
	netmsgsys  *netmsgsys.NetMsgSys
	oneSTimeID uint32  //定时器id
	oneMinuteTimeID uint32//定时器id
	hourTimeID   uint32 //小时定时器ID
}


// 程序启动
func (this *DataCenter)OnStart() {
	this.OnInit()
}

//初始化
func (this *DataCenter)OnInit() bool{
	NewData(this)
	// 加载csv配置
	// csvdata.LoadCommonCsvData()
	if appdata.WorkPool == nil {
		panic(" InitData appdata.WorkPool is nil")
	}
	dbsys.InitGameDB()
	dbsys.InitLogDB()
	dbsys.InitStatisticsDB()
	dbsys.PDBParamPool = dispatch.NewDBEventParamPool(200) //初始化对象池
	this.dispSys = dispatch.NewDispatchSys()
	this.dispSys.SetNetObserver(this)
	// 处理其他服务器的连接
	this.tcpserver = network.NewTcpServer(this.dispSys, appdata.NetConf,appdata.WorkPool)
	this.netmsgsys = netmsgsys.NewMsgHandler()
	this.tcpserver.Start()
	this.RegisterMsg()
	//this.oneSTimeID = timersys.NewWheelTimer(time.Second,this.PerOneSTimer,this.dispSys)      //每秒钟调用
	this.oneMinuteTimeID = timersys.NewWheelTimer(time.Minute,this.PerOneMinuteTimer,this.dispSys) //每分钟调用
	this.hourTimeID = timersys.NewWheelTimer(time.Hour,this.PerOneHourTimer,this.dispSys) //每小时调用
	return true
}

// 关闭
func (this *DataCenter)OnRelease(){
	this.tcpserver.Close()
	this.dispSys.Release()
	network.Release()
	ClearAllServerData()
}

func (this *DataCenter)OnNetWorkConnect(conn network.Conner) error{
	//给其他服务器发送链接服务器成功的消息
	restcode := &common_proto.RestInt32CodeMsg{
		RestCode: dcmodel.Link_Server_Succeed,
	}
	//通知其他服务器器连接成功
	erro := conn.WritePBMsg(cmdconst.Main_DataCenter, cmddatacenter.Sub_Repsn_Connet_DCSucceed, restcode)
	return  erro
}


func (this *DataCenter)OnNetWorkClose(conn network.Conner) error{
	//这里要移除掉线的服务器
	if PServerInfoMgr.RemoveServerInfo(conn.GetConnID()){
		return nil
	}
	return  errors.New("未找到相关服务器")
}

func (this *DataCenter)OnNetWorkRead(msgdata *network.MsgData) error{
	return  this.netmsgsys.OnNetWorkMsgHandle(msgdata)
}


//注册消息
func (this *DataCenter)RegisterMsg(){
	this.netmsgsys.RegisterMsgHandle(cmdconst.Main_DataCenter, cmddatacenter.Sub_Req_RegisterServer, ReqRegisterServer)
	this.netmsgsys.RegisterMsgHandle(cmdconst.Main_DataCenter, cmddatacenter.Sub_Req_Server_HeartBeat, ServerHeartBeat)
	this.netmsgsys.RegisterMsgHandle(cmdconst.Main_DataCenter, cmddatacenter.Sub_Req_Server_HeartBeat, ServerHeartBeat)
	this.netmsgsys.RegisterMsgHandle(cmdconst.Main_Account, cmdaccount.Sub_LS_DC_RegisterAccount, LSDCRegisterAccountMsgHandler)
	this.netmsgsys.RegisterMsgHandle(cmdconst.Main_Account, cmdaccount.Sub_LS_DC_LoginAccount, LSDCLoginAccountMsgHandler)
	
	this.netmsgsys.RegisterMsgHandle(cmdconst.Main_Account,cmdaccount.Sub_LS_DC_AccountOffline, LSDCAccountOffline)
	this.netmsgsys.RegisterMsgHandle(cmdconst.Main_Account,cmdaccount.Sub_LS_DC_AccountLeave, LSDCAccountLeave)

}



