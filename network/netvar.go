/*
创建时间: 2020/11/21
作者: Administrator
功能介绍:
网络相关变量
*/
package network


var connID      uint32                      // 连接自增ID给其他模块使用
//生成下一个连接ID
func  nextID() uint32 {
	connID++
	if connID == 0 {
		connID++
	}
	return connID
}
