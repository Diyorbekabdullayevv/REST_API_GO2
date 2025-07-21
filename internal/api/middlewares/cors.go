package middlewares

import (
	"fmt"
	"net/http"
)

var allowedOrigins = []string{
	"https://my-origin-url.com",
	"https://localhost:3000",
}

func Cors(next http.Handler) http.Handler {
	fmt.Println("")
	fmt.Println("CORS middleware...")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("CORS middleware being returned...")
		origin := r.Header.Get("Origin")
		fmt.Println(origin)

		if IsOriginAllowed(origin) {
			w.Header().Set("Acces-Control-Allow-Origin", origin)
		} else {
			http.Error(w, "Not allowed by CORS", http.StatusForbidden)
			return
		}

		w.Header().Set("Acces-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Acces-Control-Expose-Headers", "Authorization")
		w.Header().Set("Acces-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE")
		w.Header().Set("Acces-Control-Allow-Credentials", "true")
		w.Header().Set("Acces-Control-Max-Age", "3600")

		if r.Method == http.MethodOptions {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(w, r)
		fmt.Println("CORS middleware ends...")
	})
}

func IsOriginAllowed(origin string) bool {
	for _, allowedOrigin := range allowedOrigins {
		if origin == allowedOrigin {
			return true
		}
	}
	return false
}
