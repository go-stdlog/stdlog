package stdlog

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarning
	LevelError
	LevelFatal
)

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarning:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// FatalBehavior determines how Logger behaves when Logger.Fatal or
// Logger.FatalError is called. See Logger.SetFatalBehavior
type FatalBehavior uint8

const (
	// FatalExits makes Logger invoke os.Exit with status 1 after printing a
	// message with a Fatal log level. This is the default mode.
	FatalExits FatalBehavior = iota

	// FatalPanics makes Logger emit a panic with the same message as provided
	// to the logger method. This does not invoke os.Exit.
	FatalPanics
)

type Logger interface {
	// Named returns a new logger with its previous name followed by a dot,
	// followed by the provided name.
	Named(name string) Logger

	// SetLevel sets the logger level.
	SetLevel(level Level)

	// Leveled returns a copy of the current logger with a different level.
	Leveled(level Level) Logger

	// WithFields returns a copy of the current logger with extra provided
	// fields that will be present in every emitted log message.
	WithFields(keysAndValues ...any) Logger

	// Skipping returns a copy of the current logger that will skip a given
	// number of call stacks.
	Skipping(count uint) Logger

	// SetFatalBehavior defines how Logger behaves when Fatal or FatalError is
	// called. By default, Logger uses FatalExits, which means that any call to
	// a fatal logger level calls os.Exit with status 1. The other option,
	// FatalPanics, panics after printing the log message to the configured
	// output, allowing deferred functions to be executed.
	SetFatalBehavior(behavior FatalBehavior)

	Debug(msg string, keysAndValues ...any)
	Info(msg string, keysAndValues ...any)
	Warning(msg string, keysAndValues ...any)
	Error(err error, msg string, keysAndValues ...any)
	Fatal(msg string, keysAndValues ...any)
	FatalError(err error, msg string, keysAndValues ...any)
}
