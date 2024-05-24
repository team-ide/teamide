package maker

import "teamide/pkg/maker/modelers"

type ValueType struct {
	Name       string                `json:"name,omitempty"`
	Comment    string                `json:"comment,omitempty"`
	IsNumber   bool                  `json:"isNumber,omitempty"`
	Match      []string              `json:"match,omitempty"`
	Struct     *modelers.StructModel `json:"struct,omitempty"`
	FieldTypes map[string]*ValueType `json:"fieldTypes,omitempty"`
	isBase     bool
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

	ValueTypeAny = &ValueType{
		Name:     "any",
		Comment:  "Any",
		IsNumber: true,
		Match:    []string{"any", "Object", "interface{}"},
		isBase:   true,
	}
	ValueTypeNull = &ValueType{
		Name:     "null",
		Comment:  "NULL",
		IsNumber: true,
		Match:    []string{"null", "NULL", "nil"},
		isBase:   true,
	}
	ValueTypeError = &ValueType{
		Name:     "error",
		Comment:  "Error",
		IsNumber: true,
		Match:    []string{"error", "err"},
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
	ValueTypes = append(ValueTypes, ValueTypeAny)
	ValueTypes = append(ValueTypes, ValueTypeNull)
	ValueTypes = append(ValueTypes, ValueTypeError)

	for _, v := range ValueTypes {
		valueTypeCache[v.Name] = v
		for _, k := range v.Match {
			valueTypeCache[k] = v
		}
	}
}

var (
	valueTypeCache = make(map[string]*ValueType)
)

func GetValueType(name string) *ValueType {
	return valueTypeCache[name]
}

func GetValueTypes() []*ValueType {
	return ValueTypes
}
