/*
创建时间: 2019/12/25
作者: zjy
功能介绍:
路径相关管理
*/

package model

import (
	"path"
)

type PathModel struct {
	AppRootPath  string //  程序(main)根路径
	CsvPath      string
	LogsPath     string
	ConfPath     string
	ConfjsonPath string
}

// 创建PathModel
func NewPathModel() *PathModel {
	return new(PathModel)
}

// 路径管理相关函数
func (pthpro *PathModel) SetRootPath(pwd string ) {
	pthpro.AppRootPath = pwd
}

func (pthpro *PathModel) InitPathModel() {
	pthpro.ConfPath = path.Join(pthpro.AppRootPath, "configs")
	pthpro.CsvPath = path.Join(pthpro.AppRootPath, "csv")
	pthpro.LogsPath = path.Join(pthpro.AppRootPath, "logs")
	pthpro.ConfjsonPath = path.Join(pthpro.ConfPath, "serverconf.json")
}

