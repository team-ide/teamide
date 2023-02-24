package invokers

import (
	"errors"
	"teamide/pkg/maker/modelers"
)

type InvokeData struct {
	app      *modelers.Application
	args     []*InvokeVar
	vars     []*InvokeVar
	argCache map[string]*InvokeVar
	varCache map[string]*InvokeVar
}

func NewInvokeData(app *modelers.Application) (data *InvokeData) {
	data = &InvokeData{
		app:      app,
		argCache: make(map[string]*InvokeVar),
		varCache: make(map[string]*InvokeVar),
	}
	return
}

type InvokeVar struct {
	Name      string              `json:"name,omitempty"`
	Value     interface{}         `json:"value,omitempty"`
	ValueType *modelers.ValueType `json:"valueType,omitempty"`
}

func (this_ *InvokeData) GetArgs() (args []*InvokeVar) {
	args = this_.args
	return
}

func (this_ *InvokeData) GetVars() (vars []*InvokeVar) {
	vars = this_.vars
	return
}

func (this_ *InvokeData) AddArg(name string, value interface{}, valueType *modelers.ValueType) (err error) {
	err = this_.addArg(&InvokeVar{
		Name:      name,
		Value:     value,
		ValueType: valueType,
	})
	return
}

func (this_ *InvokeData) AddStructArg(name string, value interface{}, structName string) (err error) {
	valueType := &modelers.ValueType{
		Name: structName,
	}
	valueType.Struct = this_.app.GetStruct(valueType.Name)
	err = this_.addArg(&InvokeVar{
		Name:      name,
		Value:     value,
		ValueType: valueType,
	})
	return
}

func (this_ *InvokeData) addArg(arg *InvokeVar) (err error) {
	if this_.argCache[arg.Name] != nil {
		err = errors.New("arg [" + arg.Name + "] already exist")
		return
	}
	err = this_.formatInvokeVar(arg)
	if err != nil {
		return
	}
	this_.args = append(this_.args, arg)
	this_.argCache[arg.Name] = arg
	return
}

func (this_ *InvokeData) addVar(var_ *InvokeVar) (err error) {
	if this_.varCache[var_.Name] != nil {
		err = errors.New("arg [" + var_.Name + "] already exist")
		return
	}
	err = this_.formatInvokeVar(var_)
	if err != nil {
		return
	}
	this_.vars = append(this_.vars, var_)
	this_.varCache[var_.Name] = var_
	return
}

func (this_ *InvokeData) formatInvokeVar(invokeVar *InvokeVar) (err error) {

	return
}
