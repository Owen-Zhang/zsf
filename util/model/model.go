package model

//ResData 服务返回数据结构
type ResData struct {
	Status  int         `json:"status"`  //返回的状态码 0:失败，1:成功
	Message string      `json:"message"` //消息提示内容
	Data    interface{} `json:"data"`    //数据内容
}

const (
	Success = 1
	Fail    = 0
)
