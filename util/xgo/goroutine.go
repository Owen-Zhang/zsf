package xgo

import (
	"fmt"
	"runtime"

	"github.com/Owen-Zhang/zsf/logger"
)

//Go 包含异常处理的goroutine
func Go(fn func() error) {
	go try(fn)
}

//try 包含recover,并记录相关的日志信息
func try(fn func() error) (err error) {
	defer func() {
		if err := recover(); err != nil {
			_, file, line, _ := runtime.Caller(2)
			logger.FrameLog.Errorf("[recover] (%s:/%d): %+v", file, line, err)
			err = fmt.Errorf("%v", err)
		}
	}()
	return fn()
}
