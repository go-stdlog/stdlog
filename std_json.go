package stdlog

import (
	"encoding/json"
	"fmt"
	"io"
	"time"
)

type stdLoggerJSON struct {
	output io.Writer
	name   string
	level  Level
}

func NewStdJSON(writer io.Writer) Logger {
	return &stdLoggerJSON{output: writer}
}

func (s *stdLoggerJSON) Named(name string) Logger {
	n := new(stdLoggerJSON)
	if len(s.name) == 0 {
		n.name = name
	} else {
		n.name = s.name + "." + name
	}
	n.level = s.level
	n.output = s.output

	return n
}

func (s *stdLoggerJSON) SetLevel(level Level) { s.level = level }

func (s *stdLoggerJSON) Debug(msg string, kvs ...any) {
	if s.level != LevelDebug {
		return
	}
	s.do(getEventPool().prepare(LevelDebug, msg, kvs))
}

func (s *stdLoggerJSON) Info(msg string, kvs ...any) {
	if s.level >= LevelInfo {
		return
	}
	s.do(getEventPool().prepare(LevelInfo, msg, kvs))
}

func (s *stdLoggerJSON) Warning(msg string, kvs ...any) {
	if s.level >= LevelWarning {
		return
	}
	s.do(getEventPool().prepare(LevelWarning, msg, kvs))
}

func (s *stdLoggerJSON) Error(err error, msg string, kvs ...any) {
	if s.level >= LevelError {
		return
	}
	ev := getEventPool().prepare(LevelError, msg, kvs)
	ev.Error = err
	ev.Backtrace = stackTrace(1)
	s.do(ev)
}

func (s *stdLoggerJSON) Fatal(msg string, kvs ...any) {
	ev := getEventPool().prepare(LevelFatal, msg, kvs)
	ev.Backtrace = stackTrace(1)
	s.do(ev)
}

func (s *stdLoggerJSON) FatalError(err error, msg string, kvs ...any) {
	ev := getEventPool().prepare(LevelFatal, msg, kvs)
	ev.Error = err
	ev.Backtrace = stackTrace(1)
	s.do(ev)
}

func (s *stdLoggerJSON) do(ev *stdLoggerEvent) {
	defer putEventPool(ev)
	ensureKV(ev)
	callerLocation := caller(2)
	data := map[string]any{
		"time":  ev.Timestamp.UTC().Format(time.RFC3339Nano),
		"level": ev.Level.String(),
	}
	if len(s.name) > 0 {
		data["name"] = s.name
	}
	data["caller"] = callerLocation
	data["msg"] = ev.Message
	if ev.Error != nil {
		data["error"] = ev.Error.Error()
	}
	data["extra"] = stdJSONKVs(ev)
	if len(ev.Backtrace) > 0 {
		data["backtrace"] = ev.Backtrace
	}
	out, err := json.Marshal(data)
	if err != nil {
		s.Error(err, "Failed writing log entry")
		return
	}

	_, _ = s.output.Write([]byte(string(out) + "\n"))
}

func stdJSONKVs(ev *stdLoggerEvent) any {
	mp := map[string]any{}

	for i := 0; i < len(ev.kvs); i += 2 {
		mp[fmt.Sprintf("%v", ev.kvs[i])] = ev.kvs[i+1]
	}

	return mp
}
