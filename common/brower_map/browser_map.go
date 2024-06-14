package brower_map

import (
	"sync"
)

// GlobalMap 是一个全局的 map，用于存储 key 是 string 类型，value 是任意接口类型的数据
var globalMap = make(map[string]interface{})
var mutex = &sync.Mutex{}

// Push 将数据放入全局 map
func Push(key string, value interface{}) {
	mutex.Lock()
	defer mutex.Unlock()
	globalMap[key] = value
}

// Get 根据 key 从全局 map 中获取对应的值
func Get(key string) interface{} {
	mutex.Lock()
	defer mutex.Unlock()
	return globalMap[key]
}
