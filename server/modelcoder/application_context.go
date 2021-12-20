package modelcoder

type ApplicationContext struct {
	applicationModel *ApplicationModel
	constantMap      map[string]*ConstantModel
	dictionaryMap    map[string]*DictionaryModel
	datasourceMap    map[string]*DatasourceModel
	structMap        map[string]*StructModel
	variableMap      map[string]*VariableModel
	serverMap        map[string]*ServerModel
	serviceMap       map[string]*ServiceModel
	daoMap           map[string]*DaoModel
}

func newApplicationContext(applicationModel *ApplicationModel) *ApplicationContext {
	if applicationModel == nil {
		return nil
	}
	context := &ApplicationContext{
		applicationModel: applicationModel,
		constantMap:      map[string]*ConstantModel{},
		dictionaryMap:    map[string]*DictionaryModel{},
		datasourceMap:    map[string]*DatasourceModel{},
		structMap:        map[string]*StructModel{},
		variableMap:      map[string]*VariableModel{},
		serverMap:        map[string]*ServerModel{},
		serviceMap:       map[string]*ServiceModel{},
		daoMap:           map[string]*DaoModel{},
	}
	context.init()
	return context
}

func (this_ *ApplicationContext) init() *ApplicationContext {

	for _, one := range this_.applicationModel.Constants {
		this_.constantMap[one.Name] = one
	}
	for _, one := range this_.applicationModel.Dictionaries {
		this_.dictionaryMap[one.Name] = one
	}
	for _, one := range this_.applicationModel.Datasources {
		this_.datasourceMap[one.Name] = one
	}
	for _, one := range this_.applicationModel.Structs {
		this_.structMap[one.Name] = one
	}
	for _, one := range this_.applicationModel.Variables {
		this_.variableMap[one.Name] = one
	}
	for _, one := range this_.applicationModel.Servers {
		this_.serverMap[one.Name] = one
	}
	for _, one := range this_.applicationModel.Services {
		this_.serviceMap[one.Name] = one
	}
	for _, one := range this_.applicationModel.Daos {
		this_.daoMap[one.Name] = one
	}

	return this_
}

func (this_ *ApplicationContext) GetConstant(name string) *ConstantModel {
	model := this_.constantMap[name]
	return model
}

func (this_ *ApplicationContext) GetDictionary(name string) *DictionaryModel {
	model := this_.dictionaryMap[name]
	return model
}

func (this_ *ApplicationContext) GetDatasource(name string) *DatasourceModel {
	model := this_.datasourceMap[name]
	return model
}

func (this_ *ApplicationContext) GetStruct(name string) *StructModel {
	model := this_.structMap[name]
	return model
}

func (this_ *ApplicationContext) GetVariable(name string) *VariableModel {
	model := this_.variableMap[name]
	return model
}

func (this_ *ApplicationContext) GetServer(name string) *ServerModel {
	model := this_.serverMap[name]
	return model
}

func (this_ *ApplicationContext) GetService(name string) *ServiceModel {
	model := this_.serviceMap[name]
	return model
}

func (this_ *ApplicationContext) GetDao(name string) *DaoModel {
	model := this_.daoMap[name]
	return model
}
