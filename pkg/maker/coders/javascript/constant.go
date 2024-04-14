package javascript

import (
	"errors"
	"fmt"
	"github.com/team-ide/go-tool/util"
	"teamide/pkg/maker/base"
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
	if util.IsNotEmpty(model.Comment) {
		base.AppendLine(&code.Body, fmt.Sprintf(`// %s`, model.Comment), 0)
	}
	if util.IsNotEmpty(model.Note) {
		base.AppendLine(&code.Body, fmt.Sprintf(`// %s`, model.Note), 0)
	}
	for _, option := range model.Options {
		err = this_.appendOption(code, option)
		if err != nil {
			return
		}
	}
	return
}

func (this_ *constantCoder) appendOption(code *common.Code, model *modelers.ConstantOptionModel) (err error) {
	code.CodeType = this_.codeType
	code.Dir = this_.GetConstantDir()
	code.Name = this_.GetConstantName()
	if util.IsNotEmpty(model.Comment) {
		base.AppendLine(&code.Body, fmt.Sprintf(`// %s`, model.Comment), 0)
	}
	if util.IsNotEmpty(model.Note) {
		base.AppendLine(&code.Body, fmt.Sprintf(`// %s`, model.Note), 0)
	}
	valueType, err := this_.GetValueType(model.Type)
	if err != nil {
		err = errors.New("value type [" + model.Type + "] is not support.")
		return
	}
	if valueType == nil {
		err = errors.New("value type [" + model.Type + "] is not support.")
		return
	}
	if valueType.IsNumber {
		base.AppendLine(&code.Body, fmt.Sprintf(`const %s = %s;`, model.Name, model.Value), 0)
	} else {
		base.AppendLine(&code.Body, fmt.Sprintf(`const %s = "%s";`, model.Name, model.Value), 0)
	}
	return
}
