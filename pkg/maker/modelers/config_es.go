package modelers

type ConfigEsModel struct {
	ElementNode
	Comment  string `json:"comment,omitempty"`  // 说明
	Note     string `json:"note,omitempty"`     // 注释
	Url      string `json:"url,omitempty"`      //
	Username string `json:"username,omitempty"` //
	Password string `json:"password,omitempty"` //

}

func init() {
	addDocTemplate(&docTemplate{
		Name:    TypeConfigEsName,
		Comment: "Elasticsearch配置",
		Fields: []*docTemplateField{
			{Name: "comment", Comment: "配置说明"},
			{Name: "note", Comment: "配置源码注释"},
			{Name: "url", Comment: ""},
			{Name: "username", Comment: ""},
			{Name: "password", Comment: ""},
		},
	})
}
