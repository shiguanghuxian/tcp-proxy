package program

import "sync"

// Buffer 用于流读取

var (
	// DefaultBufferSize 缓存大小
	DefaultBufferSize = 4096
)

// 创建一个sync.Pool对象
var bufferPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, DefaultBufferSize)
	},
}
