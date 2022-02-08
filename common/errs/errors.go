package errs

import "net/http"

//// System

var ErrPanic = NewForType(ErrTypePanic,
	http.StatusInternalServerError, "uncaught panic")

//// Access

var ErrUserAccessDenied = NewForType(ErrTypeUserAccessDenied,
	http.StatusForbidden, "user forbidden")
