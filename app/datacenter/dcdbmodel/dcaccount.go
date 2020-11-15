/*
创建时间: 2020/08/2020/8/24
作者: Administrator
功能介绍:
数据服数据库需要结构体
*/
package dcdbmodel

type DB_Req_CreateAccount struct {
	Username     string
	Password     string
	ClientType   uint32
	PhoneNum     uint32
	ClientIp     string
	MacAddr      string //Mac地址
}

//账号登录
type DB_Req_LoginAccount struct {
	Username     string
	Password     string
	ClientType   uint32
	ClientIp     string
	MacAddr      string //Mac地址
}