package maker

import (
	"errors"
	"fmt"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"regexp"
	"sync"
	"teamide/pkg/maker/modelers"
)

type InvokeData struct {
	invoker      *Invoker
	args         []*InvokeVar
	vars         []*InvokeVar
	argCache     map[string]*InvokeVar
	argCacheLock sync.Locker
	varCache     map[string]*InvokeVar
	varCacheLock sync.Locker
	script       *Script
	scriptLock   sync.Locker
}

type InvokeVar struct {
	Name      string      `json:"name,omitempty"`
	Value     interface{} `json:"value,omitempty"`
	ValueType *ValueType  `json:"valueType,omitempty"`
}

func (this_ *Invoker) NewInvokeData() (data *InvokeData, err error) {
	script, err := this_.NewScriptByParent(this_.script)
	if err != nil {
		util.Logger.Error("NewInvokeData NewScript error", zap.Any("error", err))
		return
	}
	data = &InvokeData{
		invoker:      this_,
		argCache:     make(map[string]*InvokeVar),
		argCacheLock: &sync.Mutex{},
		varCache:     make(map[string]*InvokeVar),
		varCacheLock: &sync.Mutex{},
		script:       script,
		scriptLock:   &sync.Mutex{},
	}
	return
}

func (this_ *Invoker) NewInvokeDataByArgs(argModels []*modelers.ArgModel, args []interface{}) (data *InvokeData, err error) {
	data, err = this_.NewInvokeData()
	if err != nil {
		return
	}
	mSize := len(argModels)
	vSize := len(args)
	for i := 0; i < mSize; i++ {
		if i > vSize-1 {
			break
		}
		err = data.AddVar(argModels[i].Name, args[i], argModels[i].Type)
		if err != nil {
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

func (this_ *InvokeData) AddArg(name string, value interface{}, valueType *ValueType) (err error) {
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

	var valueType *ValueType = nil

	if varType != "" {
		valueType, err = this_.invoker.GetValueType(varType)
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
	util.Logger.Debug(funcInvoke.name + " start")
	if script == "" {
		return
	}

	res, err = this_.script.RunScript(script)
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

	util.Logger.Debug(funcInvoke.name + " start")
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
