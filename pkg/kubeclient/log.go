package kubeclient

import (
	"github.com/rs/zerolog"
	"os"
	"path/filepath"
)

// LogEntry defines the structure of a log entry
type LogEntry struct {
	SessionID string
	User      string
	Namespace string
	Pod       string
	Container string
	Command   string
	Timestamp string
}

// logCommand logs a command execution with context in JSON format
func logCommand(logger zerolog.Logger, entry LogEntry) {
	logger.Info().
		Str("timestamp", entry.Timestamp).
		Str("session_id", entry.SessionID).
		Str("user", entry.User).
		Str("namespace", entry.Namespace).
		Str("pod", entry.Pod).
		Str("container", entry.Container).
		Str("command", entry.Command).
		Msg("Command executed")
}

// logNewSession logs the start of a new interactive session
func logNewSession(logger zerolog.Logger, entry LogEntry) {
	logger.Info().
		Timestamp().
		Str("session_id", entry.SessionID).
		Str("user", entry.User).
		Str("namespace", entry.Namespace).
		Str("pod", entry.Pod).
		Str("container", entry.Container).
		Str("timestamp", entry.Timestamp).
		Msg("New interactive session started")
}

// createLogger initializes and returns a zerolog logger for the log file
func createLogger() (zerolog.Logger, error) {
	logFilePath := getLogFilePath()
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return zerolog.Logger{}, err
	}

	// Create a zerolog logger that writes to the log file
	logger := zerolog.New(logFile).With().Timestamp().Logger()
	return logger, nil
}

// getLogFilePath returns the path to the log file
func getLogFilePath() string {
	logFilePath := os.Getenv("KUBE_EXEC_LOG_PATH")
	if logFilePath == "" {
		logFilePath = "/tmp/kube-exec/interactive.log"
	}

	// Ensure the directory exists
	os.MkdirAll(filepath.Dir(logFilePath), 0755)
	return logFilePath
}
