package errors

import "fmt"

type Error struct {
	code int
	msg  string
}

func NewError(code int, msg string) *Error {
	return &Error{
		code: code,
		msg:  msg,
	}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%d: %s", e.code, e.msg)
}
