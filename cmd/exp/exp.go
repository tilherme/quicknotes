package main

import (
	"log/slog"
	"os"
)

func main() {
	// log.Println("Log padrao")
	// h := slog.NewTextHandler(os.Stdout, nil)
	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	})

	log := slog.New(h).With("app", "exp")
	log.Info("info message", "id", 1)
	log.Debug("debug message")
	log.Warn("warn message")
	log.Error("warn message")

}
