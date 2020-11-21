/*
创建时间: 2020/5/2
作者: zjy
功能介绍:

*/

package apploginsv

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"time"
	"github.com/zjytra/wengo/app/appdata"
	"github.com/zjytra/wengo/app/datacenter/dcmodel"
	"github.com/zjytra/wengo/cmdconst"
	"github.com/zjytra/wengo/cmdconst/cmddatacenter"
	"github.com/zjytra/wengo/network"
	"github.com/zjytra/wengo/protobuf/pb/common_proto"
	"github.com/zjytra/wengo/protobuf/pb/dc_proto"
	"github.com/zjytra/wengo/xlog"
	"github.com/zjytra/wengo/xutil/timeutil"
)

//连接到数据服务器
func ConnectDCServerRepsn(conn network.Conner,msgdata []byte) error{
	restcode := &common_proto.RestInt32CodeMsg{}
	erro := proto.Unmarshal(msgdata,restcode)
	if erro != nil {
		xlog.ErrorLogNoInScene("ConnectDCServerRepsn %v", erro)
		return erro
	}
	
	if restcode.RestCode != dcmodel.Link_Server_Succeed {
		return errors.New("连接数据中心失败")
	}
	
	serverInfo := &dc_proto.ServerInfoMsg{
		AppId:   appdata.NetConf.App_id, //serid
		AppKind: appdata.NetConf.App_kind,
		OutAddr: appdata.NetConf.Out_addr,
		OutProt: appdata.NetConf.Out_prot,
	}
	xlog.DebugLogNoInScene("连接数据中心成功 发送服务器信息到中心服",serverInfo)
	erro = conn.WritePBMsg(cmdconst.Main_DataCenter, cmddatacenter.Sub_Req_RegisterServer,serverInfo)
	if erro != nil {
		conn.Close()
		return nil
	}

	return  erro
}

//注册服务器回复
func RegisterServerRepsn(conn network.Conner,msgdata []byte) error{
	restcode := &dc_proto.RespnRegisterServerInfoMsg{}
	erro := proto.Unmarshal(msgdata,restcode)
	if erro != nil {
		xlog.ErrorLogNoInScene("RegisterServerRepsn %v", erro)
		return erro
	}
	//远程时间减当前时间
	difftime := restcode.UnixNano - time.Now().UnixNano()
	timeutil.SetDiffUnixNano(difftime)
	xlog.DebugLogNoInScene("中心服回复注册码 %v 相差时间%v",restcode.RestCode,difftime)
	if restcode.RestCode != dcmodel.Register_Server_Succeed { //Register_Server_Succeed
		conn.Close()
		return errors.New("RegisterServerRepsn 失败")
	}

	return  erro
}

//注册服务器回复
func RepsnDataCenterHeartBeat(conn network.Conner,msgdata []byte) error{
	restcode := &dc_proto.RespnServerHeartBeatMsg{}
	erro := proto.Unmarshal(msgdata,restcode)
	if erro != nil {
		xlog.ErrorLogNoInScene( "RepsnDataCenterHeartBeat %v", erro.Error())
		return erro
	}
	difftime := restcode.UnixNano - time.Now().UnixNano()
	timeutil.SetDiffUnixNano(difftime)
	//xlog.DebugLogNoInScene("心跳回复%v,相差时间%v",restcode.RestCode,difftime)
	if restcode.RestCode != 1 {
		conn.Close()
		return errors.New("RepsnDataCenterHeartBeat 失败")
	}
	
	return  erro
}

//数据中心发送关闭远端连接
func DCLSCloseFarEndConn(conn network.Conner, msgdata []byte) error {
	closeClient :=  &dc_proto.CloseClientLinkMsg{}
	erro := proto.Unmarshal(msgdata,closeClient)
	if erro != nil {
		xlog.ErrorLogNoInScene( "DCLSCloseFarEndConn %v", erro.Error())
		return erro
	}
	//连接id为零
	if closeClient.AccountId == 0 {
		xlog.WarningLogNoInScene( "中心服请求关闭连接为0")
		return nil
	}
	connID := pLoginAccountMgr.GetConnIDByAccountID(closeClient.AccountId)
	if connID == 0 {
		return nil
	}
	Client.tcpServer.CloseConnID(connID)
	return nil
}

func DCLSGateWayInfo(conn network.Conner, msgdata []byte) error {
	xlog.WarningLogNoInScene( "DCLSGateWayInfo")
	
	return nil
}