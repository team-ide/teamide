package javascript

import (
	"errors"
	"fmt"
	"teamide/pkg/maker/coders/common"
	"teamide/pkg/maker/modelers"
	"teamide/pkg/util"
)

type constantCoder struct {
	*appCoder
}

func (this_ *constantCoder) Gen(code *common.Code, model *modelers.ConstantModel) (err error) {
	code.CodeType = this_.codeType
	code.Dir = this_.GetConstantDir()
	code.Name = this_.GetConstantName()
	if util.IsNotEmpty(model.Comment) {
		util.AppendLine(&code.Body, fmt.Sprintf(`// %s`, model.Comment), 0)
	}
	if util.IsNotEmpty(model.Note) {
		util.AppendLine(&code.Body, fmt.Sprintf(`// %s`, model.Note), 0)
	}
	var valueType = this_.GetValueType(model.Type)
	if valueType == nil {
		err = errors.New("value type [" + model.Type + "] is not support.")
		return
	}
	if valueType.IsNumber {
		util.AppendLine(&code.Body, fmt.Sprintf(`const %s = %s;`, model.Name, model.Value), 0)
	} else {
		util.AppendLine(&code.Body, fmt.Sprintf(`const %s = "%s";`, model.Name, model.Value), 0)
	}
	return
}
