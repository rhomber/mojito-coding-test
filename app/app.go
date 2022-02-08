package app

import (
	"github.com/sirupsen/logrus"
	"mojito-coding-test/app/handler/health"
	"mojito-coding-test/app/service"
	"mojito-coding-test/common/chttp"
	cmiddleware "mojito-coding-test/common/chttp/middleware"
	"mojito-coding-test/common/config"
)

type Application struct {
	Config         *config.Config   `inject:""`
	Logger         *logrus.Entry    `inject:""`
	ServiceManager *service.Manager `inject:""`
}

func (a *Application) Init() (*chttp.Router, error) {
	handlers, err := a.InitHandlers()
	if err != nil {
		return nil, err
	}

	return handlers, nil
}

func (a *Application) InitHandlers() (*chttp.Router, error) {
	r := chttp.NewRouter()

	cmiddleware.Init(r, a.ServiceManager).
		Each(cmiddleware.MountPoints{
			"/":         false,
			"/api/test": true,
		}, func(r *chttp.Router, mount string, isExternal bool) {
			a.initRoutes(r, isExternal)
		})

	return r, nil
}

func (a *Application) initRoutes(r *chttp.Router, isExternal bool) {
	// Health
	for _, mp := range []string{"/status", "/_ah/health"} {
		r.Get(mp, health.Healthz())
	}

	// V1
	r.Route("/v1", func(r *chttp.Router) {
		a.initPublicRoutes(r, isExternal)
	})
}

func (a *Application) initPublicRoutes(r *chttp.Router, isExternal bool) {

}
