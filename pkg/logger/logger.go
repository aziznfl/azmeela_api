package logger

import (
	"log/slog"
	"os"
)

type LoggerWrapper struct {
	*slog.Logger
}

var Log *LoggerWrapper

func InitLogger() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	handler := slog.NewTextHandler(os.Stdout, opts)
	Log = &LoggerWrapper{slog.New(handler)}
}

func (l *LoggerWrapper) Fatal(msg string, args ...any) {
	l.Error(msg, args...)
	os.Exit(1)
}

func (l *LoggerWrapper) Sync() error {
	// slog doesn't require sync like Zap
	return nil
}
