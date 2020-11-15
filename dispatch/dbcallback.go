/*
创建时间: 2020/4/27
作者: zjy
功能介绍:
数据库的回调事件不想定义再这边,但是数据库系统依赖调度系统,所以只有定义这边,定义在数据那边会出现包的循环引用
*/

package dispatch

import (
	"database/sql"
)

//数据库线程直接查询返回结构体不过结构体的要与数据库表字段匹配
type OnDBQueryCB func(dbParam *DBEventParam) error
////数据库执行返回
type DBExcuteCallback func(dbParam  *DBEventParam,result sql.Result) error

////数据库查询返回原始数据
////@param rows 使用了记得关闭
//type DBQueryRowsCallback func(plymark PlayerMark,rows *sql.Rows) error
////数据库线程直接返回多维数组
//type DBQueryRowsToStringSlicesCb func(plymark PlayerMark,rows [][]string) error
////数据库线程返回单行数据
//type DBQueryRowToStringSliceCb func(plymark PlayerMark,row []string) error
////数据库线程返回多行结构体
//type DBQueryRowsToStructSliceCb func(plymark PlayerMark,to []interface{}) error

//向逻辑线程投递查询事件
type DBQueryData struct {
	BDParam *DBEventParam
	Cb      OnDBQueryCB //回调方法
}

//向逻辑线程投递查询事件
type DBExcuteData struct {
	BDParam  *DBEventParam //回调系统确定 回调方法的调度协程
	Cb     DBExcuteCallback //回调方法
	Result sql.Result
}
