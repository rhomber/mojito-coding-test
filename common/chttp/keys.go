package chttp

type CtxKeyType int

const (
	CtxKeyLogger CtxKeyType = iota
	CtxKeyConfig
	CtxKeyServiceManager
	CtxKeyOrigCorrelationId
	CtxKeyCorrelationId
)
