package server

//IServer 服务接口定义
type IServer interface {
	Start() error
	Stop() error
}
