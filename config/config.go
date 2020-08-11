package config

import (
	"fmt"
	"sync"

	"github.com/toolkits/pkg/logger"
	"gopkg.in/yaml.v2"
)

var def *configuration = nil

//configuration 配制只提供监控、更新提醒、获取内容([]byte)
//其缓存也是基于文件名-》文件内容
//使用方得到内容后要将内容反序列化并保存，更新时会通知
type configuration struct {
	provider IProvider
	mutex    sync.RWMutex            //读写锁
	contents map[string][]byte       //文件名->配制内容
	watcher  map[string]func([]byte) //文件名->回调方法(有更新后)
}

//Init 初始化配制信息
func Init() {
	def = &configuration{
		mutex:    sync.RWMutex{},
		contents: make(map[string][]byte),
		watcher:  make(map[string]func([]byte)),
		provider: NewFileSource(),
	}
	go wathChange()
}

//UnmarshalFile 将文件内容反序列化成struct对象
//如果反序列化成功,sValue就是返回对象
func UnmarshalFile(fileName string, sValue interface{}) error {
	contents := get(fileName)
	if contents == nil || len(contents) == 0 {
		return fmt.Errorf("获取文件内容[%s]失败", fileName)
	}
	return yaml.Unmarshal(contents, sValue)
}

//Get 通过配制文件名获取配制内容(可以返回一些类型对象)
func get(fileName string) []byte {
	def.mutex.RLock()
	result, ok := def.contents[fileName]
	def.mutex.RUnlock()
	if ok {
		return result
	}
	tmpContent := def.provider.Get(fileName)
	if tmpContent == nil || len(tmpContent) == 0 {
		return tmpContent
	}

	updateConfig(fileName, tmpContent)
	return tmpContent
}

//wathChange 监控配制对象更新
func wathChange() {
	for {
		select {
		case event, ok := <-def.provider.Notify():
			if !ok {
				logger.Error("监控配制更新的channel关闭")
				return
			}
			updateConfig(event.FileName, event.Content)
		}
	}
}

//updateConfig 更新配制信息
func updateConfig(fileName string, content []byte) {
	if content == nil || len(content) == 0 {
		return
	}
	def.mutex.Lock()
	defer def.mutex.Unlock()
	def.contents[fileName] = content
}
