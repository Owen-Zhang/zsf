package config

import (
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
	mutex    sync.RWMutex                      //读写锁
	contents map[string]map[string]interface{} //文件名->配制内容
	watcher  map[string]func()                 //文件名->回调方法(有更新后)
}

//Init 初始化配制信息
func Init() {
	def = &configuration{
		mutex:    sync.RWMutex{},
		contents: make(map[string]map[string]interface{}),
		watcher:  make(map[string]func()),
		provider: NewFileSource(),
	}
	go wathChange()
}

//Get 通过配制文件名获取配制内容(可以返回一些类型对象)
func Get(fileName string) map[string]interface{} {
	def.mutex.RLock()
	result, ok := def.contents[fileName]
	def.mutex.RUnlock()
	if ok {
		return result
	}
	tmpContent := def.provider.Get(fileName)
	if tmpContent == nil || len(tmpContent) == 0 {
		return map[string]interface{}{}
	}
	return updateConfig(fileName, tmpContent)
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
func updateConfig(fileName string, content []byte) map[string]interface{} {
	if content == nil || len(content) == 0 {
		return map[string]interface{}{}
	}
	def.mutex.Lock()
	val, ok := def.contents[fileName]
	if !ok {
		val = make(map[string]interface{})
	}
	unMarshal(fileName, content, val)
	defer def.mutex.Unlock()
	return val
}

//unMarshal yaml文件反序列化成map
func unMarshal(fileName string, content []byte, m map[string]interface{}) {
	if err := yaml.Unmarshal(content, m); err != nil {
		logger.Errorf("反序列化[%s]配制内容失败:%+v", fileName, err)
	}
}
