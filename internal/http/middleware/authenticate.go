package middleware

import (
	"net/http"

	"github.com/Iyeyasu/bingo-paste/internal/config"
	"github.com/Iyeyasu/bingo-paste/internal/session"
)

// Authenticate handles user authentication.
func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if config.Get().Authentication.Enabled && session.User(r) == nil {
			http.Redirect(w, r, "/login", http.StatusFound)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
