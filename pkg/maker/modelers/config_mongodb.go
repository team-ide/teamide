package modelers

type ConfigMongodbModel struct {
	Name     string `json:"name,omitempty"`     // 名称，同一个应用中唯一
	Comment  string `json:"comment,omitempty"`  // 说明
	Note     string `json:"note,omitempty"`     // 注释
	Address  string `json:"address,omitempty"`  //
	Username string `json:"username,omitempty"` //
	Password string `json:"password,omitempty"` //

}

var (
	docTemplateConfigMongodbName = "configMongodb"
)

func init() {
	addDocTemplate(&docTemplate{
		Name:    docTemplateConfigMongodbName,
		Comment: "Mongodb配置",
		Fields: []*docTemplateField{
			{Name: "name", Comment: "配置名称"},
			{Name: "comment", Comment: "配置说明"},
			{Name: "note", Comment: "配置源码注释"},
			{Name: "address", Comment: ""},
			{Name: "username", Comment: ""},
			{Name: "password", Comment: ""},
		},
	})
}
