package model

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"teamide/application/base"
)

type VariableModel struct {
	Name      string      `json:"name,omitempty"`                       // 名称，同一个方法模块唯一
	Comment   string      `json:"comment,omitempty"`                    // 注释说明
	Value     string      `json:"value,omitempty" yaml:"value"`         // 值脚本
	IsList    bool        `json:"isList,omitempty" yaml:"isList"`       // 是否是列表
	IsPage    bool        `json:"isPage,omitempty" yaml:"isPage"`       // 是否是列表
	DataType  string      `json:"dataType,omitempty" yaml:"dataType"`   // 数据类型
	DataPlace string      `json:"dataPlace,omitempty" yaml:"dataPlace"` // 数据位置
	TryError  *ErrorModel `json:"tryError,omitempty" yaml:"tryError"`
}

func (this_ *ModelContext) GetVariableDataType(dataType string) (res *VariableDataType) {
	if dataType == "" {
		return
	}
	res = getVariableDataType(dataType)
	if res != nil {
		return
	}
	dataStruct := this_.GetStruct(dataType)
	if dataStruct == nil {
		return
	}
	res = &VariableDataType{
		Value:      dataStruct.Name,
		Text:       dataStruct.Comment,
		DataStruct: dataStruct,
	}
	return
}

func FormatVariableModel(variable interface{}) interface{} {
	if variable == nil {
		return nil
	}
	variableMap, variableMapOk := variable.(map[interface{}]interface{})
	fmt.Println("variableMap:", variableMapOk)
	if variableMapOk {
		fmt.Println("variableMap:", variableMap)
	}
	return variable
}

type VariableDataType struct {
	Value           string            `json:"value,omitempty"`
	Text            string            `json:"text,omitempty"`
	DataStruct      *StructModel      `json:"dataStruct,omitempty"`
	DataStructField *StructFieldModel `json:"dataStructField,omitempty"`
}

var (
	dataTypes = []*VariableDataType{}

	DATA_TYPE_STRING  = newVariableDataType("string", "字符串(String,string)")
	DATA_TYPE_INT     = newVariableDataType("int", "整形(int,int32)")
	DATA_TYPE_LONG    = newVariableDataType("long", "长整型(long,int64)")
	DATA_TYPE_BOOLEAN = newVariableDataType("boolean", "布尔型(boolean,bool)")
	DATA_TYPE_BYTE    = newVariableDataType("byte", "字节型(byte,int8)")
	DATA_TYPE_DATE    = newVariableDataType("date", "字节型(Date,date)")
	DATA_TYPE_SHORT   = newVariableDataType("short", "短整型(short,int16)")
	DATA_TYPE_DOUBLE  = newVariableDataType("double", "双精度浮点型(double,float64)")
	DATA_TYPE_FLOAT   = newVariableDataType("float", "浮点型(float,float32)")
	DATA_TYPE_MAP     = newVariableDataType("map", "集合(Map,map)")
)

func newVariableDataType(value, text string) *VariableDataType {
	res := &VariableDataType{
		Value: value,
		Text:  text,
	}
	dataTypes = append(dataTypes, res)
	return res
}

func getVariableDataType(value string) *VariableDataType {
	for _, one := range dataTypes {
		if one.Value == value {
			return one
		}
	}
	return nil
}

type VariableDataPlace struct {
	Value string `json:"value,omitempty"`
	Text  string `json:"text,omitempty"`
}

var (
	dataPlaces = []*VariableDataPlace{}

	DATA_PLACE_BODY   = newVariableDataPlace("body", "Request Body")
	DATA_PLACE_HEADER = newVariableDataPlace("header", "Request Header")
	DATA_PLACE_PARAM  = newVariableDataPlace("param", "Request Param")
	DATA_PLACE_FILE   = newVariableDataPlace("file", "Request File")
	DATA_PLACE_PATH   = newVariableDataPlace("path", "Request Path")
	DATA_PLACE_FORM   = newVariableDataPlace("form", "Request Form")
)

func newVariableDataPlace(value, text string) *VariableDataPlace {
	res := &VariableDataPlace{
		Value: value,
		Text:  text,
	}
	dataPlaces = append(dataPlaces, res)
	return res
}

func GetVariableDataPlace(value string) *VariableDataPlace {
	for _, one := range dataPlaces {
		if strings.EqualFold(one.Value, value) {
			return one
		}
	}
	return nil
}

func getVariablesByValue(value interface{}) (variables []*VariableModel, err error) {
	if value == nil {
		return
	}
	values, valuesOk := value.([]interface{})
	if !valuesOk {
		err = errors.New(fmt.Sprint("value [", value, "] type [", reflect.TypeOf(value).Name(), "] to variable list error"))
		return

	}
	if len(values) == 0 {
		return
	}

	for _, valuesOne := range values {
		switch v := valuesOne.(type) {
		case map[string]interface{}:
			if len(v) == 0 {
				break
			}
			variableMap := v
			if len(v) == 1 {
				for mapKey, mapValue := range v {
					switch subV := mapValue.(type) {
					case map[string]interface{}:
						variableMap = subV
					default:
						variableMap["value"] = fmt.Sprint(subV)
					}
					if variableMap["name"] == nil {
						variableMap["name"] = mapKey
					}
				}
			}
			formatTryErrorByMap(variableMap)
			variable := &VariableModel{}
			err = base.ToBean([]byte(base.ToJSON(variableMap)), variable)
			if err != nil {
				return
			}
			variables = append(variables, variable)
			if variableMap["fields"] != nil {
				var fieldVariables []*VariableModel
				fieldVariables, err = getVariablesByValue(variableMap["fields"])
				if err != nil {
					return
				}
				for _, one := range fieldVariables {
					one.Name = variable.Name + "." + one.Name
					variables = append(variables, one)
				}
			}
			if variableMap["list"] != nil {
				switch list := variableMap["list"].(type) {
				case []interface{}:
					for index, listOne := range list {
						var listVariables []*VariableModel
						listVariables, err = getVariablesByValue(listOne)
						if err != nil {
							return
						}
						for _, listVariable := range listVariables {
							listVariable.Name = variable.Name + "[" + fmt.Sprint(index) + "]" + "." + listVariable.Name
							variables = append(variables, listVariable)
						}
					}
				}

			}
		default:
			err = errors.New(fmt.Sprint("value [", v, "] type [", reflect.TypeOf(v).Name(), "] to variable error"))
			return
		}
	}
	return
}
