package logger

import (
	"log/slog"
	"os"

	"github.com/Marlliton/slogpretty"
)

func Init(logLevel int) {
	levels := []slog.Level{
		slog.LevelError,
		slog.LevelWarn,
		slog.LevelInfo,
		slog.LevelDebug,
	}
	level := slog.LevelDebug
	if logLevel >= 0 && logLevel < len(levels) {
		level = levels[logLevel]
	}

	logOpts := &slogpretty.Options{
		Level:      level,
		AddSource:  true,
		Colorful:   true,
		Multiline:  true,
		TimeFormat: "15:04:05",
	}
	logHandler := slogpretty.New(os.Stdout, logOpts)
	slog.SetDefault(slog.New(logHandler))
	slog.Debug("Setup logging", slog.Int("Log Level", logLevel))
}
