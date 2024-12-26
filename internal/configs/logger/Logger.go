package logger

import (
	"log/slog"
	"os"
)

var logger *slog.Logger

func InitLogger(format string) *slog.Logger {
	var handler slog.Handler

	// Choose the format (e.g., "text" or "json")
	switch format {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, nil) // Use default options
	default: // Default to text format
		handler = slog.NewTextHandler(os.Stdout, nil) // Use default options
	}

	// Initialize the global logger
	logger = slog.New(handler)
	return logger
}
