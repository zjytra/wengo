/*
创建时间: 2020/2/17
作者: zjy
功能介绍:

*/

package apploginsv

import (
	"github.com/golang/protobuf/proto"
	"github.com/wengo/app/netmsgsys"
	"github.com/wengo/cmdconst"
	"github.com/wengo/cmdconst/cmdaccount"
	"github.com/wengo/model"
	"github.com/wengo/msgcode"
	"github.com/wengo/network"
	"github.com/wengo/protobuf/pb/account_proto"
	"github.com/wengo/protobuf/pb/common_proto"
	"github.com/wengo/xlog"
	"github.com/wengo/xutil/strutil"
	"strings"
)

//////////////////////////////////账号注册//////////////////////////////////

//客户端请求注册账号
func RegisterAccountMsgHandler(conn network.Conner, msgdata []byte) error {
	regAccout := &account_proto.CL_LS_ReqRegisterAccoutMsg{}
	erro := proto.Unmarshal(msgdata, regAccout)
	if erro != nil {
		xlog.ErrorLogNoInScene("RegisterAccountMsgHandler %v", erro)
		return erro
	}
	//正在注册中向中心服发送消息
	if netmsgsys.IsRegitering(regAccout.Username) {
		SendRegAccountMsgToClient(conn,  msgcode.AccountCode_IsRegistering)
		return nil
	}
	//验证帐号
	code := IsVerify(conn, regAccout.Username, regAccout.Password)
	if code !=  msgcode.AccountCode_None {
		SendRegAccountMsgToClient(conn,code)
		return nil
	}
	//TODO 验证下版本号
	//
	
	addr := strings.Split(conn.RemoteAddr().String(), ":")
	regAccoutToDc := &account_proto.LS_DC_ReqRegisterAccoutMsg{
		Username:     regAccout.Username,
		Password:     regAccout.Password,
		ClientType:   regAccout.ClientType,
		PhoneNum:     regAccout.PhoneNum,
		ClientConnID: conn.GetConnID(),
		ClientIp:     addr[0],
		MacAddr:      regAccout.GetMacAddr(),
		Version:      regAccout.Version,
	}
	//验证数据
	//验证账号是否合法
	//发送到dbServer 去验证 拉取用户账号信息
	erro = DataCenter.DataCenterConn.WritePBMsg(cmdconst.Main_Account, cmdaccount.Sub_LS_DC_RegisterAccount, regAccoutToDc)
	//向中心服投递的信息记录下
	netmsgsys.SetRegiteringAccount(regAccout.Username)
	return nil
}

//世界服返回账号注册消息
func DCLSRespnRegisterAccount(conn network.Conner, msgdata []byte) error {
	respn := &account_proto.DC_LS_RespnRegisterAccoutMsg{}
	erro := proto.UnmarshalMerge(msgdata, respn)
	if erro != nil {
		return erro
	}
	if respn.ClientConnID == 0 {
		xlog.WarningLogNoInScene("世界服返回客户端连接ID为0")
		return nil
	}
	
	restcode := &common_proto.RestInt32CodeMsg{
		RestCode: respn.GetRestCode(),
	}
	//向客户端返回消息
	erro = Client.tcpServer.WritePBMsgByConnID(respn.ClientConnID, cmdconst.Main_Account, cmdaccount.Sub_LS_C_RegisterAccount, restcode)
	netmsgsys.DelRegitering(respn.GetUsername())
	return erro
}


//给客户端回复消息
func SendRegAccountMsgToClient(conn network.Conner,  code int32) {
	restcode := &common_proto.RestInt32CodeMsg{
		RestCode: code,
	}
	//返回测试数据
	erro := conn.WritePBMsg(cmdconst.Main_Account, cmdaccount.Sub_LS_C_RegisterAccount, restcode)
	if erro != nil {
		xlog.ErrorLogNoInScene("SendRegAccountMsgToClient %v", erro)
	}
}

//////////////////////////////////账号注册结束//////////////////////////////////

//账号是否有效
func IsVerify(conn network.Conner, Username, Password string) int32 {
	//长度验证
	lencode := VerifyStrLen(conn, Username, Password)
	if lencode != msgcode.AccountCode_None{
		return lencode
	}
	//账号包含空格或者非单词字符
	isMatch := strutil.StringHasSpaceOrSpecialChar(Username)
	if isMatch {
		return msgcode.AccountCode_UserNameFormatErro
	}
	//sql注入验证
	isMatch = strutil.StringHasSqlKey(Username)
	if isMatch {
		return msgcode.AccountCode_SqlZhuRu
	}
	//sql注入验证
	isMatch = strutil.StringHasSqlKey(Password)
	if isMatch {
		return msgcode.AccountCode_SqlZhuRu
	}
	return msgcode.AccountCode_None
}

//验证长度
func VerifyStrLen(conn network.Conner, username, password string) int32 {
	strLen := len(username)
	if strLen <= 4 {
		return msgcode.AccountCode_UserNameShort
	}
	if strLen > 11 {
		return msgcode.AccountCode_UserNameLong
	}
	strLen = len(password)
	if strLen < 6 { //密码过短
		return  msgcode.AccountCode_PasswordShort
	}
	if strLen > 18 {
		return msgcode.AccountCode_PasswordLong
	}
	return msgcode.AccountCode_None
}


//////////////////////////////////账号登录//////////////////////////////////

//登陆账号
func LoginAccountMsgHandler(conn network.Conner, msgdata []byte) error {
	reqMsg := &account_proto.CL_LS_ReqLoginMsg{}
	erro := proto.Unmarshal(msgdata, reqMsg)
	if erro != nil {
		xlog.ErrorLogNoInScene("LoginAccountMsgHandler %v", erro)
		return erro
	}
	//正在登录中向中心服发送消息
	if netmsgsys.IsLogining(reqMsg.Username) {
		SendLoginAccountMsgToClient(conn,  msgcode.AccountCode_IsLogining,0,reqMsg.Username)
		return nil
	}
	//已经登录过
	account := pLoginAccountMgr.GetAccountInfoByUserName(reqMsg.Username)
	if account != nil {
		SendLoginAccountMsgToClient(conn,  msgcode.AccountCode_IsLogined,account.AccountID,reqMsg.Username)
		//conn.Close() //已经登录过这里直接 踢掉
		return nil
	}
	//验证帐号
	code := IsVerify(conn, reqMsg.Username, reqMsg.Password)
	if code !=  msgcode.AccountCode_None {
		SendLoginAccountMsgToClient(conn,code,0,reqMsg.Username)
		return nil
	}
	addr := strings.Split(conn.RemoteAddr().String(), ":")
	login := &account_proto.LS_DC_ReqLoginMsg{
		Username:     reqMsg.Username,
		Password:     reqMsg.Password,
		ClientIp:     addr[0],
		ClientType:   reqMsg.ClientType,
		MacAddr:      reqMsg.MacAddr,
		ClientConnID: conn.GetConnID(),
		Version:      reqMsg.Version,
	}
	//验证数据
	//验证账号是否合法
	//发送到dbServer 去验证 拉取用户账号信息
	erro = DataCenter.DataCenterConn.WritePBMsg(cmdconst.Main_Account, cmdaccount.Sub_LS_DC_LoginAccount, login)
	//向中心服投递的信息记录下
	netmsgsys.SetLoginingAccount(reqMsg.Username)
	return nil
}
//给客户端登录消息回复消息
func SendLoginAccountMsgToClient(conn network.Conner,  code int32,accountID  uint64,userName string) {
	restcode := &account_proto.LS_CL_RespnLoginMsg{
		RestCode: code,
		AccountID: accountID,
		Username: userName,
	}
	//返回测试数据
	erro := conn.WritePBMsg(cmdconst.Main_Account, cmdaccount.Sub_LS_C_LoginAccount, restcode)
	if erro != nil {
		xlog.ErrorLogNoInScene("SendLoginAccountMsgToClient %v", erro)
	}
}

//世界服返回账号登录消息
func DCLSRespnLoginAccount(conn network.Conner, msgdata []byte) error {
	respn := &account_proto.DC_LS_RespnLoginAccoutMsg{}
	erro := proto.UnmarshalMerge(msgdata, respn)
	if erro != nil {
		return erro
	}
	xlog.WarningLogNoInScene("世界服返回登录信息",respn)
	//移除限制信息
	netmsgsys.DelLogining(respn.GetUsername())
	if respn.ClientConnID != 0 {
		restcode := &account_proto.LS_CL_RespnLoginMsg{
			RestCode: respn.GetRestCode(),
			AccountID: respn.GetAccountID(),
			Username: respn.GetUsername(),
		}
		//向客户端返回消息
		erro = Client.tcpServer.WritePBMsgByConnID(respn.ClientConnID, cmdconst.Main_Account, cmdaccount.Sub_LS_C_LoginAccount, restcode)
	}else {
		xlog.WarningLogNoInScene("世界服返回客户端连接ID为0")
	}
	//登录成功放到缓存中
	if respn.GetRestCode() == msgcode.AccountCode_Login_Succeed  {
		logina:= new(model.LogionAccount)
		logina.ConnID = respn.GetClientConnID()
		logina.AccountID = respn.GetAccountID()
		logina.Username =respn.GetUsername()
		logina.ClientType =respn.GetClientType()
		pLoginAccountMgr.AddLoginAccount(logina)
	}
	return erro
}

//////////////////////////////////账号登录结束//////////////////////////////////

//////////////////////////////////客户端请求账号离开游戏//////////////////////////////////
func CL_LSReqLeave(conn network.Conner, msgdata []byte) error {
	reqMsg := &account_proto.ClientLeaveMsg{}
	erro := proto.UnmarshalMerge(msgdata, reqMsg)
	if erro != nil {
		return erro
	}
	account := pLoginAccountMgr.GetAccountInfoByAccountID(reqMsg.AccountID)
	if account == nil {  //不是这个号
		conn.Close() //直接踢了
		return nil
	}

	//向中心服发送离线消息
	DataCenter.DataCenterConn.WritePBMsg(cmdconst.Main_Account,cmdaccount.Sub_LS_DC_AccountLeave,reqMsg)
	return erro
}


//////////////////////////////////客户端请求账号离开游戏结束//////////////////////////////////