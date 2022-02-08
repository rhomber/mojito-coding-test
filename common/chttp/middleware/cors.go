package middleware

import (
	"github.com/go-chi/cors"
	"mojito-coding-test/common/core"
	"net/http"
)

func Cors(isExternal bool) func(next http.Handler) http.Handler {
	allowedOrigins := []string{"*"}
	if isExternal {
		if core.Config.IsEnvLocal() {
			allowedOrigins = []string{"http://localhost:3000", "http://localhost:9011", "http://localhost:9022"}
		} else {
			if appDomain := core.Config.GetString("app.domain"); appDomain != "" {
				allowedOrigins = []string{"https://" + appDomain}
			}
		}
	}

	crs := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link", "Content-Range"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	return crs.Handler
}
