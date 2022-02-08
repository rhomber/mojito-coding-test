package middleware

import (
	"mojito-coding-test/common/chttp"
	"net/http"
	"strings"
)

func RedirectHttpToHttps() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			if forProto := r.Header.Get("X-Forwarded-Proto"); forProto != "" {
				if strings.ToLower(forProto) == "http" {
					zctx := chttp.NewContext(w, r)
					zctx.Redirect("https://" + r.Host + r.RequestURI)

					return
				}
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
