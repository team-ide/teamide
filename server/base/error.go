package base

import (
	"fmt"
)

type valdateError struct {
	args []interface{}
}

var (
	valdateErrorType = GetRefValue(NewErrorValdate("")).Type()
)

func (e *valdateError) Error() string {
	return fmt.Sprint(e.args...)
}

func NewErrorValdate(args ...interface{}) *valdateError {

	return &valdateError{args}
}

func IsErrorValdate(err error) bool {
	errValue := GetRefValue(err)
	return errValue.Type() == valdateErrorType
}
