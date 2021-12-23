package modelcoder

type applicationContext struct {
	applicationModel *ApplicationModel
	constantMap      map[string]*ConstantModel
	dictionaryMap    map[string]*DictionaryModel
	datasourceMap    map[string]*DatasourceModel
	structMap        map[string]*StructModel
	serverMap        map[string]*ServerModel
	serviceMap       map[string]ServiceModel
	daoMap           map[string]DaoModel
}

func newApplicationContext(applicationModel *ApplicationModel) *applicationContext {
	if applicationModel == nil {
		return nil
	}
	context := &applicationContext{
		applicationModel: applicationModel,
		constantMap:      make(map[string]*ConstantModel),
		dictionaryMap:    make(map[string]*DictionaryModel),
		datasourceMap:    make(map[string]*DatasourceModel),
		structMap:        make(map[string]*StructModel),
		serverMap:        make(map[string]*ServerModel),
		serviceMap:       make(map[string]ServiceModel),
		daoMap:           make(map[string]DaoModel),
	}
	context.init()
	return context
}

func (this_ *applicationContext) init() *applicationContext {

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
	for _, one := range this_.applicationModel.Servers {
		this_.serverMap[one.Name] = one
	}
	for _, one := range this_.applicationModel.Services {
		this_.serviceMap[one.GetName()] = one
	}
	for _, one := range this_.applicationModel.Daos {
		this_.daoMap[one.GetName()] = one
	}

	return this_
}

func (this_ *applicationContext) GetConstant(name string) *ConstantModel {
	model := this_.constantMap[name]
	return model
}

func (this_ *applicationContext) GetDictionary(name string) *DictionaryModel {
	model := this_.dictionaryMap[name]
	return model
}

func (this_ *applicationContext) GetDatasource(name string) *DatasourceModel {
	model := this_.datasourceMap[name]
	return model
}

func (this_ *applicationContext) GetStruct(name string) *StructModel {
	model := this_.structMap[name]
	return model
}

func (this_ *applicationContext) GetServer(name string) *ServerModel {
	model := this_.serverMap[name]
	return model
}

func (this_ *applicationContext) GetService(name string) ServiceModel {
	model := this_.serviceMap[name]
	return model
}

func (this_ *applicationContext) GetDao(name string) DaoModel {
	model := this_.daoMap[name]
	return model
}
