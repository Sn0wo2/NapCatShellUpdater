package log

import (
	"encoding/json"
	"fmt"
	"github.com/Sn0wo2/NapCatShellUpdater/helper"
	"github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"io"
	"os"
	"runtime"
	"strings"
)

// Logger is an instance of the logrus logger.
var Logger = logrus.New()

// InitLogger initializes the logger.
// If logPath is empty, logs will only be output to stdout.
// Returns an error if there is a problem opening the log file.
func InitLogger(logPath string, formatter *easy.Formatter, logLevel logrus.Level) error {
	var writers []io.Writer
	writers = append(writers, os.Stdout)

	if logPath != "" {
		file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return fmt.Errorf("failed to open log file %s: %w", logPath, err)
		}
		writers = append(writers, file)
	}

	Logger.SetOutput(io.MultiWriter(writers...))
	Logger.SetLevel(logLevel)
	Logger.SetFormatter(formatter)
	Logger.SetReportCaller(true)

	return nil
}

// Log a message with a prefix at the specified level.
func logWithPrefix(level logrus.Level, prefix string, args ...interface{}) {
	entry := Logger.WithField("prefix", prefix)
	switch level {
	case logrus.TraceLevel:
		entry.Trace(args...)
	case logrus.DebugLevel:
		entry.Debug(args...)
	case logrus.InfoLevel:
		entry.Info(args...)
	case logrus.WarnLevel:
		entry.Warn(args...)
	case logrus.ErrorLevel:
		entry.Error(args...)
	case logrus.FatalLevel:
		entry.Fatal(args...)
	case logrus.PanicLevel:
		entry.Panic(args...)
	}
}

// Trace logs a message at level Trace with a prefix.
func Trace(prefix string, args ...interface{}) {
	logWithPrefix(logrus.TraceLevel, prefix, args...)
}

// Debug logs a message at level Debug with a prefix.
func Debug(prefix string, args ...interface{}) {
	logWithPrefix(logrus.DebugLevel, prefix, args...)
}

// Info logs a message at level Info with a prefix.
func Info(prefix string, args ...interface{}) {
	logWithPrefix(logrus.InfoLevel, prefix, args...)
}

// Warning logs a message at level Warning with a prefix.
func Warning(prefix string, args ...interface{}) {
	logWithPrefix(logrus.WarnLevel, prefix, args...)
}

// Warn logs a message at level Warn with a prefix.
func Warn(prefix string, args ...interface{}) {
	logWithPrefix(logrus.WarnLevel, prefix, args...)
}

// Error logs a message at level Error with a prefix.
func Error(prefix string, args ...interface{}) {
	logWithPrefix(logrus.ErrorLevel, prefix, args...)
}

// Panic logs a message at level Panic with a prefix.
func Panic(prefix string, args ...interface{}) {
	logWithPrefix(logrus.PanicLevel, prefix, args...)
}

// Fatal logs a message at level Fatal with a prefix.
func Fatal(prefix string, args ...interface{}) {
	logWithPrefix(logrus.FatalLevel, prefix, args...)
}

// RPanic logs a message at level Error with stack trace, without exiting the program.
func RPanic(data ...interface{}) {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	if n < len(buf) {
		buf = buf[:n]
	}

	Error("panic", TrimJSONArray(fmt.Sprint(data...)), "\n", helper.BytesToString(buf))
}

// TrimJSONArray trims the JSON array from the string for print
func TrimJSONArray(data string) string {
	return strings.TrimSuffix(strings.TrimPrefix(data, "["), "]")
}

// FormatInterfaceToJSON formats an object into a JSON string.
// Returns an empty string if an error occurs.
func FormatInterfaceToJSON(data ...interface{}) string {
	if len(data) == 1 {
		bytes, err := json.Marshal(data[0])
		if err != nil {
			RPanic("failed to format JSON", err)
			return ""
		}
		return helper.BytesToString(bytes)
	} else {
		bytes, err := json.Marshal(data)
		if err != nil {
			RPanic("failed to format JSON", err)
			return ""
		}
		return helper.BytesToString(bytes)
	}
}

// DefaultFormatter returns a default easy formatter.
func DefaultFormatter() *easy.Formatter {
	return &easy.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		LogFormat:       "%time% [%lvl%] (%prefix%): %msg%\n",
	}
}
