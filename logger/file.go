package logger

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

type FileLog struct {
	mu            sync.Mutex
	dir           string //directory for log files
	files         [numSeverity]syncBuffer
	flushInterval time.Duration
	rotateNum     int
	maxSize       uint64
	fall          bool
	rotateByHour  bool
	lastCheck     uint64
	reg           *regexp.Regexp // for rotatebyhour log del...
	keepHours     uint           // keep how many hours old, only make sense when rotatebyhour is T
}

//NewFileLog new
func NewFileLog(dir string) (*FileLog, error) {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	var fb FileLog
	fb.dir = dir
	for i := 0; i < numSeverity; i++ {
		fileName := path.Join(dir, severityName[i]+".log")
		f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return nil, err
		}

		count := uint64(0)
		stat, err := f.Stat()
		if err == nil {
			count = uint64(stat.Size())
		}
		fb.files[i] = syncBuffer{
			Writer:   bufio.NewWriterSize(f, bufferSize),
			file:     f,
			filePath: fileName,
			parent:   &fb,
			count:    count,
		}

	}
	// default
	fb.flushInterval = time.Second * 3
	fb.rotateNum = 20
	fb.maxSize = 1024 * 1024 * 1024
	fb.rotateByHour = false
	fb.lastCheck = 0
	// init reg to match files
	// ONLY cover this centry...
	fb.reg = regexp.MustCompile("(INFO|ERROR|WARNING|DEBUG|FATAL)\\.log\\.20[0-9]{8}")
	fb.keepHours = 24 * 7

	go fb.flushDaemon()
	go fb.monitorFiles()
	go fb.rotateByHourDaemon()
	return &fb, nil
}

//Flush 将缓冲内容写到硬盘上
func (f *FileLog) Flush() {
	f.mu.Lock()
	defer f.mu.Unlock()
	for i := 0; i < numSeverity; i++ {
		f.files[i].Flush()
		f.files[i].Sync()
	}
}

func (f *FileLog) close() {
	f.Flush()
}

func (f *FileLog) flushDaemon() {
	for {
		time.Sleep(f.flushInterval)
		f.Flush()
	}
}

//rotateByHourDaemon 是否每小时都将原文件改成一个新文件(如果文件为空就没必要新增)
//如:2020081514-INFO.log
//以及删除保存有多长时间(keephours)的日志文件
func (f *FileLog) rotateByHourDaemon() {
	for {
		time.Sleep(time.Second * 1)

		if f.rotateByHour {
			check := getLastCheck(time.Now())
			if f.lastCheck < check {
				for i := 0; i < numSeverity; i++ {
					if !isFileEmpty(f.files[i].filePath) {
						os.Rename(f.files[i].filePath, f.files[i].filePath+fmt.Sprintf(".%d", f.lastCheck))
					}
				}
				f.lastCheck = check
			}

			// also check log dir to del overtime files
			files, err := ioutil.ReadDir(f.dir)
			if err == nil {
				for _, file := range files {
					// exactly match, then we
					if file.Name() == f.reg.FindString(file.Name()) &&
						shouldDel(file.Name(), f.keepHours) {
						os.Remove(filepath.Join(f.dir, file.Name()))
					}
				}
			}
		}
	}
}

//monitorFiles 在合适的时候新增一个初始文件如:INFO.log
func (f *FileLog) monitorFiles() {
	for range time.NewTicker(time.Second * 5).C {
		for i := 0; i < numSeverity; i++ {
			fileName := path.Join(f.dir, severityName[i]+".log")
			if _, err := os.Stat(fileName); err != nil && os.IsNotExist(err) {
				if file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644); err == nil {
					f.mu.Lock()
					f.files[i].close()
					f.files[i].Writer = bufio.NewWriterSize(file, bufferSize)
					f.files[i].file = file
					f.mu.Unlock()
				}
			}
		}
	}
}

//Log 实现接口,记录日志
func (f *FileLog) Log(s Severity, msg []byte) {
	f.mu.Lock()
	switch s {
	case FATAL:
		f.files[FATAL].write(msg)
	case ERROR:
		f.files[ERROR].write(msg)
	case WARNING:
		f.files[WARNING].write(msg)
	case INFO:
		f.files[INFO].write(msg)
	case DEBUG:
		f.files[DEBUG].write(msg)
	}
	if f.fall && s < INFO {
		f.files[INFO].write(msg)
	}
	f.mu.Unlock()
	if s == FATAL {
		f.Flush()
	}
}

//Rotate 设置日志分隔条件
//Rotate(10, 1024*1024*500) //自动切分日志，保留10个文件（INFO.log.000-INFO.log.009，循环覆盖），//每个文件大小为500M, 因为dlog支持多个文件后端， 所以需要为每个file backend指定具体切分数值
func (f *FileLog) Rotate(rotateNum int, maxSize uint64) {
	f.rotateNum = rotateNum
	f.maxSize = maxSize
}

//SetRotateByHour 是否每小时分隔日志文件
func (f *FileLog) SetRotateByHour(rotateByHour bool) {
	f.rotateByHour = rotateByHour
	if f.rotateByHour {
		f.lastCheck = getLastCheck(time.Now())
	} else {
		f.lastCheck = 0
	}
}

//SetKeepHours 日志要保存的时间(多少小时)
func (f *FileLog) SetKeepHours(hours uint) {
	f.keepHours = hours
}

//Fall 是否将info级别还低的日志记录到info文件中
func (f *FileLog) Fall() {
	f.fall = true
}

//SetFlushDuration 设置将缓存写到硬盘的间隔
func (f *FileLog) SetFlushDuration(t time.Duration) {
	if t >= time.Second {
		f.flushInterval = t
	} else {
		f.flushInterval = time.Second
	}
}

//isFileEmpty 检查文件是否为空
func isFileEmpty(filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		fmt.Printf("isFileEmpty 查看文件是否为空时发生错误:%v", info)
		return true
	}
	return info.Size() == 0
}

func shouldDel(fileName string, left uint) bool {
	// tag should be like 2016071114
	tagInt, err := strconv.Atoi(strings.Split(fileName, ".")[2])
	if err != nil {
		return false
	}

	point := time.Now().Unix() - int64(left*3600)

	if getLastCheck(time.Unix(point, 0)) > uint64(tagInt) {
		return true
	}
	return false
}
