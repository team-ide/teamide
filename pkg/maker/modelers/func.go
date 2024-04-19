package modelers

type FuncModel struct {
	ElementNode
	Comment    string      `json:"comment,omitempty"`    // 说明
	Note       string      `json:"note,omitempty"`       // 注释
	Args       []*ArgModel `json:"args,omitempty"`       // 入参
	Func       string      `json:"func,omitempty"`       // 函数内容
	ReturnType string      `json:"returnType,omitempty"` // 返回类型
}

func init() {
	addDocTemplate(&docTemplate{
		Name:    TypeFuncName,
		Comment: "函数",
		Fields: []*docTemplateField{
			{Name: "comment", Comment: "函数说明"},
			{Name: "note", Comment: "函数注释"},
			{Name: "args", Comment: "参数", IsList: true, StructName: docTemplateArgName},
			{Name: "func", Comment: "函数内容"},
			{Name: "returnType", Comment: "返回类型"},
		},
	})
}
