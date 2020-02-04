package middleware

import (
	"net/http"

	"bingo/internal/config"
	"bingo/internal/http/httpext"
	"bingo/internal/session"
)

// Authorize handles user authorization.
func Authorize(next http.Handler, role config.Role) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user := session.User(r); user != nil && user.Role >= role {
			next.ServeHTTP(w, r)
		} else {
			httpext.UnauthorizedError(w)
		}
	})
}
