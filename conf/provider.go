package conf

//IProvider 配制提供接口
type IProvider interface {
	Notify() chan Event
	Get(fileName string) []byte
}
