package logger

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

func SetupLogger() {
	wd, _ := os.Getwd()
	Logfile_Path := filepath.Join(wd, "internal", "logger", "Logs", "App_Logs.log")
	file, err := os.OpenFile(Logfile_Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("Error to Open the File" + err.Error())
	}
	multiWriter := io.MultiWriter(os.Stdout, file)

	handler := slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger := slog.New(handler)

	slog.SetDefault(logger)
}
