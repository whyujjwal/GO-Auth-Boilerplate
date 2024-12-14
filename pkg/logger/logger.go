package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// LoggerConfig holds configuration for logging
type LoggerConfig struct {
	LogLevel    zerolog.Level
	LogFilePath string
	MaxSize     int64 // Maximum log file size in MB
	MaxBackups  int   // Number of old log files to keep
	MaxAge      int   // Maximum number of days to retain old log files
}

// InitFileLogger sets up file-based logging with multiple configuration options
func InitFileLogger(config LoggerConfig) (zerolog.Logger, error) {
	// Ensure log directory exists
	logDir := filepath.Dir(config.LogFilePath)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return log.Logger, fmt.Errorf("failed to create log directory: %w", err)
	}

	// Open log file with append mode and proper permissions
	logFile, err := os.OpenFile(
		config.LogFilePath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return log.Logger, fmt.Errorf("failed to open log file: %w", err)
	}

	// Create a multi-writer to log to both file and console
	multiWriter := io.MultiWriter(
		logFile,
		zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339},
	)

	// Create logger
	logger := zerolog.New(multiWriter).
		Level(config.LogLevel).
		With().
		Timestamp().
		Logger()

	return logger, nil
}

// RotateLogFile manages log file rotation
func RotateLogFile(config LoggerConfig) error {
	// Check file size
	fileInfo, err := os.Stat(config.LogFilePath)
	if err != nil {
		return err
	}

	// If file size exceeds max size, rotate
	if fileInfo.Size() > config.MaxSize*1024*1024 {
		// Generate backup filename with timestamp
		backupFilename := fmt.Sprintf(
			"%s.%s.bak",
			config.LogFilePath,
			time.Now().Format("2006-01-02-15-04-05"),
		)

		// Rename current log file to backup
		if err := os.Rename(config.LogFilePath, backupFilename); err != nil {
			return fmt.Errorf("failed to rotate log file: %w", err)
		}

		// Clean up old backup files if needed
		if err := cleanupOldLogFiles(config); err != nil {
			log.Error().Err(err).Msg("Failed to cleanup old log files")
		}
	}

	return nil
}

// cleanupOldLogFiles removes old log backup files
func cleanupOldLogFiles(config LoggerConfig) error {
	// Get directory and base filename
	logDir := filepath.Dir(config.LogFilePath)
	baseFileName := filepath.Base(config.LogFilePath)

	// Read directory contents
	files, err := os.ReadDir(logDir)
	if err != nil {
		return err
	}

	// Track backup files
	var backupFiles []os.DirEntry

	// Find backup files
	for _, file := range files {
		if matched, _ := filepath.Match(baseFileName+".*", file.Name()); matched {
			backupFiles = append(backupFiles, file)
		}
	}

	// Sort and remove excess backup files
	if len(backupFiles) > config.MaxBackups {
		// Sort files by modification time (oldest first)
		// Note: This is a simplified implementation
		for i := 0; i < len(backupFiles)-config.MaxBackups; i++ {
			fullPath := filepath.Join(logDir, backupFiles[i].Name())
			fileInfo, err := backupFiles[i].Info()
			if err != nil {
				continue
			}

			// Check file age
			if time.Since(fileInfo.ModTime()) > time.Duration(config.MaxAge)*24*time.Hour {
				os.Remove(fullPath)
			}
		}
	}

	return nil
}

// Example usage
func ExampleLogging() {
	// Configure logging
	config := LoggerConfig{
		LogLevel:    zerolog.InfoLevel,
		LogFilePath: "./logs/app.log",
		MaxSize:     10, // 10 MB
		MaxBackups:  5,  // Keep 5 backup files
		MaxAge:      30, // Keep files for 30 days
	}

	// Initialize logger
	logger, err := InitFileLogger(config)
	if err != nil {
		panic(err)
	}

	// Use logger
	logger.Info().
		Str("component", "main").
		Msg("Application started")

	// Periodic log rotation (could be in a separate goroutine)
	if err := RotateLogFile(config); err != nil {
		logger.Error().Err(err).Msg("Log rotation failed")
	}
}
