/*
创建时间: 2020/5/1
作者: zjy
功能介绍:
 根据数据库表生成gofile
*/

package main

import (
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	"github.com/wengo/csvdata"
	"github.com/wengo/dbsys"
	"github.com/wengo/xutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// 列信息
type ColumnInfo struct {
	tablName     string
	columName    string
	columType    string
	columComment string
	isNullable   string //如果可以为null就要设置成指针//add by zjy 20200826
}

var (
	sqlDB      *sql.DB
	gameDBConf *csvdata.Dbconf
	dbwg       sync.WaitGroup
)

func QueryTables() map[string][]*ColumnInfo {
	queryResult := make(map[string][]*ColumnInfo)
	//这里还是增加啊排序,虽然反射是通过字段名称读取的,为了方便查看字段位置，还是匹配下顺序为好
	rows, err := sqlDB.Query("SELECT t.TABLE_NAME,t.COLUMN_NAME,t.COLUMN_TYPE,t.COLUMN_COMMENT,t.IS_NULLABLE FROM information_schema.COLUMNS AS t WHERE t.TABLE_SCHEMA = ? ORDER BY t.TABLE_NAME,t.ORDINAL_POSITION;", gameDBConf.Dbname)
	defer rows.Close()
	if err != nil {
		panic(err)
		return nil
	}
	
	for rows.Next() {
		column := new(ColumnInfo)
		rows.Scan(&column.tablName, &column.columName, &column.columType, &column.columComment,&column.isNullable)
		_, ok := queryResult[column.tablName]
		if !ok {
			var temslice []*ColumnInfo
			temslice = append(temslice, column)
			queryResult[column.tablName] = temslice
			continue
		}
		queryResult[column.tablName] = append(queryResult[column.tablName], column)
	}
	return queryResult
}

func Report() {
	queryResult := QueryTables()
	if queryResult == nil {
		return
	}
	for tableName, columnInfos := range queryResult {
		dbwg.Add(1)
		go ParseColumn(tableName,columnInfos)
	}
}

func ParseColumn(tableName string, columnInfos []*ColumnInfo)  {
	defer  dbwg.Done()
	// 创建csv文件
	fs, err := os.Create(filepath.Join("./model/dbmodels",strings.ToLower(tableName) +"_dbfeild.go"))
	if xutil.IsError(err) {
		return
	}
	defer fs.Close()
	fs.WriteString(fmt.Sprintf("//生成的文件建议不要改动,详见mysql-struct-maker.go ParseColumn方法源码生成格式 \n"))
	fs.WriteString(fmt.Sprintf("package dbmodels \n"))
	fs.WriteString(fmt.Sprintf("\ntype %s struct {\n", xutil.Capitalize(tableName)))
	for _, info := range columnInfos {
		// translate to go struct foramt
		vname := strings.Title(info.columName) // 字段名称
		retype := DBTypeToGoT(info.columType)  //
		if strings.Compare(info.isNullable,"YES") == 0{
			//fix by zjy 20200826
			//如果字段可以为null就需要给字段设置为指针,解决数据库查询Scan给结构体地址赋值错误问题
			//sql: Scan error on column index x, name “x”: converting NULL to int64 is unsupported
			retype = fmt.Sprintf("*%s",retype) // 字段名称
		}
	
		fs.WriteString(fmt.Sprintf("\t%s %s `sql:\"%s\"` // 数据库注释:%s \n ", vname, retype,info.columName, info.columComment))
	}
	fs.WriteString("}\n")
	
}

func DBTypeToGoT(dbtype string) string {
	if dbtype == "" || strings.Compare(dbtype, "") == 0 {
		return ""
	}
	resstr := dbtype
	if strings.Contains(dbtype, "varchar") {
		resstr = "string"
	} else if strings.Contains(dbtype, "date") ||
		      strings.Contains(dbtype, "datetime") ||
		       strings.Contains(dbtype,"timestamp") {
		// resstr = "time.Time" //这里要string类型
		resstr = "string"
	} else if strings.Contains(dbtype, "tinyint") {
		resstr = "int8"
	} else if strings.Contains(dbtype, "smallint") {
		resstr = "int16"
	} else if strings.Contains(dbtype, "integer") {
		resstr = "int32"
	} else if strings.Contains(dbtype, "bigint") {
		resstr = "int64"
	} else if strings.Contains(dbtype, "int") {
		resstr = "int32"
	} else if strings.Contains(dbtype, "double") {
		resstr = "float64"
	} else if strings.Contains(dbtype, "float") {
		resstr = "float32"
	}
	
	// 查看是否是无符号类型
	if strings.Contains(dbtype, "unsigned") {
		resstr = "u" + resstr
	}
	return resstr
}

func main() {
	// set the file path that result save in
	csvdata.SetDbconfMapData("./csv")
	gameDBConf = csvdata.GetDbconfPtr("gamedb")
	if gameDBConf == nil {
		panic("conf == nil")
	}
	// connect to the database
	var  err error
	sqlDB, err = sql.Open("mysql", dbsys.GetMysqlDataSourceName(gameDBConf))
	if err != nil {
		fmt.Println("open  DB  %v", err)
		return
	}
	sqlDB.SetMaxOpenConns(20)
	sqlDB.SetMaxIdleConns(10)
	err = sqlDB.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Database Connect Scuess!")
	Report()
	sqlDB.Close()
	dbwg.Wait()
	fmt.Println("Prase Scuess!")
	
}
