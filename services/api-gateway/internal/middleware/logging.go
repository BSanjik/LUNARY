package middleware

import (
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.RequestURI)

		next.ServeHTTP(w, r)

		log.Printf("Completed %s in %v", r.RequestURI, time.Since(start))
	})
}
