package utils

import (
	"maps"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// setupLogger configures the global logger based on environment
func SetupLogger() *logrus.Logger {
	logger := logrus.New()

	// Set log level from environment
	logLevel := strings.ToLower(os.Getenv("LOG_LEVEL"))
	switch logLevel {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	// Set formatter based on environment
	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "release" {
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z",
		})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			ForceColors:     true,
		})
	}

	// Set output (optional - could write to file)
	// logFile := os.Getenv("LOG_FILE")
	// if logFile != "" {
	//     file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	//     if err == nil {
	//         logger.SetOutput(file)
	//     }
	// }

	return logger
}

// LogWithRequest logs HTTP requests with structured data
func LogWithRequest(logger *logrus.Logger, method, path, statusCode, latency, clientIP, userAgent string) {
	logger.WithFields(logrus.Fields{
		"method":      method,
		"path":        path,
		"status_code": statusCode,
		"latency":     latency,
		"client_ip":   clientIP,
		"user_agent":  userAgent,
		"service":     "swipelearn-api",
	}).Info("HTTP Request")
}

// LogError logs errors with context
func LogError(logger *logrus.Logger, err error, operation string, context map[string]any) {
	fields := logrus.Fields{
		"error":     err.Error(),
		"operation": operation,
		"service":   "swipelearn-api",
	}

	// Add context fields
	maps.Copy(fields, context)

	logger.WithFields(fields).Error("Operation failed")
}

// LogInfo logs informational messages with context
func LogInfo(logger *logrus.Logger, message string, context map[string]interface{}) {
	fields := logrus.Fields{
		"message": message,
		"service": "swipelearn-api",
	}

	maps.Copy(fields, context)

	logger.WithFields(fields).Info("Info")
}

// LogDebug logs debug messages with context
func LogDebug(logger *logrus.Logger, message string, context map[string]interface{}) {
	fields := logrus.Fields{
		"message": message,
		"service": "swipelearn-api",
	}

	maps.Copy(fields, context)

	logger.WithFields(fields).Debug("Debug")
}
