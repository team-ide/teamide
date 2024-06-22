package modelers

import (
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"strings"
)

type Type struct {
	Name     string  `json:"name"`
	Comment  string  `json:"comment"`
	IsFile   bool    `json:"isFile"`
	Children []*Type `json:"children"`

	newModel func() any
	toModel  func(name, text string) (model interface{}, err error)
	toText   func(model interface{}) (text string, err error)
}

func (this_ *Type) ToText(model interface{}) (text string, err error) {
	return this_.toText(model)
}

func (this_ *Type) ToModel(name, text string) (model interface{}, err error) {
	return this_.toModel(name, text)
}
func (this_ *Type) NewModel() any {
	return this_.newModel()
}

var (
	Types []*Type

	TypeAppName = "app"
	TypeApp     = &Type{
		Name:     TypeAppName,
		Comment:  "应用设置",
		IsFile:   true,
		newModel: func() any { return &AppModel{} },
		toModel: func(name, text string) (model interface{}, err error) {
			var app = &AppModel{}
			app.Text = text
			model = app

			err = toModel(text, TypeAppName, model)
			if err != nil {
				util.Logger.Error("text to app model error", zap.Any("text", text), zap.Error(err))
				return
			}

			data := map[string]interface{}{}
			err = yaml.Unmarshal([]byte(text), &data)
			if err != nil {
				util.Logger.Error("text to yaml data error", zap.Any("text", text), zap.Error(err))
				return
			}
			var bs []byte
			var m interface{}
			app.Other = map[string]any{}
			for k, v := range data {
				var t *Type
				var tN string
				for tK, tV := range configTypes {
					if k == tK {
						t = tV
						break
					} else if strings.HasPrefix(k, tK+"_") {
						t = tV
						tN = strings.TrimPrefix(k, tK+"_")
					}
				}
				if t != nil {
					bs, err = yaml.Marshal(v)
					if err != nil {
						util.Logger.Error("value to yaml error", zap.Any("value", v), zap.Error(err))
						return
					}
					modelName := tN
					if modelName == "" {
						modelName = "default"
					}
					m, err = t.toModel(modelName, string(bs))
					if err != nil {
						util.Logger.Error("value yaml to model error", zap.Any("type", t.Name), zap.Any("value", string(bs)), zap.Error(err))
						return
					}
					switch t {
					case TypeConfigDb:
						if app.Db == nil {
							app.Db = make(map[string]*ConfigDbModel)
						}
						app.Db[tN] = m.(*ConfigDbModel)
						break
					case TypeConfigRedis:
						if app.Redis == nil {
							app.Redis = make(map[string]*ConfigRedisModel)
						}
						app.Redis[tN] = m.(*ConfigRedisModel)
						break
					case TypeConfigZk:
						if app.Zk == nil {
							app.Zk = make(map[string]*ConfigZkModel)
						}
						app.Zk[tN] = m.(*ConfigZkModel)
						break
					case TypeConfigEs:
						if app.Es == nil {
							app.Es = make(map[string]*ConfigEsModel)
						}
						app.Es[tN] = m.(*ConfigEsModel)
						break
					case TypeConfigKafka:
						if app.Kafka == nil {
							app.Kafka = make(map[string]*ConfigKafkaModel)
						}
						app.Kafka[tN] = m.(*ConfigKafkaModel)
						break
					case TypeConfigMongodb:
						if app.Mongodb == nil {
							app.Mongodb = make(map[string]*ConfigMongodbModel)
						}
						app.Mongodb[tN] = m.(*ConfigMongodbModel)
						break
					}
				} else {
					app.Other[k] = v
				}
			}

			return
		},
		toText: func(model interface{}) (text string, err error) {
			text, err = toText(model, TypeAppName, &docOptions{
				outComment: true,
				omitEmpty:  false,
			})
			if err != nil {
				util.Logger.Error("app model to text error", zap.Any("model", model), zap.Error(err))
				return
			}
			return
		},
	}

	TypeConfigDbName = "config/database"
	TypeConfigDb     = &Type{
		Name:     TypeConfigDbName,
		Comment:  "Database",
		newModel: func() any { return &ConfigDbModel{} },
		toModel: func(name, text string) (model interface{}, err error) {
			model = &ConfigDbModel{}
			err = toModel(text, TypeConfigDbName, model)
			if err != nil {
				util.Logger.Error("text to config database model error", zap.Any("text", text), zap.Error(err))
				return
			}
			model.(*ConfigDbModel).Name = name
			return
		},
		toText: func(model interface{}) (text string, err error) {
			text, err = toText(model, TypeConfigDbName, &docOptions{
				outComment: true,
				omitEmpty:  false,
			})
			if err != nil {
				util.Logger.Error("config database model to text error", zap.Any("model", model), zap.Error(err))
				return
			}
			return
		},
	}

	TypeConfigRedisName = "config/redis"
	TypeConfigRedis     = &Type{
		Name:     TypeConfigRedisName,
		Comment:  "Redis",
		newModel: func() any { return &ConfigRedisModel{} },
		toModel: func(name, text string) (model interface{}, err error) {
			model = &ConfigRedisModel{}
			err = toModel(text, TypeConfigRedisName, model)
			if err != nil {
				util.Logger.Error("text to config redis model error", zap.Any("text", text), zap.Error(err))
				return
			}
			model.(*ConfigRedisModel).Name = name
			return
		},
		toText: func(model interface{}) (text string, err error) {
			text, err = toText(model, TypeConfigRedisName, &docOptions{
				outComment: true,
				omitEmpty:  false,
			})
			if err != nil {
				util.Logger.Error("config redis model to text error", zap.Any("model", model), zap.Error(err))
				return
			}
			return
		},
	}

	TypeConfigZkName = "config/zk"
	TypeConfigZk     = &Type{
		Name:     TypeConfigZkName,
		Comment:  "Zookeeper",
		newModel: func() any { return &ConfigZkModel{} },
		toModel: func(name, text string) (model interface{}, err error) {
			model = &ConfigZkModel{}
			err = toModel(text, TypeConfigZkName, model)
			if err != nil {
				util.Logger.Error("text to config zookeeper model error", zap.Any("text", text), zap.Error(err))
				return
			}
			model.(*ConfigZkModel).Name = name
			return
		},
		toText: func(model interface{}) (text string, err error) {
			text, err = toText(model, TypeConfigZkName, &docOptions{
				outComment: true,
				omitEmpty:  false,
			})
			if err != nil {
				util.Logger.Error("config zookeeper model to text error", zap.Any("model", model), zap.Error(err))
				return
			}
			return
		},
	}

	TypeConfigKafkaName = "config/kafka"
	TypeConfigKafka     = &Type{
		Name:     TypeConfigKafkaName,
		Comment:  "Kafka",
		newModel: func() any { return &ConfigKafkaModel{} },
		toModel: func(name, text string) (model interface{}, err error) {
			model = &ConfigKafkaModel{}
			err = toModel(text, TypeConfigKafkaName, model)
			if err != nil {
				util.Logger.Error("text to config kafka model error", zap.Any("text", text), zap.Error(err))
				return
			}
			model.(*ConfigKafkaModel).Name = name
			return
		},
		toText: func(model interface{}) (text string, err error) {
			text, err = toText(model, TypeConfigKafkaName, &docOptions{
				outComment: true,
				omitEmpty:  false,
			})
			if err != nil {
				util.Logger.Error("config kafka model to text error", zap.Any("model", model), zap.Error(err))
				return
			}
			return
		},
	}

	TypeConfigMongodbName = "config/mongodb"
	TypeConfigMongodb     = &Type{
		Name:     TypeConfigMongodbName,
		Comment:  "Mongodb",
		newModel: func() any { return &ConfigMongodbModel{} },
		toModel: func(name, text string) (model interface{}, err error) {
			model = &ConfigMongodbModel{}
			err = toModel(text, TypeConfigMongodbName, model)
			if err != nil {
				util.Logger.Error("text to config mongodb model error", zap.Any("text", text), zap.Error(err))
				return
			}
			model.(*ConfigMongodbModel).Name = name
			return
		},
		toText: func(model interface{}) (text string, err error) {
			text, err = toText(model, TypeConfigMongodbName, &docOptions{
				outComment: true,
				omitEmpty:  false,
			})
			if err != nil {
				util.Logger.Error("config mongodb model to text error", zap.Any("model", model), zap.Error(err))
				return
			}
			return
		},
	}

	TypeConfigEsName = "config/es"
	TypeConfigEs     = &Type{
		Name:     TypeConfigEsName,
		Comment:  "Elastic Search",
		newModel: func() any { return &ConfigEsModel{} },
		toModel: func(name, text string) (model interface{}, err error) {
			model = &ConfigEsModel{}
			err = toModel(text, TypeConfigEsName, model)
			if err != nil {
				util.Logger.Error("text to config elasticsearch model error", zap.Any("text", text), zap.Error(err))
				return
			}
			model.(*ConfigEsModel).Name = name
			return
		},
		toText: func(model interface{}) (text string, err error) {
			text, err = toText(model, TypeConfigEsName, &docOptions{
				outComment: true,
				omitEmpty:  false,
			})
			if err != nil {
				util.Logger.Error("config elasticsearch model to text error", zap.Any("model", model), zap.Error(err))
				return
			}
			return
		},
	}

	configTypes = map[string]*Type{
		"db":      TypeConfigDb,
		"redis":   TypeConfigRedis,
		"es":      TypeConfigEs,
		"zk":      TypeConfigZk,
		"mongodb": TypeConfigMongodb,
		"kafka":   TypeConfigKafka,
	}

	TypeConstantName = "constant"

	TypeConstant = &Type{
		Name:     TypeConstantName,
		Comment:  "常量",
		newModel: func() any { return &ConstantModel{} },
		toModel: func(name, text string) (model interface{}, err error) {
			model = &ConstantModel{}
			err = toModel(text, TypeConstantName, model)
			if err != nil {
				util.Logger.Error("text to constant model error", zap.Any("text", text), zap.Error(err))
				return
			}
			model.(*ConstantModel).Name = name
			return
		},
		toText: func(model interface{}) (text string, err error) {
			text, err = toText(model, TypeConstantName, &docOptions{
				outComment: true,
				omitEmpty:  false,
			})
			if err != nil {
				util.Logger.Error("constant model to text error", zap.Any("model", model), zap.Error(err))
				return
			}
			return
		},
	}

	TypeErrorName = "error"
	TypeError     = &Type{
		Name:     TypeErrorName,
		Comment:  "错误码",
		newModel: func() any { return &ErrorModel{} },
		toModel: func(name, text string) (model interface{}, err error) {
			model = &ErrorModel{}
			err = toModel(text, TypeErrorName, model)
			if err != nil {
				util.Logger.Error("text to error model error", zap.Any("text", text), zap.Error(err))
				return
			}
			model.(*ErrorModel).Name = name
			return
		},
		toText: func(model interface{}) (text string, err error) {
			text, err = toText(model, TypeErrorName, &docOptions{
				outComment: true,
				omitEmpty:  false,
			})
			if err != nil {
				util.Logger.Error("error model to text error", zap.Any("model", model), zap.Error(err))
				return
			}
			return
		},
	}

	TypeStructName = "struct"
	TypeStruct     = &Type{
		Name:     TypeStructName,
		Comment:  "结构体",
		newModel: func() any { return &StructModel{} },
		toModel: func(name, text string) (model interface{}, err error) {
			model = &StructModel{}
			err = toModel(text, TypeStructName, model)
			if err != nil {
				util.Logger.Error("text to struct model error", zap.Any("text", text), zap.Error(err))
				return
			}
			model.(*StructModel).Name = name
			return
		},
		toText: func(model interface{}) (text string, err error) {
			text, err = toText(model, TypeStructName, &docOptions{
				outComment: true,
				omitEmpty:  false,
			})
			if err != nil {
				util.Logger.Error("struct model to text error", zap.Any("model", model), zap.Error(err))
				return
			}
			return
		},
	}

	TypeTableName = "table"
	TypeTable     = &Type{
		Name:     TypeTableName,
		Comment:  "表",
		newModel: func() any { return &TableModel{} },
		toModel: func(name, text string) (model interface{}, err error) {
			model = &TableModel{}
			err = toModel(text, TypeTableName, model)
			if err != nil {
				util.Logger.Error("text to table model error", zap.Any("text", text), zap.Error(err))
				return
			}
			model.(*TableModel).Name = name
			return
		},
		toText: func(model interface{}) (text string, err error) {
			text, err = toText(model, TypeTableName, &docOptions{
				outComment: true,
				omitEmpty:  false,
			})
			if err != nil {
				util.Logger.Error("table model to text error", zap.Any("model", model), zap.Error(err))
				return
			}
			return
		},
	}

	TypeStorageName = "storage"
	TypeStorage     = &Type{
		Name:     TypeStorageName,
		Comment:  "数据层",
		newModel: func() any { return &StorageModel{} },
		toModel: func(name, text string) (model interface{}, err error) {
			model = &StorageModel{}
			err = toModel(text, TypeStorageName, model)
			if err != nil {
				util.Logger.Error("text to storage model error", zap.Any("text", text), zap.Error(err))
				return
			}
			model.(*StorageModel).Name = name
			return
		},
		toText: func(model interface{}) (text string, err error) {
			text, err = toText(model, TypeStorageName, &docOptions{
				outComment: true,
				omitEmpty:  false,
			})
			if err != nil {
				util.Logger.Error("storage model to text error", zap.Any("model", model), zap.Error(err))
				return
			}
			return
		},
	}

	TypeServiceName = "service"
	TypeService     = &Type{
		Name:     TypeServiceName,
		Comment:  "服务",
		newModel: func() any { return &ServiceModel{} },
		toModel: func(name, text string) (model interface{}, err error) {
			model = &ServiceModel{}
			err = toModel(text, TypeServiceName, model)
			if err != nil {
				util.Logger.Error("text to service model error", zap.Any("text", text), zap.Error(err))
				return
			}
			model.(*ServiceModel).Name = name
			return
		},
		toText: func(model interface{}) (text string, err error) {
			text, err = toText(model, TypeServiceName, &docOptions{
				outComment: true,
				omitEmpty:  false,
			})
			if err != nil {
				util.Logger.Error("service model to text error", zap.Any("model", model), zap.Error(err))
				return
			}
			return
		},
	}

	TypeFuncName = "func"
	TypeFunc     = &Type{
		Name:     TypeFuncName,
		Comment:  "函数",
		newModel: func() any { return &FuncModel{} },
		toModel: func(name, text string) (model interface{}, err error) {
			model = &FuncModel{}
			err = toModel(text, TypeFuncName, model)
			if err != nil {
				util.Logger.Error("text to func model error", zap.Any("text", text), zap.Error(err))
				return
			}
			model.(*FuncModel).Name = name
			return
		},
		toText: func(model interface{}) (text string, err error) {
			text, err = toText(model, TypeFuncName, &docOptions{
				outComment: true,
				omitEmpty:  false,
			})
			if err != nil {
				util.Logger.Error("func model to text error", zap.Any("model", model), zap.Error(err))
				return
			}
			return
		},
	}

	TypeLanguageGolangName = "language/golang"
	TypeLanguageGolang     = &Type{
		Name:     TypeLanguageGolangName,
		Comment:  "Golang",
		IsFile:   true,
		newModel: func() any { return &LanguageGolangModel{} },
		toModel: func(name, text string) (model interface{}, err error) {
			model = &LanguageGolangModel{}
			err = toModel(text, TypeLanguageGolangName, model)
			if err != nil {
				util.Logger.Error("text to language golang model error", zap.Any("text", text), zap.Error(err))
				return
			}
			return
		},
		toText: func(model interface{}) (text string, err error) {
			text, err = toText(model, TypeLanguageGolangName, &docOptions{
				outComment: true,
				omitEmpty:  false,
			})
			if err != nil {
				util.Logger.Error("language golang model to text error", zap.Any("model", model), zap.Error(err))
				return
			}
			return
		},
	}

	TypeFlowchartName = "flowchart"
	TypeFlowchart     = &Type{
		Name:     TypeFlowchartName,
		Comment:  "流程图",
		newModel: func() any { return &FlowchartModel{} },
		toModel: func(name, text string) (model interface{}, err error) {
			model = &FlowchartModel{}
			err = toModel(text, TypeFlowchartName, model)
			if err != nil {
				util.Logger.Error("text to flowchart model error", zap.Any("text", text), zap.Error(err))
				return
			}
			model.(*FlowchartModel).Name = name
			return
		},
		toText: func(model interface{}) (text string, err error) {
			text, err = toText(model, TypeFlowchartName, &docOptions{
				outComment: true,
				omitEmpty:  false,
			})
			if err != nil {
				util.Logger.Error("flowchart model to text error", zap.Any("model", model), zap.Error(err))
				return
			}
			return
		},
	}
)

func init() {
	AppendType(TypeApp)

	AppendType(TypeConstant)

	AppendType(TypeError)

	AppendType(TypeStruct)

	AppendType(TypeTable)

	AppendType(TypeStorage)

	AppendType(TypeService)

	AppendType(TypeFunc)

	AppendType(&Type{
		Name:    "language",
		Comment: "导出语言",
		Children: []*Type{
			TypeLanguageGolang,
		},
	})

	AppendType(TypeFlowchart)

}

var (
	modelTypeCache = make(map[string]*Type)
)

func GetModelType(key string) (modelType *Type) {
	modelType = modelTypeCache[key]
	return
}

func AppendType(one *Type) {

	modelTypeCache[one.Name] = one
	for _, c := range one.Children {
		modelTypeCache[c.Name] = c
	}
	Types = append(Types, one)
}

func GetTypes() []*Type {
	return Types
}

func GetTypeCache() map[string]*Type {
	return modelTypeCache
}
