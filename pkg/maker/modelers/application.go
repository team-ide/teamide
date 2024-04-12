package modelers

import (
	"errors"
	"strings"
)

type Application struct {
	ConstantList       []*ConstantModel `json:"constantList"`
	constantCache      map[string]*ConstantModel
	StructList         []*StructModel `json:"structList"`
	structCache        map[string]*StructModel
	ServiceList        []*ServiceModel `json:"serviceList"`
	serviceCache       map[string]*ServiceModel
	DaoList            []*DaoModel `json:"daoList"`
	daoCache           map[string]*DaoModel
	ErrorList          []*ErrorModel `json:"errorList"`
	errorCache         map[string]*ErrorModel
	FuncList           []*FuncModel `json:"funcList"`
	funcCache          map[string]*FuncModel
	LoadErrors         []*LoadError             `json:"loadErrors"`
	LanguageJavascript *LanguageJavascriptModel `json:"languageJavascript"`
	ConfigRedisList    []*ConfigRedisModel      `json:"configRedisList"`
	configRedisCache   map[string]*ConfigRedisModel
	ConfigDbList       []*ConfigDbModel `json:"configDbList"`
	configDbCache      map[string]*ConfigDbModel
	ConfigZkList       []*ConfigZkModel `json:"configZkList"`
	configZkCache      map[string]*ConfigZkModel
}

type LoadError struct {
	Type  *Type  `json:"type"`
	Path  string `json:"path"`
	Error string `json:"error"`
}

func (this_ *Application) AppendConstant(model *ConstantModel) (err error) {
	if this_.constantCache == nil {
		this_.constantCache = make(map[string]*ConstantModel)
	}
	if this_.constantCache[model.Name] != nil {
		err = errors.New("constant model [" + model.Name + "] already exist")
		return
	}
	this_.ConstantList = append(this_.ConstantList, model)
	this_.constantCache[model.Name] = model
	return
}

func (this_ *Application) GetConstant(name string) (model *ConstantModel) {
	if this_.constantCache != nil {
		model = this_.constantCache[name]
	}
	return
}

func (this_ *Application) AppendStruct(model *StructModel) (err error) {
	if this_.structCache == nil {
		this_.structCache = make(map[string]*StructModel)
	}
	if this_.structCache[model.Name] != nil {
		err = errors.New("struct model [" + model.Name + "] already exist")
		return
	}
	this_.StructList = append(this_.StructList, model)
	this_.structCache[model.Name] = model
	return
}

func (this_ *Application) GetStruct(name string) (model *StructModel) {
	if this_.structCache != nil {
		model = this_.structCache[name]
	}
	return
}

func (this_ *Application) AppendService(model *ServiceModel) (err error) {
	if this_.serviceCache == nil {
		this_.serviceCache = make(map[string]*ServiceModel)
	}
	if this_.serviceCache[model.Name] != nil {
		err = errors.New("service model [" + model.Name + "] already exist")
		return
	}
	this_.ServiceList = append(this_.ServiceList, model)
	this_.serviceCache[model.Name] = model
	return
}

func (this_ *Application) GetService(name string) (model *ServiceModel) {
	if this_.serviceCache != nil {
		model = this_.serviceCache[name]
	}
	return
}

func (this_ *Application) AppendDao(model *DaoModel) (err error) {
	if this_.daoCache == nil {
		this_.daoCache = make(map[string]*DaoModel)
	}
	if this_.daoCache[model.Name] != nil {
		err = errors.New("dao model [" + model.Name + "] already exist")
		return
	}
	this_.DaoList = append(this_.DaoList, model)
	this_.daoCache[model.Name] = model
	return
}

func (this_ *Application) GetDao(name string) (model *DaoModel) {
	if this_.daoCache != nil {
		model = this_.daoCache[name]
	}
	return
}

func (this_ *Application) AppendConfigRedis(model *ConfigRedisModel) (err error) {
	if this_.configRedisCache == nil {
		this_.configRedisCache = make(map[string]*ConfigRedisModel)
	}
	if this_.configRedisCache[model.Name] != nil {
		err = errors.New("redis model [" + model.Name + "] already exist")
		return
	}
	this_.ConfigRedisList = append(this_.ConfigRedisList, model)
	this_.configRedisCache[model.Name] = model
	return
}

func (this_ *Application) GetConfigRedis(name string) (model *ConfigRedisModel) {
	if this_.configRedisCache != nil {
		model = this_.configRedisCache[name]
	}
	return
}

func (this_ *Application) AppendConfigDb(model *ConfigDbModel) (err error) {
	if this_.configDbCache == nil {
		this_.configDbCache = make(map[string]*ConfigDbModel)
	}
	if this_.configDbCache[model.Name] != nil {
		err = errors.New("db model [" + model.Name + "] already exist")
		return
	}
	this_.ConfigDbList = append(this_.ConfigDbList, model)
	this_.configDbCache[model.Name] = model
	return
}

func (this_ *Application) GetConfigDb(name string) (model *ConfigDbModel) {
	if this_.configDbCache != nil {
		model = this_.configDbCache[name]
	}
	return
}

func (this_ *Application) AppendConfigZk(model *ConfigZkModel) (err error) {
	if this_.configZkCache == nil {
		this_.configZkCache = make(map[string]*ConfigZkModel)
	}
	if this_.configZkCache[model.Name] != nil {
		err = errors.New("zk model [" + model.Name + "] already exist")
		return
	}
	this_.ConfigZkList = append(this_.ConfigZkList, model)
	this_.configZkCache[model.Name] = model
	return
}

func (this_ *Application) GetConfigZk(name string) (model *ConfigZkModel) {
	if this_.configZkCache != nil {
		model = this_.configZkCache[name]
	}
	return
}

func (this_ *Application) AppendError(model *ErrorModel) (err error) {
	if this_.errorCache == nil {
		this_.errorCache = make(map[string]*ErrorModel)
	}
	if this_.errorCache[model.Name] != nil {
		err = errors.New("error model [" + model.Name + "] already exist")
		return
	}
	this_.ErrorList = append(this_.ErrorList, model)
	this_.errorCache[model.Name] = model
	return
}

func (this_ *Application) GetError(name string) (model *ErrorModel) {
	if this_.errorCache != nil {
		model = this_.errorCache[name]
	}
	return
}

func (this_ *Application) AppendFunc(model *FuncModel) (err error) {
	if this_.funcCache == nil {
		this_.funcCache = make(map[string]*FuncModel)
	}
	if this_.funcCache[model.Name] != nil {
		err = errors.New("func model [" + model.Name + "] already exist")
		return
	}
	this_.FuncList = append(this_.FuncList, model)
	this_.funcCache[model.Name] = model
	return
}

func (this_ *Application) GetFunc(name string) (model *FuncModel) {
	if this_.funcCache != nil {
		model = this_.funcCache[name]
	}
	return
}

func (this_ *Application) SetLanguageJavascript(model *LanguageJavascriptModel) {
	if model == nil {
		model = &LanguageJavascriptModel{}
	}
	this_.LanguageJavascript = model
	return
}

func (this_ *Application) GetLanguageJavascript() (model *LanguageJavascriptModel) {
	if this_.LanguageJavascript == nil {
		this_.SetLanguageJavascript(nil)
	}
	model = this_.LanguageJavascript
	return
}

func (this_ *Application) GetValueType(name string) (valueType *ValueType, err error) {
	for _, one := range ValueTypes {
		if strings.EqualFold(one.Name, name) {
			valueType = one
			return
		}
	}
	valueType = &ValueType{
		Name: name,
	}
	valueType.Struct = this_.GetStruct(valueType.Name)
	if valueType.Struct == nil {
		err = errors.New("value type and struct not found name [" + valueType.Name + "]")
		return
	}
	return
}
