package middleware

import (
	"fmt"
	"mojito-coding-test/common/chttp"
	"mojito-coding-test/common/config"
	"mojito-coding-test/common/errs"
	"net/http"
	"runtime/debug"
)

func Recoverer(config *config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil {
					ctx := chttp.NewContext(w, r)
					ctx.InternalError(errs.ErrPanic.
						WithDetails(fmt.Sprintf("%+v\n%s", rvr, string(debug.Stack()))))
				}
			}()

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}
