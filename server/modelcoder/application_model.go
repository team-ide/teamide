package modelcoder

// 应用Model
type ApplicationModel struct {
	Constants    []*ConstantModel   `json:"constants,omitempty"`    // 常量Model，一经定义不可修改，放一些全局使用的值，通过固定值、启动传参或环境变量设值
	Dictionaries []*DictionaryModel `json:"dictionaries,omitempty"` // 字典Model，定义数据字典，用户状态、XX类型值和文案映射等
	Datasources  []*DatasourceModel `json:"datasources,omitempty"`  // 数据源Model，定义数据源，Database、Redis、Kafka等
	Structs      []*StructModel     `json:"structs,omitempty"`      // 结构体Model，定义结构体，结构体字段、JSON字段、表字段等
	Variables    []*VariableModel   `json:"variables,omitempty"`    // 变量Model，调用一个接口、服务、数据访问从入口设值开始，整个线上可用
	Servers      []*ServerModel     `json:"servers,omitempty"`      // 服务器层Model，用于提供服务接口能力，HTTP、RPC等
	Services     []ServiceModel     `json:"services,omitempty"`     // 服务层Model，用于逻辑处理，验证等
	Daos         []DaoModel         `json:"daos,omitempty"`         // 数据层Model，用于处理数据存储，查询数据等
}

func (this_ *ApplicationModel) AppendConstant(model ...*ConstantModel) *ApplicationModel {
	this_.Constants = append(this_.Constants, model...)
	return this_
}

func (this_ *ApplicationModel) AppendDictionary(model ...*DictionaryModel) *ApplicationModel {
	this_.Dictionaries = append(this_.Dictionaries, model...)
	return this_
}

func (this_ *ApplicationModel) AppendDatasource(model ...*DatasourceModel) *ApplicationModel {
	this_.Datasources = append(this_.Datasources, model...)
	return this_
}

func (this_ *ApplicationModel) AppendStruct(model ...*StructModel) *ApplicationModel {
	this_.Structs = append(this_.Structs, model...)
	return this_
}

func (this_ *ApplicationModel) AppendVariable(model ...*VariableModel) *ApplicationModel {
	this_.Variables = append(this_.Variables, model...)
	return this_
}

func (this_ *ApplicationModel) AppendServer(model ...*ServerModel) *ApplicationModel {
	this_.Servers = append(this_.Servers, model...)
	return this_
}

func (this_ *ApplicationModel) AppendService(model ...ServiceModel) *ApplicationModel {
	this_.Services = append(this_.Services, model...)
	return this_
}

func (this_ *ApplicationModel) AppendDao(model ...DaoModel) *ApplicationModel {
	this_.Daos = append(this_.Daos, model...)
	return this_
}
