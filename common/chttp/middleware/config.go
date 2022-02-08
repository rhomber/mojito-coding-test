package middleware

import (
	"context"
	"mojito-coding-test/common/chttp"
	"mojito-coding-test/common/config"
	"net/http"
)

// Set Config in Context
func SetConfig(cfg *config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, chttp.CtxKeyConfig, cfg)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
