package logger

import (
	"github.com/sirupsen/logrus"
)

// Fields defines a type alias for logrus.Fields.
type Fields logrus.Fields

// Logger wraps logrus.Entry to provide context-aware logging.
type Logger struct {
	entry *logrus.Entry
}

// New creates a new Logger instance with optional default fields.
func New(fields Fields) *Logger {
	return &Logger{
		entry: logrus.WithFields(logrus.Fields(fields)),
	}
}

// WithFields adds additional fields to the logger context.
func (l *Logger) WithFields(fields Fields) *Logger {
	return &Logger{
		entry: l.entry.WithFields(logrus.Fields(fields)),
	}
}

// WithError adds an error to the logger context.
func (l *Logger) WithError(err error) *Logger {
	return &Logger{
		entry: l.entry.WithError(err),
	}
}

// Info logs an informational message.
func (l *Logger) Info(msg string) {
	l.entry.Info(msg)
}

// Warn logs a warning message.
func (l *Logger) Warn(msg string) {
	l.entry.Warn(msg)
}

// Error logs an error message.
func (l *Logger) Error(msg string) {
	l.entry.Error(msg)
}

// Debug logs a debug message.
func (l *Logger) Debug(msg string) {
	l.entry.Debug(msg)
}

// Panic logs a panic message and panics.
func (l *Logger) Panic(msg string) {
	l.entry.Panic(msg)
}

// Fatal logs a fatal message and exits the application.
func (l *Logger) Fatal(msg string) {
	l.entry.Fatal(msg)
}

// Infof logs a formatted informational message.
func (l *Logger) Infof(format string, args ...interface{}) {
	l.entry.Infof(format, args...)
}

// Errorf logs a formatted error message.
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.entry.Errorf(format, args...)
}

// Configure sets up global logging options (e.g., format, level).
func Configure(level logrus.Level, format logrus.Formatter) {
	logrus.SetLevel(level)
	if format != nil {
		logrus.SetFormatter(format)
	}
}
