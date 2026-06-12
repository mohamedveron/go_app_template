package logging

import (
	"fmt"
	"log/slog"
	"strings"
	"time"
)

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorBold   = "\033[1m"
	ColorGray   = "\033[90m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
)

// LogLevel represents a log level for pretty formatting
type LogLevel int

const (
	LevelDebug LogLevel = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelSetup
)

// LogEntry represents a log entry for pretty formatting
type LogEntry struct {
	Timestamp time.Time
	Level     LogLevel
	Message   string
	Attrs     map[string]interface{}
}

// PrettyFormatter handles pretty formatting for logs
type PrettyFormatter struct{}

// NewPrettyFormatter creates a new pretty formatter
func NewPrettyFormatter() *PrettyFormatter {
	return &PrettyFormatter{}
}

// Format formats a log entry with pretty colors and structure
func (f *PrettyFormatter) Format(entry LogEntry) string {
	// Format timestamp
	timestamp := entry.Timestamp.Format("15:04:05")

	// Get level color and text
	levelColor, levelText := f.getLevelColorAndText(entry.Level)

	// Main log line with colors
	mainLine := fmt.Sprintf("%s[%s]%s %s%s%s%s %s%s%s\n",
		ColorGray, timestamp, ColorReset,
		levelColor, ColorBold, levelText, ColorReset,
		ColorBold, entry.Message, ColorReset)

	// If no attributes, return just the main line
	if len(entry.Attrs) == 0 {
		return mainLine
	}

	// Calculate indent (same length as "[15:04:05] LEVEL ")
	indent := strings.Repeat(" ", len(timestamp)+len(levelText)+4)

	// Add attributes on separate lines with extra indent
	var result strings.Builder
	result.WriteString(mainLine)

	for key, value := range entry.Attrs {
		attrLine := fmt.Sprintf("%s  %s%s=%v%s\n",
			indent,
			ColorGray, key, value, ColorReset)
		result.WriteString(attrLine)
	}

	return result.String()
}

// getLevelColorAndText returns the color and text for a log level
func (f *PrettyFormatter) getLevelColorAndText(level LogLevel) (string, string) {
	switch level {
	case LevelDebug:
		return ColorGray, "DEBUG"
	case LevelInfo:
		return ColorGreen, "INFO"
	case LevelWarn:
		return ColorYellow, "WARN"
	case LevelError:
		return ColorRed, "ERROR"
	case LevelSetup:
		return ColorBlue, "SETUP"
	default:
		return ColorReset, "UNKNOWN"
	}
}

// SlogLevelToLogLevel converts slog.Level to LogLevel
func SlogLevelToLogLevel(level slog.Level) LogLevel {
	switch level {
	case slog.LevelDebug:
		return LevelDebug
	case slog.LevelInfo:
		return LevelInfo
	case slog.LevelWarn:
		return LevelWarn
	case slog.LevelError:
		return LevelError
	default:
		return LevelInfo
	}
}
