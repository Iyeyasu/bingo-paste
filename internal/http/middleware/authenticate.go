package middleware

import (
	"net/http"

	"bingo/internal/session"
)

// Authenticate handles user authentication.
func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if session.User(r) == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
