package middleware

import (
	"net/http"

	"github.com/Iyeyasu/bingo-paste/internal/config"
	"github.com/Iyeyasu/bingo-paste/internal/http/httpext"
	"github.com/Iyeyasu/bingo-paste/internal/session"
)

// Authorize handles user authorization.
func Authorize(next http.Handler, role config.Role) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !config.Get().Authentication.Enabled {
			next.ServeHTTP(w, r)
		} else if user := session.User(r); user != nil && user.Role >= role {
			next.ServeHTTP(w, r)
		} else {
			httpext.UnauthorizedError(w)
		}
	})
}
