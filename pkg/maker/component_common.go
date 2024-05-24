package maker

import (
	"teamide/pkg/maker/modelers"
)

func NewCommonCompiler(config *modelers.AppModel) *Component {
	component := &Component{
		Methods: []*ComponentMethod{
			{
				Name: "GenId", GetReturnTypes: func(args []interface{}) (returnTypes []*ValueType) {

					returnTypes = append(returnTypes, ValueTypeInt64)
					return
				},
			},
			{
				Name: "GenStr", GetReturnTypes: func(args []interface{}) (returnTypes []*ValueType) {

					returnTypes = append(returnTypes, ValueTypeString)
					return
				},
			},
		},
	}
	return component
}

func NewComponentCommon(config *modelers.AppModel) (res *ComponentRedis, err error) {

	res = &ComponentRedis{}
	return
}

type ComponentCommon struct {
	*Compiler
}

func (this_ *ComponentCommon) ShouldMappingFunc() bool {
	return true
}

func (this_ *ComponentCommon) GenId() (res int64, err error) {

	return
}

func (this_ *ComponentCommon) GenStr(min, max int) (res string, err error) {

	return
}
