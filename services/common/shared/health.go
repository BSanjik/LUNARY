package shared

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
)

// Статусы health check
const (
	StatusHealthy   = "healthy"
	StatusUnhealthy = "unhealthy"
	StatusDegraded  = "degraded"
)

// HealthStatus представляет статус одной зависимости
type HealthStatus struct {
	Status      string                 `json:"status"`
	Message     string                 `json:"message,omitempty"`
	LastChecked time.Time              `json:"last_checked"`
	Duration    time.Duration          `json:"duration_ms"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// HealthResponse представляет полный ответ health check
type HealthResponse struct {
	Status     string                  `json:"status"`
	Timestamp  time.Time               `json:"timestamp"`
	Service    ServiceInfo             `json:"service"`
	Checks     map[string]HealthStatus `json:"checks"`
	SystemInfo SystemInfo              `json:"system_info"`
}

// ServiceInfo содержит информацию о сервисе
type ServiceInfo struct {
	Name      string    `json:"name"`
	Version   string    `json:"version"`
	Uptime    string    `json:"uptime"`
	StartTime time.Time `json:"start_time"`
}

// SystemInfo содержит системную информацию
type SystemInfo struct {
	GoVersion     string `json:"go_version"`
	NumGoroutines int    `json:"num_goroutines"`
	MemAlloc      uint64 `json:"mem_alloc_mb"`
	NumGC         uint32 `json:"num_gc"`
}

// HealthChecker интерфейс для проверки здоровья
type HealthChecker interface {
	Check() HealthStatus
	Name() string
}

// DatabaseHealthChecker проверяет подключение к БД
type DatabaseHealthChecker struct {
	DB *sql.DB
}

func NewDatabaseHealthChecker(db *sql.DB) *DatabaseHealthChecker {
	return &DatabaseHealthChecker{DB: db}
}

func (d *DatabaseHealthChecker) Name() string {
	return "database"
}

func (d *DatabaseHealthChecker) Check() HealthStatus {
	start := time.Now()

	status := HealthStatus{
		LastChecked: start,
		Metadata:    make(map[string]interface{}),
	}

	// Проверяем ping к БД
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := d.DB.PingContext(ctx)
	status.Duration = time.Since(start)

	if err != nil {
		status.Status = StatusUnhealthy
		status.Message = "Database connection failed"
		status.Metadata["error"] = err.Error()
		return status
	}

	// Проверяем статистику соединений
	stats := d.DB.Stats()
	status.Status = StatusHealthy
	status.Message = "Database connection is healthy"
	status.Metadata["open_connections"] = stats.OpenConnections
	status.Metadata["in_use"] = stats.InUse
	status.Metadata["idle"] = stats.Idle

	// Проверяем если слишком много открытых соединений
	if stats.OpenConnections > 50 {
		status.Status = StatusDegraded
		status.Message = "High number of database connections"
	}

	return status
}

// HealthManager управляет всеми health checks
type HealthManager struct {
	serviceName string
	version     string
	startTime   time.Time
	checkers    map[string]HealthChecker
}

func NewHealthManager(serviceName, version string) *HealthManager {
	return &HealthManager{
		serviceName: serviceName,
		version:     version,
		startTime:   time.Now(),
		checkers:    make(map[string]HealthChecker),
	}
}

func (h *HealthManager) AddChecker(checker HealthChecker) {
	h.checkers[checker.Name()] = checker
}

func (h *HealthManager) RunChecks() HealthResponse {
	checks := make(map[string]HealthStatus)
	overallStatus := StatusHealthy

	// Выполняем все проверки
	for name, checker := range h.checkers {
		status := checker.Check()
		checks[name] = status

		// Определяем общий статус
		if status.Status == StatusUnhealthy {
			overallStatus = StatusUnhealthy
		} else if status.Status == StatusDegraded && overallStatus == StatusHealthy {
			overallStatus = StatusDegraded
		}
	}

	// Собираем системную информацию
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	systemInfo := SystemInfo{
		GoVersion:     runtime.Version(),
		NumGoroutines: runtime.NumGoroutine(),
		MemAlloc:      m.Alloc / 1024 / 1024, // В мегабайтах
		NumGC:         m.NumGC,
	}

	// Формируем ответ
	response := HealthResponse{
		Status:    overallStatus,
		Timestamp: time.Now(),
		Service: ServiceInfo{
			Name:      h.serviceName,
			Version:   h.version,
			Uptime:    time.Since(h.startTime).String(),
			StartTime: h.startTime,
		},
		Checks:     checks,
		SystemInfo: systemInfo,
	}

	return response
}

// HTTP handlers
func (h *HealthManager) LivenessHandler(w http.ResponseWriter, r *http.Request) {
	// Простая проверка что сервис запущен
	response := map[string]interface{}{
		"status":    StatusHealthy,
		"timestamp": time.Now(),
		"service":   h.serviceName,
		"message":   "Service is alive",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

	LogDebug("Liveness check completed", logrus.Fields{
		"service": h.serviceName,
	})
}

func (h *HealthManager) ReadinessHandler(w http.ResponseWriter, r *http.Request) {
	// Полная проверка готовности (БД, внешние сервисы)
	healthResponse := h.RunChecks()

	status := http.StatusOK
	if healthResponse.Status == StatusUnhealthy {
		status = http.StatusServiceUnavailable
	} else if healthResponse.Status == StatusDegraded {
		status = http.StatusOK // Или можно использовать 207 Multi-Status
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(healthResponse)

	LogInfo("Readiness check completed", logrus.Fields{
		"service": h.serviceName,
		"status":  healthResponse.Status,
		"checks":  len(healthResponse.Checks),
	})
}

func (h *HealthManager) HealthHandler(w http.ResponseWriter, r *http.Request) {
	// Детальная информация о здоровье сервиса
	healthResponse := h.RunChecks()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(healthResponse)

	LogDebug("Health check completed", logrus.Fields{
		"service": h.serviceName,
		"status":  healthResponse.Status,
	})
}
