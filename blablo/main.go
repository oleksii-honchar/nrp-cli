package blablo

import (
	"log/slog"
	"os"
)

type Logger struct {
	*slog.Logger // Embed the slog.Logger to extend it if needed
}

func NewLogger(prefix string) *Logger {
	opts := PrettyHandlerOptions{
		SlogOpts: slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}
	handler := NewPrettyHandler(prefix, os.Stdout, opts)
	logger := slog.New(handler)

	return &Logger{Logger: logger}
}
