/*
创建时间: 2020/2/2
作者: zjy
功能介绍:

*/

package appworldsv

import (
     "github.com/wengo/xengine"
)



type WorldSvFactory struct {

}



func (this *WorldSvFactory)CreateAppBehavor() xengine.ServerBehavior {
     ls := new(WorldServer)
     return ls
}

