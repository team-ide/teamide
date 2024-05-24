package maker

import (
	"fmt"
	"github.com/dop251/goja"
	"github.com/dop251/goja/ast"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"reflect"
	"strings"
	"teamide/pkg/maker/modelers"
)

func (this_ *Compiler) GetOrCreateSpace(space string) (res *CompilerSpace) {
	res = this_.spaceCache[space]
	if res == nil {
		res = &CompilerSpace{
			space:     space,
			Compiler:  this_,
			packCache: make(map[string]*CompilerPack),
		}
		this_.spaceList = append(this_.spaceList, res)
		this_.spaceCache[space] = res
	}
	return
}

type CompilerSpace struct {
	*Compiler
	space     string
	packList  []*CompilerPack
	packCache map[string]*CompilerPack
}

func (this_ *CompilerSpace) GetOrCreatePack(pack string) (res *CompilerPack) {
	res = this_.packCache[pack]
	if res == nil {
		res = &CompilerPack{
			pack:          pack,
			CompilerSpace: this_,
			classCache:    make(map[string]*CompilerClass),
		}
		this_.packList = append(this_.packList, res)
		this_.packCache[pack] = res
	}
	return
}

func (this_ *CompilerSpace) GetClass(path string, fileIsClass bool) (endName string, res *CompilerClass) {
	names := strings.Split(path, "/")
	size := len(names)
	classIndex := size - 1
	endNameIndex := -1
	if !fileIsClass {
		classIndex = size - 2
		endNameIndex = size - 1
	}
	var packs []string
	var class string
	if classIndex >= 0 {
		class = names[classIndex]
		if classIndex > 0 {
			packs = names[:classIndex-1]
		}
	}
	packName := strings.Join(packs, ".")
	pack := this_.GetOrCreatePack(packName)
	res = pack.GetOrCreateClass(class)
	if endNameIndex >= 0 {
		endName = names[endNameIndex]
	}
	return
}

type CompilerPack struct {
	*CompilerSpace
	pack       string
	classList  []*CompilerClass
	classCache map[string]*CompilerClass
}

func (this_ *CompilerPack) GetOrCreateClass(class string) (res *CompilerClass) {
	res = this_.classCache[class]
	if res == nil {
		res = &CompilerClass{
			class:        class,
			CompilerPack: this_,
			importCache:  make(map[*CompilerClass]*CompilerImport),
			fieldCache:   make(map[string]*CompilerField),
			methodCache:  make(map[string]*CompilerMethod),
		}
		this_.classList = append(this_.classList, res)
		this_.classCache[class] = res
	}
	return
}

type CompilerClass struct {
	*CompilerPack
	class       string
	importList  []*CompilerImport
	importCache map[*CompilerClass]*CompilerImport
	fieldList   []*CompilerField
	fieldCache  map[string]*CompilerField
	methodList  []*CompilerMethod
	methodCache map[string]*CompilerMethod
}

func (this_ *CompilerClass) GetOrCreateImport(class *CompilerClass) (res *CompilerImport) {
	res = this_.importCache[class]
	if res == nil {
		res = &CompilerImport{
			class: class,
		}
		this_.importList = append(this_.importList, res)
		this_.importCache[class] = res
	}
	return
}

type CompilerImport struct {
	class *CompilerClass
}

func (this_ *CompilerClass) GetOrCreateField(name string) (res *CompilerField) {
	res = this_.fieldCache[name]
	if res == nil {
		res = &CompilerField{
			CompilerClass: this_,
			CompilerArg: &CompilerArg{
				name: name,
			},
		}
		this_.fieldList = append(this_.fieldList, res)
		this_.fieldCache[name] = res
	}
	return
}

type CompilerField struct {
	*CompilerClass
	*CompilerArg
	value string
}

func (this_ *CompilerClass) GetMethod(name string) (res *CompilerMethod) {
	res = this_.methodCache[name]
	return
}

func (this_ *CompilerClass) CreateMethod(name string, args []*modelers.ArgModel) (res *CompilerMethod, err error) {

	res = &CompilerMethod{
		CompilerClass: this_,
		method:        name,
		paramCache:    make(map[string]*CompilerMethodParam),
		varCache:      make(map[string]*CompilerMethodVar),
	}

	for _, arg := range args {
		var t *ValueType
		t, err = this_.GetValueType(arg.Type)
		if err != nil {
			return
		}
		res.addParam(arg.Name, t)
		if err != nil {
			return
		}
	}
	this_.methodList = append(this_.methodList, res)
	this_.methodCache[name] = res
	return
}

type CompilerMethod struct {
	*CompilerClass
	method            string
	paramList         []*CompilerMethodParam
	paramCache        map[string]*CompilerMethodParam
	varList           []*CompilerMethodVar
	varCache          map[string]*CompilerMethodVar
	callComponentList []*CompilerCall
	callUtilList      []*CompilerCall
	callFuncList      []*CompilerCall
	callDaoList       []*CompilerCall
	callServiceList   []*CompilerCall
	result            *CompilerMethodResult
	program           *ast.Program
	script            *Script
}

func (this_ *CompilerMethod) GetKey() (key string) {
	key = "space:" + this_.space + ",class:" + this_.class + ",method:" + this_.method
	return
}

func (this_ *CompilerMethod) getParam(name string) (res *CompilerMethodParam) {
	res = this_.paramCache[name]
	return
}
func (this_ *CompilerMethod) addParam(name string, valueType *ValueType) (res *CompilerMethodParam) {
	res = &CompilerMethodParam{
		CompilerMethod: this_,
		CompilerArg: &CompilerArg{
			name:       name,
			valueTypes: []*ValueType{valueType},
		},
	}
	this_.paramList = append(this_.paramList, res)
	this_.paramCache[name] = res

	return
}

func (this_ *CompilerMethod) getVar(name string) (res *CompilerMethodVar) {
	res = this_.varCache[name]
	return
}
func (this_ *CompilerMethod) addVar(name string) (res *CompilerMethodVar) {
	res = &CompilerMethodVar{
		CompilerMethod: this_,
		CompilerArg: &CompilerArg{
			name: name,
		},
	}
	this_.varList = append(this_.varList, res)
	this_.varCache[name] = res

	return
}

func (this_ *CompilerMethod) getByNameScript(nameScript string) (res interface{}) {
	names := strings.Split(nameScript, `.`)
	if len(names) == 0 {
		fmt.Println("name script [", nameScript, "] is empty")
		return
	}
	v := this_.script.dataContext[names[0]]
	if v == nil {
		return
	}
	switch tV := v.(type) {
	case *CompilerMethodVar:
		res = tV
		break
	case *CompilerMethodParam:
		res = tV
		break
	default:
		panic("getByNameScript [" + reflect.TypeOf(v).String() + "] not support")
	}

	return
}
func (this_ *CompilerMethod) findType(name string) (find bool) {
	_, find = this_.varCache[name]
	if !find {
		_, find = this_.paramCache[name]
	}
	return
}

func (this_ *CompilerMethod) BindCode(script string) (err error) {
	runScript := `(function (){
` + script + `
})()`
	this_.program, err = goja.Parse("", runScript)
	if err != nil {
		util.Logger.Error("compile script parse error", zap.Any("error", err))
		return
	}

	return
}

type CompilerMethodParam struct {
	*CompilerMethod
	*CompilerArg
}

type CompilerMethodVar struct {
	*CompilerMethod
	*CompilerArg
	value string
}

type CompilerCall struct {
	*CompilerMethod
	name       string
	paramList  []*CompilerCallParam
	returnList []*ValueType
}

type CompilerCallParam struct {
	*CompilerCall
	*CompilerArg
	value string
}

type CompilerArg struct {
	name       string
	valueTypes []*ValueType
}

func (this_ *CompilerArg) addValueType(valueTypes ...*ValueType) {
	addValueTypes(&this_.valueTypes, valueTypes...)
	return
}

type CompilerMethodResult struct {
	valueTypes []*ValueType
}

func (this_ *CompilerMethodResult) addValueType(valueTypes ...*ValueType) {
	addValueTypes(&this_.valueTypes, valueTypes...)
	return
}

func addValueTypes(toList *[]*ValueType, valueTypes ...*ValueType) {
	for _, valueType := range valueTypes {
		var find bool
		for _, v := range *toList {
			if v == valueType {
				find = true
			}
		}
		if !find {
			*toList = append(*toList, valueType)
		}
	}

	return
}
