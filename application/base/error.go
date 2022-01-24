package base

import (
	"fmt"
)

type baseError interface {
	Error() string
	GetCode() string
	GetMsg() string
}

type ErrorBase struct {
	Code string
	Msg  string
}

func (err *ErrorBase) Error() string {
	return fmt.Sprint("[", err.Code, "]", err.Msg)
}
func (err *ErrorBase) GetCode() string {
	return err.Code
}
func (err *ErrorBase) GetMsg() string {
	return err.Msg
}

func NewError(code string, args ...interface{}) baseError {
	if code == "" {
		code = "-1"
	}
	return &ErrorBase{
		Code: code,
		Msg:  fmt.Sprint(args...),
	}
}

func NewErrorVariableIsNull(args ...interface{}) baseError {
	return NewError("10001", args...)
}

func NewErrorActionIsNull(args ...interface{}) baseError {
	return NewError("10002", args...)
}

func NewErrorDaoIsNull(args ...interface{}) baseError {
	return NewError("10003", args...)
}

func NewErrorActionTypeIsWrong(args ...interface{}) baseError {
	return NewError("10004", args...)
}

func NewErrorActionStepIsWrong(args ...interface{}) baseError {
	return NewError("10005", args...)
}

func NewErrorDaoTypeIsWrong(args ...interface{}) baseError {
	return NewError("10006", args...)
}

func NewErrorApiTypeIsWrong(args ...interface{}) baseError {
	return NewError("10007", args...)
}

func NewErrorStructIsNull(args ...interface{}) baseError {
	return NewError("10008", args...)
}

func NewErrorParamIsNull(args ...interface{}) baseError {
	return NewError("10009", args...)
}

func NewErrorColumnValueIsNull(args ...interface{}) baseError {
	return NewError("10010", args...)
}
