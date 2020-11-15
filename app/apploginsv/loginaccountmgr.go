/*
创建时间: 2020/09/2020/9/3
作者: Administrator
功能介绍:

*/
package apploginsv

import (
	"github.com/wengo/cmdconst"
	"github.com/wengo/cmdconst/cmdaccount"
	"github.com/wengo/model"
	"github.com/wengo/protobuf/pb/account_proto"
)

type LoginAccountMgr struct {
	loginAccount    map[uint64]*model.LogionAccount //[AccountID]Account登录的账号信息
	connMapAccount  map[uint32]uint64               //[连接id]与账号id绑定
	userNameAccount map[string]*model.LogionAccount //根据username判断
}

func NewLoginAccountMgr() *LoginAccountMgr {
	loginAccount := new(LoginAccountMgr)
	loginAccount.loginAccount = make(map[uint64]*model.LogionAccount)
	loginAccount.connMapAccount = make(map[uint32]uint64)
	loginAccount.userNameAccount =make(map[string]*model.LogionAccount)
	return loginAccount
}

func (this *LoginAccountMgr) AddLoginAccount(account *model.LogionAccount) {
	this.loginAccount[account.AccountID] = account
	//连接与accountId映射
	this.connMapAccount[account.ConnID] = account.AccountID
	this.userNameAccount[account.Username] = account
}

//根据连接获取账号信息
func (this *LoginAccountMgr) GetAccountInfoByConnID(connID uint32) *model.LogionAccount {
	accountID, ok := this.connMapAccount[connID]
	if !ok {
		return nil
	}
	return this.GetAccountInfoByAccountID(accountID)
}

//获取账号信息
func (this *LoginAccountMgr) GetAccountInfoByUserName(userName string) *model.LogionAccount {
	return this.userNameAccount[userName]
}
//获取账号信息
func (this *LoginAccountMgr) GetAccountInfoByAccountID(accountID uint64) *model.LogionAccount {
	return this.loginAccount[accountID]
}

//获取账号信息
func (this *LoginAccountMgr) GetConnIDByAccountID(accountID uint64) uint32 {
	account := this.GetAccountInfoByAccountID(accountID)
	if account == nil {
		return 0
	}
	return account.ConnID
}

//被动断开连接
func (this *LoginAccountMgr) Offline(connID uint32) {
	//移除连接信息
	//这里一定要再中心服登录成功才有数据
	accountID, ok := this.connMapAccount[connID]
	if !ok {
		return
	}
	delete(this.connMapAccount, connID) //移除数据
	paccount, ok := this.loginAccount[accountID]
	if !ok {
		return
	}
	//告诉数据中心用户离线
	offLineMsg := &account_proto.LS_DC_ClientOffLineMsg{
		Username:     paccount.Username,
		AccountID:    paccount.AccountID,
	}
	delete(this.loginAccount, accountID) //移除数据
	//向中心服发送离线消息
	DataCenter.DataCenterConn.WritePBMsg(cmdconst.Main_Account, cmdaccount.Sub_LS_DC_AccountOffline, offLineMsg)
	
	_, ok = this.userNameAccount[paccount.Username]
	if !ok {
		return
	}
	delete(this.userNameAccount, paccount.Username) //移除数据
	
}

func (this *LoginAccountMgr) Leave(accountID uint64) bool {
	paccount, ok := this.loginAccount[accountID]
	if !ok {
		return false
	}
	delete(this.loginAccount, accountID) //移除数据
	//移除连接信息
	_, ok = this.connMapAccount[paccount.ConnID]
	if !ok {
		return false
	}
	delete(this.connMapAccount, paccount.ConnID) //移除数据
	return true
}

func (this *LoginAccountMgr) Release() {
	for k, _ := range this.loginAccount {
		delete(this.loginAccount, k)
	}
	this.loginAccount = nil
	for k, _ := range this.connMapAccount {
		delete(this.connMapAccount, k)
	}
	this.connMapAccount = nil
}
