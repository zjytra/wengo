/*
创建时间: 2020/2/2
作者: zjy
功能介绍:

*/

package apploginsv

import (
     "github.com/zjytra/wengo/xengine"
)



type LoginServerFactory struct {

}



func (this *LoginServerFactory)CreateAppBehavor() xengine.ServerBehavior {
     ls := new(LogionServer)
     return ls
}

