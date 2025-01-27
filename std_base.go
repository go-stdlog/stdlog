package stdlog

import (
	"io"
	"slices"
)

type stdBase struct {
	output     io.Writer
	name       string
	level      Level
	baseFields []*kv
	handler    func(name string, out io.Writer, ev *stdLoggerEvent)
}

func (s *stdBase) dup() *stdBase {
	n := &stdBase{
		output:     s.output,
		name:       s.name,
		level:      s.level,
		baseFields: make([]*kv, 0, len(s.baseFields)),
		handler:    s.handler,
	}
	for _, v := range s.baseFields {
		n.baseFields = append(n.baseFields, v)
	}
	return n
}

func (s *stdBase) Named(name string) Logger {
	n := s.dup()
	if len(s.name) == 0 {
		n.name = name
	} else {
		n.name = s.name + "." + name
	}

	return n
}

func (s *stdBase) WithFields(keysAndValues ...any) Logger {
	if len(keysAndValues)%2 != 0 {
		panic("uneven number of keys and values")
	}

	n := s.dup()
	n.baseFields = slices.Grow(n.baseFields, len(keysAndValues)/2)
	for i := 0; i < len(keysAndValues); i += 2 {
		n.baseFields = append(n.baseFields, stdField(keysAndValues[i], keysAndValues[i+1]))
	}

	return n
}

func (s *stdBase) SetLevel(level Level) { s.level = level }

func (s *stdBase) Leveled(level Level) Logger {
	n := s.dup()
	n.level = level
	return n
}

func (s *stdBase) Debug(msg string, kvs ...any) {
	if s.level != LevelDebug {
		return
	}
	s.handler(s.name, s.output, getEventPool().prepare(LevelDebug, msg, s.baseFields, kvs))
}

func (s *stdBase) Info(msg string, kvs ...any) {
	if s.level > LevelInfo {
		return
	}
	s.handler(s.name, s.output, getEventPool().prepare(LevelInfo, msg, s.baseFields, kvs))
}

func (s *stdBase) Warning(msg string, kvs ...any) {
	if s.level > LevelWarning {
		return
	}
	s.handler(s.name, s.output, getEventPool().prepare(LevelWarning, msg, s.baseFields, kvs))
}

func (s *stdBase) Error(err error, msg string, kvs ...any) {
	if s.level > LevelError {
		return
	}
	ev := getEventPool().prepare(LevelError, msg, s.baseFields, kvs)
	ev.Error = err
	ev.Backtrace = stackTrace(3)
	s.handler(s.name, s.output, ev)
}

func (s *stdBase) Fatal(msg string, kvs ...any) {
	ev := getEventPool().prepare(LevelFatal, msg, s.baseFields, kvs)
	ev.Backtrace = stackTrace(3)
	s.handler(s.name, s.output, ev)
	stdExit(1)
}

func (s *stdBase) FatalError(err error, msg string, kvs ...any) {
	ev := getEventPool().prepare(LevelFatal, msg, s.baseFields, kvs)
	ev.Error = err
	ev.Backtrace = stackTrace(3)
	s.handler(s.name, s.output, ev)
	stdExit(1)
}
