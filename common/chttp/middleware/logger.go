package middleware

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/middleware"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"mojito-coding-test/common/chttp"
	"net/http"
	"strings"
	"time"
)

// Set Logger in Context
func SetLogger(logger *logrus.Entry) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ctx = context.WithValue(ctx, chttp.CtxKeyLogger, logger)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(fn)
	}
}

func Tracing(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		logger := chttp.GetLogger(ctx)
		logFields := logrus.Fields{}

		origCorrelationId := r.Header.Get(chttp.HeaderXOrigCorrelationId)
		if origCorrelationId != "" {
			logFields["orig_correlation_id"] = origCorrelationId
		}

		correlationId := r.Header.Get(chttp.HeaderXCorrelationId)
		if correlationId == "" {
			correlationId = uuid.New().String()
		}

		logFields["correlation_id"] = correlationId

		ctx = context.WithValue(ctx, chttp.CtxKeyCorrelationId, correlationId)

		if addrParts := strings.Split(r.RemoteAddr, ":"); len(addrParts) > 0 && addrParts[0] != "" {
			logFields["remote_ip"] = addrParts[0]
		}

		// Extra
		logFields["http_method"] = r.Method
		logFields["http_uri"] = r.RequestURI

		ctx = context.WithValue(ctx, chttp.CtxKeyLogger, logger.WithFields(logFields))
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

// StructuredLogger is a simple, but powerful implementation of a custom structured
// logger backed on logrus. I encourage users to copy it, adapt it and make it their
// own. Also take a look at https://github.com/pressly/lg for a dedicated pkg based
// on this work, designed for zcontext-based http routers.

func Logger() func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&StructuredLogger{})
}

type StructuredLogger struct {
}

func (l *StructuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	logger := chttp.GetLogger(r.Context())

	entry := &StructuredLoggerEntry{
		Logger: logger,
		l:      l,
		req:    r,
	}
	logFields := logrus.Fields{}

	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	logFields["http_scheme"] = scheme
	logFields["http_proto"] = r.Proto
	logFields["http_method"] = r.Method

	logFields["http_remote_addr"] = r.RemoteAddr
	logFields["http_user_agent"] = r.UserAgent()

	logFields["http_uri"] = r.RequestURI
	logFields["http_url"] = fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)

	entry.Logger = entry.Logger.WithFields(logFields)

	entry.Logger.Infoln(fmt.Sprintf(">>> %s %s", r.Method, r.RequestURI))

	return entry
}

type StructuredLoggerEntry struct {
	Logger *logrus.Entry
	l      *StructuredLogger
	req    *http.Request
}

func (e *StructuredLoggerEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	httpMethod := e.Logger.Data["http_method"].(string)
	httpRespElapsedMs := float64(elapsed.Nanoseconds()) / 1000000.0

	e.Logger = e.Logger.WithFields(logrus.Fields{
		"http_resp_status":       status,
		"http_resp_bytes_length": bytes,
		"http_resp_elapsed_ms":   httpRespElapsedMs,
	})

	if status < 400 {
		e.Logger.Infoln(fmt.Sprintf("<<< %s %s - %d - %fms", httpMethod,
			e.Logger.Data["http_uri"], status, httpRespElapsedMs))
	} else {
		e.Logger.Warnln(fmt.Sprintf("<<< %s %s - %d - %fms", httpMethod,
			e.Logger.Data["http_uri"], status, httpRespElapsedMs))
	}
}

func (e *StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	e.Logger = e.Logger.WithFields(logrus.Fields{
		"stack": string(stack),
		"panic": fmt.Sprintf("%+v", v),
	})
}

// Helper methods used by the application to get the request-scoped
// logger entry and set additional fields between handlers.
//
// This is a useful pattern to use to set state on the entry as it
// passes through the handler chain, which at any point can be logged
// with a call to .Print(), .InfoQueue(), etc.

func GetLogEntry(r *http.Request) logrus.FieldLogger {
	entry := middleware.GetLogEntry(r).(*StructuredLoggerEntry)
	return entry.Logger
}

func LogEntrySetField(r *http.Request, key string, value interface{}) {
	if entry, ok := r.Context().Value(middleware.LogEntryCtxKey).(*StructuredLoggerEntry); ok {
		entry.Logger = entry.Logger.WithField(key, value)
	}
}

func LogEntrySetFields(r *http.Request, fields map[string]interface{}) {
	if entry, ok := r.Context().Value(middleware.LogEntryCtxKey).(*StructuredLoggerEntry); ok {
		entry.Logger = entry.Logger.WithFields(fields)
	}
}
