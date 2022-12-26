package snsgo

import "fmt"

type ErrorCode uint

const (
	ErrorCodeUnknown ErrorCode = iota
	ErrorCodeInvalidArg
	ErrorCodeNotFound
	ErrorCodeUnauthorized
)

type Error struct {
	code    ErrorCode
	msg     string
	errOrig error
}

func WrapErrorf(orig error, code ErrorCode, msg string, args ...interface{}) error {
	return &Error{
		code:    code,
		errOrig: orig,
		msg:     fmt.Sprintf(msg, args...),
	}
}

func NewErrorf(code ErrorCode, msg string, args ...interface{}) error {
	return WrapErrorf(nil, code, msg, args...)
}

func (e *Error) Error() string {
	if e.errOrig != nil {
		return fmt.Sprintf("%s: %v", e.msg, e.errOrig)
	}

	return e.msg
}

func (e *Error) Unwrap() error {
	return e.errOrig
}

func (e *Error) Code() ErrorCode {
	return e.code
}
