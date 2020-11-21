/*
创建时间: 2020/08/2020/8/25
作者: Administrator
功能介绍:

*/
package dbsys

import "github.com/zjytra/wengo/dispatch"

//数据库查询事件对象返回结构体对象
type AsyncDBQueryData struct {
	BDParam  *dispatch.DBEventParam
	Cb        dispatch.OnDBQueryCB //回调方法
	Querystr  string               //封装好的回调语句
}

type DBExcuteEventData struct {
	BDParam  *dispatch.DBEventParam //回调系统确定 回调方法的调度协程
	Cb        dispatch.DBExcuteCallback
	Excutestr string //封装好的回调语句
}





