package middleware

import (
	"net/http"
	"net/url"
	"strings"
)

// TrimStrings handles trimming accidental whitespaces from HTTP url encoded requests.
func TrimStrings(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// if r.Method != http.MethodPost {
		// 	next.ServeHTTP(w, r)
		// 	return
		// }

		// body, err := ioutil.ReadAll(r.Body)
		// if err != nil {
		// 	next.ServeHTTP(w, r)
		// 	return
		// }

		// values, err := url.ParseQuery(string(body))
		// if err != nil {
		// 	next.ServeHTTP(w, r)
		// 	return
		// }

		// changed := trimValues(&values)
		// if !changed {
		// 	next.ServeHTTP(w, r)
		// 	return
		// }

		next.ServeHTTP(w, r)
	})
}

func trimValues(values *url.Values) bool {
	var changed bool

	for key := range *values {
		val := values.Get(key)
		trimmed := strings.TrimSpace(val)
		if len(val) != len(trimmed) {
			values.Set(key, trimmed)
			changed = true
		}
	}

	return changed
}
