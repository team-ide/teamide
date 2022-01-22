package invoke

import (
	"teamide/application/base"
	"teamide/application/common"
	"teamide/application/model"

	"github.com/wxnacy/wgo/arrays"
)

type ParseInfo struct {
	App             common.IApplication     `json:"-"`
	InvokeNamespace *common.InvokeNamespace `json:"-"`
	ParameterList   []string                `json:"parameterList"`
	UseFunctions    []string                `json:"useFunctions"`
}

var (
	parseCallMap map[string]func(parseInfo *ParseInfo, prefixName string, args []interface{}) (err error)
)

func init() {
	parseCallMap = make(map[string]func(parseInfo *ParseInfo, prefixName string, args []interface{}) (err error))
	parseCallMap["addDataInfo"] = parseAddDataInfo
	parseCallMap["service"] = parseService
}

func parseAddDataInfo(parseInfo *ParseInfo, _ string, args []interface{}) (err error) {
	name := args[0].(string)
	var dataType string = args[1].(string)
	var comment string
	var value string
	var isList bool
	var isPage bool
	if len(args) > 2 {
		comment = args[2].(string)
	}
	if len(args) > 3 {
		value = args[3].(string)
	}
	if len(args) > 4 {
		isList = args[4].(bool)
	}
	if len(args) > 5 {
		isPage = args[5].(bool)
	}
	if base.IsEmpty(dataType) && base.IsNotEmpty(value) {
		var dataInfo *common.InvokeDataInfo
		dataInfo, _ = parseInfo.InvokeNamespace.GetDataInfo(value)
		if dataInfo != nil {
			// fmt.Println("add alias [", name, "] for value [", value, "]")
			if arrays.ContainsString(dataInfo.Alias, name) < 0 {
				dataInfo.Alias = append(dataInfo.Alias, name)
			}
			return
		}
	}
	variable := &model.VariableModel{
		Name:     name,
		Comment:  comment,
		Value:    value,
		DataType: dataType,
		IsList:   isList,
		IsPage:   isPage,
	}
	err = parseInfo.InvokeNamespace.SetDataInfo(variable)
	if err != nil {
		return
	}
	return
}

func parseService(parseInfo *ParseInfo, prefixName string, args []interface{}) (err error) {
	callServiceName := args[0].(string)
	callService := parseInfo.App.GetContext().GetService(callServiceName)
	if callService == nil {
		err = base.NewErrorServiceIsNull("call service [", callServiceName, "] not defind")
		return
	}
	// for index, callVariable := range callService.InVariables {
	// 	value := args[index+1].(string)
	// 	if base.IsNotEmpty(value) && value != callVariable.Name {
	// 		fmt.Println("callService [", callServiceName, "]")
	// 		fmt.Println("call arg [", callVariable.Name, "] use [", value, "]")
	// 		var dataInfo *common.InvokeDataInfo
	// 		dataInfo, err = parseInfo.InvokeNamespace.GetDataInfo( value)
	// 		if err != nil {
	// 			return
	// 		}
	// 		if dataInfo == nil {
	// 			continue
	// 		}
	// 		callDataInfo := dataInfo
	// 		if arrays.ContainsString(callDataInfo.Alias, callVariable.Name) < 0 {
	// 			callDataInfo.Alias = append(callDataInfo.Alias, callVariable.Name)
	// 		}
	// 	}
	// }

	var javascript string
	javascript, err = common.GetServiceJavascriptByService(parseInfo.App, callService)
	if err != nil {
		return
	}
	// if callService.Name == "user/batchInsert" {
	// 	fmt.Println(javascript)
	// }
	functionParser := NewFunctionParser(javascript)
	err = functionParser.Parse(parseInfo)
	if err != nil {
		return
	}
	return
}
