package logger

import (
	"fmt"
	"time"
)

//此文件是为了不用创建对象logger对象，可以直接用包方法记录日志(使用的是本地文件记录方式), 方便使用

//业务日志记录器
var logging Logger

//FrameLog 此log框架日志记录用,其它不要用此
var FrameLog *Logger = nil

func init() {
	frameLog, err := NewFileLog("logs/frame")
	if err != nil {
		panic(fmt.Sprintf("logger/frame 生成框架日志文件时出错: %v", err))
	}
	FrameLog = NewLogger(INFO, frameLog)
}

func SetLogging(level interface{}, log ILog) {
	logging.SetLogging(level, log)
}

func SetSeverity(level interface{}) {
	logging.SetSeverity(level)
}

func Close() {
	logging.Close()
}

func LogToStderr() {
	logging.LogToStderr()
}

func Debug(args ...interface{}) {
	logging.print(DEBUG, args...)
}

func Debugf(format string, args ...interface{}) {
	logging.printf(DEBUG, format, args...)
}

func Info(args ...interface{}) {
	logging.print(INFO, args...)
}

func Infof(format string, args ...interface{}) {
	logging.printf(INFO, format, args...)
}

func Warning(args ...interface{}) {
	logging.print(WARNING, args...)
}

func Warningf(format string, args ...interface{}) {
	logging.printf(WARNING, format, args...)
}

func Error(args ...interface{}) {
	logging.print(ERROR, args...)
}

func Errorf(format string, args ...interface{}) {
	logging.printf(ERROR, format, args...)
}

func Fatal(args ...interface{}) {
	logging.print(FATAL, args...)
}

func Fatalf(format string, args ...interface{}) {
	logging.printf(FATAL, format, args...)
}

func LogDepth(s Severity, depth int, format string, args ...interface{}) {
	logging.printfDepth(s, depth+1, format, args...)
}

func Printf(format string, args ...interface{}) {
	logging.printfSimple(format, args...)
}

//---------------------------------------
// 设置文件保存日志时的一些配制信息
//--------------------------------------

//Rotate 文件大小分隔配制
func Rotate(rotateNum1 int, maxSize1 uint64) {
	logging.backend.Rotate(rotateNum1, maxSize1)
}

func SetFlushDuration(t time.Duration) {
	logging.backend.SetFlushDuration(t)

}

func SetRotateByHour(rotateByHour bool) {
	logging.backend.SetRotateByHour(rotateByHour)
}

func SetKeepHours(hours uint) {
	logging.backend.SetKeepHours(hours)
}

// func GetLogger() *Logger {
// 	return &logging
// }
