package middleware

import (
	"net/http"

	"github.com/Iyeyasu/bingo-paste/internal/session"
)

type contextKey string

var requestCacheKey contextKey = "cache"

// StartSession handles starting a session for a user.
func StartSession(next http.Handler) http.Handler {
	return session.Default.Manager.LoadAndSave(insertCache(next))
}

func insertCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = session.WithRequestCache(r)
		next.ServeHTTP(w, r)
	})
}
