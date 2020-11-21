/*
创建时间: 2020/09/2020/9/1
作者: Administrator
功能介绍:
主要防止客户端命令重复发送
*/
package netmsgsys

import "wengo/xutil/timeutil"

var(
	AccountRegiter  map[string]int64 //账号注册验证
	AccountLogin    map[string]int64 //账号登录验证
)

func InitMsgVerify(){
	AccountRegiter = make(map[string]int64)
	AccountLogin = make(map[string]int64)
}

//查看是否正在注册中
func IsRegitering(account string) bool {
	_, isOk := AccountRegiter[account]
	if !isOk {
		return false
	}
	return true
}
//向注册中写数据
func SetRegiteringAccount(account string) {
	AccountRegiter[account] = timeutil.GetCurrentTimeS() + 10
}
//删除数据
func DelRegitering(account string) {
	_, isOk := AccountRegiter[account]
	if !isOk {
		return
	}
	delete(AccountRegiter,account)
}

//查看是否正在注册中
func IsLogining(account string) bool {
	_, isOk := AccountLogin[account]
	if !isOk {
		return false
	}
	return true
}
//向注册中写数据
func SetLoginingAccount(account string) {
	AccountLogin[account] = timeutil.GetCurrentTimeS() + 10
}

//删除数据
func DelLogining(account string) {
	_, isOk := AccountLogin[account]
	if !isOk {
		return
	}
	delete(AccountLogin,account)
}

func ReleaseData(){
	for account,_ := range AccountRegiter {
		delete(AccountRegiter,account)
	}
	AccountRegiter = nil
	for account,_ := range AccountLogin {
		delete(AccountLogin,account)
	}
	AccountLogin = nil
}


func DelExpireData(){
	currentTimes := timeutil.GetCurrentTimeS()
	for account,tval := range AccountRegiter {
		if tval <= currentTimes { //过期了删除
			delete(AccountRegiter,account)
		}
	}
	for account,tval := range AccountLogin {
		if tval <= currentTimes { //过期了删除
			delete(AccountLogin,account)
		}
	}
}