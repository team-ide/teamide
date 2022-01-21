package model

// Model Context
type ModelContext struct {
	Constants            []*ConstantModel       `json:"constants,omitempty"`            // 常量Model，一经定义不可修改，放一些全局使用的值，通过固定值、启动传参或环境变量设值
	Errors               []*ErrorModel          `json:"errors,omitempty"`               // 全局错误信息
	Dictionaries         []*DictionaryModel     `json:"dictionaries,omitempty"`         // 字典Model，定义数据字典，用户状态、XX类型值和文案映射等
	DatasourceDatabases  []*DatasourceDatabase  `json:"datasourceDatabases,omitempty"`  // 数据源Model，定义数据源，Database、Redis、Kafka等
	DatasourceRedises    []*DatasourceRedis     `json:"datasourceRedises,omitempty"`    // 数据源Model，定义数据源，Database、Redis、Kafka等
	DatasourceKafkas     []*DatasourceKafka     `json:"datasourceKafkas,omitempty"`     // 数据源Model，定义数据源，Database、Redis、Kafka等
	DatasourceZookeepers []*DatasourceZookeeper `json:"datasourceZookeepers,omitempty"` // 数据源Model，定义数据源，Database、Redis、Kafka等
	Structs              []*StructModel         `json:"structs,omitempty"`              // 结构体Model，定义结构体，结构体字段、JSON字段、表字段等
	ServerWebs           []*ServerWebModel      `json:"serverWebs,omitempty"`           // 服务器层Model，用于提供服务接口能力，HTTP、RPC等
	Services             []*ServiceModel        `json:"services,omitempty"`             // 服务层Model，用于逻辑处理，验证等
	Tests                []*TestModel           `json:"tests,omitempty"`                // 测试Model，用于逻辑处理，验证等

	constantMap            map[string]*ConstantModel       `json:"-"`
	errorMap               map[string]*ErrorModel          `json:"-"`
	dictionaryMap          map[string]*DictionaryModel     `json:"-"`
	datasourceDatabaseMap  map[string]*DatasourceDatabase  `json:"-"`
	datasourceRedisMap     map[string]*DatasourceRedis     `json:"-"`
	datasourceKafkaMap     map[string]*DatasourceKafka     `json:"-"`
	datasourceZookeeperMap map[string]*DatasourceZookeeper `json:"-"`
	structMap              map[string]*StructModel         `json:"-"`
	serverWebMap           map[string]*ServerWebModel      `json:"-"`
	serviceMap             map[string]*ServiceModel        `json:"-"`
	testMap                map[string]*TestModel           `json:"-"`
}

func (this_ *ModelContext) Init() *ModelContext {
	this_.initConstant()
	this_.initError()
	this_.initDictionary()
	this_.initDatasourceDatabase()
	this_.initDatasourceRedis()
	this_.initDatasourceKafka()
	this_.initDatasourceZookeeper()
	this_.initStruct()
	this_.initServerWeb()
	this_.initService()
	this_.initTest()
	fileInfoStruct := &StructModel{
		Name: "fileInfo",
		Fields: []*StructFieldModel{
			{Name: "name", DataType: "string"},
			{Name: "type", DataType: "string"},
			{Name: "path", DataType: "string"},
			{Name: "dir", DataType: "string"},
			{Name: "size", DataType: "long"},
			{Name: "absolutePath", DataType: "string"},
		},
	}
	this_.AppendStruct(fileInfoStruct)

	pageInfoStruct := &StructModel{
		Name: "pageInfo",
		Fields: []*StructFieldModel{
			{Name: "pageNumber", DataType: "long"},
			{Name: "pageSize", DataType: "long"},
			{Name: "totalPage", DataType: "long"},
			{Name: "totalSize", DataType: "long"},
			{Name: "list", DataType: "map", IsList: true},
		},
	}
	this_.AppendStruct(pageInfoStruct)
	return this_
}

func (this_ *ModelContext) initConstant() *ModelContext {
	this_.constantMap = map[string]*ConstantModel{}
	for _, one := range this_.Constants {
		this_.constantMap[one.Name] = one
	}
	return this_
}

func (this_ *ModelContext) initError() *ModelContext {
	this_.errorMap = map[string]*ErrorModel{}
	for _, one := range this_.Errors {
		this_.errorMap[one.Name] = one
	}
	return this_
}

func (this_ *ModelContext) initDictionary() *ModelContext {
	this_.dictionaryMap = map[string]*DictionaryModel{}
	for _, one := range this_.Dictionaries {
		this_.dictionaryMap[one.Name] = one
	}
	return this_
}

func (this_ *ModelContext) initDatasourceDatabase() *ModelContext {
	this_.datasourceDatabaseMap = map[string]*DatasourceDatabase{}
	for _, one := range this_.DatasourceDatabases {
		this_.datasourceDatabaseMap[one.Name] = one
	}
	return this_
}

func (this_ *ModelContext) initDatasourceRedis() *ModelContext {
	this_.datasourceRedisMap = map[string]*DatasourceRedis{}
	for _, one := range this_.DatasourceRedises {
		this_.datasourceRedisMap[one.Name] = one
	}
	return this_
}

func (this_ *ModelContext) initDatasourceKafka() *ModelContext {
	this_.datasourceKafkaMap = map[string]*DatasourceKafka{}
	for _, one := range this_.DatasourceKafkas {
		this_.datasourceKafkaMap[one.Name] = one
	}
	return this_
}

func (this_ *ModelContext) initDatasourceZookeeper() *ModelContext {
	this_.datasourceZookeeperMap = map[string]*DatasourceZookeeper{}
	for _, one := range this_.DatasourceZookeepers {
		this_.datasourceZookeeperMap[one.Name] = one
	}
	return this_
}

func (this_ *ModelContext) initStruct() *ModelContext {
	this_.structMap = map[string]*StructModel{}
	for _, one := range this_.Structs {
		this_.structMap[one.Name] = one
	}
	return this_
}

func (this_ *ModelContext) initServerWeb() *ModelContext {
	this_.serverWebMap = map[string]*ServerWebModel{}
	for _, one := range this_.ServerWebs {
		this_.serverWebMap[one.Name] = one
	}
	return this_
}

func (this_ *ModelContext) initService() *ModelContext {
	this_.serviceMap = map[string]*ServiceModel{}
	for _, one := range this_.Services {
		this_.serviceMap[one.Name] = one
	}
	return this_
}

func (this_ *ModelContext) initTest() *ModelContext {
	this_.testMap = map[string]*TestModel{}
	for _, one := range this_.Tests {
		this_.testMap[one.Name] = one
	}
	return this_
}

func (this_ *ModelContext) GetConstant(name string) *ConstantModel {
	model := this_.constantMap[name]
	return model
}

func (this_ *ModelContext) GetError(name string) *ErrorModel {
	model := this_.errorMap[name]
	return model
}

func (this_ *ModelContext) GetDictionary(name string) *DictionaryModel {
	model := this_.dictionaryMap[name]
	return model
}

func (this_ *ModelContext) GetDatasourceDatabase(name string) *DatasourceDatabase {
	model := this_.datasourceDatabaseMap[name]
	return model
}

func (this_ *ModelContext) GetDatasourceRedis(name string) *DatasourceRedis {
	model := this_.datasourceRedisMap[name]
	return model
}

func (this_ *ModelContext) GetDatasourceKafka(name string) *DatasourceKafka {
	model := this_.datasourceKafkaMap[name]
	return model
}

func (this_ *ModelContext) GetDatasourceZookeeper(name string) *DatasourceZookeeper {
	model := this_.datasourceZookeeperMap[name]
	return model
}

func (this_ *ModelContext) GetStruct(name string) *StructModel {
	model := this_.structMap[name]
	return model
}

func (this_ *ModelContext) GetServerWeb(name string) *ServerWebModel {
	model := this_.serverWebMap[name]
	return model
}

func (this_ *ModelContext) GetService(name string) *ServiceModel {
	model := this_.serviceMap[name]
	return model
}

func (this_ *ModelContext) GetTest(name string) *TestModel {
	model := this_.testMap[name]
	return model
}

func (this_ *ModelContext) AppendConstant(model ...*ConstantModel) *ModelContext {
	this_.Constants = append(this_.Constants, model...)
	this_.initConstant()
	return this_
}

func (this_ *ModelContext) AppendError(model ...*ErrorModel) *ModelContext {
	this_.Errors = append(this_.Errors, model...)
	this_.initError()
	return this_
}

func (this_ *ModelContext) AppendDictionary(model ...*DictionaryModel) *ModelContext {
	this_.Dictionaries = append(this_.Dictionaries, model...)
	this_.initDictionary()
	return this_
}

func (this_ *ModelContext) AppendDatasourceDatabase(model ...*DatasourceDatabase) *ModelContext {
	this_.DatasourceDatabases = append(this_.DatasourceDatabases, model...)
	this_.initDatasourceDatabase()
	return this_
}

func (this_ *ModelContext) AppendDatasourceRedis(model ...*DatasourceRedis) *ModelContext {
	this_.DatasourceRedises = append(this_.DatasourceRedises, model...)
	this_.initDatasourceRedis()
	return this_
}

func (this_ *ModelContext) AppendDatasourceKafka(model ...*DatasourceKafka) *ModelContext {
	this_.DatasourceKafkas = append(this_.DatasourceKafkas, model...)
	this_.initDatasourceKafka()
	return this_
}

func (this_ *ModelContext) AppendDatasourceZookeeper(model ...*DatasourceZookeeper) *ModelContext {
	this_.DatasourceZookeepers = append(this_.DatasourceZookeepers, model...)
	this_.initDatasourceZookeeper()
	return this_
}

func (this_ *ModelContext) AppendStruct(model ...*StructModel) *ModelContext {
	this_.Structs = append(this_.Structs, model...)
	this_.initStruct()
	return this_
}

func (this_ *ModelContext) AppendServerWeb(model ...*ServerWebModel) *ModelContext {
	this_.ServerWebs = append(this_.ServerWebs, model...)
	this_.initServerWeb()
	return this_
}

func (this_ *ModelContext) AppendService(model ...*ServiceModel) *ModelContext {
	this_.Services = append(this_.Services, model...)
	this_.initService()
	return this_
}

func (this_ *ModelContext) AppendTest(model ...*TestModel) *ModelContext {
	this_.Tests = append(this_.Tests, model...)
	this_.initTest()
	return this_
}
