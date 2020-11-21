/*
创建时间: 2020/09/2020/9/1
作者: Administrator
功能介绍:

*/
package dcmodel

import "wengo/model/dbmodels"

const (
	AccountState_None  = 0
	AccountState_Offline = 1
	AccountState_Online = 2
	AccountState_Leave = 3 //主动离开
)


type Account struct {
	PDBAccountData *dbmodels.Accounts
	exprationTime  int64
	accountState   uint8  //账号状态
	clientConnID   uint32 //客户端保存的连接 有可能连接的登录服,有可能连接的是网关
	serverAppID    int32  //只保存appid
}


func AccountStateToStr(accountState   uint8 ) string {
	switch accountState {
	case AccountState_None:
		return "AccountState_None"
	case AccountState_Offline:
		return "AccountState_Offline"
	case AccountState_Online:
		return "AccountState_Online"
	case AccountState_Leave:
		return "AccountState_Leave"
	default:
		return "AccountState_Erro"
	}
	return "AccountState_Erro"
}
func (a *Account) GetServerAppID() int32 {
	return a.serverAppID
}

func (a *Account) SetServerAppID(connServerAppID int32) {
	a.serverAppID = connServerAppID
}

//设置客户端连接
func (a *Account) GetClientConnID() uint32 {
	return a.clientConnID
}

func (a *Account) SetClientConnID(clientConnID uint32) {
	a.clientConnID = clientConnID
}

//账号是否离线
func (a *Account) AccountIsOffline() bool {
	return a.accountState == AccountState_Offline
}

//账号是否离线
func (a *Account) AccountIsOnline() bool {
	return a.accountState == AccountState_Online
}

func (a *Account) AccountState() uint8 {
	return a.accountState
}


func (a *Account) SetAccountState(accountState uint8) {
	a.accountState = accountState
}


func NewAccount(dbAccountData *dbmodels.Accounts) *Account {
	return &Account{PDBAccountData: dbAccountData}
}

func (a *Account) GetExprationTime() int64 {
	return a.exprationTime
}

func (a *Account) SetExprationTime(exprationTime int64) {
	a.exprationTime = exprationTime
}


func(a *Account)SetAccountServerInfo(serverInfo *SeverInfoModel){
	if serverInfo != nil {
		a.SetServerAppID(serverInfo.AppId)
	}
}

func(a *Account)SetAccountLeaveServer(){
	a.SetServerAppID(0)
}
