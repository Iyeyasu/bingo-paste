package middleware

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	util "github.com/Iyeyasu/bingo-paste/internal/util/http"
	"github.com/Iyeyasu/bingo-paste/internal/util/log"
)

// LogMiddleware logs all incoming HTTP requests.
type LogMiddleware struct {
	next http.Handler
}

// NewLogMiddleware creates a new LogMiddleware.
func NewLogMiddleware(handler http.Handler) http.Handler {
	mw := new(LogMiddleware)
	mw.next = handler
	return mw
}

func (mw *LogMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if log.GetLevel() < log.TraceLevel {
		mw.next.ServeHTTP(w, r)
		return
	}

	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		util.WriteError(w, fmt.Sprintln("Failed to log request:", dump))
		return
	}

	log.Trace("Logging HTTP request")
	log.Trace(string(dump))
	mw.next.ServeHTTP(w, r)
}
