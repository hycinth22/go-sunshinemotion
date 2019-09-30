package ssmt

import (
	"errors"
	"fmt"
)

type httpError struct {
	msg     string
	httpErr error
}

func (e httpError) Error() string {
	return e.msg + "\n" + e.httpErr.Error()
}

func IsHTTPError(e error) bool {
	_, ok := e.(httpError)
	return ok
}

type IServiceError interface {
	error
	GetCode() int64
	GetMsg() string
	SetMsg(string)
	Equal(error) bool
}

type ServiceError struct {
	code int64
	what string
	msg  string
}

func (e *ServiceError) Error() string {
	return fmt.Sprintf("response status %d(%s) , message: %s", e.code, e.what, e.msg)
}
func (e *ServiceError) GetCode() int64 {
	return e.code
}
func (e *ServiceError) GetMsg() string {
	return e.msg
}
func (e *ServiceError) SetMsg(m string) {
	e.msg = m
}
func (e *ServiceError) Equal(err error) bool {
	var r IServiceError
	return errors.As(err, &r) && e.code == r.GetCode()
}
