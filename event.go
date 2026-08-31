package stdlog

import (
	"sync"
	"time"
)

type stdLoggerEvent struct {
	Level     Level
	StackSkip uint
	Timestamp time.Time
	Backtrace string
	Message   string
	Error     error
	kvs       []*kv
}

func (e *stdLoggerEvent) prepare(lvl Level, message string, basekvs []*kv, kvs []any, stackSkip uint) *stdLoggerEvent {
	e.Level = lvl
	e.Timestamp = time.Now()
	e.Message = message
	e.StackSkip = stackSkip
	newSize := len(basekvs) + (len(kvs) / 2)
	e.kvs = make([]*kv, 0, newSize)
	for i := range basekvs {
		e.kvs = append(e.kvs, basekvs[i])
	}
	for i := 0; i < len(kvs); i += 2 {
		e.kvs = append(e.kvs, stdField(kvs[i], kvs[i+1]))
	}
	e.Backtrace = ""
	e.Error = nil
	return e
}

var eventPool = sync.Pool{New: func() any { return &stdLoggerEvent{} }}

func getEventPool() *stdLoggerEvent {
	return eventPool.Get().(*stdLoggerEvent)
}
func putEventPool(e *stdLoggerEvent) {
	eventPool.Put(e)
}
