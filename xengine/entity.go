/*
创建时间: 2019/12/22
作者: zjy
功能介绍:

*/

package xengine

type Entity interface {
	GetComponent() (*Component, error)
	AddComponent(comp *Component) (*Component, error)
	HasComponent(comp *Component)
	IsActivity() bool
}
