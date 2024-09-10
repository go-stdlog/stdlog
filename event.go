package stdlog

import (
	"sync"
	"time"
)

type stdLoggerEvent struct {
	Level     Level
	Timestamp time.Time
	Backtrace string
	Message   string
	Error     error
	kvs       []any
}

func (e *stdLoggerEvent) prepare(lvl Level, message string, kvs []any) *stdLoggerEvent {
	e.Level = lvl
	e.Timestamp = time.Now()
	e.Message = message
	e.kvs = kvs
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
