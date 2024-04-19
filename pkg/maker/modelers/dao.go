package modelers

type DaoModel struct {
	ElementNode
	Comment string        `json:"comment,omitempty"` // 说明
	Note    string        `json:"note,omitempty"`    // 注释
	Args    []*ArgModel   `json:"args,omitempty"`    //入参
	Steps   []interface{} `json:"steps,omitempty"`   // 阶段
	Return  string        `json:"return,omitempty"`  // 返回
}

func init() {
	addDocTemplate(&docTemplate{
		Name:    TypeDaoName,
		Comment: "数据层，处理数据库等落地数据",
		Fields: []*docTemplateField{
			{Name: "comment", Comment: "数据层说明"},
			{Name: "note", Comment: "数据层源码注释"},
			{Name: "args", Comment: "参数", IsList: true, StructName: docTemplateArgName},
			{Name: "steps", Comment: "阶段", IsList: true, StructName: docTemplateStepName},
			{Name: "return", Comment: "返回值"},
		},
	})
}
