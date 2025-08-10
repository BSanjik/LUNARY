package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/BSanjik/LUNARY/services/common/shared"
)

type Proxy struct {
	Services map[string]string
}

func NewProxy(services map[string]string) *Proxy {
	return &Proxy{Services: services}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	shared.LogInfo("Входящий запрос", map[string]interface{}{
		"method":      r.Method,
		"path":        r.URL.Path,
		"remote_addr": r.RemoteAddr,
	})

	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	if len(pathParts) == 0 {
		shared.LogWarn("Пустой путь в запросе", map[string]interface{}{
			"path":        r.URL.Path,
			"remote_addr": r.RemoteAddr,
		})
		shared.WriteErrorResponse(w, shared.NewValidationError("Некорректный путь запроса"))
		return
	}

	service := pathParts[0]
	backend, ok := p.Services[service]
	if !ok {
		shared.LogWarn("Сервис не найден", map[string]interface{}{
			"service":            service,
			"available_services": getServiceNames(p.Services),
			"path":               r.URL.Path,
		})
		shared.WriteErrorResponse(w, shared.NewNotFoundError("Сервис не найден"))
		return
	}

	backendURl, err := url.Parse(backend)
	if err != nil {
		shared.LogWarn("Ошибка парсинга URL бэкенда", map[string]interface{}{
			"backend": backend,
			"servoce": service,
		})
		shared.WriteErrorResponse(w, shared.NewNotFoundError("Внутренняя ошибка сервера"))
		return
	}

	originalPath := r.URL.Path
	r.URL.Path = "/" + strings.Join(pathParts[1:], "/")
	r.Host = backendURl.Host

	shared.LogInfo("Перенаправление запроса", map[string]interface{}{
		"original_path": originalPath,
		"new_path":      r.URL.Path,
		"backend":       backend,
		"service":       service,
	})

	proxy := httputil.NewSingleHostReverseProxy(backendURl)
	proxy.ServeHTTP(w, r)

	duration := time.Since(start)
	shared.LogInfo("Запрос обработан", map[string]interface{}{
		"service":     service,
		"duration_ms": duration.Milliseconds(),
		"backend":     backend,
	})
}

// Вспомогательная функция для логирования списка сервисов
func getServiceNames(services map[string]string) []string {
	names := make([]string, 0, len(services))
	for name := range services {
		names = append(names, name)
	}
	return names
}
