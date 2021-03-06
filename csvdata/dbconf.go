//excle生成文件请勿修改
package csvdata

import (
	"fmt"
	"github.com/zjytra/wengo/csvparse"
	"github.com/zjytra/wengo/xutil"
	"github.com/zjytra/wengo/xlog"
	"sync/atomic"
)

var dbconfAtomic atomic.Value

type  Dbconf struct {
	Dbname string //#数据库名称 字段名称  dbname
	Ip string //ip地址 字段名称  ip
	Dbport string //端口号 字段名称  dbport
	Dbusername string //用户名 字段名称  dbusername
	Dbpwd string //密码 字段名称  dbpwd
	Maxopenconns int //最大链接数 字段名称  maxopenconns
	Maxidleconns int //闲置连接数 字段名称  maxidleconns
	Readnum int8 //读协程数量 字段名称  readnum
	Writenum int8 //写协程数量 字段名称  writenum
	Char_set string //字符集 字段名称  char_set
}

func SetDbconfMapData(csvpath  string ) {
  	defer xlog.RecoverToStd()
	dbconfAtomic.Store(loadDbconfUsedData(csvpath))
}

func loadDbconfUsedData(csvpath  string ) map[string]*Dbconf{
    csvmapdata := csvparse.GetCsvMapData(csvpath + "/dbconf.csv")
	tem := make(map[string]*Dbconf)
	for _, filedData := range csvmapdata {
		one := new(Dbconf)
		for filedName, filedval := range filedData {
			isok := csvparse.ReflectSetField(one, filedName, filedval)
			xutil.IsError(isok)
			if _,ok := tem[one.Dbname]; ok {
				fmt.Println(one.Dbname,"重复")
			}
		}
		tem[one.Dbname] = one
	}
	return tem
}

func GetDbconfPtr(dbname string) *Dbconf{
    alldata := GetAllDbconf()
	if alldata == nil {
		return nil
	}
	if data, ok := alldata[dbname]; ok {
		return data
	}
	return nil
}

func GetAllDbconf() map[string]*Dbconf{
    val := dbconfAtomic.Load()
	if data, ok := val.(map[string]*Dbconf); ok {
		return data
	}
	return nil
}
