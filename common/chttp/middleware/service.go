package middleware

import (
	"context"
	"mojito-coding-test/common/chttp"
	"net/http"
)

// Set ServiceManager in Context
func SetServiceManager(sm chttp.ServiceManager) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, chttp.CtxKeyServiceManager, sm)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
