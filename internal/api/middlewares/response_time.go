package middlewares

import (
	"fmt"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func ResponseTimeMiddleware(next http.Handler) http.Handler {
	fmt.Println("")
	fmt.Println("Response time middleware...")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Response time middleware being returned...")
		fmt.Println("")
		// fmt.Println("Recieved request in responseTime")
		start := time.Now()

		wrappedWriter := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		duration := time.Since(start)
		w.Header().Set("X-Response-Time", duration.String())
		next.ServeHTTP(wrappedWriter, r)
		duration = time.Since(start)

		fmt.Printf("Method: %s\nUrl: %s\nStatus: %d\nDuration: %v\n", r.Method, r.URL, wrappedWriter.status, duration.String())
		fmt.Println("Response time middleware ends...")
	})
}
