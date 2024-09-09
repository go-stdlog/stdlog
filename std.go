package stdlog

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
	"time"
)

type StdLoggerOut struct {
	output io.Writer
	name   string
	level  Level
}

func NewSTD(writer io.Writer) Logger {
	return &StdLoggerOut{output: writer}
}

func (s *StdLoggerOut) Named(name string) Logger {
	n := new(StdLoggerOut)
	*n = *s
	if n.name == "" {
		n.name = name
	} else {
		n.name = fmt.Sprintf("%s.%s", n.name, name)
	}

	return n
}

func (s *StdLoggerOut) SetLevel(level Level) { s.level = level }

func (s *StdLoggerOut) Debug(msg string, kvs ...any) {
	s.Do(now(), caller(), "", msg, LevelDebug, nil, kvs...)
}

func (s *StdLoggerOut) Info(msg string, kvs ...any) {
	s.Do(now(), caller(), "", msg, LevelInfo, nil, kvs...)
}

func (s *StdLoggerOut) Warning(msg string, kvs ...any) {
	s.Do(now(), caller(), "", msg, LevelWarning, nil, kvs...)
}

func (s *StdLoggerOut) Error(err error, msg string, kvs ...any) {
	s.Do(now(), caller(), backtrace(), msg, LevelInfo, err, kvs...)
}

func (s *StdLoggerOut) Fatal(msg string, kvs ...any) {
	s.Do(now(), caller(), backtrace(), msg, LevelFatal, nil, kvs...)
	os.Exit(1)
}

func (s *StdLoggerOut) FatalError(err error, msg string, kvs ...any) {
	s.Do(now(), caller(), backtrace(), msg, LevelFatal, err, kvs...)
}

func (s *StdLoggerOut) Do(ts, caller, bt, msg string, level Level, err error, kvs ...any) {
	if s.level > level {
		return
	}

	ensureKV(level, kvs)

	if err != nil {
		kvs = append(kvs, "error", err)
	}

	if bt != "" {
		kvs = append(kvs, "backtrace", bt)
	}

	message := fmt.Sprintf("%s %s %s %s %s %s\n", ts, level.String(), s.name, caller, msg, formatKeysValsStd(kvs))
	_, _ = fmt.Fprintf(s.output, message)
}

func formatKeysValsStd(kvs []any) string {
	res := make([]string, len(kvs)/2)
	idx := 0
	for i := 0; i < len(kvs)-1; i += 2 {
		res[i] = fmt.Sprintf("%v=%v", kvs[i], kvs[i+1])
		res[idx] = fmt.Sprintf("%v=%v", kvs[i], kvs[i+1])
		idx++
	}
	return strings.Join(res, " ")
}

func ensureKV(level Level, kvs []any) {
	if l := len(kvs); l != 0 && l%2 != 0 {
		panic(fmt.Errorf("uneven keys and values passed to %s", level.String()))
	}
}

func now() string { return time.Now().Format(time.RFC3339Nano) }

func caller() string {
	_, file, line, ok := runtime.Caller(2)
	if ok {
		return fmt.Sprintf("%s:%d", file, line)
	}
	return "unknown"
}

func backtrace() string {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}
