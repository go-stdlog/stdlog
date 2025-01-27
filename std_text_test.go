package stdlog

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
)

func TestStdText(t *testing.T) {
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
		l := NewStd(buf)
		return l, buf
	}

	messageRegexp := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}\s\d+:\d+:\d+.\d+(?:-\d+:\d+)?\s\[([^]]+)]\s(?:(\S+)\s)?[^:]+:\d+\s([^$]+)\n$`)
	assertMessage := func(t *testing.T, buf *bytes.Buffer, msg string) {
		t.Helper()
		require.Truef(t, messageRegexp.MatchString(buf.String()), "expected %q to match %s", buf.String(), messageRegexp.String())

		groups := messageRegexp.FindStringSubmatch(buf.String())
		if len(groups) != 4 {
			return
		}
		assert.Equal(t, msg, groups[3])
	}
	assertName := func(t *testing.T, buf *bytes.Buffer, name string) {
		t.Helper()
		require.Truef(t, messageRegexp.MatchString(buf.String()), "expected %q to match %s", buf.String(), messageRegexp.String())

		groups := messageRegexp.FindStringSubmatch(buf.String())
		if len(groups) != 4 {
			return
		}
		assert.Equal(t, name, groups[2])
	}
	assertLevel := func(t *testing.T, buf *bytes.Buffer, level string) {
		t.Helper()
		require.Truef(t, messageRegexp.MatchString(buf.String()), "expected %q to match %s", buf.String(), messageRegexp.String())

		groups := messageRegexp.FindStringSubmatch(buf.String())
		if len(groups) != 4 {
			return
		}
		assert.Equal(t, level, groups[1])
	}
	assertPattern := func(t *testing.T, buf *bytes.Buffer, pattern string) {
		t.Helper()
		require.Truef(t, regexp.MustCompile(pattern).MatchString(buf.String()), "expected %q to match %s", buf.String(), pattern)
	}

	t.Run("Named", func(t *testing.T) {
		t.Run("Single", func(t *testing.T) {
			drainStdExit()
			l, out := makeLogger()
			l = l.Named("a")
			l.Info("Hello, World!")
			assertMessage(t, out, "Hello, World!")
			assertName(t, out, "a")

		})
		t.Run("Multiple", func(t *testing.T) {
			drainStdExit()
			l, out := makeLogger()
			l = l.Named("a").Named("b")
			l.Info("Hello, World!")
			assertMessage(t, out, "Hello, World!")
			assertName(t, out, "a.b")
		})
	})
	t.Run("WithFields", func(t *testing.T) {
		drainStdExit()
		l, out := makeLogger()
		l = l.Named("a").Named("b").WithFields("foo", "bar")
		l.Info("Hello, World!")
		assertMessage(t, out, "Hello, World! foo=\"bar\"")
	})
	t.Run("Leveled", func(t *testing.T) {
		drainStdExit()
		l, out := makeLogger()
		l = l.Named("a").Named("b").WithFields("foo", "bar").Leveled(LevelInfo)
		l.Debug("DO NOT APPEAR")
		l.Info("Hello, World!")
		assertMessage(t, out, "Hello, World! foo=\"bar\"")
	})
	t.Run("SetLevel", func(t *testing.T) {
		drainStdExit()
		l, out := makeLogger()
		l = l.Named("a").Named("b").WithFields("foo", "bar")
		l.SetLevel(LevelInfo)
		l.Debug("DO NOT APPEAR")
		l.Info("Hello, World!")
		assertMessage(t, out, "Hello, World! foo=\"bar\"")
	})

	t.Run("Debug", func(t *testing.T) {
		drainStdExit()
		l, out := makeLogger()
		l = l.Named("a").Named("b").WithFields("foo", "bar")
		l.SetLevel(LevelDebug)
		l.Debug("Hello, World!")
		assertMessage(t, out, "Hello, World! foo=\"bar\"")
		assertLevel(t, out, "DEBUG")
	})
	t.Run("Info", func(t *testing.T) {
		drainStdExit()
		l, out := makeLogger()
		l = l.Named("a").Named("b").WithFields("foo", "bar")
		l.Info("Hello, World!")
		assertMessage(t, out, "Hello, World! foo=\"bar\"")
		assertLevel(t, out, "INFO")
	})
	t.Run("Warning", func(t *testing.T) {
		drainStdExit()
		l, out := makeLogger()
		l = l.Named("a").Named("b").WithFields("foo", "bar")
		l.Warning("Hello, World!")
		assertMessage(t, out, "Hello, World! foo=\"bar\"")
		assertLevel(t, out, "WARN")
	})
	t.Run("Error", func(t *testing.T) {
		drainStdExit()
		l, out := makeLogger()
		l = l.Named("a").Named("b").WithFields("foo", "bar")
		err := errors.New("error")
		l.Error(err, "Hello, World!")
		assertPattern(t, out, `Hello, World! foo="bar" error="error"\n(\t[a-zA-Z\.0-9/-]+\n\t{2}[/a-zA-Z0-9_.]+(?::\d+)?\n?)*`)
		assertLevel(t, out, "ERROR")
	})
	t.Run("Error, nil error", func(t *testing.T) {
		drainStdExit()
		l, out := makeLogger()
		l = l.Named("a").Named("b").WithFields("foo", "bar")
		l.Error(nil, "Hello, World!")
		assertPattern(t, out, `Hello, World! foo="bar"\n(\t[a-zA-Z\.0-9/-]+\n\t{2}[/a-zA-Z0-9_.]+(?::\d+)?\n?)*`)
		assertLevel(t, out, "ERROR")
	})
	t.Run("Fatal", func(t *testing.T) {
		drainStdExit()
		l, out := makeLogger()
		l = l.Named("a").Named("b").WithFields("foo", "bar")
		l.Fatal("Hello, World!")
		assertPattern(t, out, `Hello, World! foo="bar"\n(\t[a-zA-Z\.0-9/-]+\n\t{2}[/a-zA-Z0-9_.]+(?::\d+)?\n?)*`)
		assertLevel(t, out, "FATAL")
		assert.NotEmpty(t, stdExitCh)
	})
	t.Run("FatalError", func(t *testing.T) {
		drainStdExit()
		l, out := makeLogger()
		l = l.Named("a").Named("b").WithFields("foo", "bar")
		err := errors.New("error")
		l.FatalError(err, "Hello, World!")
		assertPattern(t, out, `Hello, World! foo="bar" error="error"\n(\t[a-zA-Z\.0-9/-]+\n\t{2}[/a-zA-Z0-9_.]+(?::\d+)?\n?)*`)
		assertLevel(t, out, "FATAL")
	})
}
