package modelers

import (
	"errors"
	"github.com/team-ide/go-tool/util"
	"os"
	"reflect"
	"strings"
	"sync"
)

type Application struct {
	Dir string `json:"dir"`

	Children []*Element `json:"children"`

	constantList []*ConstantModel
	structList   []*StructModel
	serviceList  []*ServiceModel
	daoList      []*DaoModel
	errorList    []*ErrorModel
	funcList     []*FuncModel

	configRedisList         []*ConfigRedisModel
	configDbList            []*ConfigDbModel
	configZkList            []*ConfigZkModel
	configKafkaList         []*ConfigKafkaModel
	configMongodbList       []*ConfigMongodbModel
	configElasticsearchList []*ConfigEsModel

	languageJavascript *LanguageJavascriptModel

	modelTypeCaches     map[*Type]*util.Cache
	modelTypeCachesLock sync.Mutex

	modelTypeCache map[string]*Type

	LoadErrors []*LoadError `json:"loadErrors"`
}

type Element struct {
	Key           string     `json:"key,omitempty"`
	Text          string     `json:"text,omitempty"`
	IsType        bool       `json:"isType,omitempty"`
	IsPack        bool       `json:"isPack,omitempty"`
	ModelName     string     `json:"modelName,omitempty"`
	ModelType     string     `json:"modelType,omitempty"`
	ModelTypeText string     `json:"modelTypeText,omitempty"`
	Children      []*Element `json:"children,omitempty"`
	Pack          *Pack      `json:"pack,omitempty"`
}

type Pack struct {
	Name    string `json:"name,omitempty"`    // 名称，同一个应用中唯一
	Comment string `json:"comment,omitempty"` // 说明
	Note    string `json:"note,omitempty"`    // 注释
}

type LoadError struct {
	Type  *Type  `json:"type"`
	Path  string `json:"path"`
	Error string `json:"error"`
}

func (this_ *Application) appendType(parent *Element, modelType *Type) (element *Element) {
	element = &Element{
		Key:  modelType.Name,
		Text: modelType.Comment,
	}
	if parent != nil {
		parent.Children = append(parent.Children, element)
	} else {
		this_.Children = append(this_.Children, element)
	}
	element.IsType = true
	if this_.modelTypeCache == nil {
		this_.modelTypeCache = make(map[string]*Type)
	}
	this_.modelTypeCache[element.Key] = modelType
	return
}

func (this_ *Application) GetModelType(key string) (modelType *Type) {
	modelType = this_.modelTypeCache[key]
	return
}

func (this_ *Application) GetModelTypeModel(key string, name string) (model interface{}) {
	modelType := this_.modelTypeCache[key]
	if modelType == nil {
		util.Logger.Error("model type " + key + " not found")
		return
	}
	cache := this_.getModelTypeCache(modelType)
	model, _ = cache.Get(name)
	return
}

func (this_ *Application) getModePath(modelType *Type, name string) (path string) {
	path = this_.Dir + modelType.Name + "/" + name + ".yml"
	return
}

func (this_ *Application) Remove(key string, name string) (err error) {
	if name == "" {
		err = errors.New("model name is empty")
		return
	}
	modelType := this_.modelTypeCache[key]
	if modelType == nil {
		err = errors.New("model type " + key + " not found")
		return
	}
	path := this_.getModePath(modelType, name)
	exist, err := util.PathExists(path)
	if err != nil {
		return
	}
	if exist {
		err = os.Remove(path)
	}
	return
}

func (this_ *Application) Save(key string, name string, model interface{}, isNew bool) (err error) {
	if name == "" {
		err = errors.New("model name is empty")
		return
	}
	modelType := this_.modelTypeCache[key]
	if modelType == nil {
		err = errors.New("model type " + key + " not found")
		return
	}
	path := this_.getModePath(modelType, name)
	exist, err := util.PathExists(path)
	if err != nil {
		return
	}
	if isNew && exist {
		err = errors.New("model [" + key + "] [" + name + "] is exist")
		return
	}
	dir := path[:strings.LastIndex(path, "/")]

	if dirExist, _ := util.PathExists(dir); !dirExist {
		if err = os.MkdirAll(dir, os.ModePerm); err != nil {
			return
		}
	}
	if model != nil {
		var text string
		if text, err = modelType.toText(model); err != nil {
			return
		}
		var f *os.File
		if f, err = os.Create(path); err != nil {
			return
		}
		defer func() { _ = f.Close() }()
		if _, err = f.WriteString(text); err != nil {
			return
		}
	} else {
		if err = os.MkdirAll(path, os.ModePerm); err != nil {
			return
		}
	}
	return
}

func (this_ *Application) GetModelTypeModels(key string) (models interface{}) {
	modelType := this_.modelTypeCache[key]
	if modelType == nil {
		util.Logger.Error("model type " + key + " not found")
		return
	}
	switch modelType {
	case TypeConstant:
		models = this_.constantList
		break
	case TypeError:
		models = this_.errorList
		break
	case TypeStruct:
		models = this_.structList
		break
	case TypeDao:
		models = this_.daoList
		break
	case TypeService:
		models = this_.serviceList
		break
	case TypeFunc:
		models = this_.funcList
		break
	case TypeConfigDb:
		models = this_.configDbList
		break
	case TypeConfigRedis:
		models = this_.configRedisList
		break
	case TypeConfigZk:
		models = this_.configZkList
		break
	case TypeConfigKafka:
		models = this_.configKafkaList
		break
	case TypeConfigMongodb:
		models = this_.configMongodbList
		break
	case TypeConfigElasticsearch:
		models = this_.configElasticsearchList
		break

	default:
		util.Logger.Error("model type " + key + " not support get model list")
		return
	}
	return
}

func (this_ *Application) appendPack(parent *Element, modelType *Type, name string) (element *Element) {
	element = &Element{
		Key:  name,
		Text: name,
	}
	if parent != nil {
		element.Key = parent.Key + "/" + element.Key
		parent.Children = append(parent.Children, element)
	} else {
		this_.Children = append(this_.Children, element)
	}
	element.IsPack = true
	element.Pack = &Pack{
		Name: name,
	}
	element.ModelType = modelType.Name
	element.ModelTypeText = modelType.Comment
	return
}

func (this_ *Application) appendModel(parent *Element, modelType *Type, fileName string, name string, model interface{}) (err error) {
	element := &Element{
		Key:  fileName,
		Text: fileName,
	}
	if parent != nil {
		element.Key = parent.Key + "/" + element.Key
		parent.Children = append(parent.Children, element)
	} else {
		this_.Children = append(this_.Children, element)
	}

	element.ModelType = modelType.Name
	element.ModelTypeText = modelType.Comment
	element.ModelName = name

	this_.modelTypeCachesLock.Lock()
	defer this_.modelTypeCachesLock.Unlock()

	if this_.modelTypeCaches == nil {
		this_.modelTypeCaches = make(map[*Type]*util.Cache)
	}
	cache := this_.modelTypeCaches[modelType]
	if cache == nil {
		cache = util.NewCache()
		this_.modelTypeCaches[modelType] = cache
	}

	if _, ok := cache.Get(name); ok {
		err = errors.New("type " + modelType.Name + " model [" + name + "] already exist")
		return
	}

	cache.Put(name, model)
	switch mV := model.(type) {
	case *ConstantModel:
		this_.constantList = append(this_.constantList, mV)
		break
	case *ErrorModel:
		this_.errorList = append(this_.errorList, mV)
		break
	case *StructModel:
		this_.structList = append(this_.structList, mV)
		break
	case *DaoModel:
		this_.daoList = append(this_.daoList, mV)
		break
	case *ServiceModel:
		this_.serviceList = append(this_.serviceList, mV)
		break
	case *FuncModel:
		this_.funcList = append(this_.funcList, mV)
		break
	case *ConfigDbModel:
		this_.configDbList = append(this_.configDbList, mV)
		break
	case *ConfigKafkaModel:
		this_.configKafkaList = append(this_.configKafkaList, mV)
		break
	case *ConfigMongodbModel:
		this_.configMongodbList = append(this_.configMongodbList, mV)
		break
	case *ConfigZkModel:
		this_.configZkList = append(this_.configZkList, mV)
		break
	case *ConfigRedisModel:
		this_.configRedisList = append(this_.configRedisList, mV)
		break
	case *ConfigEsModel:
		this_.configElasticsearchList = append(this_.configElasticsearchList, mV)
		break
	case *LanguageJavascriptModel:
		this_.languageJavascript = mV
		break
	default:
		err = errors.New("type " + modelType.Name + " model [" + name + "] [" + reflect.TypeOf(model).String() + "] not support")
		return
	}

	return
}

func (this_ *Application) getModelTypeCache(modelType *Type) (cache *util.Cache) {
	this_.modelTypeCachesLock.Lock()
	defer this_.modelTypeCachesLock.Unlock()

	if this_.modelTypeCaches == nil {
		this_.modelTypeCaches = make(map[*Type]*util.Cache)
	}
	cache = this_.modelTypeCaches[modelType]
	if cache == nil {
		cache = util.NewCache()
		this_.modelTypeCaches[modelType] = cache
	}
	return
}

func (this_ *Application) GetConstantList() []*ConstantModel {
	return this_.constantList
}

func (this_ *Application) GetConstant(name string) (model *ConstantModel) {
	cache := this_.getModelTypeCache(TypeConstant)
	find, _ := cache.Get(name)
	if find != nil {
		model = find.(*ConstantModel)
	}
	return
}

func (this_ *Application) GetStruct(name string) (model *StructModel) {
	cache := this_.getModelTypeCache(TypeStruct)
	find, _ := cache.Get(name)
	if find != nil {
		model = find.(*StructModel)
	}
	return
}

func (this_ *Application) GetService(name string) (model *ServiceModel) {
	cache := this_.getModelTypeCache(TypeService)
	find, _ := cache.Get(name)
	if find != nil {
		model = find.(*ServiceModel)
	}
	return
}

func (this_ *Application) GetDao(name string) (model *DaoModel) {
	cache := this_.getModelTypeCache(TypeDao)
	find, _ := cache.Get(name)
	if find != nil {
		model = find.(*DaoModel)
	}
	return
}

func (this_ *Application) GetConfigRedisList() []*ConfigRedisModel {
	return this_.configRedisList
}

func (this_ *Application) GetConfigRedis(name string) (model *ConfigRedisModel) {
	cache := this_.getModelTypeCache(TypeConfigRedis)
	find, _ := cache.Get(name)
	if find != nil {
		model = find.(*ConfigRedisModel)
	}
	return
}

func (this_ *Application) GetConfigDbList() []*ConfigDbModel {
	return this_.configDbList
}

func (this_ *Application) GetConfigDb(name string) (model *ConfigDbModel) {
	cache := this_.getModelTypeCache(TypeConfigDb)
	find, _ := cache.Get(name)
	if find != nil {
		model = find.(*ConfigDbModel)
	}
	return
}

func (this_ *Application) GetConfigZkList() []*ConfigZkModel {
	return this_.configZkList
}

func (this_ *Application) GetConfigZk(name string) (model *ConfigZkModel) {
	cache := this_.getModelTypeCache(TypeConfigZk)
	find, _ := cache.Get(name)
	if find != nil {
		model = find.(*ConfigZkModel)
	}
	return
}

func (this_ *Application) GetConfigKafkaList() []*ConfigKafkaModel {
	return this_.configKafkaList
}

func (this_ *Application) GetConfigKafka(name string) (model *ConfigKafkaModel) {
	cache := this_.getModelTypeCache(TypeConfigKafka)
	find, _ := cache.Get(name)
	if find != nil {
		model = find.(*ConfigKafkaModel)
	}
	return
}

func (this_ *Application) GetConfigMongodbList() []*ConfigMongodbModel {
	return this_.configMongodbList
}

func (this_ *Application) GetConfigMongodb(name string) (model *ConfigMongodbModel) {
	cache := this_.getModelTypeCache(TypeConfigMongodb)
	find, _ := cache.Get(name)
	if find != nil {
		model = find.(*ConfigMongodbModel)
	}
	return
}

func (this_ *Application) GetConfigElasticsearchList() []*ConfigEsModel {
	return this_.configElasticsearchList
}

func (this_ *Application) GetConfigElasticsearch(name string) (model *ConfigEsModel) {
	cache := this_.getModelTypeCache(TypeConfigElasticsearch)
	find, _ := cache.Get(name)
	if find != nil {
		model = find.(*ConfigEsModel)
	}
	return
}

func (this_ *Application) GetErrorList() []*ErrorModel {
	return this_.errorList
}

func (this_ *Application) GetError(name string) (model *ErrorModel) {
	cache := this_.getModelTypeCache(TypeError)
	find, _ := cache.Get(name)
	if find != nil {
		model = find.(*ErrorModel)
	}
	return
}

func (this_ *Application) GetFuncList() []*FuncModel {
	return this_.funcList
}

func (this_ *Application) GetFunc(name string) (model *FuncModel) {
	cache := this_.getModelTypeCache(TypeFunc)
	find, _ := cache.Get(name)
	if find != nil {
		model = find.(*FuncModel)
	}
	return
}

func (this_ *Application) GetLanguageJavascript() (model *LanguageJavascriptModel) {
	if this_.languageJavascript == nil {
		this_.languageJavascript = &LanguageJavascriptModel{}
	}
	model = this_.languageJavascript
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
