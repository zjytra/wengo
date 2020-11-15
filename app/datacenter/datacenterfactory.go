/*
创建时间: 2020/2/2
作者: zjy
功能介绍:
数据服务器工厂
*/

package datacenter

import (
     "github.com/wengo/xengine"
)



type DataCenterFactory struct {

}



func (this *DataCenterFactory)CreateAppBehavor() xengine.ServerBehavior {
     ls := new(DataCenter)
     return ls
}

