package maker

import (
	"github.com/team-ide/go-tool/db"
	"teamide/pkg/maker/modelers"
)

func NewDbCompiler(config *modelers.ConfigDbModel) *Component {
	component := &Component{
		Methods: []*ComponentMethod{
			{
				Name: "SelectOne", GetReturnTypes: func(args []interface{}) (returnType *ValueType) {
					if len(args) == 4 {
						returnType = args[3].(*ValueType)
					} else {
						returnType = ValueTypeMap
					}
					return
				},
			},
			{
				Name: "Insert", GetReturnTypes: func(args []interface{}) (returnType *ValueType) {
					returnType = ValueTypeInt64
					return
				},
			},
			{
				Name: "Update", GetReturnTypes: func(args []interface{}) (returnType *ValueType) {
					returnType = ValueTypeInt64
					return
				},
			},
			{
				Name: "Delete", GetReturnTypes: func(args []interface{}) (returnType *ValueType) {
					returnType = ValueTypeInt64
					return
				},
			},
		},
	}
	return component
}

func NewComponentDb(config *modelers.ConfigDbModel) (res *ComponentDb, err error) {
	ser, err := db.New(&db.Config{
		Type:     config.Type,
		Host:     config.Host,
		Port:     config.Port,
		Username: config.Username,
		Password: config.Password,
		Database: config.Database,
	})
	if err != nil {
		return
	}
	res = &ComponentDb{
		ser: ser,
	}
	return
}

type ComponentDb struct {
	ser db.IService
}

func (this_ *ComponentDb) ShouldMappingFunc() bool {
	return true
}

func (this_ *ComponentDb) SelectOne(columns string, table string, where string, valueType *ValueType) (res any, err error) {
	find := map[string]interface{}{
		"userId": 1,
		"name":   "张三",
	}
	res = find
	return
}
