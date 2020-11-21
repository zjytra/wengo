package dcmodel

import (
	"container/list"
	"wengo/model/dbmodels"
	"wengo/xlog"
	"wengo/xutil/timeutil"
)

//账号管理
type AccountMgr struct {
	UserNameAccountMap map[string]*Account //[使用用户名作为Key]使用账号
	AccountIDMap  map[uint64]*Account      //[使用账号作为key]
	AccountsByMacMap   map[string]uint8    //单台机器只能注册10个号
	pDelAccountsList   *list.List
	pServerInfoMgr     *ServerInfoMgr
}

func NewAccountsMgr(serInfoMgr *ServerInfoMgr) *AccountMgr {
	if serInfoMgr == nil {
		panic("serInfoMgr is nil")
	}
	pAccountMgr := new(AccountMgr)
	pAccountMgr.UserNameAccountMap = make(map[string]*Account)
	pAccountMgr.AccountIDMap = make(map[uint64]*Account)
	pAccountMgr.AccountsByMacMap = make(map[string]uint8)
	pAccountMgr.pDelAccountsList = list.New()
	pAccountMgr.pServerInfoMgr = serInfoMgr
	return pAccountMgr
}

//账号注册向管理类添加
func (this *AccountMgr) AddAccunts(dbAccounts *dbmodels.Accounts) *Account {
	if dbAccounts == nil  || dbAccounts.LoginName == "" || dbAccounts.AccountID == 0 {
		return nil
	}
	paccount := NewAccount(dbAccounts)
	this.UserNameAccountMap[dbAccounts.LoginName] = paccount
	this.AccountIDMap[dbAccounts.AccountID] = paccount
	return  paccount
}

//获取账号信息
func (this *AccountMgr) GetAccountByUserName(username string) *Account {
	return this.UserNameAccountMap[username]
}

//获取账号信息
func (this *AccountMgr) GetAccountByAccountID(accountID uint64) *Account {
	return this.AccountIDMap[accountID]
}

//更具Mac地址获取账号数量
func (this *AccountMgr) GetMacCreateAccount(macStr string) uint8 {
	num,ok := this.AccountsByMacMap[macStr]
	if !ok {
		return 0
	}
	return num
}
//更具ip获取账号数量
func (this *AccountMgr)SetMacCreateAccountNum(macStr string,num uint8)  {
	this.AccountsByMacMap[macStr] = num
}

//定时删除
func (this *AccountMgr)DelAccountOnTimer()  {
	//没有数据就不处理
	if this.pDelAccountsList.Len() == 0 {
		return
	}
	var n *list.Element  //下一个数据的变量临时存放
	for item := this.pDelAccountsList.Front();nil != item ;item = n {
		account,ok :=  item.Value.(*Account)
		if !ok {
			return
		}
		expTime := account.GetExprationTime()
		if expTime == 0 {  //重新登录了
			n = item.Next() //保存下一个数据
			this.pDelAccountsList.Remove(item)
			continue
		}
		if expTime > timeutil.GetCurrentTimeS()  { //这个没有过期证明后面的也没有过期直接返回
			return
		}
		//没有在线 才清除
		if !account.AccountIsOnline() {
			n = item.Next() //保存下一个数据
			this.pDelAccountsList.Remove(item)
		}
	}
}

//定时删除
func (this *AccountMgr)ReleaseData()  {
	var n *list.Element  //清空所有数据
	for item := this.pDelAccountsList.Front();nil != item ;item = n {
		n = item.Next() //保存下一个数据
		this.pDelAccountsList.Remove(item)
	}
	
	for k,_ :=range  this.UserNameAccountMap{
		delete(this.UserNameAccountMap,k)
	}
	this.UserNameAccountMap = nil
	
	for k,_ :=range  this.AccountIDMap{
		delete(this.AccountIDMap,k)
	}
	this.AccountIDMap = nil
	
	for k,_ :=range  this.AccountsByMacMap{
		delete(this.AccountsByMacMap,k)
	}
	this.AccountsByMacMap = nil
	
	xlog.DebugLogNoInScene("ReleaseData 释放账号数据")
}