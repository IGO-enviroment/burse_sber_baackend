package middleware

import (
	"boilerplate/gen"
	"net/http"
)

func CORS() gen.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "OPTIONS" {
				r.Header.Add("Access-Control-Allow-Origin", "*")
				r.Header.Add("Access-Control-Allow-Methods", "*")
				r.Header.Add("Access-Control-Allow-Headers", "*")
				r.Header.Add("Content-Type", "application/json")
				next.ServeHTTP(w, r)
			} else {
				w.WriteHeader(http.StatusOK)
				return
			}
		})
	}
}
