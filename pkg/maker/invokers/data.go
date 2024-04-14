package invokers

import (
	"errors"
	"fmt"
	"github.com/team-ide/go-tool/javascript"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"regexp"
	"sync"
	"teamide/pkg/maker/modelers"
)

type InvokeData struct {
	app          *modelers.Application
	args         []*InvokeVar
	vars         []*InvokeVar
	argCache     map[string]*InvokeVar
	argCacheLock sync.Locker
	varCache     map[string]*InvokeVar
	varCacheLock sync.Locker
	script       *javascript.Script
	scriptLock   sync.Locker
}

type InvokeVar struct {
	Name      string              `json:"name,omitempty"`
	Value     interface{}         `json:"value,omitempty"`
	ValueType *modelers.ValueType `json:"valueType,omitempty"`
}

func NewInvokeData(app *modelers.Application) (data *InvokeData, err error) {
	script, err := javascript.NewScript()
	if err != nil {
		util.Logger.Error("NewInvokeData NewScript error", zap.Any("error", err))
		return
	}
	data = &InvokeData{
		app:          app,
		argCache:     make(map[string]*InvokeVar),
		argCacheLock: &sync.Mutex{},
		varCache:     make(map[string]*InvokeVar),
		varCacheLock: &sync.Mutex{},
		script:       script,
		scriptLock:   &sync.Mutex{},
	}
	err = data.init()
	if err != nil {
		return nil, err
	}
	return
}

func (this_ *InvokeData) init() (err error) {
	// 将 常量 error func 填充 至 script 变量域中
	for _, one := range this_.app.GetConstantList() {
		for _, option := range one.Options {
			err = this_.scriptSet(option.Name, option.Value)
			if err != nil {
				util.Logger.Error("invoke data init set constant value error", zap.Any("name", option.Name), zap.Any("value", option.Value), zap.Any("error", err))
				return
			}
		}
	}

	for _, one := range this_.app.GetErrorList() {
		for _, option := range one.Options {
			err = this_.scriptSet(option.Name, option)
			if err != nil {
				util.Logger.Error("invoke data init set error value error", zap.Any("name", option.Name), zap.Any("error", option), zap.Any("error", err))
				return
			}
		}
	}

	for _, one := range this_.app.GetFuncList() {
		err = this_.scriptSet(one.Name, func(args ...interface{}) {
			util.Logger.Debug("func "+one.Name+" run start", zap.Any("func", one))
		})
		if err != nil {
			util.Logger.Error("invoke data init set func value error", zap.Any("name", one.Name), zap.Any("func", one), zap.Any("error", err))
			return
		}
	}

	return
}

func (this_ *InvokeData) GetArgs() (args []*InvokeVar) {
	args = this_.args
	return
}

func (this_ *InvokeData) GetVars() (vars []*InvokeVar) {
	vars = this_.vars
	return
}

func (this_ *InvokeData) scriptSet(name string, value interface{}) (err error) {
	this_.scriptLock.Lock()
	defer this_.scriptLock.Unlock()

	err = this_.script.Set(name, value)
	if err != nil {
		util.Logger.Error("script set value error", zap.Any("name", name), zap.Any("value", value), zap.Any("error", err))
		return
	}
	return
}

func (this_ *InvokeData) AddArg(name string, value interface{}, valueType *modelers.ValueType) (err error) {
	this_.argCacheLock.Lock()
	defer this_.argCacheLock.Unlock()

	err = this_.addArg(&InvokeVar{
		Name:      name,
		Value:     value,
		ValueType: valueType,
	})

	err = this_.scriptSet(name, value)
	if err != nil {
		return
	}
	return
}

func (this_ *InvokeData) AddVar(name string, value interface{}, varType string) (err error) {
	this_.varCacheLock.Lock()
	defer this_.varCacheLock.Unlock()

	var valueType *modelers.ValueType = nil

	if varType != "" {
		valueType, err = this_.app.GetValueType(varType)
		if err != nil {
			util.Logger.Error("invoke set var get value type error", zap.Any("varType", varType), zap.Any("error", err))
			return
		}
		if valueType.Struct != nil {
			strValue := util.GetStringValue(value)
			value = nil
			if strValue != "" {
				value = map[string]interface{}{}
				err = util.JSONDecodeUseNumber([]byte(strValue), &value)
				if err != nil {
					util.Logger.Error("invoke set var to json error", zap.Any("strValue", strValue), zap.Any("error", err))
					return
				}
			}
		}
	}

	err = this_.addVar(&InvokeVar{
		Name:      name,
		Value:     value,
		ValueType: valueType,
	})

	err = this_.scriptSet(name, value)
	if err != nil {
		return
	}
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
	//if this_.varCache[var_.Name] != nil {
	//	err = errors.New("var [" + var_.Name + "] already exist")
	//	return
	//}
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

func (this_ *InvokeData) InvokeScript(script string) (res interface{}, err error) {
	funcInvoke := invokeStart("invoke script", this_)
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(funcInvoke.name + " error:" + fmt.Sprint(e))
			util.Logger.Error(funcInvoke.name+" error", zap.Any("script", script), zap.Any("error", err))
		}
		funcInvoke.end(err)
		util.Logger.Debug(funcInvoke.name+" end", zap.Any("use", funcInvoke.use()))
	}()
	util.Logger.Debug(funcInvoke.name+" start", zap.Any("script", script))
	if script == "" {
		return
	}

	res, err = this_.script.GetScriptValue(script)
	if err != nil {
		util.Logger.Error(funcInvoke.name+" error", zap.Any("script", script), zap.Any("error", err))
		return
	}
	return
}

func (this_ *InvokeData) InvokeScriptStringValue(script string) (res string, err error) {
	funcInvoke := invokeStart("invoke script string value", this_)
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(funcInvoke.name + " error:" + fmt.Sprint(e))
			util.Logger.Error(funcInvoke.name+" error", zap.Any("script", script), zap.Any("error", err))
		}
		funcInvoke.end(err)
		util.Logger.Debug(funcInvoke.name+" end", zap.Any("use", funcInvoke.use()))
	}()

	util.Logger.Debug(funcInvoke.name+" start", zap.Any("script", script))
	if script == "" {
		return
	}

	v, err := this_.script.GetScriptValue(script)
	if err != nil {
		return
	}
	if v != nil {
		res = util.GetStringValue(v)
	}
	return
}

func (this_ *InvokeData) InvokeByStringRule(stringRule string) (res string, err error) {
	funcInvoke := invokeStart("invoke by string rule", this_)
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(funcInvoke.name + " error:" + fmt.Sprint(e))
			util.Logger.Error(funcInvoke.name+" error", zap.Any("stringRule", stringRule), zap.Any("error", err))
		}
		funcInvoke.end(err)
		util.Logger.Debug(funcInvoke.name+" end", zap.Any("use", funcInvoke.use()))
	}()

	util.Logger.Debug(funcInvoke.name+" start", zap.Any("stringRule", stringRule))
	if stringRule == "" {
		return
	}

	re := regexp.MustCompile(`\$({[^}]+})`)
	matches := re.FindAllStringIndex(stringRule, -1)
	var lastMatch []int
	for _, match := range matches {
		if lastMatch == nil {
			res += stringRule[:match[0]]
		} else {
			res += stringRule[lastMatch[1]:match[0]]
		}
		matchStr := stringRule[match[0]+2 : match[1]-1]

		var matchValue string
		matchValue, err = this_.InvokeScriptStringValue(matchStr)
		if err != nil {
			util.Logger.Error(funcInvoke.name+" get match value error", zap.Any("matchStr", matchStr), zap.Any("error", err))
			return
		}
		res += matchValue

		lastMatch = match
	}
	if lastMatch == nil {
		res += stringRule[lastMatch[1]:]
	}

	if err != nil {
		util.Logger.Error(funcInvoke.name+" error", zap.Any("stringRule", stringRule), zap.Any("error", err))
		return
	}
	return
}
