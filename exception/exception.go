package exception

import (
	"fmt"
)

type Exception interface {
	Error() string
	Code() int32
	WithCode(code int32) Exception
	WithExtendMsg(ext string) Exception
	WithError(err error) Exception
	WithMsg(msg string) Exception
	ToError() error
}

type exception struct {
	code int32
	msg  string
}

func (e *exception) Error() string {
	return fmt.Sprintf("[%d]%s", e.code, e.msg)
}

func (e *exception) Code() int32 {
	return e.code
}

func (e *exception) WithCode(code int32) Exception {
	e.code = code
	return e
}

func (e *exception) WithExtendMsg(ext string) Exception {
	e.msg += fmt.Sprintf(",%s", ext)
	return e
}

func (e *exception) ToError() error {
	return e
}

func (e *exception) WithError(err error) Exception {
	if err != nil {
		e.msg = err.Error()
	}
	return e
}
func (e *exception) WithMsg(msg string) Exception {
	e.msg = msg
	return e
}

func NewByErr(code int32, e error) Exception {
	return &exception{
		code: code,
		msg:  e.Error(),
	}
}

func New(code int32, msg string) Exception {
	return &exception{
		code: code,
		msg:  msg,
	}
}
