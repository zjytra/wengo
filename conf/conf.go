/*
创建时间: 2020/2/6
作者: zjy
功能介绍:
配置相关功能
*/

package conf

import (
	"encoding/json"
	"fmt"
	"wengo/xlog"
	"os"
)

type ServerJson struct {
	LogConf  xlog.VolatileLogModel  `json:"LogConf"`
	WorldServerId  int32   `json:WorldServerId`
}

var (
	// conf          *goini.Config
	 SvJson *ServerJson
)

func init()  {
	SvJson = new(ServerJson)
	SvJson.WorldServerId = 15 // 默认值
}

func ReadIni(iniPath string)  {
	// conf = goini.SetConfig(iniPath)
	// conf.ReadList()
	// tem, isok := strconv.Atoi(conf.GetValue("LogConf", "LogQueueCap"))
	// if isok == nil {
	// 	VolatileModel.LogQueueCap = tem
	// }
	// tem, isok = strconv.Atoi(conf.GetValue("LogConf", "ShowLvl"))
	// if isok == nil {
	// 	VolatileModel.ShowLvl = uint16(tem)
	// }
	// isOutStd, isok := strconv.ParseBool(conf.GetValue("LogConf", "IsOutStd"))
	// if isok == nil {
	// 	VolatileModel.IsOutStd = isOutStd
	// }
	// tem, isok = strconv.Atoi(conf.GetValue("LogConf", "FileTimeSpan"))
	// if isok == nil {
	// 	VolatileModel.FileTimeSpan = tem
	// }
}

func ReadJson(iniPath string) error  {
	filePtr, err := os.Open(iniPath)
	if err != nil {
		return err
	}
	defer filePtr.Close()
	// 创建json解码器
	decoder := json.NewDecoder(filePtr)
	err = decoder.Decode(SvJson)
	if err != nil {
		return err
	}
	fmt.Println(SvJson)
	return nil
}