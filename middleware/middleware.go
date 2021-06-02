package middleware

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"
)

const (
	ClientIPContextKey ContextKey = "ClientIPContextKey"
)

type (
	ContextKey string
)

func Register(mux http.Handler, middlewares ...func(http.Handler) http.Handler) (handler http.Handler) {
	handler = mux
	for i := len(middlewares) - 1; i >= 0; i-- {
		middleware := middlewares[i]
		handler = middleware(handler)
	}
	return
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(start)
		clientIP := ClientIPFromContext(r.Context())
		log.Printf(
			"%s %s %s",
			clientIP,
			r.URL.String(),
			duration.String(),
		)
	})
}

func ClientIPFromContext(ctx context.Context) net.IP {
	clientIP, ok := ctx.Value(ClientIPContextKey).(net.IP)
	if !ok {
		clientIP = net.IPv4zero
	}
	return clientIP
}
