package httpmw

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/Iyeyasu/bingo-paste/internal/http/httpext"
	"github.com/Iyeyasu/bingo-paste/internal/util/log"
)

// LogMiddleware handles logging HTTP requests.
func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if log.GetLevel() >= log.TraceLevel {
			dump, err := httputil.DumpRequest(r, true)
			if err != nil {
				httpext.WriteError(w, fmt.Sprintln("Failed to log request:", dump))
				return
			}

			log.Trace("Logging HTTP request")
			log.Trace(string(dump))
		}

		next.ServeHTTP(w, r)
	})
}
