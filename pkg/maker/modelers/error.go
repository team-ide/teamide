package modelers

type ErrorModel struct {
	Name    string              `json:"name,omitempty"`    // 字段名称，同一个结构体中唯一
	Comment string              `json:"comment,omitempty"` // 说明
	Note    string              `json:"note,omitempty"`    // 注释
	Options []*ErrorOptionModel `json:"options,omitempty"`
}

type ErrorOptionModel struct {
	Name    string `json:"name,omitempty"`    // 字段名称，同一个结构体中唯一
	Comment string `json:"comment,omitempty"` // 说明
	Note    string `json:"note,omitempty"`    // 注释
	Code    string `json:"code,omitempty"`    // 错误码
	Msg     string `json:"msg,omitempty"`     // 错误信息
}

var (
	docTemplateErrorOptionName = "error/option"
)

func init() {
	addDocTemplate(&docTemplate{
		Name:    TypeErrorName,
		Comment: "异常配置",
		Fields: []*docTemplateField{
			{Name: "name", Comment: "配置名称"},
			{Name: "comment", Comment: "配置说明"},
			{Name: "note", Comment: "配置源码注释"},
			{Name: "options", Comment: "异常项", IsList: true, StructName: docTemplateErrorOptionName},
		},
	})

	addDocTemplate(&docTemplate{
		Name:    docTemplateErrorOptionName,
		Comment: "异常项配置",
		Fields: []*docTemplateField{
			{Name: "name", Comment: "配置名称"},
			{Name: "comment", Comment: "配置说明"},
			{Name: "note", Comment: "配置源码注释"},
			{Name: "code", Comment: "错误码"},
			{Name: "msg", Comment: "错误信息"},
		},
		newModel: func() interface{} {
			return &ErrorOptionModel{}
		},
		newModels: func() interface{} {
			var vs []*ErrorOptionModel
			return vs
		},
		appendModel: func(values interface{}, value interface{}) (res interface{}) {
			vs := values.([]*ErrorOptionModel)
			vs = append(vs, value.(*ErrorOptionModel))
			return vs
		},
	})
}
