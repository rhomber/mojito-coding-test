package origin

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
)

type Context struct {
	AppKey            string `json:"app_key" yaml:"app_key"`
	HttpMethod        string `json:"http_method" yaml:"http_method"`
	HttpUri           string `json:"http_uri" yaml:"http_uri"`
	RemoteIP          string `json:"remote_ip" yaml:"remote_ip"`
	CorrelationId     string `json:"correlation_id" yaml:"correlation_id"`
	OrigCorrelationId string `json:"orig_correlation_id" yaml:"orig_correlation_id"`
	jwt.StandardClaims
}

func (c Context) IsZero() bool {
	if c.AppKey != "" && c.HttpMethod != "" && c.HttpUri != "" {
		return false
	}

	return true
}

func (c Context) DebugLabel() string {
	return fmt.Sprintf("%s (%s %s)", c.AppKey, c.HttpMethod, c.HttpUri)
}
