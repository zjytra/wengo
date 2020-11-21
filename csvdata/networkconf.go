//excle生成文件请勿修改
package csvdata

import (
	"fmt"
	"github.com/zjytra/wengo/csvparse"
	"github.com/zjytra/wengo/xutil"
	"github.com/zjytra/wengo/xlog"
	"sync/atomic"
)

var networkconfAtomic atomic.Value

type  Networkconf struct {
	App_id int32 //#服务器id 字段名称  app_id
	App_kind int32 //服务器类型 字段名称  app_kind
	App_name string //服务器名称 字段名称  app_name
	Out_addr string //外部连接的地址 字段名称  out_addr
	Out_prot string //外部连接端口 字段名称  out_prot
	Max_connect int //最大连接数 字段名称  max_connect
	Goroutines_size int //go协程池数量,连接数的两倍多一点 字段名称  goroutines_size
	Msglen_size uint8 //消息包长字节大小2 字段名称  msglen_size
	Max_msglen uint32 //消息最大长度 字段名称  max_msglen
	Write_cap_num int //连接写的包队列大小 字段名称  write_cap_num
	Checklink_s int //检查连接存活时间间隔秒 字段名称  checklink_s
	Max_rec_msg_ps int //每秒最大收包数量 字段名称  max_rec_msg_ps
	Sys_stauts_port int //系统状态端口 字段名称  sys_stauts_port
	Msg_isencrypt bool //消息是否加密 字段名称  msg_isencrypt
}

func SetNetworkconfMapData(csvpath  string ) {
  	defer xlog.RecoverToStd()
	networkconfAtomic.Store(loadNetworkconfUsedData(csvpath))
}

func loadNetworkconfUsedData(csvpath  string ) map[int32]*Networkconf{
    csvmapdata := csvparse.GetCsvMapData(csvpath + "/networkconf.csv")
	tem := make(map[int32]*Networkconf)
	for _, filedData := range csvmapdata {
		one := new(Networkconf)
		for filedName, filedval := range filedData {
			isok := csvparse.ReflectSetField(one, filedName, filedval)
			xutil.IsError(isok)
			if _,ok := tem[one.App_id]; ok {
				fmt.Println(one.App_id,"重复")
			}
		}
		tem[one.App_id] = one
	}
	return tem
}

func GetNetworkconfPtr(app_id int32) *Networkconf{
    alldata := GetAllNetworkconf()
	if alldata == nil {
		return nil
	}
	if data, ok := alldata[app_id]; ok {
		return data
	}
	return nil
}

func GetAllNetworkconf() map[int32]*Networkconf{
    val := networkconfAtomic.Load()
	if data, ok := val.(map[int32]*Networkconf); ok {
		return data
	}
	return nil
}
