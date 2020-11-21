/*
创建时间: 2020/5/17
作者: zjy
功能介绍:
服务器之间的消息通信
*/

package datacenter

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"wengo/app/datacenter/dcmodel"
	"wengo/cmdconst"
	"wengo/cmdconst/cmddatacenter"
	"wengo/csvdata"
	"wengo/model"
	"wengo/network"
	"wengo/protobuf/pb/dc_proto"
	"wengo/xlog"
	"time"
)

// 注册服务器
func ReqRegisterServer(conn network.Conner, msgdata []byte) error {
	serverInfo := &dc_proto.ServerInfoMsg{}
	erro := proto.Unmarshal(msgdata, serverInfo)
	if erro != nil {
		xlog.ErrorLogNoInScene( "ReqRegisterServer %v", erro)
		return erro
	}
	// 配置没有找到证明没有配置这个服务器
	netConf := csvdata.GetNetworkconfPtr(serverInfo.AppId)
	if netConf == nil {
		SendRegisterServerRest(conn,dcmodel.Register_Server_NotFindConf)
		conn.Close()
		return errors.New(fmt.Sprintf(" 未找到 serverInfo.AppId = %v的配置",serverInfo.AppId))
	}

	serverInfoModel := &dcmodel.SeverInfoModel{
		AppId:   serverInfo.AppId,
		AppKind: serverInfo.AppKind,
		OutAddr: serverInfo.OutAddr,
		OutProt: serverInfo.OutProt,
		ConnID:  conn.GetConnID(),
	}
	// 查询内存中是否有这个注册用户
	isOk := PServerInfoMgr.AddServerInfo( serverInfoModel)
	if !isOk {
		// 服务器已经存在
		// 同一個ServerID已经被注册
		SendRegisterServerRest(conn,dcmodel.Register_Server_Exist )
	} else {
		//1成功
		SendRegisterServerRest(conn,dcmodel.Register_Server_Succeed)
	}
	//先发送消息在关闭连接
	if !isOk {
		conn.Close()
		return  errors.New(fmt.Sprintf("同一個ServerID %v已经被注册",serverInfo.AppId))
	}
	// 如果连接进来的是游戏服需要把游戏服的信息发送给网关
	if serverInfo.AppKind == model.APP_GameServer {
		
	}
	
	return erro
}

//发送注册服务器结果
func SendRegisterServerRest(conn network.Conner,restCode  int32)  {
	sendMsg := &dc_proto.RespnRegisterServerInfoMsg{
		RestCode: restCode,
		UnixNano: time.Now().UnixNano() , //时间na秒
	}
	erro := conn.WritePBMsg( cmdconst.Main_DataCenter, cmddatacenter.Sub_Repsn_RegisterServer, sendMsg)
	if erro != nil {
		xlog.ErrorLogNoInScene( "SendRegisterServerRest  WritePBMsg %v", erro)
	}
}
// 接收其他服務器心跳
func ServerHeartBeat(conn network.Conner, msgdata []byte) error {
	syerSysInfo := &dc_proto.ServerSysInfo{}
	erro := proto.Unmarshal(msgdata, syerSysInfo)
	if erro != nil {
		xlog.ErrorLogNoInScene( "ServerHeartBeat %v", erro)
		return erro
	}
	sendMsg := &dc_proto.RespnServerHeartBeatMsg{
		RestCode: 1,
		UnixNano: time.Now().UnixNano() , //时间na秒
	}
	erro = conn.WritePBMsg(cmdconst.Main_DataCenter, cmddatacenter.Sub_Repsn_Server_HeartBeat, sendMsg)
	return erro
}


//关闭远端客户端端连接
func CloseFarEndClientConn(severConnID uint32,accountID uint64){
	closeClient :=  &dc_proto.CloseClientLinkMsg{
		AccountId: accountID,
	}
	PdataCenter.tcpserver.WritePBMsgByConnID(severConnID,cmdconst.Main_DataCenter, cmddatacenter.Sub_DC_Close_FarEnd_Conn,closeClient)
}