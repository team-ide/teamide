package base

import (
	"fmt"
)

type baseError struct {
	Code string
	Msg  string
}

func (e *baseError) Error() string {
	return e.Msg
}

var (
	baseErrorType = GetRefValue(NewBaseError("")).Type()
)

const (
	shouldLoginErrCode      = "100"
	noPowerErrCode          = "101"
	validateErrCode         = "200"
	FileSizeOversizeErrCode = "5001"
)

var (
	ShouldLoginError = NewBaseError(shouldLoginErrCode, "未检出到登录用户,请先登录!")
	NoPowerError     = NewBaseError(noPowerErrCode, "暂无权限操作!")
)

func NewBaseError(code string, args ...interface{}) *baseError {

	return &baseError{Code: code, Msg: fmt.Sprint(args...)}
}

func NewValidateError(args ...interface{}) *baseError {

	return NewBaseError(validateErrCode, args...)
}

func IsBaseError(err error) bool {
	errValue := GetRefValue(err)
	return errValue.Type() == baseErrorType
}

func ToBaseError(err error) *baseError {
	if !IsBaseError(err) {
		return nil
	}
	return err.(*baseError)
}

func IsValidateError(err error) bool {
	baseErr := ToBaseError(err)
	if baseErr == nil {
		return false
	}
	return baseErr.Code == validateErrCode
}
