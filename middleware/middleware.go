package middleware

import (
	"log"
	"net/http"
	"time"
)

func Register(mux http.Handler, middlewares ...func(http.Handler) http.Handler) (handler http.Handler) {
	handler = mux
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		log.Printf(
			"%s %s",
			r.URL.String(),
			duration.String(),
		)
	})
}
