/*
创建时间: 2020/7/6
作者: zjy
功能介绍:

*/

package appclient

import (
	"github.com/golang/protobuf/proto"
	"wengo/network"
	"wengo/protobuf/pb/account_proto"
	"wengo/protobuf/pb/common_proto"
	"wengo/xlog"
)

//注册消息回复
func OnRegisterAccountHanlder(conn network.Conner,msgdata []byte) error{
	restCode := &common_proto.RestInt32CodeMsg{}
	erro := proto.Unmarshal(msgdata,restCode)
	if erro != nil {
		xlog.ErrorLogNoInScene( "OnRegisterAccountHanlder %v", erro)
		return erro
	}
	// 查询内存中是否有这个注册用户
	xlog.DebugLogNoInScene( "OnRegisterAccountHanlder code =%v", restCode.RestCode)
	return  nil
}


//登录消息回复
func OnRespnLoginAccountHanlder(conn network.Conner,msgdata []byte) error{
	restCode := &account_proto.LS_CL_RespnLoginMsg{}
	erro := proto.Unmarshal(msgdata,restCode)
	if erro != nil {
		xlog.ErrorLogNoInScene( "OnRespnLoginAccountHanlder %v", erro)
		return erro
	}
	// 查询内存中是否有这个注册用户
	xlog.DebugLogNoInScene( "OnRespnLoginAccountHanlder code =%v username =%v accountID =%v", restCode.RestCode,restCode.Username,restCode.AccountID)
	return  nil
}