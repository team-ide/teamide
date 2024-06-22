package modelers

type AppModel struct {
	ElementNode
	Comment string                         `json:"comment,omitempty"`
	Db      map[string]*ConfigDbModel      `json:"db,omitempty"`
	Redis   map[string]*ConfigRedisModel   `json:"redis,omitempty"`
	Zk      map[string]*ConfigZkModel      `json:"zk,omitempty"`
	Kafka   map[string]*ConfigKafkaModel   `json:"kafka,omitempty"`
	Es      map[string]*ConfigEsModel      `json:"es,omitempty"`
	Mongodb map[string]*ConfigMongodbModel `json:"mongodb,omitempty"`
	Other   map[string]any                 `json:"other,omitempty"`
	Text    string                         `json:"-"`
}

func init() {
	addDocTemplate(&docTemplate{
		Name:    TypeAppName,
		Comment: "应用",
		Fields:  []*docTemplateField{},
	})
}

type ConfigDbModel struct {
	ElementNode
	Comment         string `json:"comment,omitempty"` // 说明
	Note            string `json:"note,omitempty"`    // 注释
	Type            string `json:"type,omitempty"`
	Host            string `json:"host,omitempty"`
	Port            int    `json:"port,omitempty"`
	Database        string `json:"database,omitempty"`
	DbName          string `json:"dbName,omitempty"`
	Username        string `json:"username,omitempty"`
	Password        string `json:"password,omitempty"`
	OdbcDsn         string `json:"odbcDsn,omitempty"`
	OdbcDialectName string `json:"odbcDialectName,omitempty"`

	Schema       string `json:"schema,omitempty"`
	Sid          string `json:"sid,omitempty"`
	MaxIdleConn  int    `json:"maxIdleConn,omitempty"`
	MaxOpenConn  int    `json:"maxOpenConn,omitempty"`
	DatabasePath string `json:"databasePath,omitempty"`
}

func init() {
	addDocTemplate(&docTemplate{
		Name:    TypeConfigDbName,
		Comment: "Database配置",
		Fields: []*docTemplateField{
			{Name: "comment", Comment: "配置说明"},
			{Name: "note", Comment: "配置源码注释"},
			{Name: "type", Comment: ""},
			{Name: "host", Comment: ""},
			{Name: "port", Comment: ""},
			{Name: "database", Comment: ""},
			{Name: "dbName", Comment: ""},
			{Name: "username", Comment: ""},
			{Name: "password", Comment: ""},
			{Name: "odbcDsn", Comment: ""},
			{Name: "odbcDialectName", Comment: ""},
			{Name: "schema", Comment: ""},
			{Name: "sid", Comment: ""},
			{Name: "maxIdleConn", Comment: ""},
			{Name: "maxOpenConn", Comment: ""},
			{Name: "databasePath", Comment: ""},
		},
	})
}

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

type ConfigKafkaModel struct {
	ElementNode
	Comment  string `json:"comment,omitempty"`  // 说明
	Note     string `json:"note,omitempty"`     // 注释
	Address  string `json:"address,omitempty"`  //
	Username string `json:"username,omitempty"` //
	Password string `json:"password,omitempty"` //

}

func init() {
	addDocTemplate(&docTemplate{
		Name:    TypeConfigKafkaName,
		Comment: "Kafka配置",
		Fields: []*docTemplateField{
			{Name: "comment", Comment: "配置说明"},
			{Name: "note", Comment: "配置源码注释"},
			{Name: "address", Comment: ""},
			{Name: "username", Comment: ""},
			{Name: "password", Comment: ""},
		},
	})
}

type ConfigMongodbModel struct {
	ElementNode
	Comment  string `json:"comment,omitempty"`  // 说明
	Note     string `json:"note,omitempty"`     // 注释
	Address  string `json:"address,omitempty"`  //
	Username string `json:"username,omitempty"` //
	Password string `json:"password,omitempty"` //

}

func init() {
	addDocTemplate(&docTemplate{
		Name:    TypeConfigMongodbName,
		Comment: "Mongodb配置",
		Fields: []*docTemplateField{
			{Name: "comment", Comment: "配置说明"},
			{Name: "note", Comment: "配置源码注释"},
			{Name: "address", Comment: ""},
			{Name: "username", Comment: ""},
			{Name: "password", Comment: ""},
		},
	})
}

type ConfigRedisModel struct {
	ElementNode
	Comment  string `json:"comment,omitempty"`  // 说明
	Note     string `json:"note,omitempty"`     // 注释
	Address  string `json:"address,omitempty"`  //
	Username string `json:"username,omitempty"` //
	Auth     string `json:"auth,omitempty"`     //
	CertPath string `json:"certPath,omitempty"` //

}

func init() {
	addDocTemplate(&docTemplate{
		Name:    TypeConfigRedisName,
		Comment: "Redis配置",
		Fields: []*docTemplateField{
			{Name: "comment", Comment: "配置说明"},
			{Name: "note", Comment: "配置源码注释"},
			{Name: "address", Comment: ""},
			{Name: "username", Comment: ""},
			{Name: "auth", Comment: ""},
			{Name: "certPath", Comment: ""},
		},
	})
}

type ConfigZkModel struct {
	ElementNode
	Comment  string `json:"comment,omitempty"`  // 说明
	Note     string `json:"note,omitempty"`     // 注释
	Address  string `json:"address,omitempty"`  //
	Username string `json:"username,omitempty"` //
	Password string `json:"password,omitempty"` //

}

func init() {
	addDocTemplate(&docTemplate{
		Name:    TypeConfigZkName,
		Comment: "Zookeeper配置",
		Fields: []*docTemplateField{
			{Name: "comment", Comment: "配置说明"},
			{Name: "note", Comment: "配置源码注释"},
			{Name: "address", Comment: ""},
			{Name: "username", Comment: ""},
			{Name: "password", Comment: ""},
		},
	})
}
