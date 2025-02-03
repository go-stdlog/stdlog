package stdlog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

type jKV struct {
	Name  string
	Value any
}

func (j *jKV) Marshal() ([]byte, error) {
	d, err := json.Marshal(j.Value)
	if err != nil {
		return nil, err
	}
	return []byte(fmt.Sprintf("%q: %s", j.Name, d)), nil
}

func NewStdJSON(writer io.Writer) Logger {
	b := &stdBase{
		output: writer,
	}
	b.handler = func(name string, out io.Writer, ev *stdLoggerEvent) {
		defer putEventPool(ev)
		callerLocation := caller(3)
		size := 4
		if len(name) > 0 {
			size++
		}
		if ev.Error != nil {
			size++
		}
		if len(ev.Backtrace) > 0 {
			size++
		}
		if len(ev.kvs) > 0 {
			size++
		}
		data := make([]jKV, 0, size)
		data = append(data, jKV{"time", ev.Timestamp.Format(time.RFC3339Nano)})
		data = append(data, jKV{"level", ev.Level.String()})

		if len(name) > 0 {
			data = append(data, jKV{"name", name})
		}

		data = append(data, jKV{"caller", callerLocation})
		data = append(data, jKV{"msg", ev.Message})
		if ev.Error != nil {
			data = append(data, jKV{"error", ev.Error.Error()})
		}
		if len(ev.kvs) > 0 {
			mp := make(map[string]any, len(ev.kvs))
			for _, kv := range ev.kvs {
				mp[kv.k] = kv.v
			}
			data = append(data, jKV{"extra", mp})
		}

		if len(ev.Backtrace) > 0 {
			data = append(data, jKV{"backtrace", ev.Backtrace})
		}
		buf := bytes.NewBuffer(nil)
		buf.WriteByte('{')
		fields := make([]string, 0, len(data))
		var err error
		for _, datum := range data {
			var d []byte
			d, err = datum.Marshal()
			if err != nil {
				break
			}
			fields = append(fields, string(d))
		}
		if err != nil {
			b.Error(err, "Failed writing log entry")
			return
		}
		buf.Write([]byte(strings.Join(fields, ", ")))
		buf.Write([]byte("}\n"))
		_, _ = out.Write(buf.Bytes())
	}
	return b
}
