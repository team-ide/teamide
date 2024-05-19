package maker

import (
	"github.com/dop251/goja"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"reflect"
)

type Script struct {
	dataContext map[string]interface{}
	vm          *goja.Runtime
}

type ShouldMappingFunc interface {
	ShouldMappingFunc() bool
}
type MappingFunc struct {
}

func (this_ *Script) Set(name string, value interface{}) (err error) {

	//util.Logger.Debug("script set var", zap.Any("name", name))
	var setValue = value
	if shouldMappingFunc, ok := value.(ShouldMappingFunc); ok && shouldMappingFunc.ShouldMappingFunc() {
		mappingFunc := map[string]interface{}{}
		mappingFunc["_bind_obj"] = value
		vOf := reflect.ValueOf(value)
		tOf := reflect.TypeOf(value)
		num := vOf.NumMethod()
		for i := 0; i < num; i++ {
			tM := tOf.Method(i)
			if tM.Name == "ShouldMappingFunc" {
				continue
			}
			vM := vOf.Method(i)
			mappingFunc[tM.Name] = vM.Interface()
			mappingFunc[util.FirstToLower(tM.Name)] = vM.Interface()
		}

		setValue = mappingFunc
	}

	this_.dataContext[name] = setValue
	err = this_.vm.Set(name, setValue)
	if err != nil {
		util.Logger.Error("script set var error", zap.Any("name", name), zap.Any("error", err))
		return
	}

	return
}
func (this_ *Script) GetScriptValue(script string) (interface{}, error) {
	if script == "" {
		return nil, nil
	}

	var scriptValue goja.Value
	scriptValue, err := this_.vm.RunString(script)
	if err != nil {
		util.Logger.Error("表达式执行异常", zap.Any("script", script), zap.Error(err))
		return nil, err
	}
	return scriptValue.Export(), nil
}

func (this_ *Script) CompileScript(script string) (*goja.Program, error) {
	runScript := `(function (){` + script + `})()
`
	p, err := goja.Compile("", runScript, false)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (this_ *Script) RunScript(script string) (interface{}, error) {
	if script == "" {
		return nil, nil
	}

	runScript := `(function (){
` + script + `
})()
`
	scriptValue, err := this_.vm.RunScript("", runScript)
	if err != nil {
		return nil, err
	}
	return scriptValue.Export(), nil
}

func (this_ *Script) GetStringScriptValue(script string) (value string, err error) {

	var scriptValue interface{}
	scriptValue, err = this_.GetScriptValue(script)
	if scriptValue != nil {
		value = util.GetStringValue(scriptValue)
		return
	}
	return
}
func NewScript() (script *Script, err error) {

	return NewScriptByParent(nil)
}
func NewScriptByParent(parent *Script) (script *Script, err error) {
	script = &Script{}
	script.vm = goja.New()
	script.dataContext = make(map[string]interface{})
	if parent != nil {
		for key, value := range parent.dataContext {
			err = script.Set(key, value)
			if err != nil {
				return
			}
		}
	}

	return
}
