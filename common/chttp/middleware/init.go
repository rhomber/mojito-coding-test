package middleware

import (
	chim "github.com/go-chi/chi/middleware"
	"mojito-coding-test/common/chttp"
	"mojito-coding-test/common/chttp/handler"
	"mojito-coding-test/common/core"
	"net/http"
)

type MountPoints map[string]bool

type initializer struct {
	r           *chttp.Router
	middlewares []func(next http.Handler) http.Handler
}

func Init(r *chttp.Router, sm chttp.ServiceManager) *initializer {
	//// Global Middleware
	var middlewares []func(next http.Handler) http.Handler
	middlewares = append(middlewares, SetLogger(core.Logger))
	middlewares = append(middlewares, SetConfig(core.Config))
	middlewares = append(middlewares, SetServiceManager(sm))
	middlewares = append(middlewares, chim.RealIP)
	middlewares = append(middlewares, Tracing)
	if core.Config.GetBool("http.logging") {
		middlewares = append(middlewares, Logger())
	}
	middlewares = append(middlewares, chim.NoCache)
	middlewares = append(middlewares, Recoverer(core.Config))
	middlewares = append(middlewares, chim.Timeout(core.Config.GetDuration("http.timeout")))

	return &initializer{
		r:           r,
		middlewares: middlewares,
	}
}

func (i *initializer) Each(mps MountPoints, handlerCb func(r *chttp.Router, mount string, isExternal bool)) {
	for mount, isExternal := range mps {
		i.r.Route(mount, func(r *chttp.Router) {
			//// Apply Globals
			for _, mid := range i.middlewares {
				r.Use(mid)
			}

			//// Apply Externals
			if isExternal {
				r.Use(RedirectHttpToHttps())
			}
			r.Use(Cors(isExternal))

			//// Error Handlers
			r.NotFound(handler.NotFound)
			r.MethodNotAllowed(handler.MethodNotAllowed)

			// Handler
			handlerCb(r, mount, isExternal)
		})
	}
}
