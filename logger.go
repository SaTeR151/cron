package cron

import "log/slog"

type Logger interface {
	Info(msg string, keysAndValues ...any)
	Error(err error, msg string, keysAndValues ...any)
}

type cronLogger struct {
	logger *slog.Logger
}

// NewLogger возвращает логер для cron.Option
func NewLogger(base *slog.Logger) Logger {
	return &cronLogger{
		logger: base,
	}
}

func (l *cronLogger) Info(msg string, keysAndValues ...any) {
	l.logger.Info(msg, keysAndValues...)
}

func (l *cronLogger) Error(err error, msg string, keysAndValues ...any) {
	args := make([]any, 0, len(keysAndValues)+2)
	args = append(args, "error", err)
	args = append(args, keysAndValues...)

	l.logger.Error(msg, args...)
}
