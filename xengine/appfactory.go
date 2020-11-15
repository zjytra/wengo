/*
创建时间: 2020/2/2
作者: zjy
功能介绍:

*/

package xengine

type  AppFactory interface {
	CreateAppBehavor() ServerBehavior
	// CreateConfer() Confer
}


