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
		return "WARNING"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

type Logger interface {
	// Named returns a new logger with its previous name followed by a dot,
	// followed by the provided name.
	Named(name string) Logger
	SetLevel(level Level)

	Debug(msg string, keysAndValues ...any)
	Info(msg string, keysAndValues ...any)
	Warning(msg string, keysAndValues ...any)
	Error(err error, msg string, keysAndValues ...any)
	Fatal(msg string, keysAndValues ...any)
	FatalError(err error, msg string, keysAndValues ...any)
}
