package errs

import "net/http"

//// System

var ErrPanic = NewForType(ErrTypePanic,
	http.StatusInternalServerError, "uncaught panic")

//// Access
// I usually create more codes / better messages for this. I also return more verbose information that is hidden
// internally (but I removed that logic for the test).

var ErrUserAccessDenied = NewForType(ErrTypeUserAccessDenied,
	http.StatusForbidden, "user forbidden")
var ErrUserAuthMissingHeader = NewForType(ErrTypeUserAccessDenied,
	http.StatusBadRequest, "missing auth header")
var ErrUserAuthFailed = NewForType(ErrTypeUserAccessDenied,
	http.StatusUnauthorized, "authentication failed")
