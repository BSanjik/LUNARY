// Общий логгер
package shared

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger

// Ключи для контекста
type ContextKey string

const (
	RequestIDKey ContextKey = "request_id"
	UserIDKey    ContextKey = "user_id"
	ServiceKey   ContextKey = "service"
)

// Инициализация логгера
func InitLogger(serviceName string, level string) {
	Logger = logrus.New()

	// JSON формат для продакшена
	Logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
		},
	})

	// Устанавливаем уровень логирования
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	Logger.SetLevel(logLevel)

	// Вывод в stdout
	Logger.SetOutput(os.Stdout)

	// Добавляем глобальные поля
	Logger = Logger.WithFields(logrus.Fields{
		"service": serviceName,
	}).Logger
}

// Создание логгера с контекстом
func WithContext(ctx context.Context) *logrus.Entry {
	entry := Logger.WithFields(logrus.Fields{})

	// Добавляем поля из контекста
	if requestID := ctx.Value(RequestIDKey); requestID != nil {
		entry = entry.WithField("request_id", requestID)
	}

	if userID := ctx.Value(UserIDKey); userID != nil {
		entry = entry.WithField("user_id", userID)
	}

	return entry
}

// Логгер с дополнительными полями
func WithFields(fields logrus.Fields) *logrus.Entry {
	return Logger.WithFields(fields)
}

// WithRequestID добавляет request ID в контекст
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// WithUserID добавляет user ID в контекст
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// WithService добавляет имя сервиса в контекст
func WithService(ctx context.Context, serviceName string) context.Context {
	return context.WithValue(ctx, ServiceKey, serviceName)
}

// GetRequestID извлекает request ID из контекста
func GetRequestID(ctx context.Context) string {
	if requestID := ctx.Value(RequestIDKey); requestID != nil {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}

// GetUserID извлекает user ID из контекста
func GetUserID(ctx context.Context) string {
	if userID := ctx.Value(UserIDKey); userID != nil {
		if id, ok := userID.(string); ok {
			return id
		}
	}
	return ""
}

// Быстрые методы для часто используемых случаев
func LogInfo(msg string, fields ...logrus.Fields) {
	if len(fields) > 0 {
		Logger.WithFields(fields[0]).Info(msg)
	} else {
		Logger.Info(msg)
	}
}

// LogInfoWithContext логирует с автоматическим добавлением полей из контекста
func LogInfoWithContext(ctx context.Context, msg string, fields ...logrus.Fields) {
	entry := WithContext(ctx)
	if len(fields) > 0 {
		entry = entry.WithFields(fields[0])
	}
	entry.Info(msg)
}

// LogErrorWithContext логирует ошибку с контекстом
func LogErrorWithContext(ctx context.Context, err error, msg string, fields ...logrus.Fields) {
	entry := WithContext(ctx)
	if err != nil {
		entry = entry.WithField("error", err.Error())
	}
	if len(fields) > 0 {
		entry = entry.WithFields(fields[0])
	}
	entry.Error(msg)
}

// LogWarnWithContext логирует предупреждение с контекстом
func LogWarnWithContext(ctx context.Context, msg string, fields ...logrus.Fields) {
	entry := WithContext(ctx)
	if len(fields) > 0 {
		entry = entry.WithFields(fields[0])
	}
	entry.Warn(msg)
}

func LogError(err error, msg string, fields ...logrus.Fields) {
	entry := Logger.WithField("error", err.Error())
	if len(fields) > 0 {
		entry = entry.WithFields(fields[0])
	}
	entry.Error(msg)
}

func LogWarn(msg string, fields ...logrus.Fields) {
	if len(fields) > 0 {
		Logger.WithFields(fields[0]).Warn(msg)
	} else {
		Logger.Warn(msg)
	}
}

func LogDebug(msg string, fields ...logrus.Fields) {
	if len(fields) > 0 {
		Logger.WithFields(fields[0]).Debug(msg)
	} else {
		Logger.Debug(msg)
	}
}

// Структурированные логи для HTTP запросов
func LogHTTPRequest(method, path, userAgent, clientIP string, statusCode int, duration int64) {
	Logger.WithFields(logrus.Fields{
		"method":      method,
		"path":        path,
		"status_code": statusCode,
		"duration_ms": duration,
		"user_agent":  userAgent,
		"client_ip":   clientIP,
		"type":        "http_request",
	}).Info("HTTP request processed")
}

// Логи для базы данных
func LogDBOperation(operation, table string, duration int64, rowsAffected int64, err error) {
	fields := logrus.Fields{
		"operation":     operation,
		"table":         table,
		"duration_ms":   duration,
		"rows_affected": rowsAffected,
		"type":          "database",
	}

	if err != nil {
		fields["error"] = err.Error()
		Logger.WithFields(fields).Error("Database operation failed")
	} else {
		Logger.WithFields(fields).Debug("Database operation completed")
	}
}

// Логи для аутентификации
func LogAuthEvent(event, phone string, userID int64, success bool, reason string) {
	fields := logrus.Fields{
		"event":   event, // "login", "register", "logout"
		"phone":   phone,
		"success": success,
		"type":    "auth",
	}

	if userID > 0 {
		fields["user_id"] = userID
	}

	if reason != "" {
		fields["reason"] = reason
	}

	if success {
		Logger.WithFields(fields).Info("Authentication event")
	} else {
		Logger.WithFields(fields).Warn("Authentication failed")
	}
}
