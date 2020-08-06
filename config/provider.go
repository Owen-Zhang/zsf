package config

//IProvider 配制提供接口
type IProvider interface {
	Notify() chan Event
}
