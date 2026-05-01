package middleware

import (
	"log"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter        // embedding: ganha todos os métodos do original
	status int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status                    // salva o status
	rw.ResponseWriter.WriteHeader(status) // repassa para o original
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)
		log.Printf("%s %s %d — %s", r.Method, r.URL.Path, rw.status, time.Since(start))
	})
}
