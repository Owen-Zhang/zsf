package xerrors

//BusinessError 业务错误(需要展示到前台)
type BusinessError struct {
	Message string
}

func (b BusinessError) Error() string {
	return b.Message
}

func (b BusinessError) String() string {
	return b.Message
}
