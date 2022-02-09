package chttp

type CtxKeyType int

const (
	CtxKeyLogger CtxKeyType = iota
	CtxKeyConfig
	CtxKeyDb
	CtxKeyServiceManager
	CtxKeyOrigCorrelationId
	CtxKeyCorrelationId
)
