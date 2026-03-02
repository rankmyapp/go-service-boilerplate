package logging

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"
)

type Config struct {
	Level     string
	Format    string
	AddSource bool
}

// New creates a logger writing to stdout with the configured level/format.
func New(cfg Config) (*slog.Logger, error) {
	return NewWithWriter(cfg, os.Stdout)
}

// NewWithWriter creates a logger writing to the provided writer.
func NewWithWriter(cfg Config, w io.Writer) (*slog.Logger, error) {
	level, err := parseLevel(cfg.Level)
	if err != nil {
		return nil, err
	}

	handlerOptions := &slog.HandlerOptions{
		Level:     level,
		AddSource: cfg.AddSource,
	}

	format := strings.ToLower(strings.TrimSpace(cfg.Format))
	var handler slog.Handler
	switch format {
	case "", "json":
		handler = slog.NewJSONHandler(w, handlerOptions)
	case "text":
		handler = slog.NewTextHandler(w, handlerOptions)
	default:
		return nil, fmt.Errorf("unsupported log format %q", cfg.Format)
	}

	return slog.New(handler), nil
}

func parseLevel(raw string) (slog.Level, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "debug":
		return slog.LevelDebug, nil
	case "", "info":
		return slog.LevelInfo, nil
	case "warn":
		return slog.LevelWarn, nil
	case "error":
		return slog.LevelError, nil
	default:
		return 0, fmt.Errorf("unsupported log level %q", raw)
	}
}
