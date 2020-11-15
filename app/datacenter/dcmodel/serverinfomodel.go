/*
创建时间: 2020/5/17
作者: zjy
功能介绍:

*/

package dcmodel

//传递的服务器信息
type SeverInfoModel struct {
	AppId     int32     //服务器id 每个服务器id 为一个,同组服务器id 不能重复
	AppKind   int32
	OutAddr   string
	OutProt   string
	ConnID    uint32 // 服务器创建连接时生成id
}

const (
	Link_Server_Succeed = 1 // 连接服务器成功
)



const (
	Register_Server_Succeed   = 1 //注册服务器成功
	Register_Server_Exist   = 2   //对应的连接已经在服务器存在
	Register_Server_NotFindConf = 3 //未找到服务器相关配置
)