package errors

import (
	"errors"
	"fmt"

	pkgErrors "github.com/pkg/errors"
)

// Define alias
var (
	New          = pkgErrors.New
	Unwrap       = pkgErrors.Unwrap
	Wrap         = pkgErrors.Wrap
	Wrapf        = pkgErrors.Wrapf
	WithStack    = pkgErrors.WithStack
	WithMessage  = pkgErrors.WithMessage
	WithMessagef = pkgErrors.WithMessagef
)

var (
	// ErrInvalidToken invalid
	ErrInvalidToken = NewResponse(401, 401, "invalid X-Authorization-Token")
	// ErrOpHomeNotLogin not login
	ErrOpHomeNotLogin = NewResponse(401, 401, "please login first")
	// ErrNoPerm no permission
	ErrNoPerm = NewResponse(403, 403, "no permission")
	// ErrNotFound not found
	ErrNotFound = NewResponse(404, 404, "url not found")
	// ErrMethodNotAllow method not allowed
	ErrMethodNotAllow = NewResponse(405, 405, "method not allowed")
	// ErrTooManyRequests too many requests
	ErrTooManyRequests = NewResponse(429, 429, "too many requests")
	// ErrInternalServer internal server error
	ErrInternalServer = NewResponse(500, 500, "internal server error")
	// ErrBadRequest bad request
	ErrBadRequest = New400Response("bad request")
	// ErrUserDisable user forbidden
	ErrUserDisable = New400Response("user disabled")
)

// ResponseError 定义响应错误
type ResponseError struct {
	Code    int    // 错误码
	Message string // 错误消息
	Status  int    // 响应状态码
	ERR     error  // 响应错误
}

func (r *ResponseError) Error() string {
	if r.ERR != nil {
		return r.ERR.Error()
	}
	return r.Message
}

// UnWrapResponse ...
func UnWrapResponse(err error) *ResponseError {
	if v, ok := err.(*ResponseError); ok {
		return v
	}
	return nil
}

// WrapResponse ..
func WrapResponse(err error, code, status int, msg string, args ...interface{}) error {
	res := &ResponseError{
		Code:    code,
		Message: fmt.Sprintf(msg, args...),
		ERR:     err,
		Status:  status,
	}
	return res
}

// Wrap400Response .
func Wrap400Response(err error, msg string, args ...interface{}) error {
	return WrapResponse(err, 0, 400, msg, args...)
}

// Wrap500Response .
func Wrap500Response(err error, msg string, args ...interface{}) error {
	return WrapResponse(err, 0, 500, msg, args...)
}

// NewResponse .
func NewResponse(code, status int, msg string, args ...interface{}) error {
	res := &ResponseError{
		Code:    code,
		Message: fmt.Sprintf(msg, args...),
		Status:  status,
	}
	return res
}

// New400Response .
func New400Response(msg string, args ...interface{}) error {
	return NewResponse(0, 400, msg, args...)
}

// New500Response .
func New500Response(msg string, args ...interface{}) error {
	return NewResponse(0, 500, msg, args...)
}

// Append return multi error
func Append(e, e2 error) error {
	if e2 == nil {
		return e
	}
	return errors.Join(e, e2)
}
