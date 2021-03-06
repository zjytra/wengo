/*
创建时间: 2020/3/3
作者: zjy
功能介绍:
数据库逻辑封装读写分离
*/

package dbsys

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/zjytra/wengo/csvdata"
	"github.com/zjytra/wengo/dispatch"
	"github.com/zjytra/wengo/xlog"
	"github.com/zjytra/wengo/xutil"
	"github.com/zjytra/wengo/xutil/osutil"
	"github.com/zjytra/wengo/xutil/strutil"
)

const (
	DBQueryRowsCB_Event              = 1 // 查询返回原始的
	DBExcuteEvent                    = 2 // 写事件
	DBQueryRowsToStrSlicesCb_Event   = 3 // 查询返回二维字符串切片多行数据
	DBQueryRowToStrSliceCb_Event     = 4 // 查询返回字符串切片单行数据
	DBQueryRowToStructCb_Event       = 5 // 查询返回单行结构体事件
	DBQueryRowsToStructSliceCb_Event = 6 // 查询多行结构体事件
	DBQueryCustomOneRowQuery_Event   = 7 // 自定义查询
	DBEvent_max                      = 8
)

// 封装数据库处理
type MySqlDBStore struct {
	db            *sql.DB
	dbConf        *csvdata.Dbconf
	quereyEvent   *dispatch.QueueEvent // 数据库查询队列
	writeEvent   *dispatch.QueueEvent // 数据库写队列
	eventDatapool *dispatch.EventDataPool
}

func NewMySqlDBStore(dbconf *csvdata.Dbconf) *MySqlDBStore {

	dbstore := new(MySqlDBStore)
	dbstore.dbConf = dbconf
	
	if erro := dbstore.OpenDB(); erro != nil {
		xlog.ErrorLogNoInScene("open db error: %v ", erro)
		return nil
	}
	dbstore.eventDatapool = dispatch.NewEventDataPool()
	dbstore.quereyEvent = dispatch.NewQueueEvent()
	dbstore.writeEvent = dispatch.NewQueueEvent()
	//开多个线程操作数据库
	i := dbconf.Readnum
	for ; i > 0;i -- {
		dbstore.quereyEvent.AddEventDealer(dbstore.OnDBQuereyEvent)
	}
	//多个线程写
	w:= dbconf.Writenum
	for ; w > 0;w -- {
		dbstore.writeEvent.AddEventDealer(dbstore.OnDBWriteEvent)
	}
	xlog.ErrorLogNoInScene("open db %v succeed", dbconf.Dbname)
	return dbstore
}

func (this *MySqlDBStore) OpenDB() error {
	if this.dbConf == nil {
		return errors.New("dbconf is nil")
	}
	DataSoureName := GetMysqlDataSourceName(this.dbConf)
	if strutil.StringIsNil(DataSoureName) {
		return errors.New(fmt.Sprintf("%v 数据库连接信息为nil",this.dbConf.Dbname))
	}
	var Erro error
	this.db, Erro = sql.Open("mysql", DataSoureName)
	if xutil.IsError(Erro) {
		return Erro
	}
	this.db.SetMaxOpenConns(this.dbConf.Maxopenconns)
	this.db.SetMaxIdleConns(this.dbConf.Maxidleconns)
	if erro := this.db.Ping(); xutil.IsError(erro) {
		this.CloseDB()
		return erro
	}
	return Erro
}

// 关闭数据库
func (this *MySqlDBStore) CloseDB() {
	 erro := this.db.Close()
	if erro != nil {
		xlog.ErrorLogNoInScene("CloseDB %v",erro)
	}
	 //退出
	this.quereyEvent.Release()
	this.writeEvent.Release()
	xlog.ErrorLogNoInScene("CloseDB %v",this.dbConf)
}

func (this *MySqlDBStore) Query(query string, args ...interface{}) (row *sql.Rows, erro error) {
	row, erro = this.db.Query(query, args ...)
	if erro != nil {
		xlog.ErrorLogNoInScene( "%s \n db.Query sql =%s \n erro %v", osutil.GetRuntimeFileAndLineStr(1),query, erro)
		if row != nil {
			erro = row.Close()
		}
		return
	}
	// erro = row.Close() 解析完数据才能关闭
	return
}

func (this *MySqlDBStore) Excute(query string, args ...interface{}) (result sql.Result, erro error) {
	result, erro = this.db.Exec(query, args ...)
	if erro != nil {
		xlog.ErrorLogNoInScene( "%s \n db.Excute sql =%s \n erro %v", osutil.GetRuntimeFileAndLineStr(1),query, erro)
	}
	return
}

func (this *MySqlDBStore) CheckTableExists(tableName string) bool {
	if this.db == nil {
		fmt.Println("CheckTableExists this.db is nil")
		return false
	}
	rows, erro := this.db.Query("SELECT t.TABLE_NAME FROM information_schema.TABLES AS t WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? ", this.dbConf.Dbname, tableName)
	if xutil.IsError(erro) {
		return false
	}
	if rows.Next() {
		return true
	}
	return false
}
