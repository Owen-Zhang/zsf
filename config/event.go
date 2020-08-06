package config

//Event 通知信息实体
type Event struct {
	FileName string //文件名
	Content  []byte //变动后的文件内容
}
