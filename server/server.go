package server

//IServer 服务接口定义
type IServer interface {
	Start() error
	Stop() error
	Info() *ServiceInfo //服务信息
}

type ServiceInfo struct {
	Name     string            `json:"name"`
	AppID    string            `json:"appId"`
	Scheme   string            `json:"scheme"`
	Address  string            `json:"address"`
	Port     int32             `json:"port"`
	Weight   float64           `json:"weight"`
	Enable   bool              `json:"enable"`
	Healthy  bool              `json:"healthy"`
	Metadata map[string]string `json:"metadata"`
	// Group 流量组: 流量在Group之间进行负载均衡
	Group string `json:"group"` //相当于服务组名称(如:商品服务,库存服务)
}
