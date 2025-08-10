# Graceful Shutdown

## Обзор

Graceful shutdown обеспечивает корректное завершение сервиса без потери данных и прерывания активных запросов.

## Компоненты

### 1. GracefulShutdown Manager
- Управляет процессом завершения
- Ожидает сигналы SIGTERM/SIGINT
- Выполняет функции очистки
- Устанавливает таймауты

### 2. Функции очистки
- `DatabaseCleanup` - закрытие соединений с БД
- `ContextCleanup` - пользовательская логика
- `WaitForActiveConnections` - ожидание завершения запросов

## Последовательность shutdown

1. **Получение сигнала** - SIGTERM/SIGINT
2. **Остановка HTTP сервера** - `server.Shutdown()`
3. **Выполнение cleanup функций** - параллельно
4. **Закрытие БД соединений**
5. **Финальное логирование**

## Конфигурация

### Таймауты
- **Shutdown timeout**: 30 секунд
- **Active connections wait**: 10 секунд
- **HTTP timeouts**: Read 15s, Write 15s, Idle 60s

### Docker Compose
```yaml
healthcheck:
  test: ["CMD", "curl", "-f", "http://localhost:8081/health/ready"]
  interval: 30s
  timeout: 10s
  retries: 3
  start_period: 40s
```

### Kubernetes
```yaml
terminationGracePeriodSeconds: 45
lifecycle:
  preStop:
    exec:
      command: ["/bin/sh", "-c", "sleep 5"]
```

## Логирование

Graceful shutdown генерирует структурированные логи:

```json
{
  "level": "info",
  "message": "Received shutdown signal",
  "signal": "terminated"
}
```

```json
{
  "level": "info", 
  "message": "Starting graceful shutdown",
  "timeout": "30s"
}
```

```json
{
  "level": "info",
  "message": "Cleanup function completed",
  "cleanup_index": 0,
  "duration_ms": 150
}
```

## Тестирование

Используйте скрипт для тестирования:

```bash
bash test_graceful_shutdown.sh
```

Проверьте логи на наличие:
- ✅ "Received shutdown signal"
- ✅ "Starting graceful shutdown" 
- ✅ "HTTP server stopped gracefully"
- ✅ "Database connection closed successfully"
- ✅ "Graceful shutdown completed successfully"

## Best Practices

1. **Всегда используйте контекст** с таймаутом
2. **Логируйте каждый этап** shutdown
3. **Тестируйте в продакшене** с реальной нагрузкой
4. **Настройте мониторинг** времени shutdown
5. **Учитывайте зависимости** между сервисами

## Устранение проблем

### Долгий shutdown
- Проверьте активные goroutines
- Увеличьте таймауты
- Проверьте блокирующие операции

### Потеря соединений БД
- Убедитесь что `defer db.Close()` вызывается
- Проверьте connection pool настройки
- Мониторьте database connections

### Kubernetes pod termination
- Проверьте `terminationGracePeriodSeconds`
- Убедитесь что preStop hook выполняется
- Мониторьте pod lifecycle события
