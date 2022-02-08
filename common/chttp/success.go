package chttp

type SuccessBody struct {
	Code          int    `json:"code" yaml:"code"`
	Message       string `json:"message" yaml:"message"`
	CorrelationId string `json:"correlation_id,omitempty" yaml:"correlation_id,omitempty"`
}
