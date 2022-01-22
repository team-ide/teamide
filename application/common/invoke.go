package common

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"teamide/application/base"
	"teamide/application/model"

	"github.com/gin-gonic/gin"
	clone "github.com/huandu/go-clone"
	"github.com/wxnacy/wgo/arrays"
)

type InvokeNamespace struct {
	App            IApplication           `json:"-"`
	Datas          []*InvokeData          `json:"datas,omitempty"`
	DataInfos      []*InvokeDataInfo      `json:"dataInfos,omitempty"`
	TempMap        map[string]*InvokeData `json:"tempMap,omitempty"`
	RequestContext *gin.Context           `json:"-"`
	RequestBody    interface{}            `json:"requestBody,omitempty"`
	ServerWebToken *model.ServerWebToken  `json:"-"`
}

func NewInvokeNamespace(app IApplication) (res *InvokeNamespace, err error) {
	res = &InvokeNamespace{App: app}
	err = res.init()
	if err != nil {
		return
	}
	return
}

func (this_ *InvokeNamespace) init() (err error) {
	for _, one := range this_.App.GetContext().Constants {
		dataType := one.DataType
		if base.IsEmpty(dataType) {
			dataType = "string"
		}
		err = this_.SetDataInfo(&model.VariableModel{
			Name:     "$" + one.Name,
			DataType: dataType,
		})
		if err != nil {
			return
		}
		var stringValue string = one.Value
		if base.IsNotEmpty(one.EnvironmentVariable) {
			environmentVariableValue := os.Getenv(one.EnvironmentVariable)
			if base.IsNotEmpty(environmentVariableValue) {
				stringValue = environmentVariableValue
			}
		}
		var value interface{} = stringValue
		dataType_ := this_.App.GetContext().GetVariableDataType(dataType)
		switch dataType_ {
		case model.DATA_TYPE_BOOLEAN:
			value = this_.App.GetScript().IsTrue(stringValue)
		case model.DATA_TYPE_BYTE:
			var num int
			num, err = strconv.Atoi(stringValue)
			if err != nil {
				return
			}
			value = int8(num)
		case model.DATA_TYPE_FLOAT:
			var num int
			num, err = strconv.Atoi(stringValue)
			if err != nil {
				return
			}
			value = int16(num)
		case model.DATA_TYPE_INT:
			var num int
			num, err = strconv.Atoi(stringValue)
			if err != nil {
				return
			}
			value = int(num)
		case model.DATA_TYPE_LONG:
			var num int
			num, err = strconv.Atoi(stringValue)
			if err != nil {
				return
			}
			value = int64(num)
		case model.DATA_TYPE_MAP:
			value, err = this_.App.GetScript().JSONToData(stringValue)
			if err != nil {
				return
			}
		}
		err = this_.SetData("$"+one.Name, value, nil)
		if err != nil {
			return
		}
	}
	return
}

type InvokeData struct {
	Name      string        `json:"name,omitempty"`
	DotName   string        `json:"dotName,omitempty"`
	Datas     []*InvokeData `json:"datas,omitempty"`
	Parent    *InvokeData   `json:"-"`
	Value     interface{}   `json:"value,omitempty"`
	ListIndex int           `json:"listIndex,omitempty"`
}

type InvokeDataInfo struct {
	Name       string                  `json:"name,omitempty"`
	Comment    string                  `json:"comment,omitempty"`
	Alias      []string                `json:"alias,omitempty"`
	Value      string                  `json:"value,omitempty"`
	DotName    string                  `json:"dotName,omitempty"`
	DataInfos  []*InvokeDataInfo       `json:"dataInfos,omitempty"`
	Parent     *InvokeDataInfo         `json:"-"`
	DataType   *model.VariableDataType `json:"-"`
	IsList     bool                    `json:"isList,omitempty"` // 是否是列表
	IsPage     bool                    `json:"isPage,omitempty"` // 是否是列表
	Validatas  []*model.ValidateModel  `json:"-"`
	IsUse      bool                    `json:"isUse,omitempty"`      // 是否 使用
	IsSetValue bool                    `json:"isSetValue,omitempty"` // 是否 设值
}

func CloneDataInfo(dataInfo *InvokeDataInfo) (res *InvokeDataInfo) {
	Parent_ := dataInfo.Parent
	dataInfo.Parent = nil

	res = clone.Clone(dataInfo).(*InvokeDataInfo)

	dataInfo.Parent = Parent_

	return
}
func CloneData(data *InvokeData) (res *InvokeData) {
	Parent_ := data.Parent
	data.Parent = nil
	Datas_ := data.Datas
	data.Datas = nil

	res = clone.Clone(data).(*InvokeData)

	data.Parent = Parent_
	data.Datas = Datas_

	res.Datas = []*InvokeData{}
	for _, one := range data.Datas {
		res.Datas = append(res.Datas, CloneData(one))
	}
	return
}

func (this_ *InvokeNamespace) GetData(name string) (res *InvokeData, err error) {
	if strings.HasPrefix(name, "$invoke_temp.") {
		if this_.TempMap == nil {
			this_.TempMap = make(map[string]*InvokeData)
		}
		res = this_.TempMap[name]
		if res == nil {
			res = &InvokeData{
				Name:    name,
				DotName: name,
			}
			this_.TempMap[name] = res
		}
		return
	}

	res, err = GetData(this_.App, &this_.Datas, name)
	return
}

func (this_ *InvokeNamespace) SetData(name string, value interface{}, dataInfo *InvokeDataInfo) (err error) {
	if strings.HasPrefix(name, "$invoke_temp.") {
		if this_.TempMap == nil {
			this_.TempMap = make(map[string]*InvokeData)
		}
		data := this_.TempMap[name]
		if data == nil {
			data = &InvokeData{
				Name:    name,
				DotName: name,
			}
			this_.TempMap[name] = data
		}
		data.Value = value
		return
	}
	if dataInfo == nil {
		dataInfo, err = this_.GetDataInfo(name)
		if err != nil {
			return
		}
	}
	err = SetData(this_.App, this_, nil, &this_.Datas, name, value, dataInfo)
	return
}

func (this_ *InvokeNamespace) GetDataInfo(name string) (res *InvokeDataInfo, err error) {
	if strings.HasPrefix(name, "$invoke_temp.") {
		res = &InvokeDataInfo{
			Name:    name,
			DotName: name,
		}
		return
	}

	res, err = GetDataInfo(this_.App, this_, &this_.DataInfos, name)
	return
}

func (this_ *InvokeNamespace) SetDataInfo(variable *model.VariableModel) (err error) {
	if strings.HasPrefix(variable.Name, "$invoke_temp.") {

		return
	}
	err = SetDataInfo(this_.App, this_, nil, &this_.DataInfos, variable)
	return
}

func GetFieldNameAndSubName(name string) (fieldName string, fieldIndex int, subName string) {
	fieldIndex = -1
	if strings.Contains(name, ".") {
		fieldName = name[:strings.Index(name, ".")]

		subName = name[strings.Index(name, ".")+1:]
	} else {
		fieldName = name
	}
	if strings.Contains(fieldName, "[") {
		fieldIndexStr := fieldName[strings.Index(fieldName, "[")+1 : strings.Index(fieldName, "]")]
		fieldIndex, _ = strconv.Atoi(fieldIndexStr)
		fieldName = fieldName[:strings.Index(fieldName, "[")]
	}
	return
}

func GetData(app IApplication, datas *[]*InvokeData, name string) (res *InvokeData, err error) {
	fieldName, fieldIndex, subName := GetFieldNameAndSubName(name)
	var findOne *InvokeData
	for _, one := range *datas {

		if one.Name == fieldName {
			if fieldIndex >= 0 {
				if one.ListIndex != fieldIndex {
					continue
				}
			}
			findOne = one
			break
		}
	}
	if base.IsNotEmpty(subName) {
		if findOne == nil {
			err = base.NewError("", "get data [", fieldName, "] list index [", fieldIndex, "] not defind")
			return
		}
		res, err = GetData(app, &findOne.Datas, subName)
	} else {
		res = findOne
	}
	return
}

func SetData(app IApplication, invokeNamespace *InvokeNamespace, parent *InvokeData, datas *[]*InvokeData, name string, value interface{}, dataInfo *InvokeDataInfo) (err error) {

	fieldName, fieldIndex, subName := GetFieldNameAndSubName(name)
	var findOne *InvokeData
	for _, one := range *datas {
		if one.Name == fieldName && one.ListIndex == fieldIndex {
			findOne = one
			break
		}
	}

	if findOne == nil {
		if fieldIndex >= 0 {
			var listRootOne *InvokeData
			for _, one := range *datas {
				if one.Name == fieldName && one.ListIndex == -1 {
					listRootOne = one
					break
				}
			}
			if listRootOne != nil {

				findOne_ := CloneData(listRootOne)
				findOne_.ListIndex = fieldIndex

				if base.IsEmpty(subName) && value != nil {
					findOne_.Value = value
					for _, one := range findOne_.Datas {
						one.Value = value.(map[string]interface{})[one.Name]
					}
				} else {
					findOne_.Value = map[string]interface{}{}
					for _, one := range findOne_.Datas {
						one.Value = nil
					}
				}

				if listRootOne.Value == nil {
					listRootOne.Value = []map[string]interface{}{}
				}
				listRootOne.Value = append(listRootOne.Value.([]map[string]interface{}), findOne_.Value.(map[string]interface{}))

				*datas = append(*datas, findOne_)
				findOne = findOne_
			}
		}
	}
	// fmt.Println("SetData name [", name, "] value "+base.ToJSON(value)+"")
	if base.IsNotEmpty(subName) {
		if findOne == nil {
			err = base.NewError("", "data [", fieldName, "] not defind")
			return
		}
		err = SetData(app, invokeNamespace, findOne, &findOne.Datas, subName, value, dataInfo)
	} else {
		if parent != nil {
			var parentValue interface{}

			if parent.Value == nil {
				parent.Value = map[string]interface{}{}
			}
			parentValue = parent.Value

			switch m := parentValue.(type) {
			case map[string]interface{}:
				m[fieldName] = value
			default:
				err = base.NewError("", "invoke data [", parent.DotName, "] value can not to map[string]interface{}")
				return
			}
		}
		if findOne == nil {
			dotName := fieldName
			if parent != nil {
				dotName = parent.DotName + "." + dotName
			}
			findOne = &InvokeData{
				Name:      fieldName,
				DotName:   dotName,
				Parent:    parent,
				Datas:     []*InvokeData{},
				ListIndex: -1,
			}
			*datas = append(*datas, findOne)
		}
		if value == nil {
			if dataInfo != nil && dataInfo.DataType != nil {
				switch dataInfo.DataType {
				case model.DATA_TYPE_LONG, model.DATA_TYPE_INT, model.DATA_TYPE_SHORT, model.DATA_TYPE_BYTE:
					value = 0
				case model.DATA_TYPE_BOOLEAN:
					value = false
				case model.DATA_TYPE_DOUBLE, model.DATA_TYPE_FLOAT:
					value = 0.0
				case model.DATA_TYPE_MAP:
					// value = map[string]interface{}{}
				default:
					if dataInfo.DataType.DataStruct != nil {
						// value = map[string]interface{}{}
					} else {
						value = ""
					}
				}
			}
		}
		findOne.Value = value

		if dataInfo != nil && dataInfo.DataType != nil && dataInfo.DataType.DataStruct != nil {

			if !dataInfo.IsPage {
				for _, field := range dataInfo.DataType.DataStruct.Fields {
					if value != nil {
						switch m := value.(type) {
						case map[string]interface{}:
							fieldValue := m[field.Name]
							var fieldDataInfo *InvokeDataInfo
							fieldDataInfo, err = invokeNamespace.GetDataInfo(findOne.DotName + "." + field.Name)
							if err != nil {
								return
							}
							err = SetData(app, invokeNamespace, findOne, &findOne.Datas, field.Name, fieldValue, fieldDataInfo)
							if err != nil {
								return
							}
						default:
							// err = base.NewError("", "invoke data [", findOne.DotName, "] value can not to map[string]interface{}")
							// return
						}
					}
				}
			}
		}

		if value != nil && fieldIndex == -1 && base.IsEmpty(subName) {
			listMap, listMapOk := value.([]map[string]interface{})
			if listMapOk && len(listMap) > 0 {
				for index, one := range listMap {
					oneName := fieldName + "[" + fmt.Sprint(index) + "]"
					err = SetData(app, invokeNamespace, parent, datas, oneName, one, dataInfo)
					if err != nil {
						return
					}
				}
			}
		}
	}
	return
}

func GetDataInfo(app IApplication, invokeNamespace *InvokeNamespace, dataInfos *[]*InvokeDataInfo, name string) (res *InvokeDataInfo, err error) {
	fieldName, _, subName := GetFieldNameAndSubName(name)
	var findOne *InvokeDataInfo
	for _, one := range *dataInfos {
		if one.Name == fieldName || arrays.ContainsString(one.Alias, fieldName) >= 0 {
			findOne = one
			break
		}
	}
	if base.IsNotEmpty(subName) {
		if findOne == nil {
			err = base.NewError("", "get [", name, "] data info [", fieldName, "] not defind")
			return
		}
		res, err = GetDataInfo(app, invokeNamespace, &findOne.DataInfos, subName)
	} else {
		res = findOne
	}
	return
}

func SetDataInfo(app IApplication, invokeNamespace *InvokeNamespace, parent *InvokeDataInfo, dataInfos *[]*InvokeDataInfo, variable *model.VariableModel) (err error) {
	fieldName, _, subName := GetFieldNameAndSubName(variable.Name)
	var findOne *InvokeDataInfo
	for _, one := range *dataInfos {
		if one.Name == fieldName || arrays.ContainsString(one.Alias, fieldName) >= 0 {
			findOne = one
			break
		}
	}
	if base.IsNotEmpty(subName) {
		if findOne == nil {
			err = base.NewError("", "set [", variable.Name, "] data info [", fieldName, "] not defind")
			return
		}
		var subVariable = &model.VariableModel{
			Name:      subName,
			Value:     variable.Value,
			Comment:   variable.Comment,
			DataType:  variable.DataType,
			DataPlace: variable.DataPlace,
			IsList:    variable.IsList,
		}

		err = SetDataInfo(app, invokeNamespace, findOne, &findOne.DataInfos, subVariable)
	} else {
		if findOne == nil {
			dotName := fieldName
			if parent != nil {
				dotName = parent.DotName + "." + dotName
			}
			findOne = &InvokeDataInfo{
				Name:      fieldName,
				DotName:   dotName,
				Parent:    parent,
				Value:     variable.Value,
				Comment:   variable.Comment,
				IsList:    variable.IsList,
				IsPage:    variable.IsPage,
				DataInfos: []*InvokeDataInfo{},
			}
			*dataInfos = append(*dataInfos, findOne)
		}
		if findOne.IsPage {
			pageDataType := app.GetContext().GetVariableDataType("pageInfo")
			for _, field := range pageDataType.DataStruct.Fields {
				var subVariable = &model.VariableModel{
					Name:     field.Name,
					Comment:  field.Comment,
					DataType: field.DataType,
					IsList:   field.IsList,
				}

				if field.Name == "list" {
					subVariable.DataType = variable.DataType
				}
				err = SetDataInfo(app, invokeNamespace, findOne, &findOne.DataInfos, subVariable)
				if err != nil {
					return
				}
			}
		} else {
			dataType := app.GetContext().GetVariableDataType(variable.DataType)
			if dataType != nil {
				findOne.DataType = dataType
				if dataType.DataStruct != nil {
					for _, field := range dataType.DataStruct.Fields {
						var subVariable = &model.VariableModel{
							Name:     field.Name,
							Comment:  field.Comment,
							DataType: field.DataType,
							IsList:   field.IsList,
						}
						err = SetDataInfo(app, invokeNamespace, findOne, &findOne.DataInfos, subVariable)
						if err != nil {
							return
						}
					}
				}
			}
		}
	}
	return
}
