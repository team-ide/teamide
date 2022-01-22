package model

type StructModel struct {
	Name     string              `json:"name,omitempty" yaml:"name,omitempty"`         // 名称，同一个应用中唯一
	Comment  string              `json:"comment,omitempty" yaml:"comment,omitempty"`   // 注释
	Table    string              `json:"table,omitempty" yaml:"table,omitempty"`       // 表名
	Database string              `json:"database,omitempty" yaml:"database,omitempty"` // 库名
	Parent   string              `json:"parent,omitempty" yaml:"parent,omitempty"`     // 父结构体
	Fields   []*StructFieldModel `json:"fields,omitempty" yaml:"fields,omitempty"`     // 结构体字段
	Indexs   []*StructIndexModel `json:"indexs,omitempty" yaml:"indexs,omitempty"`     // 结构体字段
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
	Name          string `json:"name,omitempty" yaml:"name,omitempty"`                   // 字段名称，同一个结构体中唯一
	Comment       string `json:"comment,omitempty" yaml:"comment,omitempty"`             // 注释
	Column        string `json:"column,omitempty" yaml:"column,omitempty"`               // 映射 数据库 字段 默认和字段名称一致
	JsonName      string `json:"jsonName,omitempty" yaml:"jsonName,omitempty"`           // 映射 JSON 字段 默认和字段名称一致
	JsonOmitempty bool   `json:"jsonOmitempty,omitempty" yaml:"jsonOmitempty,omitempty"` // 映射 JSON 字段 省略空值
	IsList        bool   `json:"isList,omitempty" yaml:"isList,omitempty"`               // 是否是列表
	DataType      string `json:"dataType,omitempty" yaml:"dataType,omitempty"`           // 数据类型
	ColumnType    string `json:"columnType,omitempty" yaml:"columnType,omitempty"`       // 字段类型
	ColumnLength  int    `json:"columnLength,omitempty" yaml:"columnLength,omitempty"`   // 字段长度
	ColumnDecimal int    `json:"columnDecimal,omitempty" yaml:"columnDecimal,omitempty"` // 字段小数点
	PrimaryKey    bool   `json:"primaryKey,omitempty" yaml:"primaryKey,omitempty"`       // 主键
	NotNull       bool   `json:"notNull,omitempty" yaml:"notNull,omitempty"`             // 不能为空
	Default       string `json:"default,omitempty" yaml:"default,omitempty"`             // 默认值
}

type StructIndexModel struct {
	Name    string `json:"name,omitempty" yaml:"name,omitempty"`       // 字段名称，同一个结构体中唯一
	Type    string `json:"type,omitempty" yaml:"type,omitempty"`       // 索引类型
	Columns string `json:"columns,omitempty" yaml:"columns,omitempty"` // 字段
	Comment string `json:"comment,omitempty" yaml:"comment,omitempty"` // 注释
}

func TextToStructModel(namePath string, text string) (model *StructModel, err error) {
	var name string
	model = &StructModel{}
	name, err = TextToModel(namePath, text, model)
	if err != nil {
		return
	}
	model.Name = name
	return
}
