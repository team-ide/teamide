package model

type ModelType struct {
	Name   string           `json:"name,omitempty"`
	Text   string           `json:"text,omitempty"`
	Dir    string           `json:"dir,omitempty"`
	Fields []ModelTypeField `json:"fields,omitempty"`
}

type ModelTypeField struct {
	Name                 string           `json:"name,omitempty"`
	Text                 string           `json:"text,omitempty"`
	Comment              string           `json:"comment,omitempty"`
	Type                 string           `json:"type,omitempty"`
	Readonly             bool             `json:"readonly,omitempty"`
	Width                int              `json:"width,omitempty"`
	Fields               []ModelTypeField `json:"fields,omitempty"`
	IsNumber             bool             `json:"isNumber,omitempty"`
	IsList               bool             `json:"isList,omitempty"`
	IsDataTypeOption     bool             `json:"isDataTypeOption,omitempty"`
	IsStructOption       bool             `json:"isStructOption,omitempty"`
	IsIndexTypeOption    bool             `json:"isIndexTypeOption,omitempty"`
	IsActionOption       bool             `json:"isActionOption,omitempty"`
	IsDatabaseTypeOption bool             `json:"isDatabaseTypeOption,omitempty"`
	IsColumnTypeOption   bool             `json:"isColumnTypeOption,omitempty"`
	IfScript             string           `json:"ifScript,omitempty"`
}

var (
	MODEL_TYPES []*ModelType

	MODEL_TYPE_STRUCT = appendModelType(&ModelType{
		Name: "structs", Text: "结构体", Dir: "struct",
		Fields: []ModelTypeField{
			{Name: "name", Text: "名称", Readonly: true},
			{Name: "comment", Text: "注释"},
			{Name: "table", Text: "表"},
			{Name: "parent", Text: "父", Type: "select", IsStructOption: true},

			{Name: "fields", Text: "字段", IsList: true,
				Fields: []ModelTypeField{
					{Name: "name", Text: "名称"},
					{Name: "comment", Text: "注释"},
					{Name: "dataType", Text: "数据类型", Type: "select", IsDataTypeOption: true},
					{Name: "isList", Text: "是List", Type: "switch", Width: 40},
					{Name: "column", Text: "字段", IfScript: "tool.isNotEmpty(model.table)"},
					{Name: "columnType", Text: "字段类型", Type: "select", Width: 80, IsColumnTypeOption: true, IfScript: "tool.isNotEmpty(model.table)"},
					{Name: "columnLength", Text: "长度", Width: 60, IsNumber: true, IfScript: "tool.isNotEmpty(model.table)"},
					{Name: "columnDecimal", Text: "小数", Width: 60, IsNumber: true, IfScript: "tool.isNotEmpty(model.table)"},
					{Name: "primaryKey", Text: "主键", Type: "switch", Width: 60, IfScript: "tool.isNotEmpty(model.table)"},
					{Name: "notNull", Text: "不为空", Type: "switch", Width: 60, IfScript: "tool.isNotEmpty(model.table)"},
					{Name: "default", Text: "默认", Width: 60, IfScript: "tool.isNotEmpty(model.table)"},
				},
			},

			{Name: "indexs", Text: "索引", IsList: true, IfScript: "tool.isNotEmpty(model.table)",
				Fields: []ModelTypeField{
					{Name: "name", Text: "名称"},
					{Name: "comment", Text: "注释"},
					{Name: "type", Text: "类型", Type: "select", Width: 80, IsIndexTypeOption: true},
					{Name: "columns", Text: "字段"},
				},
			},
		},
	})

	MODEL_TYPE_ACTION = appendModelType(&ModelType{
		Name: "actions", Text: "服务接口", Dir: "action",
		Fields: []ModelTypeField{
			{Name: "name", Text: "名称", Readonly: true},
			{Name: "comment", Text: "注释"},
		},
	})

	MODEL_TYPE_CONSTANT = appendModelType(&ModelType{
		Name: "constants", Text: "常量", Dir: "constant",
		Fields: []ModelTypeField{
			{Name: "name", Text: "名称", Readonly: true},
			{Name: "comment", Text: "注释"},
			{Name: "dataType", Text: "数据类型", Type: "select", IsDataTypeOption: true},
			{Name: "value", Text: "值"},
			{Name: "environmentVariable", Text: "环境变量", Comment: "优先取环境变量中的值"},
		},
	})

	MODEL_TYPE_ERROR = appendModelType(&ModelType{
		Name: "errors", Text: "错误码", Dir: "error",
		Fields: []ModelTypeField{
			{Name: "name", Text: "名称", Readonly: true},
			{Name: "comment", Text: "注释"},
			{Name: "code", Text: "错误码"},
			{Name: "msg", Text: "错误信息"},
		},
	})

	MODEL_TYPE_TEST = appendModelType(&ModelType{
		Name: "tests", Text: "测试", Dir: "test",
		Fields: []ModelTypeField{
			{Name: "name", Text: "名称", Readonly: true},
			{Name: "comment", Text: "注释"},
		},
	})

	MODEL_TYPE_DICTIONARY = appendModelType(&ModelType{
		Name: "dictionaries", Text: "数据字典", Dir: "dictionary",
		Fields: []ModelTypeField{
			{Name: "name", Text: "名称", Readonly: true},
			{Name: "comment", Text: "注释"},

			{Name: "options", Text: "选项", IsList: true,
				Fields: []ModelTypeField{
					{Name: "text", Text: "文案"},
					{Name: "value", Text: "值"},
					{Name: "comment", Text: "注释"},
				},
			},
		},
	})

	MODEL_TYPE_SERVER_WEB = appendModelType(&ModelType{
		Name: "serverWebs", Text: "数据字典", Dir: "server/web",
		Fields: []ModelTypeField{
			{Name: "name", Text: "名称", Readonly: true},
			{Name: "comment", Text: "注释"},
			{Name: "host", Text: "Host"},
			{Name: "port", Text: "Port", IsNumber: true},
			{Name: "contextPath", Text: "ContextPath"},

			{Name: "token", Text: "Token",
				Fields: []ModelTypeField{
					{Name: "include", Text: "验证路径"},
					{Name: "exclude", Text: "忽略路径"},
					{Name: "createAction", Text: "创建Token操作", Type: "select", IsActionOption: true},
					{Name: "validateAction", Text: "验证Token操作", Type: "select", IsActionOption: true},
					{Name: "variableName", Text: "变量名称"},
					{Name: "variableDataType", Text: "变量数据类型", Type: "select", IsDataTypeOption: true},
				},
			},
		},
	})

	MODEL_TYPE_DATASOURCE_DATABASE = appendModelType(&ModelType{
		Name: "datasourceDatabases", Text: "Database数据源", Dir: "datasource/database",
		Fields: []ModelTypeField{
			{Name: "name", Text: "名称", Readonly: true},
			{Name: "comment", Text: "注释"},
			{Name: "type", Text: "类型", Type: "select", IsDatabaseTypeOption: true},
			{Name: "host", Text: "Host"},
			{Name: "port", Text: "Port", IsNumber: true},
			{Name: "database", Text: "Database"},
			{Name: "username", Text: "Username"},
			{Name: "password", Text: "Password"},
		},
	})

	MODEL_TYPE_DATASOURCE_REDIS = appendModelType(&ModelType{
		Name: "datasourceRedises", Text: "Redis数据源", Dir: "datasource/redis",
		Fields: []ModelTypeField{
			{Name: "name", Text: "名称", Readonly: true},
			{Name: "comment", Text: "注释"},
			{Name: "address", Text: "Redis地址"},
			{Name: "auth", Text: "密码"},
			{Name: "prefix", Text: "前缀", Comment: "如果配置，所有key将自动拼接该前缀"},
		},
	})

	MODEL_TYPE_DATASOURCE_KAFKA = appendModelType(&ModelType{
		Name: "datasourceKafkas", Text: "Kafka数据源", Dir: "datasource/kafka",
		Fields: []ModelTypeField{
			{Name: "name", Text: "名称", Readonly: true},
			{Name: "comment", Text: "注释"},
			{Name: "address", Text: "Kafka地址"},
			{Name: "prefix", Text: "前缀", Comment: "如果配置，所有topic将自动拼接该前缀"},
		},
	})

	MODEL_TYPE_DATASOURCE_ZOOKEEPER = appendModelType(&ModelType{
		Name: "datasourceZookeepers", Text: "Zookeeper数据源", Dir: "datasource/zookeeper",
		Fields: []ModelTypeField{
			{Name: "name", Text: "名称", Readonly: true},
			{Name: "comment", Text: "注释"},
			{Name: "address", Text: "Zookeeper地址"},
			{Name: "namespace", Text: "命名空间", Comment: "如果配置，则所有路径将放在该命名空间下"},
		},
	})
)

func appendModelType(modelType *ModelType) *ModelType {
	MODEL_TYPES = append(MODEL_TYPES, modelType)
	return modelType
}

func GetModelType(name string) *ModelType {
	for _, one := range MODEL_TYPES {
		if one.Name == name {
			return one
		}
	}
	return nil
}
