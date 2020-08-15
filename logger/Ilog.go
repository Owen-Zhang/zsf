package logger

import "time"

//ILog 日志记录接口
type ILog interface {
	Log(s Severity, msg []byte)
	close()
	Rotate(rotateNum int, maxSize uint64)
	SetFlushDuration(t time.Duration)
	SetRotateByHour(rotateByHour bool)
	SetKeepHours(hours uint)
}
