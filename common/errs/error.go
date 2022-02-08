package errs

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

func NewForType(errorType ErrorType, statusCode int, details string) ErrorBody {
	return New(statusCode, details, errorType.ErrorCode, errorType.Message)
}

func New(statusCode int, details string, errorCode ErrorCode, message string) ErrorBody {
	return ErrorBody{
		StatusCode: statusCode,
		Details:    details,
		ErrorCode:  errorCode,
		Message:    message,
	}
}

type ErrorBody struct {
	StatusCode    int           `json:"status_code" yaml:"status_code"`
	Details       string        `json:"details,omitempty" yaml:"details,omitempty"`
	Message       string        `json:"message" yaml:"message"`
	ErrorCode     ErrorCode     `json:"error_code,omitempty" yaml:"error_code,omitempty"`
	CorrelationId string        `json:"correlation_id,omitempty" yaml:"correlation_id,omitempty"`
	Fields        logrus.Fields `json:"fields,omitempty" yaml:"fields,omitempty"`
}

func (e ErrorBody) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%s] %v: %v", e.ErrorCode, e.Message, e.Details)
	} else {
		return fmt.Sprintf("[%s] %v", e.ErrorCode, e.Message)
	}
}

func (e ErrorBody) WithStatusCode(statusCode int) ErrorBody {
	e.StatusCode = statusCode

	return e
}

func (e ErrorBody) WithMessage(message string) ErrorBody {
	e.Message = message

	return e
}

func (e ErrorBody) WithDetails(details string) ErrorBody {
	e.Details = details

	return e
}

func (e ErrorBody) WithDetailsF(format string, params ...interface{}) ErrorBody {
	e.Details = fmt.Sprintf(format, params...)

	return e
}

func (e ErrorBody) WithPrependedDetails(details string) ErrorBody {
	if e.Details != "" {
		e.Details = details + " " + e.Details
	} else {
		e.Details = details
	}

	return e
}

func (e ErrorBody) WithFields(fields logrus.Fields) ErrorBody {
	e.Fields = fields

	return e
}

func ToErrorBody(err error) ErrorBody {
	if errorBody, ok := err.(ErrorBody); ok {
		return errorBody
	}

	return ErrorBody{
		Message: err.Error(),
	}
}

func IsOfType(err error, cmpErrBodies ...ErrorBody) bool {
	if errBody, ok := err.(ErrorBody); ok {
		return BodyIsOfType(errBody, cmpErrBodies...)
	}

	return false
}

func BodyIsOfType(errBody ErrorBody, cmpErrBodies ...ErrorBody) bool {
	for _, cmp := range cmpErrBodies {
		if errBody.ErrorCode == cmp.ErrorCode {
			return true
		}
	}

	return false
}

func HasField(err error, fieldName string) bool {
	if errBody, ok := err.(ErrorBody); ok {
		if _, ok := errBody.Fields[fieldName]; ok {
			return true
		}
	}

	return false
}
