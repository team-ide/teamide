package modelers

type ConstantModel struct {
	ElementNode
	Comment string                 `json:"comment,omitempty"` // 说明
	Note    string                 `json:"note,omitempty"`    // 注释
	Options []*ConstantOptionModel `json:"options,omitempty"`
}

type ConstantOptionModel struct {
	Name    string `json:"name,omitempty"`    // 字段名称，同一个结构体中唯一
	Comment string `json:"comment,omitempty"` // 说明
	Note    string `json:"note,omitempty"`    // 注释
	Type    string `json:"type,omitempty"`    // 类型
	Value   string `json:"value,omitempty"`   // 值
}

var (
	docTemplateConstantOptionName = "constant/option"
)

func init() {
	addDocTemplate(&docTemplate{
		Name:    TypeConstantName,
		Comment: "常量配置",
		Fields: []*docTemplateField{
			{Name: "comment", Comment: "配置说明"},
			{Name: "note", Comment: "配置源码注释"},
			{Name: "options", Comment: "配置项", IsList: true, StructName: docTemplateConstantOptionName},
		},
	})

	addDocTemplate(&docTemplate{
		Name:    docTemplateConstantOptionName,
		Comment: "常量项配置",
		Fields: []*docTemplateField{
			{Name: "name", Comment: "配置名称"},
			{Name: "comment", Comment: "配置说明"},
			{Name: "note", Comment: "配置源码注释"},
			{Name: "type", Comment: "类型"},
			{Name: "value", Comment: "值"},
		},
		newModel: func() interface{} {
			return &ConstantOptionModel{}
		},
		newModels: func() interface{} {
			var vs []*ConstantOptionModel
			return vs
		},
		appendModel: func(values interface{}, value interface{}) (res interface{}) {
			vs := values.([]*ConstantOptionModel)
			vs = append(vs, value.(*ConstantOptionModel))
			return vs
		},
	})
}
