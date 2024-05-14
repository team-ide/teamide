package modelers

type ServiceModel struct {
	ElementNode
	Comment string      `json:"comment,omitempty"` // 说明
	Note    string      `json:"note,omitempty"`    // 注释
	Args    []*ArgModel `json:"args,omitempty"`    //入参
	Func    string      `json:"func,omitempty"`    // 函数内容
	//Steps   []interface{} `json:"steps,omitempty"`   // 阶段
	Return string `json:"return,omitempty"` // 返回
}

type ArgModel struct {
	Name    string `json:"name,omitempty"`    // 名称，同一个应用中唯一
	Comment string `json:"comment,omitempty"` // 说明
	Note    string `json:"note,omitempty"`    // 注释
	Type    string `json:"type,omitempty"`    // 类型
}

var (
	docTemplateArgName = "arg"
)

func init() {
	addDocTemplate(&docTemplate{
		Name:    TypeServiceName,
		Comment: "服务文件，该文件用于动作处理，如数据库、redis、文件等地方",
		Fields: []*docTemplateField{
			{Name: "comment", Comment: "结构体说明"},
			{Name: "note", Comment: "结构体源码注释"},
			{Name: "args", Comment: "参数", IsList: true, StructName: docTemplateArgName},
			{Name: "func", Comment: "函数内容"},
			//{Name: "steps", Comment: "阶段", IsList: true, StructName: docTemplateStepName},
			{Name: "return", Comment: "返回值"},
		},
	})
}

func init() {
	addDocTemplate(&docTemplate{
		Name:         docTemplateArgName,
		Abbreviation: "name",
		Fields: []*docTemplateField{
			{Name: "name", Comment: "参数名称"},
			{Name: "comment", Comment: "参数说明"},
			{Name: "note", Comment: "参数源码注释"},
			{Name: "type", Comment: "参数类型"},
		},
		newModel: func() interface{} {
			return &ArgModel{}
		},
		newModels: func() interface{} {
			var vs []*ArgModel
			return vs
		},
		appendModel: func(values interface{}, value interface{}) (res interface{}) {
			vs := values.([]*ArgModel)
			vs = append(vs, value.(*ArgModel))
			return vs
		},
	})
}
