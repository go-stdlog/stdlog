package stdlog

import (
	"slices"
	"sync"
	"time"
)

type stdLoggerEvent struct {
	Level     Level
	Timestamp time.Time
	Backtrace string
	Message   string
	Error     error
	kvs       []*kv
}

func (e *stdLoggerEvent) prepare(lvl Level, message string, basekvs []*kv, kvs []any) *stdLoggerEvent {
	e.Level = lvl
	e.Timestamp = time.Now()
	e.Message = message
	newSize := len(basekvs) + (len(kvs) / 2)
	e.kvs = slices.Grow(e.kvs, newSize)
	for i := 0; i < len(basekvs); i++ {
		e.kvs = append(e.kvs, basekvs[i])
	}
	for i := 0; i < len(kvs); i += 2 {
		e.kvs = append(e.kvs, stdField(kvs[i], kvs[i+1]))
	}
	e.kvs = e.kvs[:newSize]
	e.Backtrace = ""
	e.Error = nil
	return e
}

var eventPool = sync.Pool{New: func() interface{} { return &stdLoggerEvent{} }}

func getEventPool() *stdLoggerEvent {
	return eventPool.Get().(*stdLoggerEvent)
}
func putEventPool(e *stdLoggerEvent) {
	eventPool.Put(e)
}
