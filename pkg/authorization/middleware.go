package authorization

import (
	"net/http"
	"os"
)

func Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if len(auth) <= 8 {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		tok := auth[7:]

		_, err := ParseJWT(tok, []byte(os.Getenv("SECRET_STR")))
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		next(w, r)
	}
}
