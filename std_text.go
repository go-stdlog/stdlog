package stdlog

import (
	"fmt"
	"io"
	"slices"
	"strings"
)

var timeFormat = "2006-01-02 15:04:05.000000-07:00"

func NewStd(writer io.Writer) Logger {
	return &stdBase{
		output: writer,
		handler: func(name string, out io.Writer, ev *stdLoggerEvent) {
			defer putEventPool(ev)
			callerLocation := caller(3)
			comps := []string{
				ev.Timestamp.Format(timeFormat),
				"[" + ev.Level.String() + "]",
				name,
				callerLocation,
				ev.Message,
				formatKeysValsStd(ev),
			}
			comps = slices.DeleteFunc(comps, func(s string) bool { return len(s) == 0 })
			_, _ = fmt.Fprintf(out, "%s\n", strings.Join(comps, " "))
			if ev.Backtrace != "" {
				lines := strings.Split(ev.Backtrace, "\n")
				for i, line := range lines {
					lines[i] = "\t" + line
				}
				_, _ = fmt.Fprintf(out, "%s\n", strings.Join(lines, "\n"))
			}
		},
	}
}

func formatKeysValsStd(ev *stdLoggerEvent) string {
	size := len(ev.kvs)
	if ev.Error != nil {
		size++
	}

	res := make([]string, size)
	for i, kv := range ev.kvs {
		switch t := kv.v.(type) {
		case string:
			res[i] = fmt.Sprintf("%s=%q", kv.k, t)
		case fmt.Stringer:
			res[i] = fmt.Sprintf("%s=%s", kv.k, t.String())
		default:
			res[i] = fmt.Sprintf("%s=%v", kv.k, kv.v)
		}
	}

	if ev.Error != nil {
		res[len(ev.kvs)] = fmt.Sprintf("error=%q", ev.Error.Error())
	}

	return strings.Join(res, " ")
}
