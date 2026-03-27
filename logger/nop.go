package logger

import (
	"io"
	"log/slog"
	"os"
)

// NewNop returns a Logger that discards all output.
// Intended for unit tests where log output is irrelevant.
func NewNop() *Logger {
	h := slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	return &Logger{
		handler:     h,
		serviceName: "nop",
		level:       LevelDebug,
		exitFunc:    os.Exit,
	}
}
