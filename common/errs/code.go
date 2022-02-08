package errs

type ErrorCode string

func NewType(errorCode ErrorCode, message string) ErrorType {
	return ErrorType{
		ErrorCode: errorCode,
		Message:   message,
	}
}

type ErrorType struct {
	ErrorCode ErrorCode `json:"error_code,omitempty" yaml:"error_code,omitempty"`
	Message   string    `json:"message" yaml:"message"`
}

const (
	//// System

	ErrCodePanic ErrorCode = "SY000"

	//// User

	ErrCodeUserAccessDenied ErrorCode = "US000"
)

//// System

var ErrTypePanic = NewType(ErrCodePanic,
	"Unexpected Fault")

//// User

var ErrTypeUserAccessDenied = NewType(ErrCodeUserAccessDenied,
	"User Access Denied")
