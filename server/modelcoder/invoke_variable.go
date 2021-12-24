package modelcoder

type VariableData struct {
	Name       string       `json:"name,omitempty"`     // 变量名称
	Data       interface{}  `json:"data,omitempty"`     // 变量值
	DataType   *DataType    `json:"dataType,omitempty"` // 数据类型
	DataStruct *StructModel `json:"-"`                  // 数据结构体
}

type invokeVariable struct {
	Parent        *invokeVariable        `json:"-"`
	VariableDatas []*VariableData        `json:"variableDatas,omitempty"`
	InvokeCache   map[string]interface{} `json:"invokeCache,omitempty"`
}

func (this_ *invokeVariable) GetVariableData(name string) *VariableData {
	if len(this_.VariableDatas) == 0 {
		return nil
	}
	for _, one := range this_.VariableDatas {
		if one.Name == name {
			return one
		}
	}
	return nil
}

func (this_ *invokeVariable) AddVariableData(list ...*VariableData) *invokeVariable {
	this_.VariableDatas = append(this_.VariableDatas, list...)
	return this_
}

func (this_ *invokeVariable) Clone() *invokeVariable {
	res := &invokeVariable{
		VariableDatas: []*VariableData{},
		Parent:        this_,
		InvokeCache:   map[string]interface{}{},
	}
	return res
}

func processVariableDatas(application *Application, parameters []*VariableModel, variable *invokeVariable) (err error) {
	if len(parameters) == 0 {
		return
	}
	for _, one := range parameters {
		err = processVariableData(application, one, variable)
		if err != nil {
			return
		}
	}
	return
}

func processVariableData(application *Application, parameter *VariableModel, variable *invokeVariable) (err error) {
	if parameter == nil {
		return
	}
	if application.OutDebug() {
		application.Debug("process variable data:", ToJSON(parameter))
	}
	variableData := variable.GetVariableData(parameter.Name)
	if variableData == nil {
		err = newErrorParamIsNull("variable [" + parameter.Name + "] not defind")
		return
	}

	data := variableData.Data

	dataType := GetDataType(parameter.DataType)
	var dataStruct *StructModel

	// 数据类型为空 表示结构体
	if dataType == nil {
		dataStruct = application.context.GetStruct(parameter.DataType)
		if dataStruct == nil {
			err = newErrorStructIsNull("struct [" + parameter.DataType + "] not defind")
			return
		}
	}

	// 根据类型转换位新数据
	variableData = &VariableData{
		Name:       parameter.Name,
		Data:       data,
		DataType:   dataType,
		DataStruct: dataStruct,
	}

	if application.OutDebug() {
		application.Debug("process variable [", variableData.Name, "] data:", ToJSON(variableData.Data))
	}
	variable.AddVariableData(variableData)

	return
}

func newCallInvokeVariable(application *Application, variable *invokeVariable, callParameters []string, targetParameters []*VariableModel) (callVariable *invokeVariable, err error) {
	// 调用外部Model 需要重置invokeVariable
	callVariable = variable.Clone()

	// 根据Call传参 解析当前变量载体应该传输哪些

	//
	if len(callParameters) == 0 {
		if len(targetParameters) > 0 {
			for _, targetParameter := range targetParameters {
				callVariableData := variable.GetVariableData(targetParameter.Name)
				if callVariableData != nil {
					callParameters = append(callParameters, targetParameter.Name)
				}
			}
		}
	}
	if len(callParameters) > 0 {
		for _, callParameterName := range callParameters {
			callVariableData := variable.GetVariableData(callParameterName)
			if callVariableData == nil {
				err = newErrorParamIsNull("call parameter [" + callParameterName + "] not defind")
				return
			}
			callVariable.AddVariableData(callVariableData)
		}
	}
	return
}
