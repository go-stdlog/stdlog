package stdlog

import (
	"encoding/json"
	"io"
	"time"
)

func NewStdJSON(writer io.Writer) Logger {
	b := &stdBase{
		output: writer,
	}
	b.handler = func(name string, out io.Writer, ev *stdLoggerEvent) {
		defer putEventPool(ev)
		callerLocation := caller(2)
		data := map[string]any{
			"time":  ev.Timestamp.UTC().Format(time.RFC3339Nano),
			"level": ev.Level.String(),
		}
		if len(name) > 0 {
			data["name"] = name
		}
		data["caller"] = callerLocation
		data["msg"] = ev.Message
		if ev.Error != nil {
			data["error"] = ev.Error.Error()
		}
		if len(ev.kvs) > 0 {
			mp := make(map[string]any, len(ev.kvs))
			for _, kv := range ev.kvs {
				mp[kv.k] = kv.v
			}
			data["extra"] = mp
		}

		if len(ev.Backtrace) > 0 {
			data["backtrace"] = ev.Backtrace
		}
		logMsg, err := json.Marshal(data)
		if err != nil {
			b.Error(err, "Failed writing log entry")
			return
		}

		_, _ = out.Write(logMsg)
		_, _ = out.Write([]byte{'\n'})
	}
	return b
}
