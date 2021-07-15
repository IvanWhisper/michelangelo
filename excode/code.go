package excode

import "github.com/IvanWhisper/michelangelo/exception"

type ExCode int32

func (c ExCode) ToNum() int32 {
	return int32(c)
}

func (c ExCode) ToEx(msg string) exception.Exception {
	return exception.New(c.ToNum(), msg)
}
