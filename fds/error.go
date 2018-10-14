package fds

import (
	"errors"
	"fmt"
	"runtime"
	"time"
)

// Errors
var (
	ErrorEndpoint = errors.New("Wrong endpoint")
)

// ServerError is a common structure for FDS client error
type ServerError struct {
	code     int
	time     time.Time
	msg      string
	funcName string
}

// Error makes ServerError a string
func (e *ServerError) Error() string {
	return fmt.Sprintf("%s %s Code: [%d] Msg: %s", e.time.Format(time.ANSIC), e.funcName, e.code, e.msg)
}

// Code is the code of ServerError
func (e *ServerError) Code() int {
	return e.code
}

// Message is the msg of ServerError
func (e *ServerError) Message() string {
	return e.msg
}

// newServerError new a ServerError struct
func newServerError(msg string, code int) *ServerError {

	pc, _, _, _ := runtime.Caller(1)

	return &ServerError{
		code:     code,
		msg:      msg,
		time:     time.Now(),
		funcName: runtime.FuncForPC(pc).Name(),
	}
}
