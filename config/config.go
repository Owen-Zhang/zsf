package config

import (
	"fmt"
	"sync"
)

var def *configuration = nil

//configuration 配制只提供监控、更新提醒、获取内容([]byte)
//其缓存也是基于文件名-》文件内容
//使用方得到内容后要将内容反序列化并保存，更新时会通知
type configuration struct {
	provider IProvider
	mutex    *sync.RWMutex     //读写锁
	contents map[string][]byte //文件->配制内容
	watcher  map[string]func() //文件->回调方法(有更新后)
}

//Init 初始化配制信息
func Init(enableWatch bool) {
	def = &configuration{
		mutex:    &sync.RWMutex{},
		contents: make(map[string][]byte),
		watcher:  make(map[string]func()),
		provider: NewFileSource(enableWatch),
	}
	go wathChange()
}

//Get 通过配制文件名获取配制内容
func Get(fileName string) []byte {
	def.mutex.RLock()
	defer def.mutex.RUnlock()
	return def.contents[fileName]
}

//wathChange 监控配制对象更新
func wathChange() {
	for {
		select {
		case event, ok := <-def.provider.Notify():
			if !ok {
				fmt.Println("监控配制更新的channel关闭")
				return
			}
			updateConfig(event)
		}
	}
}

//updateConfig 更新配制信息
func updateConfig(event Event) {
	def.mutex.Lock()
	def.contents[event.FileName] = event.Content
	defer def.mutex.Unlock()
}
