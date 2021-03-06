package middleware

import (
	"net/http"
	"net/http/httputil"

	"bingo/internal/util/log"
)

// Log handles logging HTTP requests.
func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if log.GetLevel() >= log.TraceLevel {
			dump, err := httputil.DumpRequest(r, true)
			if err != nil {
				log.Trace("Failed to log HTTP request")
				return
			}
			log.Trace(string(dump))
		}

		next.ServeHTTP(w, r)
	})
}
