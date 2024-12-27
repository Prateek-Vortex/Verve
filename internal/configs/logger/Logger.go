package logger

import (
	"log/slog"
	"os"
)

var logger *slog.Logger

// InitLogger initializes the logger with the specified format ("text" or "json")
func InitLogger(format string) *slog.Logger {
	var handler slog.Handler

	// Choose the format (e.g., "text" or "json")
	switch format {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo, // Set the log level to info
		})
	default: // Default to text format
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo, // Set the log level to info
		})
	}

	// Initialize the global logger with the chosen handler
	logger = slog.New(handler)

	// Return the logger
	return logger
}
