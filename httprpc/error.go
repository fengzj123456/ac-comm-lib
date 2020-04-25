package httprpc

import (
	"errors"
	"fmt"

	"git.ablecloud.cn/ablecloud/ac-comm-lib/httprpc/codes"
)

var StackTrace = false

type Error struct {
	code  codes.Code
	cause error
	stack []byte
}

func NewError(code codes.Code, cause error) error {
	var stack []byte
	if StackTrace {
		stack = runtimeStack()
	}
	return &Error{code: code, cause: cause, stack: stack}
}

func Errorf(code codes.Code, format string, a ...interface{}) error {
	var stack []byte
	if StackTrace {
		stack = runtimeStack()
	}
	return &Error{
		code:  code,
		cause: fmt.Errorf(format, a...),
		stack: stack,
	}
}

func clientError(code codes.Code, cause, stack string) error {
	return &Error{
		code:  code,
		cause: errors.New(cause),
		stack: []byte(stack),
	}
}

func (e *Error) Code() codes.Code {
	return e.code
}

func (e *Error) Cause() error {
	return e.cause
}

func (e *Error) Stack() []byte {
	return e.stack
}

func (e *Error) Error() string {
	return fmt.Sprintf("{code: %d, desc: %s, cause: %v}", e.code, e.code.String(), e.cause)
}

type ErrorCode interface {
	Code() codes.Code
}

type ErrorCause interface {
	Cause() error
}

type ErrorStack interface {
	Stack() []byte
}

func GetErrorCode(err error) codes.Code {
	if e, ok := err.(ErrorCode); ok {
		return e.Code()
	}
	return codes.Unknown
}

func GetErrorCause(err error) error {
	if e, ok := err.(ErrorCause); ok {
		return e.Cause()
	}
	return err
}

func GetLastErrorCause(err error) error {
	if e, ok := err.(ErrorCause); ok {
		return GetLastErrorCause(e.Cause())
	}
	return err
}

func GetErrorStack(err error) []byte {
	if e, ok := err.(ErrorStack); ok {
		return e.Stack()
	}
	return nil
}
