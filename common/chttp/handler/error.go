package handler

import (
	"mojito-coding-test/common/chttp"
	"net/http"
)

func NotFound(res http.ResponseWriter, req *http.Request, ctx *chttp.Context) {
	ctx.Error(http.StatusNotFound, "Not Found")
}

func MethodNotAllowed(res http.ResponseWriter, req *http.Request, ctx *chttp.Context) {
	ctx.Error(http.StatusMethodNotAllowed, "Method Not Allowed")
}
