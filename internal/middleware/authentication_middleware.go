package middleware

import (
	"net/http"
)

// AuthenticationMiddleware handles user authentication.
type AuthenticationMiddleware struct {
	next http.Handler
}

// NewAuthenticationMiddleware creates a new AuthenticationMiddleware.
func NewAuthenticationMiddleware(handler http.Handler) http.Handler {
	mw := new(AuthenticationMiddleware)
	mw.next = handler
	return mw
}

func (mw *AuthenticationMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mw.next.ServeHTTP(w, r)
}
