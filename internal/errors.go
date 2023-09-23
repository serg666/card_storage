package internal

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var (
	httpErrInternalServer   = &httpErrorWithStatus{500, "Internal server error"}
	httpErrValidationFailed = &httpErrorWithStatus{400, "Validation failed"}
	httpErrModelNotFound    = &httpErrorWithStatus{404, "Model not found"}
)

type httpError interface {
	getHttpStatus() int
	toJson() gin.H
}

type httpErrorWithStatus struct {
	status  int
	message string
}

func (err *httpErrorWithStatus) Error() string {
	return err.message
}
func (err *httpErrorWithStatus) toJson() gin.H {
	return gin.H{
		"message": err.message,
	}
}
func (err *httpErrorWithStatus) getHttpStatus() int {
	return err.status
}

type httpValidationError struct {
	Fields map[string]string
}

func (err httpValidationError) Error() string {
	return "Validation failed"
}
func (err httpValidationError) toJson() gin.H {
	return gin.H{
		"message": err.Error(),
		"errors":  err.Fields,
	}
}
func (err httpValidationError) getHttpStatus() int {
	return 400
}

func httpValidationErrorByValidationErrors(errs validator.ValidationErrors) *httpValidationError {
	fields := make(map[string]string, len(errs))
	for _, e := range errs {
		fields[e.Field()] = e.Error()
	}

	return &httpValidationError{
		Fields: fields,
	}
}
