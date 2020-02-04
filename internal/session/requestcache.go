package session

import (
	"context"
	"net/http"
)

type contextKey string

var cacheKey contextKey = "k"

// RequestCache is a memory cache that is faster than the persistant session but only
// lasts for the duration of the request. The cache is used to cache values
// when they are accessed the first time and are needed later during the request.
type RequestCache map[string]interface{}

// WithRequestCache adds a RequestCache for the given request.
func WithRequestCache(r *http.Request) *http.Request {
	cache := RequestCache{}
	ctx := context.WithValue(r.Context(), cacheKey, cache)
	return r.WithContext(ctx)
}

// GetRequestValue returns a value from the RequestCache.
func GetRequestValue(r *http.Request, key string) interface{} {
	cache := r.Context().Value(cacheKey)
	if cache == nil {
		return nil
	}

	return cache.(RequestCache)[key]
}

// SetRequestValue sets a value in the RequestCache.
func SetRequestValue(r *http.Request, key string, value interface{}) {
	cache := r.Context().Value(cacheKey)
	if cache != nil {
		cache.(RequestCache)[key] = value
	}
}
