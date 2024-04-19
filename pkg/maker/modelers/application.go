package modelers

import (
	"errors"
	"github.com/team-ide/go-tool/util"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
)

type Application struct {
	dir string

	Children []*Element `json:"children"`

	elementCache map[string]*Element

	modelTypeCaches map[*Type]*util.Cache
	modelTypeItems  map[*Type][]ElementIFace

	doLocker    sync.Mutex
	cacheLocker sync.Mutex

	LoadErrors []*LoadError `json:"loadErrors"`
}

func (this_ *Application) GetDir() string {
	return this_.dir
}

type ElementNode struct {
	Name    string `json:"name,omitempty"` // 名称，同一个应用中唯一
	element *Element
}

func (this_ *ElementNode) GetName() string {
	if this_ == nil {
		return ""
	}
	return this_.Name
}

func (this_ *ElementNode) SetName(name string) {
	if this_ == nil {
		return
	}
	this_.Name = name
}

func (this_ *ElementNode) GetElement() *Element {
	if this_ == nil {
		return nil
	}
	return this_.element
}

func (this_ *ElementNode) SetElement(element *Element) {
	if this_ == nil {
		return
	}
	this_.element = element
}

type ElementIFace interface {
	GetName() string
	SetName(name string)
	GetElement() *Element
	SetElement(element *Element)
}

type Element struct {
	Key      string     `json:"key,omitempty"`
	Text     string     `json:"text,omitempty"`
	IsType   bool       `json:"isType,omitempty"`
	IsPack   bool       `json:"isPack,omitempty"`
	IsModel  bool       `json:"isModel,omitempty"`
	Children []*Element `json:"children,omitempty"`
	Pack     *Pack      `json:"pack,omitempty"`
	parent   *Element
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

func (this_ *Application) getElement(key string) (element *Element) {
	this_.cacheLocker.Lock()
	defer this_.cacheLocker.Unlock()
	element = this_.elementCache[key]
	return
}

func (this_ *Application) setElement(element *Element) {
	this_.cacheLocker.Lock()
	defer this_.cacheLocker.Unlock()
	this_.elementCache[element.Key] = element
	return
}

func (this_ *Application) removeElement(modelType *Type, key string) {
	items := this_.getModelTypeItems(modelType)
	cache := this_.getModelTypeCache(modelType)

	this_.cacheLocker.Lock()
	defer this_.cacheLocker.Unlock()

	var deleteKeys []string

	itemCache := make(map[string]ElementIFace)
	for _, one := range items {
		itemCache[one.GetElement().Key] = one
	}
	findElement := this_.elementCache[key]
	for k := range this_.elementCache {
		if k == key || strings.HasPrefix(k, key+"/") {
			deleteKeys = append(deleteKeys, k)

			model := itemCache[k]
			if model != nil {
				delete(itemCache, k)
				cache.Delete(model.GetName())
			}
		}
	}
	for _, k := range deleteKeys {
		delete(this_.elementCache, k)
	}
	var newItems []ElementIFace
	for _, one := range itemCache {
		newItems = append(newItems, one)
	}
	this_.modelTypeItems[modelType] = newItems

	if findElement != nil && findElement.parent != nil {
		findElement.parent.Children = removeElement(findElement.parent.Children, key)
	} else {
		this_.Children = removeElement(this_.Children, key)
	}

	return
}

func removeElement(list []*Element, removeKey string) (newList []*Element) {
	for _, one := range list {
		if one.Key != removeKey {
			newList = append(newList, one)
		}
	}
	return
}

func (this_ *Application) GetModelTypeModel(key string, name string) (model interface{}) {
	modelType := GetModelType(key)
	if modelType == nil {
		util.Logger.Error("model type " + key + " not found")
		return
	}
	cache := this_.getModelTypeCache(modelType)
	if modelType.IsFile {
		model, _ = cache.Get("")
	} else {
		model, _ = cache.Get(name)
	}
	return
}

func (this_ *Application) getModeTypePath(modelType *Type) (path string) {
	path = this_.dir + modelType.Name
	if !modelType.IsFile {
		path += "/"
	}
	return
}

func (this_ *Application) getModePath(modelType *Type, name string, isPack bool) (path string, err error) {
	var modelTypePath = this_.getModeTypePath(modelType)
	if modelType.IsFile {
		path = modelTypePath
	} else {
		if strings.Contains(name, "./") {
			err = errors.New("model name[" + name + "] has [./] is error")
			return
		}
		path = modelTypePath + name
		path = util.FormatPath(path)
		var isSub bool
		isSub, err = util.IsSubPath(modelTypePath, path)
		if err != nil {
			return
		}
		if !isSub {
			err = errors.New("model name[" + name + "] is error")
			return
		}
	}
	if !isPack {
		path += ".yml"
	}
	return
}

func (this_ *Application) getModeKey(modelType *Type, name string) (key string) {
	key = modelType.Name
	if !modelType.IsFile {
		key += "/" + name
	}
	return
}

func FormatName(name string) (res string) {
	res = strings.ReplaceAll(name, "\\", "/")
	re, _ := regexp.Compile("/+")
	res = re.ReplaceAllLiteralString(res, "/")
	if strings.HasSuffix(res, "/") {
		res = res[0 : len(res)-1]
	}
	if strings.HasPrefix(res, "/") {
		res = res[1:]
	}
	return
}

func (this_ *Application) Remove(modelType *Type, modelName string, isPack bool) (element *Element, err error) {
	modelName = FormatName(modelName)
	if modelName == "" {
		err = errors.New("model name is empty")
		return
	}
	modelKey := this_.getModeKey(modelType, modelName)

	this_.doLocker.Lock()
	defer this_.doLocker.Unlock()

	element = this_.getElement(modelKey)
	if element == nil {
		err = errors.New("model [" + modelKey + "] not found")
		return
	}
	modelPath, err := this_.getModePath(modelType, modelName, isPack)
	if err != nil {
		return
	}

	exist, err := util.PathExists(modelPath)
	if err != nil {
		return
	}
	if exist {
		err = os.RemoveAll(modelPath)
	}

	this_.removeElement(modelType, element.Key)

	return
}

func (this_ *Application) Rename(modelType *Type, oldModelName string, newModelName string, isPack bool) (oldElement *Element, newElement *Element, err error) {
	oldModelName = FormatName(oldModelName)
	if oldModelName == "" {
		err = errors.New("old model name is empty")
		return
	}
	newModelName = FormatName(newModelName)
	if newModelName == "" {
		err = errors.New("new model name is empty")
		return
	}
	oldModelKey := this_.getModeKey(modelType, oldModelName)
	newModelKey := this_.getModeKey(modelType, newModelName)

	this_.doLocker.Lock()
	defer this_.doLocker.Unlock()

	oldElement = this_.getElement(oldModelKey)
	if oldElement == nil {
		err = errors.New("old model [" + oldModelKey + "] element not found")
		return
	}
	newElement = this_.getElement(newModelKey)
	if newElement != nil {
		err = errors.New("new model [" + newModelKey + "] element is exist")
		return
	}
	oldModelPath, err := this_.getModePath(modelType, oldModelName, isPack)
	if err != nil {
		return
	}
	newModelPath, err := this_.getModePath(modelType, newModelName, isPack)
	if err != nil {
		return
	}

	oldExist, err := util.PathExists(oldModelPath)
	if err != nil {
		return
	}
	if !oldExist {
		err = errors.New("old model [" + oldModelKey + "] file not found")
		return
	}

	newExist, err := util.PathExists(newModelPath)
	if err != nil {
		return
	}
	if newExist {
		err = errors.New("new model [" + newModelKey + "] file is exist")
		return
	}
	err = os.Rename(oldModelPath, newModelPath)
	if err != nil {
		return
	}

	this_.removeElement(modelType, oldElement.Key)

	packFullName := newModelName
	if !isPack {
		if strings.LastIndex(packFullName, "/") > 0 {
			packFullName = packFullName[:strings.LastIndex(packFullName, "/")]
		} else {
			packFullName = ""
		}
	}
	if modelType.IsFile {
		packFullName = ""
	}

	parentElement, err := this_.appendPackByName(modelType, packFullName)
	if err != nil {
		return
	}

	if isPack {
		newElement = parentElement
	} else {
		_, newElement = this_.loadFile(parentElement, modelType, newModelPath)
	}

	return
}

func (this_ *Application) appendPackByName(modelType *Type, packFullName string) (element *Element, err error) {
	path := this_.getModeTypePath(modelType)

	var exist bool
	if exist, err = util.PathExists(path); err != nil {
		return
	}
	if !exist {
		if err = os.MkdirAll(path, os.ModePerm); err != nil {
			return
		}
	}
	var parentElement = this_.getElement(modelType.Name)
	if packFullName != "" {
		names := strings.Split(packFullName, "/")
		key := modelType.Name
		for i := 0; i < len(names); i++ {
			key += "/" + names[i]
			path += "/" + names[i]
			element = this_.getElement(key)
			if element == nil {
				parentElement = this_.appendPack(parentElement, modelType, names[i])
			} else {
				parentElement = element
			}
			if exist, err = util.PathExists(path); err != nil {
				return
			}
			if !exist {
				if err = os.MkdirAll(path, os.ModePerm); err != nil {
					return
				}
			}
		}
	}
	element = parentElement
	return
}
func (this_ *Application) Save(modelType *Type, modelName string, model interface{}, isPack, isNew bool) (res interface{}, element *Element, err error) {
	modelName = FormatName(modelName)
	if modelName == "" {
		err = errors.New("model name is empty")
		return
	}
	modelPath, err := this_.getModePath(modelType, modelName, isPack)
	if err != nil {
		return
	}

	this_.doLocker.Lock()
	defer this_.doLocker.Unlock()

	exist, err := util.PathExists(modelPath)
	if err != nil {
		return
	}
	if isNew && exist {
		err = errors.New("model [" + modelType.Name + "] [" + modelName + "] is exist")
		return
	}

	packFullName := modelName
	if !isPack {
		if strings.LastIndex(modelName, "/") > 0 {
			packFullName = modelName[:strings.LastIndex(modelName, "/")]
		} else {
			packFullName = ""
		}
	}
	if modelType.IsFile {
		packFullName = ""
	}

	parentElement, err := this_.appendPackByName(modelType, packFullName)
	if err != nil {
		return
	}

	// 添加 pack

	if !isPack {
		var text string
		if text, err = modelType.toText(model); err != nil {
			return
		}
		var f *os.File
		if f, err = os.Create(modelPath); err != nil {
			return
		}
		defer func() { _ = f.Close() }()
		if _, err = f.WriteString(text); err != nil {
			return
		}
		res, element = this_.loadFile(parentElement, modelType, modelPath)
	} else {
		element = parentElement
	}
	return
}

func (this_ *Application) GetModelTypeModels(key string) (models interface{}) {
	modelType := GetModelType(key)
	if modelType == nil {
		util.Logger.Error("model type " + key + " not found")
		return
	}

	this_.cacheLocker.Lock()
	defer this_.cacheLocker.Unlock()

	models = this_.modelTypeItems[modelType]

	return
}

func (this_ *Application) appendType(parent *Element, modelType *Type) (element *Element) {

	element = &Element{
		Key:  modelType.Name,
		Text: modelType.Comment,
	}
	if parent != nil {
		parent.Children = append(parent.Children, element)
		element.parent = parent
	} else {
		this_.Children = append(this_.Children, element)
	}
	element.IsType = true
	this_.setElement(element)
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
		element.parent = parent
	} else {
		this_.Children = append(this_.Children, element)
	}
	element.IsPack = true
	element.Pack = &Pack{
		Name: name,
	}
	this_.setElement(element)
	return
}

func (this_ *Application) appendModel(parent *Element, modelType *Type, fileName string, name string, model interface{}) (element *Element, err error) {
	model_, ok := model.(ElementIFace)
	if !ok {
		err = errors.New("type " + modelType.Name + " model [" + name + "] can not to ElementIFace")
		return
	}
	var key = parent.Key
	if !modelType.IsFile {
		key = parent.Key + "/" + fileName
	}
	element = parent
	if !modelType.IsFile {
		element = &Element{
			Key:  key,
			Text: fileName,
		}
		element.parent = parent
	}
	element.IsModel = true
	model_.SetElement(element)

	cache := this_.modelTypeCaches[modelType]
	if cache == nil {
		cache = util.NewCache()
		this_.modelTypeCaches[modelType] = cache
	}

	//if _, ok := cache.Get(name); ok {
	//	err = errors.New("type " + modelType.Name + " model [" + name + "] already exist")
	//	return
	//}

	if modelType.IsFile {
		cache.Put("", model)
	} else {
		cache.Put(name, model)
	}
	this_.setElement(element)

	var list = this_.modelTypeItems[modelType]
	var newList []ElementIFace
	for _, one := range list {
		if one.GetElement().Key != model_.GetElement().Key {
			newList = append(newList, one)
		}
	}
	newList = append(newList, model_)
	sort.SliceStable(newList, func(i, j int) bool {
		return newList[i].GetElement().Key < newList[j].GetElement().Key
	})
	this_.modelTypeItems[modelType] = newList

	var Children = element.parent.Children
	var newChildren []*Element
	for _, one := range Children {
		if one.Key != model_.GetElement().Key {
			newChildren = append(newChildren, one)
		}
	}
	newChildren = append(newChildren, element)
	sort.SliceStable(newChildren, func(i, j int) bool {
		return newChildren[i].Key < newChildren[j].Key
	})
	element.parent.Children = newChildren
	return
}

func (this_ *Application) getModelTypeCache(modelType *Type) (cache *util.Cache) {
	this_.cacheLocker.Lock()
	defer this_.cacheLocker.Unlock()

	cache = this_.modelTypeCaches[modelType]
	if cache == nil {
		cache = util.NewCache()
		this_.modelTypeCaches[modelType] = cache
	}
	return
}

func (this_ *Application) getModelTypeItems(modelType *Type) (items []ElementIFace) {
	this_.cacheLocker.Lock()
	defer this_.cacheLocker.Unlock()

	items = this_.modelTypeItems[modelType]
	return
}

func (this_ *Application) GetConstantList() (res []*ConstantModel) {
	items := this_.getModelTypeItems(TypeConstant)
	for _, one := range items {
		res = append(res, one.(*ConstantModel))
	}
	return
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

func (this_ *Application) GetConfigRedisList() (res []*ConfigRedisModel) {
	items := this_.getModelTypeItems(TypeConfigRedis)
	for _, one := range items {
		res = append(res, one.(*ConfigRedisModel))
	}
	return
}

func (this_ *Application) GetConfigRedis(name string) (model *ConfigRedisModel) {
	cache := this_.getModelTypeCache(TypeConfigRedis)
	find, _ := cache.Get(name)
	if find != nil {
		model = find.(*ConfigRedisModel)
	}
	return
}

func (this_ *Application) GetConfigDbList() (res []*ConfigDbModel) {
	items := this_.getModelTypeItems(TypeConfigDb)
	for _, one := range items {
		res = append(res, one.(*ConfigDbModel))
	}
	return
}

func (this_ *Application) GetConfigDb(name string) (model *ConfigDbModel) {
	cache := this_.getModelTypeCache(TypeConfigDb)
	find, _ := cache.Get(name)
	if find != nil {
		model = find.(*ConfigDbModel)
	}
	return
}

func (this_ *Application) GetConfigZkList() (res []*ConfigZkModel) {
	items := this_.getModelTypeItems(TypeConfigZk)
	for _, one := range items {
		res = append(res, one.(*ConfigZkModel))
	}
	return
}

func (this_ *Application) GetConfigZk(name string) (model *ConfigZkModel) {
	cache := this_.getModelTypeCache(TypeConfigZk)
	find, _ := cache.Get(name)
	if find != nil {
		model = find.(*ConfigZkModel)
	}
	return
}

func (this_ *Application) GetConfigKafkaList() (res []*ConfigKafkaModel) {
	items := this_.getModelTypeItems(TypeConfigKafka)
	for _, one := range items {
		res = append(res, one.(*ConfigKafkaModel))
	}
	return
}

func (this_ *Application) GetConfigKafka(name string) (model *ConfigKafkaModel) {
	cache := this_.getModelTypeCache(TypeConfigKafka)
	find, _ := cache.Get(name)
	if find != nil {
		model = find.(*ConfigKafkaModel)
	}
	return
}

func (this_ *Application) GetConfigMongodbList() (res []*ConfigMongodbModel) {
	items := this_.getModelTypeItems(TypeConfigMongodb)
	for _, one := range items {
		res = append(res, one.(*ConfigMongodbModel))
	}
	return
}

func (this_ *Application) GetConfigMongodb(name string) (model *ConfigMongodbModel) {
	cache := this_.getModelTypeCache(TypeConfigMongodb)
	find, _ := cache.Get(name)
	if find != nil {
		model = find.(*ConfigMongodbModel)
	}
	return
}

func (this_ *Application) GetConfigElasticsearchList() (res []*ConfigEsModel) {
	items := this_.getModelTypeItems(TypeConfigElasticsearch)
	for _, one := range items {
		res = append(res, one.(*ConfigEsModel))
	}
	return
}

func (this_ *Application) GetConfigElasticsearch(name string) (model *ConfigEsModel) {
	cache := this_.getModelTypeCache(TypeConfigElasticsearch)
	find, _ := cache.Get(name)
	if find != nil {
		model = find.(*ConfigEsModel)
	}
	return
}

func (this_ *Application) GetErrorList() (res []*ErrorModel) {
	items := this_.getModelTypeItems(TypeError)
	for _, one := range items {
		res = append(res, one.(*ErrorModel))
	}
	return
}

func (this_ *Application) GetError(name string) (model *ErrorModel) {
	cache := this_.getModelTypeCache(TypeError)
	find, _ := cache.Get(name)
	if find != nil {
		model = find.(*ErrorModel)
	}
	return
}

func (this_ *Application) GetFuncList() (res []*FuncModel) {
	items := this_.getModelTypeItems(TypeFunc)
	for _, one := range items {
		res = append(res, one.(*FuncModel))
	}
	return
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
	items := this_.getModelTypeItems(TypeFunc)
	if len(items) > 0 {
		model = items[0].(*LanguageJavascriptModel)
	}
	if model == nil {
		model = &LanguageJavascriptModel{}
	}
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
