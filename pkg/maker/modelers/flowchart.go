package modelers

type FlowchartModel struct {
	ElementNode
	Comment string `json:"comment,omitempty"` // 说明
	Note    string `json:"note,omitempty"`    // 注释
	Content string `json:"content,omitempty"` // 内容
}

func init() {
	addDocTemplate(&docTemplate{
		Name:    TypeFlowchartName,
		Comment: "流程图",
		Fields: []*docTemplateField{
			{Name: "comment", Comment: "流程图说明"},
			{Name: "note", Comment: "流程图源码注释"},
			{Name: "content", Comment: "内容"},
		},
	})
}
