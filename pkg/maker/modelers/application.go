package modelers

import "errors"

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
