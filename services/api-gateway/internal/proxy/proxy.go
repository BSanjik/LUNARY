package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type Proxy struct {
	Services map[string]string
}

func NewProxy(services map[string]string) *Proxy {
	return &Proxy{Services: services}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("🔍 Proxy: получен запрос %s %s", r.Method, r.URL.Path)

	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	log.Printf("🔍 Proxy: pathParts = %v", pathParts)

	if len(pathParts) == 0 {
		log.Printf("❌ Proxy: пустой путь")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	service := pathParts[0]
	log.Printf("🔍 Proxy: service = '%s'", service)

	backend, ok := p.Services[service]
	if !ok {
		log.Printf("❌ Proxy: сервис '%s' не найден в списке: %v", service, p.Services)
		http.Error(w, "Service Not Found", http.StatusNotFound)
		return
	}

	log.Printf("🔍 Proxy: backend = '%s'", backend)

	backendURl, err := url.Parse(backend)
	if err != nil {
		log.Printf("❌ Proxy: ошибка парсинга URL '%s': %v", backend, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	originalPath := r.URL.Path
	r.URL.Path = "/" + strings.Join(pathParts[1:], "/")
	r.Host = backendURl.Host

	log.Printf("🔍 Proxy: перенаправляем %s → %s%s", originalPath, backend, r.URL.Path)

	proxy := httputil.NewSingleHostReverseProxy(backendURl)
	proxy.ServeHTTP(w, r)

	log.Printf("✅ Proxy: запрос обработан")
}
