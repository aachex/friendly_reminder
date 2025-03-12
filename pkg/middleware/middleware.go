package middleware

import "net/http"

func RequireAuthorization(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			http.Error(w, "Unautorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}
