package stdlog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStdJSON(t *testing.T) {
	stdExitCh := make(chan int, 100)
	stdExit = func(code int) {
		stdExitCh <- code
	}
	drainStdExit := func() {
		for {
			select {
			case <-stdExitCh:
			default:
				return
			}
		}
	}

	makeLogger := func() (Logger, *bytes.Buffer) {
		buf := new(bytes.Buffer)
		l := NewStdJSON(buf)
		return l, buf
	}

	decode := func(t *testing.T, buf *bytes.Buffer) map[string]any {
		t.Helper()
		v := make(map[string]any)
		err := json.Unmarshal(buf.Bytes(), &v)
		require.NoError(t, err)
		return v
	}

	t.Run("Named", func(t *testing.T) {
		t.Run("Single", func(t *testing.T) {
			drainStdExit()
			l, out := makeLogger()
			l = l.Named("a")
			l.Info("Hello, World!")
			d := decode(t, out)["name"]
			assert.Equal(t, "a", d)

		})
		t.Run("Multiple", func(t *testing.T) {
			drainStdExit()
			l, out := makeLogger()
			l = l.Named("a").Named("b")
			l.Info("Hello, World!")
			d := decode(t, out)["name"]
			assert.Equal(t, "a.b", d)
		})
	})
	t.Run("WithFields", func(t *testing.T) {
		drainStdExit()
		l, out := makeLogger()
		l = l.Named("a").Named("b").WithFields("foo", "bar")
		l.Info("Hello, World!")
		d := decode(t, out)["extra"].(map[string]any)["foo"]
		assert.Equal(t, "bar", d)
	})
	t.Run("Leveled", func(t *testing.T) {
		drainStdExit()
		l, out := makeLogger()
		l = l.Named("a").Named("b").WithFields("foo", "bar").Leveled(LevelInfo)
		l.Debug("DO NOT APPEAR")
		l.Info("Hello, World!")
		dec := decode(t, out)
		assert.Equal(t, "Hello, World!", dec["msg"])
		assert.Equal(t, "INFO", dec["level"])
	})
	t.Run("SetLevel", func(t *testing.T) {
		drainStdExit()
		l, out := makeLogger()
		l = l.Named("a").Named("b").WithFields("foo", "bar")
		l.SetLevel(LevelInfo)
		l.Debug("DO NOT APPEAR")
		l.Info("Hello, World!")
		dec := decode(t, out)
		assert.Equal(t, "Hello, World!", dec["msg"])
		assert.Equal(t, "INFO", dec["level"])
	})

	t.Run("Debug", func(t *testing.T) {
		drainStdExit()
		l, out := makeLogger()
		l = l.Named("a").Named("b").WithFields("foo", "bar")
		l.Debug("Hello, World!")
		dec := decode(t, out)
		assert.Equal(t, "Hello, World!", dec["msg"])
		assert.Equal(t, "DEBUG", dec["level"])
	})
	t.Run("Info", func(t *testing.T) {
		drainStdExit()
		l, out := makeLogger()
		l = l.Named("a").Named("b").WithFields("foo", "bar")
		l.Info("Hello, World!")
		dec := decode(t, out)
		assert.Equal(t, "Hello, World!", dec["msg"])
		assert.Equal(t, "INFO", dec["level"])
	})
	t.Run("Warning", func(t *testing.T) {
		drainStdExit()
		l, out := makeLogger()
		l = l.Named("a").Named("b").WithFields("foo", "bar")
		l.Warning("Hello, World!")
		dec := decode(t, out)
		assert.Equal(t, "Hello, World!", dec["msg"])
		assert.Equal(t, "WARN", dec["level"])
	})
	t.Run("Error", func(t *testing.T) {
		drainStdExit()
		l, out := makeLogger()
		l = l.Named("a").Named("b").WithFields("foo", "bar")
		err := fmt.Errorf("error")
		l.Error(err, "Hello, World!")
		dec := decode(t, out)
		assert.Equal(t, "Hello, World!", dec["msg"])
		assert.Equal(t, "ERROR", dec["level"])
		assert.NotNil(t, dec["backtrace"])
		assert.Equal(t, "error", dec["error"])
	})
	t.Run("Error, nil error", func(t *testing.T) {
		drainStdExit()
		l, out := makeLogger()
		l = l.Named("a").Named("b").WithFields("foo", "bar")
		l.Error(nil, "Hello, World!")
		dec := decode(t, out)
		assert.Equal(t, "Hello, World!", dec["msg"])
		assert.Equal(t, "ERROR", dec["level"])
		assert.NotNil(t, dec["backtrace"])
		assert.Nil(t, dec["error"])
	})
	t.Run("Fatal", func(t *testing.T) {
		drainStdExit()
		l, out := makeLogger()
		l = l.Named("a").Named("b").WithFields("foo", "bar")
		l.Fatal("Hello, World!")
		dec := decode(t, out)
		assert.Equal(t, "Hello, World!", dec["msg"])
		assert.Equal(t, "FATAL", dec["level"])
		assert.NotEmpty(t, stdExitCh)
	})
	t.Run("FatalError", func(t *testing.T) {
		drainStdExit()
		l, out := makeLogger()
		l = l.Named("a").Named("b").WithFields("foo", "bar")
		err := fmt.Errorf("error")
		l.FatalError(err, "Hello, World!")
		dec := decode(t, out)
		assert.Equal(t, "Hello, World!", dec["msg"])
		assert.Equal(t, "FATAL", dec["level"])
		assert.NotNil(t, dec["backtrace"])
		assert.Equal(t, "error", dec["error"])
	})
}
