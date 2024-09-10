package stdlog

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

func ensureKV(ev *stdLoggerEvent) {
	if l := len(ev.kvs); l != 0 && l%2 != 0 {
		panic(fmt.Errorf("uneven keys and values passed to %s", ev.Level.String()))
	}
}

func caller(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if ok {
		fileDir := filepath.Base(filepath.Dir(file))
		fileName := filepath.Base(file)
		return fmt.Sprintf("%s/%s:%d", fileDir, fileName, line)
	} else {
		return "unknown:?"
	}
}

func stackTrace(skip int) string {
	pcs := make([]uintptr, 20)
	n := runtime.Callers(skip, pcs)
	pcs = pcs[:n]
	frames := runtime.CallersFrames(pcs)
	var data []string
	for {
		frame, more := frames.Next()
		data = append(data, fmt.Sprintf("%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line))
		if !more {
			break
		}
	}
	return strings.Join(data, "\n")
}
