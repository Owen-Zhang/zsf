package logger

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func getLastCheck(now time.Time) uint64 {
	return uint64(now.Year())*1000000 + uint64(now.Month())*10000 + uint64(now.Day())*100 + uint64(now.Hour())
}

type syncBuffer struct {
	*bufio.Writer
	file     *os.File
	count    uint64
	cur      int
	filePath string
	parent   *FileLog
}

func (self *syncBuffer) Sync() error {
	return self.file.Sync()
}

func (self *syncBuffer) close() {
	self.Flush()
	self.Sync()
	self.file.Close()
}

func (self *syncBuffer) write(b []byte) {
	if !self.parent.rotateByHour && self.parent.maxSize > 0 && self.parent.rotateNum > 0 && self.count+uint64(len(b)) >= self.parent.maxSize {
		os.Rename(self.filePath, self.filePath+fmt.Sprintf(".%03d", self.cur))
		self.cur++
		if self.cur >= self.parent.rotateNum {
			self.cur = 0
		}
		self.count = 0
	}
	self.count += uint64(len(b))
	self.Writer.Write(b)
}
