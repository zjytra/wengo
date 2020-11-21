package datacenter

import (
	"github.com/golang/protobuf/proto"
	"wengo/app/datacenter/dcdbmodel"
	"wengo/cmdconst"
	"wengo/cmdconst/cmdaccount"
	"wengo/dbsys"
	"wengo/msgcode"
	"wengo/network"
	"wengo/protobuf/pb/account_proto"
	"wengo/xlog"
)

//////////////////////////////////账号注册//////////////////////////////////

//登录向中心服务器请求账号创建
func LSDCRegisterAccountMsgHandler(conn network.Conner, msgdata []byte) error {
	regAccout := &account_proto.LS_DC_ReqRegisterAccoutMsg{}
	erro := proto.Unmarshal(msgdata, regAccout)
	if erro != nil {
		xlog.ErrorLogNoInScene( "ServerHeartBeat %v", erro)
		return erro
	}
	Paccounts := PaccountMgr.GetAccountByUserName(regAccout.Username)
	if Paccounts != nil  {	//账号已经存在
		reqCode := &account_proto.DC_LS_RespnRegisterAccoutMsg{
			ClientConnID: regAccout.ClientConnID,
			RestCode: msgcode.AccountCode_IsExsist,
			Username: regAccout.GetUsername(),
		}
		conn.WritePBMsg(cmdconst.Main_Account,cmdaccount.Sub_DC_LS_RegisterAccount,reqCode)
	 return erro
	}
	//单台机器超过注册码10个
	accountNumByMac  := PaccountMgr.GetMacCreateAccount(regAccout.MacAddr)
	if accountNumByMac >= 9 {
		reqCode := &account_proto.DC_LS_RespnRegisterAccoutMsg{
			ClientConnID: regAccout.ClientConnID,
			RestCode: msgcode.AccountCode_MACAccountNumIsMore,
			Username: regAccout.GetUsername(),
		}
		conn.WritePBMsg(cmdconst.Main_Account,cmdaccount.Sub_DC_LS_RegisterAccount,reqCode)
		return nil
	}
	param := dbsys.PDBParamPool.Pop()
	account := PAccountPool.Pop()
	param.ClientConnID = regAccout.GetClientConnID()
	param.ServerConnID = conn.GetConnID()
	param.CbDispSys = PdataCenter.dispSys
	param.ReqParam = &dcdbmodel.DB_Req_CreateAccount{
		Username:regAccout.GetUsername(),
		Password:regAccout.GetPassword(),
		ClientIp: regAccout.GetClientIp(),
		PhoneNum:regAccout.GetPhoneNum(),
		ClientType: regAccout.GetClientType(),
		MacAddr: regAccout.GetMacAddr(),
	}
	account.SetDBEventParam(param)
	dbsys.GameDB.AsyncCustomOneRowQuery(account)
	return erro
}
//////////////////////////////////账号注册结束//////////////////////////////////

//////////////////////////////////账号登录//////////////////////////////////

