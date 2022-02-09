package chttp

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"mojito-coding-test/common/config"
	"mojito-coding-test/common/data/dto"
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

func GetDb(ctx context.Context) *gorm.DB {
	if ctx == nil {
		return nil
	}
	if v, ok := ctx.Value(CtxKeyDb).(*gorm.DB); ok {
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

func GetAuth(ctx context.Context) dto.Auth {
	if ctx == nil {
		return dto.Auth{}
	}
	if v, ok := ctx.Value(CtxKeyAuth).(dto.Auth); ok {
		return v
	}
	return dto.Auth{}
}
