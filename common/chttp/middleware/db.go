package middleware

import (
	"context"
	"gorm.io/gorm"
	"mojito-coding-test/common/chttp"
	"net/http"
)

// Set Db in Context
func SetDb(db *gorm.DB) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, chttp.CtxKeyDb, db)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}
