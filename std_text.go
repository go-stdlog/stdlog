package stdlog

import (
	"fmt"
	"io"
	"slices"
	"strings"
)

type stdLoggerText struct {
	output io.Writer
	name   string
	level  Level
}

var timeFormat = "2006-01-02 15:04:05.000000-07:00"

func NewStd(writer io.Writer) Logger {
	return &stdLoggerText{output: writer}
}

func (s *stdLoggerText) Named(name string) Logger {
	n := new(stdLoggerText)
	if len(s.name) == 0 {
		n.name = name
	} else {
		n.name = s.name + "." + name
	}
	n.level = s.level
	n.output = s.output

	return n
}

func (s *stdLoggerText) SetLevel(level Level) { s.level = level }

func (s *stdLoggerText) Debug(msg string, kvs ...any) {
	if s.level != LevelDebug {
		return
	}
	s.do(getEventPool().prepare(LevelDebug, msg, kvs))
}

func (s *stdLoggerText) Info(msg string, kvs ...any) {
	if s.level >= LevelInfo {
		return
	}
	s.do(getEventPool().prepare(LevelInfo, msg, kvs))
}

func (s *stdLoggerText) Warning(msg string, kvs ...any) {
	if s.level >= LevelWarning {
		return
	}
	s.do(getEventPool().prepare(LevelWarning, msg, kvs))
}

func (s *stdLoggerText) Error(err error, msg string, kvs ...any) {
	if s.level >= LevelError {
		return
	}
	ev := getEventPool().prepare(LevelError, msg, kvs)
	ev.Error = err
	ev.Backtrace = stackTrace(3)
	s.do(ev)
}

func (s *stdLoggerText) Fatal(msg string, kvs ...any) {
	ev := getEventPool().prepare(LevelFatal, msg, kvs)
	ev.Backtrace = stackTrace(3)
	s.do(ev)
}

func (s *stdLoggerText) FatalError(err error, msg string, kvs ...any) {
	ev := getEventPool().prepare(LevelFatal, msg, kvs)
	ev.Error = err
	ev.Backtrace = stackTrace(3)
	s.do(ev)
}

func (s *stdLoggerText) do(ev *stdLoggerEvent) {
	defer putEventPool(ev)
	ensureKV(ev)
	callerLocation := caller(3)
	comps := []string{
		ev.Timestamp.Format(timeFormat),
		"[" + ev.Level.String() + "]",
		s.name,
		callerLocation,
		ev.Message,
		formatKeysValsStd(ev),
	}
	comps = slices.DeleteFunc(comps, func(s string) bool { return len(s) == 0 })
	_, _ = fmt.Fprintf(s.output, "%s\n", strings.Join(comps, " "))
	if ev.Backtrace != "" {
		lines := strings.Split(ev.Backtrace, "\n")
		for i, line := range lines {
			lines[i] = "\t" + line
		}
		_, _ = fmt.Fprintf(s.output, "%s\n", strings.Join(lines, "\n"))
	}
}

func formatKeysValsStd(ev *stdLoggerEvent) string {
	size := len(ev.kvs) / 2
	if ev.Error != nil {
		size++
	}
	res := make([]string, size)
	idx := 0
	for i := 0; i < len(ev.kvs)-1; i += 2 {
		if s, ok := ev.kvs[i+1].(string); ok {
			res[idx] = fmt.Sprintf("%v=%q", ev.kvs[i], s)
		} else {
			res[idx] = fmt.Sprintf("%v=%v", ev.kvs[i], ev.kvs[i+1])
		}
		idx++
	}
	if ev.Error != nil {
		res[idx] = fmt.Sprintf("error=%q", ev.Error.Error())
	}
	return strings.Join(res, " ")
}
