package httpext

import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// ParseID parses object id from HTTP request.
func ParseID(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	return strconv.ParseInt(params.ByName("id"), 10, 64)
}

// ParseFilter parses filter string from the HTTP request.
func ParseFilter(r *http.Request) string {
	query := r.URL.Query()
	return query.Get("search")
}

// ParseRange parses limits for a list HTTP request.
func ParseRange(r *http.Request) (int64, int64) {
	query := r.URL.Query()
	limitParam := query.Get("limit")
	offsetParam := query.Get("offset")

	limit, err := strconv.ParseInt(limitParam, 10, 64)
	if err != nil {
		limit = 10
	}

	offset, err := strconv.ParseInt(offsetParam, 10, 64)
	if err != nil {
		offset = 0
	}

	return limit, offset
}
