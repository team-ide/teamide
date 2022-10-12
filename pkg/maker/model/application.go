package model

import "errors"

type Application struct {
	StructList   []*StructModel `json:"structList"`
	structCache  map[string]*StructModel
	ServiceList  []*ServiceModel `json:"serviceList"`
	serviceCache map[string]*ServiceModel
	ErrorList    []*ErrorModel `json:"errorList"`
	errorCache   map[string]*ErrorModel
	LoadErrors   []*LoadError `json:"loadErrors"`
}

type LoadError struct {
	Type  *Type  `json:"type"`
	Path  string `json:"path"`
	Error string `json:"error"`
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
