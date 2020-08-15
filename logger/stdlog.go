package logger

import (
	"os"
	"time"
)

type stdLog struct{}

func (f *stdLog) Log(s Severity, msg []byte) {
	os.Stdout.Write(msg)
}

func (f *stdLog) close() {

}

func (f *stdLog) Rotate(rotateNum int, maxSize uint64) {

}
func (f *stdLog) SetFlushDuration(t time.Duration) {

}
func (f *stdLog) SetRotateByHour(rotateByHour bool) {

}
func (f *stdLog) SetKeepHours(hours uint) {

}
