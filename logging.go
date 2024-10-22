package main

import (
	"log"
	"log/slog"
	"os"
)

var loggerLevels = map[string]slog.Level{
	"info":  slog.LevelInfo,
	"debug": slog.LevelDebug,
	"error": slog.LevelError,
	"warn":  slog.LevelWarn,
}

type sloggerConfig struct {
	level slog.Level
	style string
}

func newSlogConfig(loggerLevel, loggerStyle string) sloggerConfig {
	logLevel, ok := loggerLevels[loggerLevel]
	if !ok {
		log.Fatalf("invalid log level provided: %s", logLevel)
	}
	return sloggerConfig{level: logLevel, style: loggerStyle}
}

func (s *sloggerConfig) Create() *slog.Logger {
	var handler slog.Handler
	switch s.style {
	case "text":
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: s.level})
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: s.level})
	case "dev":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: s.level, AddSource: true})
	default:
		log.Fatalf("invalid log style provided: %s", s.style)
	}

	return slog.New(handler)
}
