package conf

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/Owen-Zhang/zsf/logger"
	"github.com/fsnotify/fsnotify"
)

const fileConfigDir = "etc"

//FileSource 文件配制实例
type FileSource struct {
	changed chan Event //变化提示channel
}

//NewFileSource 实例化文件配制
//会读取etc目录下面的文件信息
func NewFileSource() IProvider {
	fileSource := &FileSource{
		changed: make(chan Event, 6),
	}
	go fileSource.watch()
	return fileSource
}

//watch 监控文件变化
func (f FileSource) watch() {
	fileWather, err := fsnotify.NewWatcher()
	if err != nil {
		logger.FrameLog.Errorf("监控配制文件变化出现错误: %+v", err)
		return
	}
	defer fileWather.Close()
	//这里只能watch文件夹，watch文件不能正常的收到更新信息,而且windows下有两次更新提醒
	//所以这个只能用在linux下面,其它平台没有测试
	if err := fileWather.Add("etc"); err != nil {
		logger.FrameLog.Errorf("增加监控目录[etc]出现错误: %+v", err)
		return
	}
	for {
		select {
		case event, ok := <-fileWather.Events:
			if !ok {
				return
			}
			//如果新增配制文件需要保存一下，如果监控create会造成在修改文件时同时会
			//触发write和create
			if event.Op&fsnotify.Write == fsnotify.Write && !strings.HasSuffix(event.Name, ".swp") {
				//fileName := strings.Split(event.Name, "\\")[1]
				fileName := strings.Split(event.Name, "/")[1]
				content := f.Get(fileName)
				if err != nil {
					continue
				}
				f.changed <- Event{
					FileName: fileName,
					Content:  content,
				}
			}
		case err, ok := <-fileWather.Errors:
			if !ok {
				return
			}
			logger.FrameLog.Errorf("fileWather返回错误: %+v", err)
		}
	}
}

//Get 返回单个文件配制信息
func (f FileSource) Get(fileName string) (result []byte) {
	result, err := ioutil.ReadFile(filepath.Join(fileConfigDir, fileName))
	if err != nil {
		logger.FrameLog.Errorf("返回单个文件配制信息出错: %+v", err)
		return
	}
	return
}

//Notify 更新通知
func (f FileSource) Notify() chan Event {
	return f.changed
}
