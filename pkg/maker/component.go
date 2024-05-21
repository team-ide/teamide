package maker

import (
	"github.com/team-ide/go-tool/util"
	"teamide/pkg/maker/modelers"
)

type Component struct {
	Fields  []*ComponentField  `json:"fields"`
	Methods []*ComponentMethod `json:"methods"`
}

type ComponentField struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type ComponentMethod struct {
	Name           string                                                       `json:"name"`
	Args           []*ComponentField                                            `json:"args"`
	ThrowErrors    []string                                                     `json:"throwErrors"`
	GetReturnTypes func(args []interface{}) (returnTypes []*modelers.ValueType) `json:"-"`
}

func (this_ *Component) ToContext() (ctx map[string]interface{}) {
	ctx = make(map[string]interface{})
	for _, method := range this_.Methods {
		name := method.Name
		ctx[name] = method
		ctx[util.FirstToLower(name)] = method
	}

	return ctx
}
