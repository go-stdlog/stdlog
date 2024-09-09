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

type Logger[T any] interface {
	// Named returns a new logger with its previous name followed by a dot,
	// followed by the provided name.
	Named(name string) Logger[T]
	SetLevel(level Level)

	Debug(msg string, keysAndValues ...T)
	Info(msg string, keysAndValues ...T)
	Warning(msg string, keysAndValues ...T)
	Error(err error, msg string, keysAndValues ...T)
	Fatal(msg string, keysAndValues ...T)
	FatalError(err error, msg string, keysAndValues ...T)
}
