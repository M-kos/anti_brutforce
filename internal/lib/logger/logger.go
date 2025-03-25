package logger

import (
	"log/slog"
	"os"
)

type Logger struct {
	log *slog.Logger
}

func NewLogger() *Logger {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	return &Logger{
		log: logger,
	}
}

func (l *Logger) Info(msg string, args ...any) {
	l.log.Info(msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.log.Error(msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.log.Warn(msg, args...)
}
