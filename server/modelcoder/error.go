package modelcoder

import "fmt"

type ErrorVariableIsNull struct {
	Msg string
}

func (err *ErrorVariableIsNull) Error() string {
	return err.Msg
}

func newErrorVariableIsNull(args ...interface{}) *ErrorVariableIsNull {
	return &ErrorVariableIsNull{
		Msg: fmt.Sprint(args...),
	}
}

type ErrorServiceIsNull struct {
	Msg string
}

func (err *ErrorServiceIsNull) Error() string {
	return err.Msg
}

func newErrorServiceIsNull(args ...interface{}) *ErrorServiceIsNull {
	return &ErrorServiceIsNull{
		Msg: fmt.Sprint(args...),
	}
}

type ErrorDaoIsNull struct {
	Msg string
}

func (err *ErrorDaoIsNull) Error() string {
	return err.Msg
}

func newErrorDaoIsNull(args ...interface{}) *ErrorDaoIsNull {
	return &ErrorDaoIsNull{
		Msg: fmt.Sprint(args...),
	}
}
