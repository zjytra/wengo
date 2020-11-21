/*
创建时间: 2020/4/28
作者: zjy
功能介绍:

*/

package csvdata

import "github.com/panjf2000/ants/v2"

//初始化登陆服数据
func LoadCommonCsvData()  {
	SetNetworkconfMapData(csvPath)
	SetDbconfMapData(csvPath)
}


func ReLoadCommonCsvData(workPool *ants.Pool)  {
	workPool.Submit(LoadCommonCsvData)
}