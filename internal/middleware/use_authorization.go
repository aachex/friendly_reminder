package middleware

import (
	"net/http"
	"os"

	"github.com/artemwebber1/friendly_reminder/pkg/jwtservice"
)

func UseAuthorization(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if len(auth) <= 8 {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		tok := auth[7:]

		_, err := jwtservice.Parse(tok, []byte(os.Getenv("SECRET_STR")))
		if err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		next(w, r)
	}
}
