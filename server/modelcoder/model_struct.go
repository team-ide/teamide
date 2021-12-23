package modelcoder

type StructModel struct {
	Name       string              `json:"name,omitempty"`       // 名称，同一个应用中唯一
	TableName  string              `json:"tableName,omitempty"`  // 表名
	ParentName string              `json:"parentName,omitempty"` // 父结构体
	Fields     []*StructFieldModel `json:"fields,omitempty"`     // 结构体字段
}

func (this_ *StructModel) GetField(name string) *StructFieldModel {
	if len(this_.Fields) == 0 {
		return nil
	}
	for _, one := range this_.Fields {
		if one.Name == name {
			return one
		}
	}
	return nil
}

type StructFieldModel struct {
	Name          string `json:"name,omitempty"`          // 字段名称，同一个结构体中唯一
	Comment       string `json:"comment,omitempty"`       // 注释
	ColumnName    string `json:"columnName,omitempty"`    // 映射 数据库 字段 默认和字段名称一致
	JsonName      string `json:"jsonName,omitempty"`      // 映射 JSON 字段 默认和字段名称一致
	JsonOmitempty bool   `json:"jsonOmitempty,omitempty"` // 映射 JSON 字段 省略空值
	IsList        string `json:"isList,omitempty"`        // 是否是列表
	DataType      string `json:"dataType,omitempty"`      // 数据类型
}
