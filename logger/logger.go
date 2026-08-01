// Package logger provides a simple logging utility for the Power Query Converter.
// It allows logging messages with different severity levels (Info, Warning, Error, Debug)
// and stores the logs in memory for later retrieval.
package logger

import (
	"fmt"
	"os"
	"time"
)

// Logger is a simple logging utility that allows logging messages with different severity levels.
type Logger struct {
	prefix      string
	logs        []string
	isDebugging bool
}

// NewLogger creates a new Logger instance with the specified prefix.
// The prefix is used to identify the source of the log messages.
func NewLogger(prefix string, isDebugging bool) *Logger {
	return &Logger{
		prefix:      prefix,
		isDebugging: isDebugging,
		logs:        []string{},
	}
}

// Info logs an informational message to stdout.
func (l *Logger) Info(message string) {
	msg := l.generateLogMessage(message, "INFO")
	l.logs = append(l.logs, msg)
	fmt.Println(msg)
}

// Warning logs a warning message to stderr.
func (l *Logger) Warning(message string) {
	msg := l.generateLogMessage(message, "WARNING")
	l.logs = append(l.logs, msg)
	fmt.Fprintln(os.Stderr, msg)
}

// Error logs an error message to stderr.
func (l *Logger) Error(message string) {
	msg := l.generateLogMessage(message, "ERROR")
	l.logs = append(l.logs, msg)
	fmt.Fprintln(os.Stderr, msg)
}

// Debug logs a debug message to stderr when debugging is enabled.
func (l *Logger) Debug(message string) {
	msg := l.generateLogMessage(message, "DEBUG")
	l.logs = append(l.logs, msg)
	if l.isDebugging {
		fmt.Fprintln(os.Stderr, msg)
	}
}

// GetLogs returns all the logs stored in the Logger instance.
// This can be useful for testing or debugging purposes.
func (l *Logger) GetLogs() []string {
	return l.logs
}

// SetDebugging enables or disables debug logging.
// When debugging is enabled, debug messages will be printed to stderr.
func (l *Logger) SetDebugging(isDebugging bool) {
	l.isDebugging = isDebugging
}

// generateLogMessage creates a formatted log message with the current time, log level, prefix, and the actual message.
// The format is: "YYYY-MM-DD HH:MM:SS [LEVEL][PREFIX]: message".
func (l *Logger) generateLogMessage(message string, level string) string {
	return fmt.Sprintf("%s [%s][%s]: %s", l.getCurrentTime(), level, l.prefix, message)
}

// getCurrentTime returns the current time formatted as "YYYY-MM-DD HH:MM:SS".
// This is used to timestamp the log messages.
// The format is defined by the time package in Go.
func (l *Logger) getCurrentTime() string {
	timeFormat := "2006-01-02 15:04:05"
	currentTime := time.Now()
	return currentTime.Format(timeFormat)
}
