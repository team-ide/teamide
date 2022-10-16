package javascript

import (
	"teamide/pkg/maker/coders/common"
	"teamide/pkg/maker/modelers"
)

type constantCoder struct {
	*appCoder
}

func (this_ *constantCoder) Gen(code *common.Code, model *modelers.ConstantModel) (err error) {
	code.CodeType = this_.codeType
	code.Dir = this_.GetConstantDir()
	code.Name = this_.GetConstantName()
	return
}
