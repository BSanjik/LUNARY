package proxy

import (
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
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) == 0 {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	service := pathParts[0]
	backend, ok := p.Services[service]
	if !ok {
		http.Error(w, "Service Not Found", http.StatusNotFound)
		return
	}

	backendURl, err := url.Parse(backend)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	r.URL.Path = "/" + strings.Join(pathParts[1:], "/")
	r.Host = backendURl.Host

	proxy := httputil.NewSingleHostReverseProxy(backendURl)
	proxy.ServeHTTP(w, r)
}
