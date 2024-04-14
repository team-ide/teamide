package modelers

import (
	"errors"
	"strings"
)

type Application struct {
	Dir string `json:"dir"`

	Children []*Element `json:"children"`

	Constant      []*ConstantModel `json:"constant"`
	constantCache map[string]*ConstantModel
	Struct        []*StructModel `json:"struct"`
	structCache   map[string]*StructModel
	Service       []*ServiceModel `json:"service"`
	serviceCache  map[string]*ServiceModel
	Dao           []*DaoModel `json:"dao"`
	daoCache      map[string]*DaoModel
	Error         []*ErrorModel `json:"error"`
	errorCache    map[string]*ErrorModel
	Func          []*FuncModel `json:"func"`
	funcCache     map[string]*FuncModel

	ConfigRedis        []*ConfigRedisModel `json:"configRedis"`
	configRedisCache   map[string]*ConfigRedisModel
	ConfigDb           []*ConfigDbModel `json:"configDb"`
	configDbCache      map[string]*ConfigDbModel
	ConfigZk           []*ConfigZkModel `json:"configZk"`
	configZkCache      map[string]*ConfigZkModel
	ConfigKafka        []*ConfigKafkaModel `json:"configKafka"`
	configKafkaCache   map[string]*ConfigKafkaModel
	ConfigMongodb      []*ConfigMongodbModel `json:"configMongodb"`
	configMongodbCache map[string]*ConfigMongodbModel

	LanguageJavascript *LanguageJavascriptModel `json:"languageJavascript"`

	LoadErrors []*LoadError `json:"loadErrors"`
}

type Element struct {
	Key      string      `json:"key,omitempty"`
	Text     string      `json:"text,omitempty"`
	IsType   bool        `json:"isType,omitempty"`
	IsPack   bool        `json:"isPack,omitempty"`
	Type     string      `json:"type,omitempty"`
	Children []*Element  `json:"children,omitempty"`
	Model    interface{} `json:"model,omitempty"`
}

type LoadError struct {
	Type  *Type  `json:"type"`
	Path  string `json:"path"`
	Error string `json:"error"`
}

func (this_ *Application) AppendType(parent *Element, modelType *Type) (err error) {
	element := &Element{
		Key:  modelType.Name,
		Text: modelType.Comment,
	}
	return
}

func (this_ *Application) AppendPack(model *ConstantModel) (err error) {
	if this_.constantCache == nil {
		this_.constantCache = make(map[string]*ConstantModel)
	}
	if this_.constantCache[model.Name] != nil {
		err = errors.New("constant model [" + model.Name + "] already exist")
		return
	}
	this_.Constant = append(this_.Constant, model)
	this_.constantCache[model.Name] = model
	return
}

func (this_ *Application) AppendConstant(model *ConstantModel) (err error) {
	if this_.constantCache == nil {
		this_.constantCache = make(map[string]*ConstantModel)
	}
	if this_.constantCache[model.Name] != nil {
		err = errors.New("constant model [" + model.Name + "] already exist")
		return
	}
	this_.Constant = append(this_.Constant, model)
	this_.constantCache[model.Name] = model
	return
}

func (this_ *Application) GetConstantList() []*ConstantModel {
	return this_.Constant
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
	this_.Struct = append(this_.Struct, model)
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
	this_.Service = append(this_.Service, model)
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
	this_.Dao = append(this_.Dao, model)
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
	this_.ConfigRedis = append(this_.ConfigRedis, model)
	this_.configRedisCache[model.Name] = model
	return
}
func (this_ *Application) GetConfigRedisList() []*ConfigRedisModel {
	return this_.ConfigRedis
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
	this_.ConfigDb = append(this_.ConfigDb, model)
	this_.configDbCache[model.Name] = model
	return
}
func (this_ *Application) GetConfigDbList() []*ConfigDbModel {
	return this_.ConfigDb
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
	this_.ConfigZk = append(this_.ConfigZk, model)
	this_.configZkCache[model.Name] = model
	return
}
func (this_ *Application) GetConfigZkList() []*ConfigZkModel {
	return this_.ConfigZk
}

func (this_ *Application) GetConfigZk(name string) (model *ConfigZkModel) {
	if this_.configZkCache != nil {
		model = this_.configZkCache[name]
	}
	return
}

func (this_ *Application) AppendConfigKafka(model *ConfigKafkaModel) (err error) {
	if this_.configKafkaCache == nil {
		this_.configKafkaCache = make(map[string]*ConfigKafkaModel)
	}
	if this_.configKafkaCache[model.Name] != nil {
		err = errors.New("kafka model [" + model.Name + "] already exist")
		return
	}
	this_.ConfigKafka = append(this_.ConfigKafka, model)
	this_.configKafkaCache[model.Name] = model
	return
}
func (this_ *Application) GetConfigKafkaList() []*ConfigKafkaModel {
	return this_.ConfigKafka
}

func (this_ *Application) GetConfigKafka(name string) (model *ConfigKafkaModel) {
	if this_.configKafkaCache != nil {
		model = this_.configKafkaCache[name]
	}
	return
}

func (this_ *Application) AppendConfigMongodb(model *ConfigMongodbModel) (err error) {
	if this_.configMongodbCache == nil {
		this_.configMongodbCache = make(map[string]*ConfigMongodbModel)
	}
	if this_.configMongodbCache[model.Name] != nil {
		err = errors.New("mongodb model [" + model.Name + "] already exist")
		return
	}
	this_.ConfigMongodb = append(this_.ConfigMongodb, model)
	this_.configMongodbCache[model.Name] = model
	return
}
func (this_ *Application) GetConfigMongodbList() []*ConfigMongodbModel {
	return this_.ConfigMongodb
}

func (this_ *Application) GetConfigMongodb(name string) (model *ConfigMongodbModel) {
	if this_.configMongodbCache != nil {
		model = this_.configMongodbCache[name]
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
	this_.Error = append(this_.Error, model)
	this_.errorCache[model.Name] = model
	return
}

func (this_ *Application) GetErrorList() []*ErrorModel {
	return this_.Error
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
	this_.Func = append(this_.Func, model)
	this_.funcCache[model.Name] = model
	return
}

func (this_ *Application) GetFuncList() []*FuncModel {
	return this_.Func
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
