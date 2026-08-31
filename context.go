package stdlog

import "context"

type contextKey int

const (
	contextKLogger contextKey = iota
)

func IntoContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, contextKLogger, logger)
}

func FromContext(ctx context.Context) Logger {
	logger, ok := ctx.Value(contextKLogger).(Logger)
	if !ok {
		return Discard
	}
	return logger
}
