package middleware

import (
	"log"
	"net/http"
)

func HTTPLogger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)
		next(w, r)
		log.Printf("Handled request: %s %s", r.Method, r.URL.Path)
	}
}
