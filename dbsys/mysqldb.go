// 创建时间: 2019/10/17
// 作者: zjy
// 功能介绍:
// 1.处理mysql相关逻辑
// 2.简单orm操作
// 3.
package dbsys

import (
	"database/sql"
	"fmt"
	"github.com/zjytra/wengo/csvdata"
	"github.com/zjytra/wengo/dispatch"
	"github.com/zjytra/wengo/xlog"
	"reflect"
	"strings"
	"sync"
)

const(
	GameDBName  = "gamedb"
	GameLogDBName = "gamelogdb"
	GameStatisticsDBName = "game_statisticsdb"
	SchameDBName = "information_schema"
)






var (
	GameDB *MySqlDBStore //游戏库
	LogDB *MySqlDBStore  //日志库
	Game_statisticsDB *MySqlDBStore  //分析库
	PDBParamPool *dispatch.DBEventParamPool
	AccountMutex sync.Mutex   //账号锁保证创建账号唯一性
)


func InitGameDB() {
	gameDBCf := csvdata.GetDbconfPtr(GameDBName)
	if gameDBCf == nil {
		panic("gamedb 数据库配置为null")
	}
	if GameDB != nil {
		GameDB.CloseDB()
	}
	// 创建数据库相关操作
	GameDB = NewMySqlDBStore(gameDBCf)
}


func InitLogDB() {
	gameDBCf := csvdata.GetDbconfPtr(GameLogDBName)
	if gameDBCf == nil {
		panic("logdb 数据库配置为null")
	}
	if LogDB != nil {
		LogDB.CloseDB()
	}
	// 创建数据库相关操作
	LogDB = NewMySqlDBStore(gameDBCf)
}

func InitStatisticsDB() {
	gameDBCf := csvdata.GetDbconfPtr(GameStatisticsDBName)
	if gameDBCf == nil {
		panic("StatisticsDB 数据库配置为null")
	}
	if Game_statisticsDB != nil {
		Game_statisticsDB.CloseDB()
	}
	// 创建数据库相关操作
	Game_statisticsDB = NewMySqlDBStore(gameDBCf)
}

//获取连接字符串
func GetMysqlDataSourceName(dbinfo *csvdata.Dbconf) string {
	if dbinfo == nil {
		fmt.Println("dbinfo is nil")
		return ""
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&multiStatements=true",
		dbinfo.Dbusername,
		dbinfo.Dbpwd,
		dbinfo.Ip,
		dbinfo.Dbport,
		dbinfo.Dbname,
		dbinfo.Char_set)
}

//将多行查询的结果转换为string 数组
func RowsToStringsSlices(rows *sql.Rows) [][]string{
	if rows == nil{
		xlog.ErrorLogNoInScene("RowsToStringsSlices rows is nil ")
		return nil
	}
	columns,erro := rows.Columns()  //获取查询出的字段
	if erro != nil {
		xlog.ErrorLogNoInScene("RowsToStringsSlices Columns %v",erro)
		return nil
	}
	columnsCount := len(columns) //字段个数
	if columnsCount == 0 {
		xlog.ErrorLogNoInScene("RowsToStringsSlices columnsCount == 0")
		return nil
	}
	//拼接[][]string 这里作为桶装东西
	values := make([]string, columnsCount)  //每行的值
	scans := make([]interface{}, columnsCount)
	for i := range values {
		scans[i] = &values[i]  //这里存储装数据的地址 给扫描的时候使用
	}
	var strrows [][]string
	for rows.Next() {
		erro = rows.Scan(scans...) //传入values的地址进行赋值
		if erro != nil {
			xlog.ErrorLogNoInScene("RowsToStringsSlices Scan %v",erro)
		    continue
		}
		onerow := make([]string,columnsCount)
		copy(onerow,values) //赋值给新空间,下次扫描的时候要覆盖values
		strrows = append(strrows, onerow)
	}
	return strrows
}

//将单行查询的结果转换为string 数组
func RowToStringSlice(rows *sql.Rows) []string{
	if rows == nil{
		xlog.ErrorLogNoInScene("RowToStringSlice rows is nil ")
		return nil
	}
	if !rows.Next()  {
		xlog.ErrorLogNoInScene("RowToStringSlice 没有查询数据")
		return nil
	}
	columns,erro := rows.Columns()  //获取查询出的字段
	if erro != nil {
		xlog.ErrorLogNoInScene("RowToStringSlice Columns %v",erro)
		return nil
	}
	columnsCount := len(columns) //字段个数
	if columnsCount == 0 {
		xlog.ErrorLogNoInScene("RowToStringSlice columnsCount == 0")
		return nil
	}
	//拼接[][]string 这里作为桶装东西
	values := make([]string, columnsCount)  //每行的值
	scans := make([]interface{}, columnsCount)
	for i := range values {
		scans[i] = &values[i]
	}
	erro = rows.Scan(scans...)
	if erro != nil {
		xlog.ErrorLogNoInScene("RowToStringSlice Scan %v",erro)
		return nil
	}
	return values
}


//查询结果转为结构体
//@param to 可以获取地址的类型,给传入的对象直接赋值
//@param rows 数据库查询结果
func RowToStruct(rows *sql.Rows, to interface{})  {
	if rows == nil{
		xlog.ErrorLogNoInScene("RowToStruct rows is nil ")
		return
	}
	if to == nil {
		xlog.ErrorLogNoInScene("RowToStruct to= nil")
		return
	}
	val := reflect.ValueOf(to)
	if val.Kind() != reflect.Ptr || val.CanAddr() {
		xlog.ErrorLogNoInScene("赋值对象不是指针 %v",val.Kind())
		return
	}
	valElem := val.Elem()
	if valElem.Kind() != reflect.Struct{
		xlog.ErrorLogNoInScene("赋值对象不是结构体 %v,%v",valElem.Kind(),valElem)
		return
	}
	column_names,column_types := GetRowsColsAndTypes(rows)
	if column_names == nil || column_types == nil {
		return
	}
	colLen := len(column_names)
	valElemLen := valElem.NumField();
	if colLen != valElemLen {
		xlog.ErrorLogNoInScene("数据库列数 %v 与结构体列数 %v 不匹配 ",colLen,valElemLen)
		return
	}
	scan_dest := GetStructScanArr(&valElem,column_names,column_types) //获取结构体的扫描地址
	if  scan_dest == nil{
		return
	}
	for rows.Next() {
		erro := rows.Scan(scan_dest...) //传入values的地址进行赋值
		if erro != nil {
			xlog.ErrorLogNoInScene("RowToStruct Scan %v",erro)
			continue
		}
	}
}


//查询结果转为结构体
//@param to 可以是结构体类型 也可以是指针类型,主要创建实例使用
//@return []interface{}  遍历时必须用指针断言
func RowsToStructSlice(rows *sql.Rows,to interface{})[]interface{}  {
	if rows == nil{
		xlog.ErrorLogNoInScene("RowsToStructSlice rows is nil ")
		return nil
	}
	
	if to == nil {
		xlog.ErrorLogNoInScene("RowsToStructSlice to= nil")
		return nil
	}
	tp := reflect.TypeOf(to)
	if tp.Kind() == reflect.Ptr {
		tp = tp.Elem()
	}

	column_names,column_types := GetRowsColsAndTypes(rows)
	if column_names == nil || column_types == nil {
		return  nil
	}
	colLen := len(column_names)
	valElemLen := tp.NumField()
	if colLen != valElemLen {
		xlog.ErrorLogNoInScene("数据库列数 %v 与结构体列数 %v 不匹配 ",colLen,valElemLen)
		return nil
	}
	
	var allData  []interface{}
	for rows.Next() {
		val := reflect.New(tp)//创建个对象
		if val.Kind() != reflect.Ptr || val.CanAddr() {
			xlog.ErrorLogNoInScene("赋值对象不是指针 %v",val.Kind())
			return nil
		}
		// xlog.DebugLogNoInScene("val.Elem() = %v",val.Elem())
		valElem := val.Elem()
		if valElem.Kind() != reflect.Struct {
			xlog.ErrorLogNoInScene("valElem.Kind() %v",valElem.Kind())
			continue
		}
		scan_dest := GetStructScanArr(&valElem,column_names,column_types)
		if scan_dest == nil {
			continue
		}
		erro := rows.Scan(scan_dest...) //传入values的地址进行赋值
		if erro != nil {
			xlog.ErrorLogNoInScene("RowsToStringsSlices Scan %v",erro)
			continue
		}
		// xlog.ErrorLogNoInScene("data = %v ",val)
		allData = append(allData,val.Interface())
	}
	return allData
}

//获取db结果集 的列与类型
func GetRowsColsAndTypes(rows *sql.Rows)(column_names []string,column_types []*sql.ColumnType){
	column_names, erro := rows.Columns()
	if erro != nil {
		xlog.ErrorLogNoInScene("GetRowsColsAndTypes获取数据库列错误 %v",erro)
		return nil,nil
	}
	column_types,erro = rows.ColumnTypes()
	if erro != nil {
		xlog.ErrorLogNoInScene("GetRowsColsAndTypes获取数据库类型错误 %v",erro)
		return  nil,nil
	}
	// 打印数据库类型
	// for _,v := range column_types{
	// 	xlog.DebugLogNoInScene("db =%v,scan =%v",v.DatabaseTypeName() ,v.ScanType())
	// }
	return
}

//获取获取结构体地址为rows.Scan提供容器
func GetStructScanArr(valElem *reflect.Value,column_names []string,column_types []*sql.ColumnType)[]interface{}{
	if valElem == nil {
		xlog.DebugLogNoInScene("valElem = %v",valElem)
		return nil
	}
	scan_dest := []interface{}{} //扫描赋值的
	colLen := len(column_names)
	//结构体字段的地址
	for i := 0; i < colLen; i++ {
		dbColName := strings.Title(column_names[i]) //数据库字段生成结构体首字母大写
		one_value := valElem.FieldByName(dbColName) //用名字去找结构体字段可以更好的匹配
		if !DBTypeMatchFieldType(column_types[i].ScanType().String(),one_value.Kind().String()) {
			xlog.ErrorLogNoInScene("字段 = %v 数据库字段类型 %v 与结构体字段类型 %v不匹配 ",dbColName,column_types[i].ScanType().String(),one_value.Kind().String())
			// continue
			return nil //这里没有匹配上字段就直接返回了,因为这里出错scan也会出错的
		}
		//将结构体的地址赋值给
		scan_dest = append(scan_dest,  one_value.Addr().Interface())
	}
	return  scan_dest
}

//数据库结果字段扫描类型与结构类型进行匹配
func DBTypeMatchFieldType(dbscantp,fieldTp string) bool{
	switch dbscantp {
	case "sql.RawBytes","mysql.NullTime": //fix by zjy 20200826 这里增加ptr 判断
		if fieldTp != "string" && fieldTp != "ptr" {
			return false
		}
	case "uint64","int64","uint32","int32","uint16","int16","uint8","int8","float64","float32":
		if dbscantp != fieldTp && fieldTp != "ptr" {
			return false
		}
	default:
		xlog.ErrorLogNoInScene("DBTypeMatchFieldType 未处理类型%v", dbscantp)
		return false
	}
	return true
}