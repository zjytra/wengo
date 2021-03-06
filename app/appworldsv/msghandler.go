/*
创建时间: 2020/5/17
作者: zjy
功能介绍:

*/

package appworldsv

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/zjytra/wengo/app/appdata"
	"github.com/zjytra/wengo/app/datacenter/dcmodel"
	"github.com/zjytra/wengo/cmdconst"
	"github.com/zjytra/wengo/cmdconst/cmddatacenter"
	"github.com/zjytra/wengo/csvdata"
	"github.com/zjytra/wengo/model"
	"github.com/zjytra/wengo/network"
	"github.com/zjytra/wengo/protobuf/pb/common_proto"
	"github.com/zjytra/wengo/protobuf/pb/datacenter_proto"
	"github.com/zjytra/wengo/xlog"
)

// 注册服务器
func ReqRegisterServer(conn network.Conner, msgdata []byte) error {
	serverInfo := &datacenter_proto.ServerInfoMsg{}
	erro := proto.Unmarshal(msgdata, serverInfo)
	if erro != nil {
		xlog.ErrorLog(appdata.GetSecenName(), "ReqRegisterServer %v", erro.Error())
		return erro
	}
	// 配置没有找到证明没有配置这个服务器
	restcode := &common_proto.RestInt32CodeMsg{}
	netConf := csvdata.GetNetworkconfPtr(serverInfo.AppId)
	if netConf == nil {
		restcode.RestCode = dcmodel.Register_Server_NotFindConf
		erro = conn.WritePBMsg(cmdconst.Main_DataCenter, cmddatacenter.Sub_Req_RegisterServer, restcode)
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
	isOk := AddServerInfo( serverInfoModel)
	if !isOk {
		// 同一個ServerID已经被注册
		restcode.RestCode = dcmodel.Register_Server_Exist // 服务器已经存在
	} else {
		restcode.RestCode = dcmodel.Register_Server_Succeed //1成功
	}
	erro = conn.WritePBMsg( cmdconst.Main_DataCenter, cmddatacenter.Sub_Req_RegisterServer, restcode)
	if !isOk {
		conn.Close()
		return  errors.New(fmt.Sprintf("同一個ServerID %v已经被注册",serverInfo.AppId))
	}
	// 如果连接进来的是登陆服需要发送 gatesver 给他
	if serverInfo.AppKind == model.APP_LoginServer {
		
	}
	
	return erro
}

// 接收其他服務器心跳
func ServerHeartBeat(conn network.Conner, msgdata []byte) error {
	syerSysInfo := &datacenter_proto.ServerSysInfo{}
	erro := proto.Unmarshal(msgdata, syerSysInfo)
	if erro != nil {
		xlog.ErrorLog(appdata.GetSecenName(), "OnNetWorkConnect %v", erro.Error())
		return erro
	}
	restcode := &common_proto.RestInt32CodeMsg{}
	restcode.RestCode = 0
	erro = conn.WritePBMsg(cmdconst.Main_DataCenter, cmddatacenter.Sub_Req_Server_HeartBeat, restcode)
	return erro
}
