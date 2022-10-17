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
	Name     string `json:"name"`
	Comment  string `json:"comment"`
	IsNumber bool   `json:"isNumber"`
}

var (
	ValueTypes []*ValueType

	ValueTypeString = &ValueType{
		Name:    "string",
		Comment: "string",
	}

	ValueTypeInt = &ValueType{
		Name:     "i",
		Comment:  "int",
		IsNumber: true,
	}

	ValueTypeInt8 = &ValueType{
		Name:     "i8",
		Comment:  "int8",
		IsNumber: true,
	}

	ValueTypeInt16 = &ValueType{
		Name:     "i16",
		Comment:  "int16",
		IsNumber: true,
	}

	ValueTypeInt32 = &ValueType{
		Name:     "i32",
		Comment:  "int32",
		IsNumber: true,
	}

	ValueTypeInt64 = &ValueType{
		Name:     "i64",
		Comment:  "int64",
		IsNumber: true,
	}
)

func init() {
	ValueTypes = append(ValueTypes, ValueTypeString)
	ValueTypes = append(ValueTypes, ValueTypeInt)
	ValueTypes = append(ValueTypes, ValueTypeInt64)
	ValueTypes = append(ValueTypes, ValueTypeInt8)
	ValueTypes = append(ValueTypes, ValueTypeInt16)
	ValueTypes = append(ValueTypes, ValueTypeInt32)
}
