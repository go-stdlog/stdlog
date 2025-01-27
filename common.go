package stdlog

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type kv struct {
	k string
	v any
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
		if frame.Function == "" {
			data = append(data, fmt.Sprintf("%s:%d\n", frame.File, frame.Line))
		} else {
			data = append(data, fmt.Sprintf("%s\n\t%s:%d\n", frame.Function, frame.File, frame.Line))
		}

		if !more {
			break
		}
	}
	return strings.Join(data, "")
}

func stdField(k, v any) *kv {
	return &kv{fmt.Sprintf("%v", k), v}
}

var stdExit func(int) = os.Exit
