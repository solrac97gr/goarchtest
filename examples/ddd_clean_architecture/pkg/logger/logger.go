package logger

// Logger defines a simple logging interface
type Logger interface {
	Info(message string, args ...interface{})
	Error(message string, args ...interface{})
	Debug(message string, args ...interface{})
	Warn(message string, args ...interface{})
}

// ConsoleLogger implements Logger for console output
type ConsoleLogger struct{}

// NewConsoleLogger creates a new console logger
func NewConsoleLogger() *ConsoleLogger {
	return &ConsoleLogger{}
}

// Info logs an info message
func (l *ConsoleLogger) Info(message string, args ...interface{}) {
	// Implementation would log to console
}

// Error logs an error message
func (l *ConsoleLogger) Error(message string, args ...interface{}) {
	// Implementation would log to console
}

// Debug logs a debug message
func (l *ConsoleLogger) Debug(message string, args ...interface{}) {
	// Implementation would log to console
}

// Warn logs a warning message
func (l *ConsoleLogger) Warn(message string, args ...interface{}) {
	// Implementation would log to console
}
