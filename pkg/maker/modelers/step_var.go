package modelers

type StepVarModel struct {
	*StepModel `json:",inline"`

	Var   string `json:"var,omitempty"`   // 定义变量
	Value string `json:"value,omitempty"` // 值
	Type  string `json:"type,omitempty"`  // 值类型
}

var (
	docTemplateStepVarName = "step_var"
)

func init() {
	addDocTemplate(&docTemplate{
		Name: docTemplateStepVarName,
		Fields: []*docTemplateField{
			{Name: "var", Comment: "定义变量操作"},
			{Name: "value", Comment: "变量值"},
			{Name: "type", Comment: "变量值类型"},
		},
		Inline: "StepModel",
		inlineNewModel: func() interface{} {
			return &StepModel{}
		},
		newModel: func() interface{} {
			return &StepVarModel{}
		},
	})
}

type ValueType struct {
	Name     string       `json:"name,omitempty"`
	Comment  string       `json:"comment,omitempty"`
	IsNumber bool         `json:"isNumber,omitempty"`
	Match    []string     `json:"match,omitempty"`
	Struct   *StructModel `json:"struct,omitempty"`
	isBase   bool
}

func (this_ *ValueType) IsBase() bool {
	return this_.isBase
}

var (
	ValueTypes []*ValueType

	ValueTypeString = &ValueType{
		Name:    "string",
		Comment: "string",
		Match:   []string{"s", "string", "String"},
		isBase:  true,
	}

	ValueTypeInt = &ValueType{
		Name:     "i",
		Comment:  "int",
		IsNumber: true,
		Match:    []string{"i", "int", "Integer"},
		isBase:   true,
	}

	ValueTypeInt8 = &ValueType{
		Name:     "i8",
		Comment:  "int8",
		IsNumber: true,
		Match:    []string{"i8", "int8", "byte", "Byte"},
		isBase:   true,
	}

	ValueTypeInt16 = &ValueType{
		Name:     "i16",
		Comment:  "int16",
		IsNumber: true,
		Match:    []string{"i16", "int16", "short", "Short"},
		isBase:   true,
	}

	ValueTypeInt32 = &ValueType{
		Name:     "i32",
		Comment:  "int32",
		IsNumber: true,
		Match:    []string{"i32", "int32"},
		isBase:   true,
	}

	ValueTypeInt64 = &ValueType{
		Name:     "i64",
		Comment:  "int64",
		IsNumber: true,
		Match:    []string{"i64", "int64", "long", "Long"},
		isBase:   true,
	}

	ValueTypeFloat32 = &ValueType{
		Name:     "f32",
		Comment:  "float32",
		IsNumber: true,
		Match:    []string{"f32", "float32"},
		isBase:   true,
	}
	ValueTypeFloat64 = &ValueType{
		Name:     "f64",
		Comment:  "float64",
		IsNumber: true,
		Match:    []string{"f", "f64", "float64", "float", "Float"},
		isBase:   true,
	}

	ValueTypeBool = &ValueType{
		Name:     "bool",
		Comment:  "bool",
		IsNumber: true,
		Match:    []string{"bool", "boolean", "Boolean"},
		isBase:   true,
	}

	ValueTypeMap = &ValueType{
		Name:     "map",
		Comment:  "Map",
		IsNumber: true,
		Match:    []string{"map", "Map"},
		isBase:   true,
	}
)

func init() {
	ValueTypes = append(ValueTypes, ValueTypeString)
	ValueTypes = append(ValueTypes, ValueTypeBool)
	ValueTypes = append(ValueTypes, ValueTypeInt)
	ValueTypes = append(ValueTypes, ValueTypeInt64)
	ValueTypes = append(ValueTypes, ValueTypeFloat64)
	ValueTypes = append(ValueTypes, ValueTypeInt8)
	ValueTypes = append(ValueTypes, ValueTypeInt16)
	ValueTypes = append(ValueTypes, ValueTypeInt32)
	ValueTypes = append(ValueTypes, ValueTypeFloat32)
	ValueTypes = append(ValueTypes, ValueTypeMap)
}

func GetValueTypes() []*ValueType {
	return ValueTypes
}
