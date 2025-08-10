package shared

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

// GracefulShutdown управляет корректным завершением сервиса
type GracefulShutdown struct {
	server          *http.Server
	shutdownTimeout time.Duration
	cleanup         []CleanupFunc
	mu              sync.Mutex
}

// CleanupFunc функция очистки ресурсов
type CleanupFunc func(ctx context.Context) error

// NewGracefulShutdown создает новый менеджер graceful shutdown
func NewGracefulShutdown(server *http.Server, timeout time.Duration) *GracefulShutdown {
	return &GracefulShutdown{
		server:          server,
		shutdownTimeout: timeout,
		cleanup:         make([]CleanupFunc, 0),
	}
}

// AddCleanup добавляет функцию очистки
func (g *GracefulShutdown) AddCleanup(fn CleanupFunc) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.cleanup = append(g.cleanup, fn)
}

// Start запускает сервер и ожидает сигналы завершения
func (g *GracefulShutdown) Start() error {
	// Канал для получения сигналов завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Канал для ошибок сервера
	serverErrors := make(chan error, 1)

	// Запускаем сервер в горутине
	go func() {
		LogInfo("Starting HTTP server", logrus.Fields{
			"address": g.server.Addr,
		})
		
		if err := g.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrors <- err
		}
	}()

	// Ожидаем сигнал завершения или ошибку сервера
	select {
	case err := <-serverErrors:
		LogError(err, "Server failed to start")
		return err
	case sig := <-quit:
		LogInfo("Received shutdown signal", logrus.Fields{
			"signal": sig.String(),
		})
		return g.shutdown()
	}
}

// shutdown выполняет корректное завершение
func (g *GracefulShutdown) shutdown() error {
	LogInfo("Starting graceful shutdown", logrus.Fields{
		"timeout": g.shutdownTimeout.String(),
	})

	// Создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), g.shutdownTimeout)
	defer cancel()

	// Останавливаем HTTP сервер
	if err := g.server.Shutdown(ctx); err != nil {
		LogError(err, "Server shutdown failed")
		return err
	}

	LogInfo("HTTP server stopped gracefully")

	// Выполняем функции очистки
	g.runCleanup(ctx)

	LogInfo("Graceful shutdown completed successfully")
	return nil
}

// runCleanup выполняет все функции очистки
func (g *GracefulShutdown) runCleanup(ctx context.Context) {
	g.mu.Lock()
	cleanupFuncs := make([]CleanupFunc, len(g.cleanup))
	copy(cleanupFuncs, g.cleanup)
	g.mu.Unlock()

	var wg sync.WaitGroup
	
	for i, cleanupFunc := range cleanupFuncs {
		wg.Add(1)
		go func(index int, fn CleanupFunc) {
			defer wg.Done()
			
			start := time.Now()
			if err := fn(ctx); err != nil {
				LogError(err, "Cleanup function failed", logrus.Fields{
					"cleanup_index": index,
					"duration_ms":   time.Since(start).Milliseconds(),
				})
			} else {
				LogInfo("Cleanup function completed", logrus.Fields{
					"cleanup_index": index,
					"duration_ms":   time.Since(start).Milliseconds(),
				})
			}
		}(i, cleanupFunc)
	}

	// Ожидаем завершения всех функций очистки
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		LogInfo("All cleanup functions completed")
	case <-ctx.Done():
		LogWarn("Cleanup timeout exceeded, forcing shutdown")
	}
}

// DatabaseCleanup создает функцию очистки для базы данных
func DatabaseCleanup(db interface{ Close() error }, name string) CleanupFunc {
	return func(ctx context.Context) error {
		LogInfo("Closing database connection", logrus.Fields{
			"database": name,
		})
		
		start := time.Now()
		err := db.Close()
		
		LogDBOperation("close", name, time.Since(start).Milliseconds(), 0, err)
		
		if err != nil {
			return err
		}
		
		LogInfo("Database connection closed successfully", logrus.Fields{
			"database": name,
		})
		
		return nil
	}
}

// ContextCleanup создает функцию очистки с пользовательской логикой
func ContextCleanup(name string, fn func(ctx context.Context) error) CleanupFunc {
	return func(ctx context.Context) error {
		LogInfo("Running custom cleanup", logrus.Fields{
			"name": name,
		})
		
		start := time.Now()
		err := fn(ctx)
		
		if err != nil {
			LogError(err, "Custom cleanup failed", logrus.Fields{
				"name":        name,
				"duration_ms": time.Since(start).Milliseconds(),
			})
		} else {
			LogInfo("Custom cleanup completed", logrus.Fields{
				"name":        name,
				"duration_ms": time.Since(start).Milliseconds(),
			})
		}
		
		return err
	}
}

// WaitForActiveConnections ожидает завершения активных соединений
func WaitForActiveConnections(server *http.Server, maxWait time.Duration) CleanupFunc {
	return func(ctx context.Context) error {
		LogInfo("Waiting for active connections to finish", logrus.Fields{
			"max_wait": maxWait.String(),
		})
		
		// Создаем контекст с максимальным временем ожидания
		waitCtx, cancel := context.WithTimeout(ctx, maxWait)
		defer cancel()
		
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		
		for {
			select {
			case <-waitCtx.Done():
				LogWarn("Timeout waiting for connections to finish")
				return waitCtx.Err()
			case <-ticker.C:
				// В реальном приложении здесь можно проверить активные соединения
				// Для примера просто ждем немного
				LogDebug("Checking active connections...")
				return nil
			}
		}
	}
}
