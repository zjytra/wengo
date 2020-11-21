/*
创建时间: 2020/2/2
作者: zjy
功能介绍:

*/

package appgatesv

import (
     "github.com/zjytra/wengo/xengine"
)



type GateSvFactory struct {

}



func (this *GateSvFactory)CreateAppBehavor() xengine.ServerBehavior {
     ls := new(GateServer)
     return ls
}

