package modelers

type StructModel struct {
	ElementNode
	Comment string         `json:"comment,omitempty"` // 说明
	Note    string         `json:"note,omitempty"`    // 注释
	Parent  string         `json:"parent,omitempty"`  // 父结构体
	Fields  []*StructField `json:"fields,omitempty"`  // 结构体字段
}

type StructField struct {
	Name          string `json:"name,omitempty"`          // 字段名称，同一个结构体中唯一
	Comment       string `json:"comment,omitempty"`       // 说明
	Note          string `json:"note,omitempty"`          // 注释
	JsonName      string `json:"jsonName,omitempty"`      // 映射 JSON 字段 默认和字段名称一致
	JsonOmitempty bool   `json:"jsonOmitempty,omitempty"` // 映射 JSON 字段 省略空值
	IsList        bool   `json:"isList,omitempty"`        // 是否是列表
	Type          string `json:"type,omitempty"`          // 数据类型
	Default       string `json:"default,omitempty"`       // 默认值
	Column        string `json:"column,omitempty"`        // 默认值
}

var (
	docTemplateStructFieldName = "struct_field"
)

func init() {
	addDocTemplate(&docTemplate{
		Name:    TypeStruct.Name,
		Comment: "结构体文件，该文件用于入参、出参、函数调用、数据存储等地方",
		Fields: []*docTemplateField{
			{Name: "comment", Comment: "结构体说明"},
			{Name: "note", Comment: "结构体源码注释"},
			{Name: "parent", Comment: "父级结构体，源码将继承该结构体"},
			{Name: "fields", Comment: "这是结构体字段", IsList: true, StructName: docTemplateStructFieldName},
		},
		newModel: func() interface{} {
			return &StructModel{}
		},
		newModels: func() interface{} {
			var vs []*StructModel
			return vs
		},
		appendModel: func(values interface{}, value interface{}) (res interface{}) {
			vs := values.([]*StructModel)
			vs = append(vs, value.(*StructModel))
			return vs
		},
	})
	addDocTemplate(&docTemplate{
		Comment:      "结构体字段",
		Abbreviation: "name",
		Name:         docTemplateStructFieldName,
		Fields: []*docTemplateField{
			{Name: "name", Comment: "字段名称"},
			{Name: "comment", Comment: "字段说明"},
			{Name: "note", Comment: "字段源码注释"},
			{Name: "type", Comment: "字段类型", Default: "string"},
			{Name: "jsonName", Comment: "序列化JSON名称"},
			{Name: "jsonOmitempty", Comment: "序列化JSON，省略空值"},
			{Name: "isList", Comment: "是集合"},
			{Name: "default", Comment: "创建对象该字段默认的值"},
			{Name: "column", Comment: "字段"},
		},
		newModel: func() interface{} {
			return &StructField{}
		},
		newModels: func() interface{} {
			var vs []*StructField
			return vs
		},
		appendModel: func(values interface{}, value interface{}) (res interface{}) {
			vs := values.([]*StructField)
			vs = append(vs, value.(*StructField))
			return vs
		},
	})
}
