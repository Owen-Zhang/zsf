package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
)

const fileConfigDir = "etc"

//FileSource 文件配制实例
type FileSource struct {
	enableWatch bool       //是否监控文件变化
	changed     chan Event //变化提示channel
}

//NewFileSource 实例化文件配制
//会读取etc目录下面的文件信息
func NewFileSource(enableWatch bool) IProvider {
	fileSource := &FileSource{
		enableWatch: enableWatch,
		changed:     make(chan Event, 6),
	}
	if fileSource.enableWatch {
		go fileSource.watch()
	}
	fileSource.obtainAll()
	return fileSource
}

//watch 监控文件变化
func (f FileSource) watch() {
	fileWather, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println(err)
	}
	defer fileWather.Close()
	//这里只能watch文件夹，watch文件不能正常的收到更新信息,而且windows下有两次更新提醒
	//所以这个只能用在linux下面,其它平台没有测试
	if err := fileWather.Add("etc"); err != nil {
		fmt.Printf("add watch file has wrong: %v\n", err)
	}
	for {
		select {
		case event, ok := <-fileWather.Events:
			if !ok {
				return
			}
			//如果新增配制文件需要保存一下，如果监控create会造成在修改文件时同时会
			//触发write和create
			if event.Op&fsnotify.Write == fsnotify.Write {
				content, err := f.get(event.Name)
				if err != nil {
					fmt.Println(err)
					continue
				}
				f.changed <- Event{
					FileName: strings.Split(event.Name, "/")[1],
					Content:  content,
				}
			}
		case err, ok := <-fileWather.Errors:
			if !ok {
				return
			}
			fmt.Printf("watch file has err:%v\n", err)
		}
	}
}

//返回单个文件配制信息
func (f FileSource) get(fileName string) ([]byte, error) {
	return ioutil.ReadFile(filepath.Join(fileConfigDir, fileName))
}

//Notify 更新通知
func (f FileSource) Notify() chan Event {
	return f.changed
}

//obtainAll 加载所有的配制文件
func (f FileSource) obtainAll() {
	files, err := ioutil.ReadDir(fileConfigDir)
	if err != nil {
		panic(err)
	}
	if len(files) == 0 {
		return
	}
	for _, file := range files {
		fileName := file.Name()
		content, err := f.get(fileName)
		if err != nil {
			panic(err)
		}
		f.changed <- Event{
			FileName: fileName,
			Content:  content,
		}
	}
}
