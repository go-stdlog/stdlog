package stdlog

type noop string

const Discard noop = "noop"

func (n noop) Named(string) Logger { return n }

func (n noop) SetLevel(Level) {}

func (n noop) Leveled(Level) Logger { return n }

func (n noop) Skipping(uint) Logger { return n }

func (n noop) SetFatalBehavior(FatalBehavior) {}

func (n noop) Debug(_ string, keysAndValues ...any) {
	assertKvs(LevelDebug.String(), keysAndValues...)
}

func (n noop) Info(_ string, keysAndValues ...any) {
	assertKvs(LevelInfo.String(), keysAndValues...)
}

func (n noop) Warning(_ string, keysAndValues ...any) {
	assertKvs(LevelWarning.String(), keysAndValues...)
}

func (n noop) Error(_ error, _ string, keysAndValues ...any) {
	assertKvs(LevelError.String(), keysAndValues...)
}

func (n noop) Fatal(msg string, keysAndValues ...any) {
	assertKvs(LevelFatal.String(), keysAndValues...)
	panic(msg)
}

func (n noop) FatalError(msg error, _ string, keysAndValues ...any) {
	assertKvs(LevelFatal.String()+"_ERROR", keysAndValues...)
	panic(msg)
}

func (n noop) WithFields(keysAndValues ...any) Logger {
	assertKvs("WithFields", keysAndValues...)
	return n
}
