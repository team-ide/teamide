package model

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"teamide/pkg/application/base"
)

type VariableModel struct {
	Name      string      `json:"name,omitempty" yaml:"name,omitempty"`           // 名称，同一个方法模块唯一
	Comment   string      `json:"comment,omitempty" yaml:"comment,omitempty"`     // 注释说明
	Value     string      `json:"value,omitempty" yaml:"value,omitempty"`         // 值脚本
	IsList    bool        `json:"isList,omitempty" yaml:"isList,omitempty"`       // 是否是列表
	IsPage    bool        `json:"isPage,omitempty" yaml:"isPage,omitempty"`       // 是否是列表
	DataType  string      `json:"dataType,omitempty" yaml:"dataType,omitempty"`   // 数据类型
	DataPlace string      `json:"dataPlace,omitempty" yaml:"dataPlace,omitempty"` // 数据位置
	TryError  *ErrorModel `json:"tryError,omitempty" yaml:"tryError,omitempty"`
}

func (this_ *ModelContext) GetVariableDataType(dataType string) (res *DataType) {
	if dataType == "" {
		return
	}
	res = getDataType(dataType)
	if res != nil {
		return
	}
	dataStruct := this_.GetStruct(dataType)
	if dataStruct == nil {
		return
	}
	res = &DataType{
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

type DataPlace struct {
	Value string `json:"value,omitempty"`
	Text  string `json:"text,omitempty"`
}

var (
	DATA_PLACES = []*DataPlace{}

	DATA_PLACE_BODY   = newDataPlace("body", "Request Body")
	DATA_PLACE_HEADER = newDataPlace("header", "Request Header")
	DATA_PLACE_PARAM  = newDataPlace("param", "Request Param")
	DATA_PLACE_FILE   = newDataPlace("file", "Request File")
	DATA_PLACE_PATH   = newDataPlace("path", "Request Path")
	DATA_PLACE_FORM   = newDataPlace("form", "Request Form")
)

func newDataPlace(value, text string) *DataPlace {
	res := &DataPlace{
		Value: value,
		Text:  text,
	}
	DATA_PLACES = append(DATA_PLACES, res)
	return res
}

func GetDataPlace(value string) *DataPlace {
	for _, one := range DATA_PLACES {
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
