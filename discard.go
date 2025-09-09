package stdlog

import "os"

type noop string

const Discard noop = "noop"

func (n noop) Named(name string) Logger { return n }

func (n noop) SetLevel(level Level) {}

func (n noop) Leveled(level Level) Logger { return n }

func (n noop) Skipping(count uint) Logger { return n }

func (n noop) Debug(msg string, keysAndValues ...any) {
	assertKvs(LevelDebug.String(), keysAndValues...)
}

func (n noop) Info(msg string, keysAndValues ...any) {
	assertKvs(LevelInfo.String(), keysAndValues...)
}

func (n noop) Warning(msg string, keysAndValues ...any) {
	assertKvs(LevelWarning.String(), keysAndValues...)
}

func (n noop) Error(err error, msg string, keysAndValues ...any) {
	assertKvs(LevelError.String(), keysAndValues...)
}

func (n noop) Fatal(msg string, keysAndValues ...any) {
	assertKvs(LevelFatal.String(), keysAndValues...)
	os.Exit(1)
}

func (n noop) FatalError(err error, msg string, keysAndValues ...any) {
	assertKvs(LevelFatal.String()+"_ERROR", keysAndValues...)
	os.Exit(1)
}

func (n noop) WithFields(keysAndValues ...any) Logger {
	assertKvs("WithFields", keysAndValues...)
	return n
}
