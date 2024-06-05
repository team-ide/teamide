package golang

import (
	"github.com/team-ide/go-tool/util"
	"teamide/pkg/maker"
)

var (
	typeStr = map[*maker.ValueType]string{}
)

func init() {
	typeStr[maker.ValueTypeString] = "string"
	typeStr[maker.ValueTypeInt] = "int"
	typeStr[maker.ValueTypeInt8] = "int8"
	typeStr[maker.ValueTypeInt16] = "int16"
	typeStr[maker.ValueTypeInt32] = "int32"
	typeStr[maker.ValueTypeInt64] = "int64"
	typeStr[maker.ValueTypeFloat32] = "float32"
	typeStr[maker.ValueTypeFloat64] = "float64"
	typeStr[maker.ValueTypeBool] = "bool"
	typeStr[maker.ValueTypeMap] = "map[string]any"
	typeStr[maker.ValueTypeContext] = "context.Context"
}

// GetTypeStr 获取  类型 字符串 如 string、int
func (this_ *Generator) GetTypeStr(valueType *maker.ValueType) (str string, err error) {
	str = typeStr[valueType]
	if str == "" {
		// 获取对象类型
		if valueType.Struct != nil {
			structName := util.FirstToUpper(valueType.Struct.Name)
			structPack := this_.golang.GetStructPack()
			str = "*" + structPack + "." + structName
		}
	}
	return
}

// GetTypeQuoteStr 获取  类型 引用 字符串 如 string 类型为 *string
func (this_ *Generator) GetTypeQuoteStr(valueType *maker.ValueType) (str string, err error) {
	str = typeStr[valueType]
	if str == "" {
		// 获取对象类型
	} else {
		str = "*" + str
	}
	return
}
