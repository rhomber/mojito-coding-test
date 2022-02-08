package chttp

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"mojito-coding-test/common/config"
)

func GetChiContext(ctx context.Context) *chi.Context {
	if ctx == nil {
		return nil
	}
	if chiCtx, ok := ctx.Value(chi.RouteCtxKey).(*chi.Context); ok {
		return chiCtx
	}
	return nil
}

func GetLogger(ctx context.Context) *logrus.Entry {
	if ctx == nil {
		return nil
	}
	if logger, ok := ctx.Value(CtxKeyLogger).(*logrus.Entry); ok {
		return logger
	}
	return nil
}

func GetConfig(ctx context.Context) *config.Config {
	if ctx == nil {
		return nil
	}
	if v, ok := ctx.Value(CtxKeyConfig).(*config.Config); ok {
		return v
	}
	return nil
}

func GetServiceManager(ctx context.Context) ServiceManager {
	if ctx == nil {
		return nil
	}
	if v, ok := ctx.Value(CtxKeyServiceManager).(ServiceManager); ok {
		return v
	}
	return nil
}

func GetOrigCorrelationId(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if v, ok := ctx.Value(CtxKeyOrigCorrelationId).(string); ok {
		return v
	}
	return ""
}

func GetCorrelationId(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if v, ok := ctx.Value(CtxKeyCorrelationId).(string); ok {
		return v
	}
	return ""
}
