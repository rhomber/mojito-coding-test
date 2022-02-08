package health

import (
	"encoding/json"
	"mojito-coding-test/common/chttp"
	"mojito-coding-test/common/chttp/handler/health"
	"net/http"
)

func Healthz() func(res http.ResponseWriter, req *http.Request, ctx *chttp.Context) {
	return func(res http.ResponseWriter, req *http.Request, ctx *chttp.Context) {
		overall := health.NewOverall(ctx.GetConfig())
		overall.AddItem(health.TimedCheck(ctx, health.CheckDbConn()))
		overall.SetStatus()

		if overall.Status != health.StatusOK {
			result, _ := json.Marshal(overall)
			ctx.GetLogger().Errorf("health probe failed: %s", string(result))
		}

		ctx.Respond(overall.GetHttpStatus(), overall)
	}
}
