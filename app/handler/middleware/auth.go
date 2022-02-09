package middleware

import (
	"context"
	"mojito-coding-test/app/service"
	"mojito-coding-test/common/chttp"
	"mojito-coding-test/common/errs"
	"net/http"
)

const (
	HeaderXUserAuth = "X-User-Auth"
)

func AuthRequired(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		zctx := chttp.NewContext(w, r)

		sm := zctx.GetServiceManager().(*service.Manager)

		//// User Auth
		// Very primitive, just provide basic authentication.

		if userAuthHeader := r.Header.Get(HeaderXUserAuth); userAuthHeader != "" {
			authDTO, err := sm.UserAuth.Authenticate(zctx.GetDb(), userAuthHeader)
			if err != nil {
				zctx.InternalError(err)
				return
			}

			ctx = context.WithValue(ctx, chttp.CtxKeyAuth, authDTO)
		} else {
			zctx.InternalError(errs.ErrUserAuthMissingHeader)
			return
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	}

	return http.HandlerFunc(fn)
}
