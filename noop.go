package stdlog

type noop struct{}

func (n noop) Named(name string) Logger[any] {
	return n
}

func (n noop) SetLevel(level Level) {}

func (n noop) Debug(msg string, keysAndValues ...any) {

}

func (n noop) Info(msg string, keysAndValues ...any) {}

func (n noop) Warning(msg string, keysAndValues ...any) {}

func (n noop) Error(err error, msg string, keysAndValues ...any) {}

func (n noop) Fatal(msg string, keysAndValues ...any) {}

func (n noop) FatalError(err error, msg string, keysAndValues ...any) {}