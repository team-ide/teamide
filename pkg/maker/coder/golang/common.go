package golang

import "strings"

var (
	commonCode = `package {pack}

import "fmt"

type Error struct {
	code string
	msg  string
}

func (this_ *Error) Error() string {
	return fmt.Sprintf("code:%s , msg:%s", this_.code, this_.msg)
}

// NewError 构造异常对象，code为错误码，msg为错误信息
func NewError(code string, msg string) *Error {
	err := &Error{
		code: code,
		msg:  msg,
	}
	return err
}

func GenId() (res int64){
	return
}

func GenStr(min int, max int) (res string){
	return
}
`
)

func (this_ *Generator) GenCommon() (err error) {
	dir := this_.golang.GetCommonDir(this_.Dir)
	if err = this_.Mkdir(dir); err != nil {
		return
	}
	path := dir + "common.go"
	builder, err := this_.NewBuilder(path)
	if err != nil {
		return
	}
	defer builder.Close()

	code := strings.ReplaceAll(commonCode, "{pack}", this_.golang.GetCommonPack())

	builder.AppendCode(code)
	return
}
