package chttp

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"gorm.io/gorm"
	"mojito-coding-test/common/config"
	"mojito-coding-test/common/errs"
	"net/http"
	"strconv"
)

const (
	ContentTypeYAML = "application/x-yaml"
)

func NewContext(res http.ResponseWriter, req *http.Request) *Context {
	return &Context{req: req, res: res}
}

type Context struct {
	req *http.Request
	res http.ResponseWriter
}

// Transform

func (c *Context) WithContext(ctx context.Context) *Context {
	return &Context{
		req: c.req.WithContext(ctx),
		res: c.res,
	}
}

// Utility

func (c *Context) GetChiContext() *chi.Context {
	return GetChiContext(c.req.Context())
}

func (c *Context) GetServiceManager() ServiceManager {
	return GetServiceManager(c.req.Context())
}

func (c *Context) GetLogger() *logrus.Entry {
	return GetLogger(c.req.Context())
}

func (c *Context) GetConfig() *config.Config {
	return GetConfig(c.req.Context())
}

func (c *Context) GetDb() *gorm.DB {
	return GetDb(c.req.Context())
}

func (c *Context) GetOrigCorrelationId() string {
	return GetOrigCorrelationId(c.req.Context())
}

func (c *Context) GetCorrelationId() string {
	return GetCorrelationId(c.req.Context())
}

func (c *Context) URLParamInt(key string) (int, error) {
	str := chi.URLParam(c.req, key)
	if str != "" {
		return strconv.Atoi(str)
	}

	return 0, nil
}

func (c *Context) URLParamInt64(key string) (int64, error) {
	str := chi.URLParam(c.req, key)
	if str != "" {
		return strconv.ParseInt(str, 10, 64)
	}

	return 0, nil
}

func (c *Context) URLParamString(key string) string {
	return chi.URLParam(c.req, key)
}

// Request

func (c *Context) IsYAMLRequest() bool {
	if cType := c.req.Header.Get("Content-Type"); cType != "" {
		return cType == ContentTypeYAML
	}

	return false
}

func (c *Context) Decode(v interface{}) error {
	if c.IsYAMLRequest() {
		decoder := yaml.NewDecoder(c.req.Body)
		return decoder.Decode(v)
	} else {
		decoder := json.NewDecoder(c.req.Body)
		return decoder.Decode(v)
	}
}

func (c *Context) MustDecode(v interface{}) bool {
	if err := c.Decode(v); err != nil {
		c.InternalError(fmt.Errorf("invalid %s body data: %v", c.req.Method, err))
		return false
	}

	return true
}

// Response

func (c *Context) Redirect(location string) {
	c.res.Header().Set("Location", location)
	c.res.WriteHeader(http.StatusMovedPermanently)
}

func (c *Context) Respond(code int, payload interface{}) {
	logger := c.GetLogger()

	if c.IsYAMLRequest() {
		response, err := yaml.Marshal(payload)
		if err != nil {
			errMsg := fmt.Sprintf("error encoding response - %v", err)
			code = http.StatusInternalServerError
			logger.Error(errMsg)

			response, _ = yaml.Marshal(errs.ErrorBody{StatusCode: code, Message: errMsg,
				CorrelationId: c.GetCorrelationId()})
		}

		c.res.Header().Set("Content-Type", ContentTypeYAML)
		c.res.WriteHeader(code)
		c.res.Write(response)
	} else {
		response, err := json.Marshal(payload)
		if err != nil {
			errMsg := fmt.Sprintf("error encoding response - %v", err)
			code = http.StatusInternalServerError
			logger.Error(errMsg)

			response, _ = json.Marshal(errs.ErrorBody{StatusCode: code, Message: errMsg,
				CorrelationId: c.GetCorrelationId()})
		}

		c.res.Header().Set("Content-Type", "application/json")
		c.res.WriteHeader(code)
		c.res.Write(response)
	}
}

func (c *Context) RespondWithError(error errs.ErrorBody) {
	c.Respond(error.StatusCode, error)
}

func (c *Context) ErrorCode(errCode errs.ErrorCode, code int, message string) {
	c.RespondWithError(errs.ErrorBody{ErrorCode: errCode, StatusCode: code, Message: message,
		CorrelationId: c.GetCorrelationId()})
}

func (c *Context) Error(code int, message string) {
	c.RespondWithError(errs.ErrorBody{StatusCode: code, Message: message,
		CorrelationId: c.GetCorrelationId()})
}

func (c *Context) InternalError(err error) {
	logger := c.GetLogger()

	if errorBody, ok := err.(errs.ErrorBody); ok {
		// Update logger
		if errorBody.Fields != nil && len(errorBody.Fields) > 0 {
			logger = logger.WithFields(errorBody.Fields)
		}
		if errorBody.ErrorCode != "" {
			logger = logger.WithField("error_code", errorBody.ErrorCode)
		}

		// Default status
		if errorBody.StatusCode < 1 {
			errorBody.StatusCode = http.StatusInternalServerError
		}

		// Assign correlation id.
		errorBody.CorrelationId = c.GetCorrelationId()

		switch errorBody.StatusCode {
		case http.StatusUnauthorized:
			logger.Warnf("unauthorized - %v", errorBody.Error())
			break
		case http.StatusForbidden:
			logger.Warnf("forbidden - %v", errorBody.Error())
			break
		case http.StatusNotFound:
			logger.Warnf("not found - %v", errorBody.Error())
			break
		case http.StatusInternalServerError:
			logger.Errorf("internal error - %v", errorBody.Error())
			break
		case http.StatusBadRequest:
			logger.Warnf("bad request - %v", errorBody.Error())
			break
		case http.StatusServiceUnavailable:
			logger.Warnf("service unavailable - %v", errorBody.Error())
			break
		default:
			logger.Warnf("%d - %v", errorBody.StatusCode, errorBody.Error())
			break
		}

		c.RespondWithError(errorBody)
	} else {
		logger.Errorf("internal error - %v", err)
		c.Error(http.StatusInternalServerError, err.Error())
	}
}

func (c *Context) Success(code int, message string) {
	c.Respond(code, SuccessBody{Code: code, Message: message, CorrelationId: c.GetCorrelationId()})
}
