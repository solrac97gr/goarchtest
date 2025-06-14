package utils

import (
	"crypto/md5"
	"fmt"
	"strings"
)

// StringHelper provides string utility functions
type StringHelper struct{}

// ToUpper converts string to uppercase
func (h *StringHelper) ToUpper(s string) string {
	return strings.ToUpper(s)
}

// HashGenerator provides hashing utilities
type HashGenerator struct{}

// MD5Hash generates MD5 hash of input
func (h *HashGenerator) MD5Hash(input string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(input)))
}

// Logger provides logging utilities
type Logger struct {
	level string
}

// Info logs info message
func (l *Logger) Info(message string) {
	fmt.Printf("INFO: %s\n", message)
}

// ConfigManager manages configuration
type ConfigManager struct {
	config map[string]string
}

// GetConfig retrieves configuration value
func (c *ConfigManager) GetConfig(key string) string {
	return c.config[key]
}
