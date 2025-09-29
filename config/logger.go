package config

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *log.Logger

func SetupLogger() {
	// Create logs directory if not exists
	logDir := filepath.Dir(AppConfig.LogPath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatal("Cannot create log directory:", err)
	}

	// Setup log rotation
	logFile := &lumberjack.Logger{
		Filename:   AppConfig.LogPath,
		MaxSize:    10, // MB
		MaxBackups: 3,
		MaxAge:     28, // days
		Compress:   true,
	}

	// Create multi-writer (file + stdout)
	multiWriter := io.MultiWriter(os.Stdout, logFile)

	Logger = log.New(multiWriter, "", log.LstdFlags|log.Lshortfile)
}