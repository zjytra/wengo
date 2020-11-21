/*
创建时间: 2020/09/2020/9/18
作者: Administrator
功能介绍:

*/
package datacenter

import (
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/zjytra/wengo/app/datacenter/dcmodel"
	"github.com/zjytra/wengo/network"
	"github.com/zjytra/wengo/protobuf/pb/account_proto"
	"github.com/zjytra/wengo/xlog"
	"github.com/zjytra/wengo/xutil/osutil"
	"github.com/zjytra/wengo/xutil/timeutil"
)

//////////////////////////////////账号离线//////////////////////////////////
func LSDCAccountOffline(conn network.Conner, msgdata []byte) error {
	reqMsg := &account_proto.LS_DC_ClientOffLineMsg{}
	erro := proto.UnmarshalMerge(msgdata, reqMsg)
	if erro != nil {
		return erro
	}
	xlog.DebugLogNoInScene("账号%v 离线",reqMsg)
	account := PaccountMgr.GetAccountByAccountID(reqMsg.GetAccountID())
	//没有账号信息直接返回
	if account == nil {
		return nil
	}
	//服务器信息都为nil
	serverInfo := PServerInfoMgr.GetServerInfoByAppID(account.GetServerAppID())
	if serverInfo == nil {
		return errors.New(fmt.Sprintf("%v,服务器信息为nil",osutil.GetRuntimeFileAndLineStr(0)))
	}
	conserverInfo := PServerInfoMgr.GetServerInfoByConnID( conn.GetConnID())
	if conserverInfo == nil { //服务器信息不存在
		return errors.New(fmt.Sprintf("%v,当前连接的服务器信息为nil",osutil.GetRuntimeFileAndLineStr(0)))
	}
	//不是同一个服务器连接
	if serverInfo.ConnID != conserverInfo.ConnID  {
		xlog.DebugLogNoInScene("账号%v不在%v,在%v服务器连接%v,账号连接%v",reqMsg.GetAccountID(),serverInfo.AppId,account.GetServerAppID(),serverInfo.ConnID,conn.GetConnID() )
		return nil
	}
	//账号没有在线也不处理
	if !account.AccountIsOnline()  {
		return nil
	}
	
	account.SetAccountState(dcmodel.AccountState_Offline)
	now := timeutil.NowAddDate(0,0,7) //设置7天后过期
	account.SetExprationTime(now.Unix())
	account.SetClientConnID(0)
	newaccount := PaccountMgr.GetAccountByAccountID(reqMsg.GetAccountID())
	xlog.DebugLogNoInScene("离线",newaccount.PDBAccountData,"账号新状态",newaccount)
	return nil
}

//////////////////////////////////账号离线结束//////////////////////////////////



//////////////////////////////////账号离开//////////////////////////////////
func LSDCAccountLeave(conn network.Conner, msgdata []byte) error {
	reqMsg := &account_proto.ClientLeaveMsg{}
	erro := proto.UnmarshalMerge(msgdata, reqMsg)
	if erro != nil {
		return erro
	}
	account := PaccountMgr.GetAccountByAccountID(reqMsg.GetAccountID())
	//没有账号信息直接返回
	if account == nil {
		return nil
	}
	//服务器信息都为nil
	serverInfo := PServerInfoMgr.GetServerInfoByAppID(account.GetServerAppID())
	if serverInfo == nil {
		return errors.New(fmt.Sprintf("%v,服务器信息为nil",osutil.GetRuntimeFileAndLineStr(0)))
	}
	//不是同一个服务器
	if  serverInfo.ConnID != conn.GetConnID()   {
		xlog.DebugLogNoInScene("服务器连接%v,账号连接%v不匹配",conn.GetConnID(),serverInfo.ConnID)
		return nil
	}
	//可以是在线 可以是离线
	//if !account.AccountIsOnline() {
	//	return nil
	//}
	
	account.SetAccountState(dcmodel.AccountState_Leave)
	now := timeutil.NowAddDate(0,0,7) //设置7天后过期
	account.SetExprationTime(now.Unix())
	account.SetClientConnID(0)
	account.SetAccountLeaveServer()
	
	return nil
}

//////////////////////////////////账号离线结束//////////////////////////////////