package maker

import (
	"errors"
	"github.com/team-ide/go-tool/util"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"teamide/pkg/maker/modelers"
)

func newApplication() (app *Application) {
	app = &Application{
		elementCache:    make(map[string]*modelers.Element),
		modelTypeCaches: make(map[*modelers.Type]*util.Cache[any]),
		modelTypeItems:  make(map[*modelers.Type][]modelers.ElementIFace),

		constantContext: make(map[string]interface{}),
		errorContext:    make(map[string]interface{}),
		strictContext:   make(map[string]interface{}),
		tableContext:    make(map[string]interface{}),

		daoContext: make(map[string]interface{}),
		daoProgram: make(map[string]*CompileProgram),

		serviceContext: make(map[string]interface{}),
		serviceProgram: make(map[string]*CompileProgram),

		funcContext: make(map[string]interface{}),
		funcProgram: make(map[string]*CompileProgram),

		typeContext: make(map[string]*ValueType),
	}
	return
}

type Application struct {
	dir string

	Children []*modelers.Element `json:"children"`

	elementCache map[string]*modelers.Element

	modelTypeCaches map[*modelers.Type]*util.Cache[any]
	modelTypeItems  map[*modelers.Type][]modelers.ElementIFace

	doLocker    sync.Mutex
	cacheLocker sync.Mutex

	LoadErrors []*LoadError `json:"loadErrors"`

	constantContext map[string]interface{}
	errorContext    map[string]interface{}
	strictContext   map[string]interface{}
	tableContext    map[string]interface{}

	daoContext map[string]interface{}
	daoProgram map[string]*CompileProgram

	serviceContext map[string]interface{}
	serviceProgram map[string]*CompileProgram

	funcContext map[string]interface{}
	funcProgram map[string]*CompileProgram

	typeContext map[string]*ValueType
}

func (this_ *Application) GetDir() string {
	return this_.dir
}

type LoadError struct {
	Type  *modelers.Type `json:"type"`
	Path  string         `json:"path"`
	Error string         `json:"error"`
}

func (this_ *Application) getElement(key string) (element *modelers.Element) {
	this_.cacheLocker.Lock()
	defer this_.cacheLocker.Unlock()
	element = this_.elementCache[key]
	return
}

func (this_ *Application) setElement(element *modelers.Element) {
	this_.cacheLocker.Lock()
	defer this_.cacheLocker.Unlock()
	this_.elementCache[element.Key] = element
	return
}

func (this_ *Application) removeElement(modelType *modelers.Type, key string) {
	items := this_.getModelTypeItems(modelType)
	cache := this_.getModelTypeCache(modelType)

	this_.cacheLocker.Lock()
	defer this_.cacheLocker.Unlock()

	var deleteKeys []string

	itemCache := make(map[string]modelers.ElementIFace)
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
	var newItems []modelers.ElementIFace
	for _, one := range itemCache {
		newItems = append(newItems, one)
	}
	this_.modelTypeItems[modelType] = newItems

	if findElement != nil && findElement.GetParent() != nil {
		findElement.GetParent().Children = removeElement(findElement.GetParent().Children, key)
	} else {
		this_.Children = removeElement(this_.Children, key)
	}

	return
}

func removeElement(list []*modelers.Element, removeKey string) (newList []*modelers.Element) {
	for _, one := range list {
		if one.Key != removeKey {
			newList = append(newList, one)
		}
	}
	return
}

func (this_ *Application) GetModelTypeModel(key string, name string) (model interface{}) {
	modelType := modelers.GetModelType(key)
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

func (this_ *Application) getModeTypePath(modelType *modelers.Type) (path string) {
	path = this_.dir + modelType.Name
	if !modelType.IsFile {
		path += "/"
	}
	return
}

func (this_ *Application) getModePath(modelType *modelers.Type, name string, isPack bool) (path string, err error) {
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

func (this_ *Application) getModeKey(modelType *modelers.Type, name string) (key string) {
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

func (this_ *Application) Remove(modelType *modelers.Type, modelName string, isPack bool) (element *modelers.Element, err error) {
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

func (this_ *Application) Rename(modelType *modelers.Type, oldModelName string, newModelName string, isPack bool) (oldElement *modelers.Element, newElement *modelers.Element, err error) {
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

func (this_ *Application) appendPackByName(modelType *modelers.Type, packFullName string) (element *modelers.Element, err error) {
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
func (this_ *Application) Save(modelType *modelers.Type, modelName string, model interface{}, isPack, isNew bool) (res interface{}, element *modelers.Element, err error) {
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
		if text, err = modelType.ToText(model); err != nil {
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
	modelType := modelers.GetModelType(key)
	if modelType == nil {
		util.Logger.Error("model type " + key + " not found")
		return
	}

	this_.cacheLocker.Lock()
	defer this_.cacheLocker.Unlock()

	models = this_.modelTypeItems[modelType]

	return
}

func (this_ *Application) appendType(parent *modelers.Element, modelType *modelers.Type) (element *modelers.Element) {

	element = &modelers.Element{
		Key:  modelType.Name,
		Text: modelType.Comment,
	}
	if parent != nil {
		parent.Children = append(parent.Children, element)
		element.SetParent(parent)
	} else {
		this_.Children = append(this_.Children, element)
	}
	element.IsType = true
	this_.setElement(element)
	return
}

func (this_ *Application) appendPack(parent *modelers.Element, modelType *modelers.Type, name string) (element *modelers.Element) {

	element = &modelers.Element{
		Key:  name,
		Text: name,
	}
	if parent != nil {
		element.Key = parent.Key + "/" + element.Key
		parent.Children = append(parent.Children, element)
		element.SetParent(parent)
	} else {
		this_.Children = append(this_.Children, element)
	}
	element.IsPack = true
	element.Pack = &modelers.Pack{
		Name: name,
	}
	this_.setElement(element)
	return
}

func (this_ *Application) appendModel(parent *modelers.Element, modelType *modelers.Type, fileName string, name string, model interface{}) (element *modelers.Element, err error) {
	model_, ok := model.(modelers.ElementIFace)
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
		element = &modelers.Element{
			Key:  key,
			Text: fileName,
		}
		element.SetParent(parent)
	}
	element.IsModel = true
	model_.SetElement(element)

	cache := this_.modelTypeCaches[modelType]
	if cache == nil {
		cache = util.NewCache(modelType.NewModel())
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
	var newList []modelers.ElementIFace
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

	var Children = this_.Children
	if element.GetParent() != nil {
		Children = element.GetParent().Children
	}
	var newChildren []*modelers.Element
	for _, one := range Children {
		if one.Key != model_.GetElement().Key {
			newChildren = append(newChildren, one)
		}
	}
	newChildren = append(newChildren, element)
	sort.SliceStable(newChildren, func(i, j int) bool {
		return newChildren[i].Key < newChildren[j].Key
	})
	if element.GetParent() != nil {
		element.GetParent().Children = newChildren
	} else {
		this_.Children = newChildren
	}
	return
}

func (this_ *Application) getModelTypeCache(modelType *modelers.Type) (cache *util.Cache[any]) {
	this_.cacheLocker.Lock()
	defer this_.cacheLocker.Unlock()

	cache = this_.modelTypeCaches[modelType]
	if cache == nil {
		cache = util.NewCache(modelType.NewModel())
		this_.modelTypeCaches[modelType] = cache
	}
	return
}

func (this_ *Application) getModelTypeItems(modelType *modelers.Type) (items []modelers.ElementIFace) {
	this_.cacheLocker.Lock()
	defer this_.cacheLocker.Unlock()

	items = this_.modelTypeItems[modelType]
	return
}

func (this_ *Application) GetConstantList() (res []*modelers.ConstantModel) {
	items := this_.getModelTypeItems(modelers.TypeConstant)
	for _, one := range items {
		res = append(res, one.(*modelers.ConstantModel))
	}
	return
}

func (this_ *Application) GetConstant(name string) (model *modelers.ConstantModel) {
	cache := this_.getModelTypeCache(modelers.TypeConstant)
	find, _ := cache.Get(name)
	if find != nil {
		model = find.(*modelers.ConstantModel)
	}
	return
}

func (this_ *Application) GetStruct(name string) (model *modelers.StructModel) {
	cache := this_.getModelTypeCache(modelers.TypeStruct)
	find, _ := cache.Get(name)
	if find != nil {
		model = find.(*modelers.StructModel)
	}
	return
}

func (this_ *Application) GetStructList() (res []*modelers.StructModel) {
	items := this_.getModelTypeItems(modelers.TypeStruct)
	for _, one := range items {
		res = append(res, one.(*modelers.StructModel))
	}
	return
}

func (this_ *Application) GetDao(name string) (model *modelers.DaoModel) {
	cache := this_.getModelTypeCache(modelers.TypeDao)
	find, _ := cache.Get(name)
	if find != nil {
		model = find.(*modelers.DaoModel)
	}
	return
}

func (this_ *Application) GetDaoList() (res []*modelers.DaoModel) {
	items := this_.getModelTypeItems(modelers.TypeDao)
	for _, one := range items {
		res = append(res, one.(*modelers.DaoModel))
	}
	return
}

func (this_ *Application) GetService(name string) (model *modelers.ServiceModel) {
	cache := this_.getModelTypeCache(modelers.TypeService)
	find, _ := cache.Get(name)
	if find != nil {
		model = find.(*modelers.ServiceModel)
	}
	return
}

func (this_ *Application) GetServiceList() (res []*modelers.ServiceModel) {
	items := this_.getModelTypeItems(modelers.TypeService)
	for _, one := range items {
		res = append(res, one.(*modelers.ServiceModel))
	}
	return
}

func (this_ *Application) GetConfigDbList() (res []*modelers.ConfigDbModel) {
	app := this_.GetApp()
	if app.Db != nil {
		res = append(res, app.Db)
	}
	if app.DbOther != nil {
		for _, one := range app.DbOther {
			res = append(res, one)
		}
	}
	return
}

func (this_ *Application) GetConfigRedisList() (res []*modelers.ConfigRedisModel) {
	app := this_.GetApp()
	if app.Redis != nil {
		res = append(res, app.Redis)
	}
	if app.RedisOther != nil {
		for _, one := range app.RedisOther {
			res = append(res, one)
		}
	}
	return
}

func (this_ *Application) GetConfigZkList() (res []*modelers.ConfigZkModel) {
	app := this_.GetApp()
	if app.Zk != nil {
		res = append(res, app.Zk)
	}
	if app.ZkOther != nil {
		for _, one := range app.ZkOther {
			res = append(res, one)
		}
	}
	return
}

func (this_ *Application) GetConfigKafkaList() (res []*modelers.ConfigKafkaModel) {
	app := this_.GetApp()
	if app.Kafka != nil {
		res = append(res, app.Kafka)
	}
	if app.KafkaOther != nil {
		for _, one := range app.KafkaOther {
			res = append(res, one)
		}
	}
	return
}

func (this_ *Application) GetConfigMongodbList() (res []*modelers.ConfigMongodbModel) {
	app := this_.GetApp()
	if app.Mongodb != nil {
		res = append(res, app.Mongodb)
	}
	if app.MongodbOther != nil {
		for _, one := range app.MongodbOther {
			res = append(res, one)
		}
	}
	return
}

func (this_ *Application) GetConfigEsList() (res []*modelers.ConfigEsModel) {
	app := this_.GetApp()
	if app.Es != nil {
		res = append(res, app.Es)
	}
	if app.EsOther != nil {
		for _, one := range app.EsOther {
			res = append(res, one)
		}
	}
	return
}

func (this_ *Application) GetErrorList() (res []*modelers.ErrorModel) {
	items := this_.getModelTypeItems(modelers.TypeError)
	for _, one := range items {
		res = append(res, one.(*modelers.ErrorModel))
	}
	return
}

func (this_ *Application) GetError(name string) (model *modelers.ErrorModel) {
	cache := this_.getModelTypeCache(modelers.TypeError)
	find, _ := cache.Get(name)
	if find != nil {
		model = find.(*modelers.ErrorModel)
	}
	return
}

func (this_ *Application) GetFuncList() (res []*modelers.FuncModel) {
	items := this_.getModelTypeItems(modelers.TypeFunc)
	for _, one := range items {
		res = append(res, one.(*modelers.FuncModel))
	}
	return
}

func (this_ *Application) GetFunc(name string) (model *modelers.FuncModel) {
	cache := this_.getModelTypeCache(modelers.TypeFunc)
	find, _ := cache.Get(name)
	if find != nil {
		model = find.(*modelers.FuncModel)
	}
	return
}

func (this_ *Application) GetLanguageJavascript() (model *modelers.LanguageJavascriptModel) {
	items := this_.getModelTypeItems(modelers.TypeLanguageJavascript)
	model = items[0].(*modelers.LanguageJavascriptModel)
	return
}

func (this_ *Application) GetLanguageGolang() (model *modelers.LanguageGolangModel) {
	items := this_.getModelTypeItems(modelers.TypeLanguageGolang)
	model = items[0].(*modelers.LanguageGolangModel)
	return
}

func (this_ *Application) GetApp() (model *modelers.AppModel) {
	items := this_.getModelTypeItems(modelers.TypeApp)
	model = items[0].(*modelers.AppModel)
	return
}

func (this_ *Application) GetValueType(name string) (valueType *ValueType, err error) {
	if name == "" {
		name = "string"
	}
	valueType = GetValueType(name)
	if valueType != nil {
		return
	}

	valueType = this_.typeContext[name]
	if valueType != nil {
		return
	}
	valueType = &ValueType{
		Name:       name,
		FieldTypes: make(map[string]*ValueType),
	}
	valueType.Struct = this_.GetStruct(valueType.Name)
	if valueType.Struct == nil {
		err = errors.New("value type and struct not found name [" + valueType.Name + "]")
		return
	}
	this_.typeContext[name] = valueType
	for _, field := range valueType.Struct.Fields {
		if field.Type == "" {
			field.Type = "string"
		}
		valueType.FieldTypes[field.Name], err = this_.GetValueType(field.Type)
		if err != nil {
			return
		}
	}
	return
}
