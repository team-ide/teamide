package golang

import "teamide/pkg/maker"

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
}

// GetTypeStr 获取  类型 字符串 如 string、int
func (this_ *Generator) GetTypeStr(name string) (str string, err error) {
	valueType, err := this_.GetValueType(name)
	if err != nil {
		return
	}
	str = typeStr[valueType]
	if str == "" {
		// 获取对象类型
	}
	return
}

// GetTypeQuoteStr 获取  类型 引用 字符串 如 string 类型为 *string
func (this_ *Generator) GetTypeQuoteStr(name string) (str string, err error) {
	valueType, err := this_.GetValueType(name)
	if err != nil {
		return
	}
	str = typeStr[valueType]
	if str == "" {
		// 获取对象类型
	} else {
		str = "*" + str
	}
	return
}
