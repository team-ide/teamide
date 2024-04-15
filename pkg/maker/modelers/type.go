package modelers

import (
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
)

type Type struct {
	Name     string  `json:"name"`
	Comment  string  `json:"comment"`
	Dir      string  `json:"dir"`
	FileName string  `json:"fileName"`
	Children []*Type `json:"children"`

	toModel func(name, text string) (model interface{}, err error)
	toText  func(model interface{}) (text string, err error)
}

var (
	Types []*Type

	TypeConstantName = "constant"

	TypeConstant = &Type{
		Name:    TypeConstantName,
		Comment: "常量",
		Dir:     "constant",
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
		Name:    TypeErrorName,
		Comment: "错误码",
		Dir:     "error",
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
		Name:    TypeStructName,
		Comment: "结构体",
		Dir:     "struct",
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

	TypeServiceName = "service"
	TypeService     = &Type{
		Name:    TypeServiceName,
		Comment: "服务",
		Dir:     "service",
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

	TypeDaoName = "dao"
	TypeDao     = &Type{
		Name:    TypeDaoName,
		Comment: "数据层",
		Dir:     "dao",
		toModel: func(name, text string) (model interface{}, err error) {
			model = &DaoModel{}
			err = toModel(text, TypeDaoName, model)
			if err != nil {
				util.Logger.Error("text to dao model error", zap.Any("text", text), zap.Error(err))
				return
			}
			model.(*DaoModel).Name = name
			return
		},
		toText: func(model interface{}) (text string, err error) {
			text, err = toText(model, TypeDaoName, &docOptions{
				outComment: true,
				omitEmpty:  false,
			})
			if err != nil {
				util.Logger.Error("dao model to text error", zap.Any("model", model), zap.Error(err))
				return
			}
			return
		},
	}

	TypeFuncName = "func"
	TypeFunc     = &Type{
		Name:    TypeFuncName,
		Comment: "函数",
		Dir:     "func",
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

	TypeConfigDbName = "config/database"
	TypeConfigDb     = &Type{
		Name:    TypeConfigDbName,
		Comment: "Database",
		Dir:     "database",
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
		Name:    TypeConfigRedisName,
		Comment: "Redis",
		Dir:     "redis",
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

	TypeConfigZkName = "config/zookeeper"
	TypeConfigZk     = &Type{
		Name:    TypeConfigZkName,
		Comment: "Zookeeper",
		Dir:     "zookeeper",
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
		Name:    TypeConfigKafkaName,
		Comment: "Kafka",
		Dir:     "kafka",
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
		Name:    TypeConfigMongodbName,
		Comment: "Mongodb",
		Dir:     "mongodb",
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

	TypeConfigElasticsearchName = "config/elasticsearch"
	TypeConfigElasticsearch     = &Type{
		Name:    TypeConfigElasticsearchName,
		Comment: "Elastic Search",
		Dir:     "elasticsearch",
		toModel: func(name, text string) (model interface{}, err error) {
			model = &ConfigMongodbModel{}
			err = toModel(text, TypeConfigElasticsearchName, model)
			if err != nil {
				util.Logger.Error("text to config elasticsearch model error", zap.Any("text", text), zap.Error(err))
				return
			}
			model.(*ConfigMongodbModel).Name = name
			return
		},
		toText: func(model interface{}) (text string, err error) {
			text, err = toText(model, TypeConfigElasticsearchName, &docOptions{
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

	TypeLanguageJavascriptName = "language/javascript"
	TypeLanguageJavascript     = &Type{
		Name:     TypeLanguageJavascriptName,
		Comment:  "JavaScript",
		Dir:      "",
		FileName: "javascript",
		toModel: func(name, text string) (model interface{}, err error) {
			model = &LanguageJavascriptModel{}
			err = toModel(text, TypeLanguageJavascriptName, model)
			if err != nil {
				util.Logger.Error("text to language javascript model error", zap.Any("text", text), zap.Error(err))
				return
			}
			return
		},
		toText: func(model interface{}) (text string, err error) {
			text, err = toText(model, TypeLanguageJavascriptName, &docOptions{
				outComment: true,
				omitEmpty:  false,
			})
			if err != nil {
				util.Logger.Error("language javascript model to text error", zap.Any("model", model), zap.Error(err))
				return
			}
			return
		},
	}
)

func init() {
	AppendType(TypeConstant)

	AppendType(TypeError)

	AppendType(TypeStruct)

	AppendType(TypeDao)

	AppendType(TypeService)

	AppendType(TypeFunc)

	AppendType(&Type{
		Name:    "config",
		Comment: "配置",
		Dir:     "config",
		Children: []*Type{

			TypeConfigDb,
			TypeConfigRedis,
			TypeConfigZk,
			TypeConfigKafka,
			TypeConfigMongodb,
			TypeConfigElasticsearch,
		},
	})

	AppendType(&Type{
		Name:    "language",
		Comment: "导出语言",
		Dir:     "language",
		Children: []*Type{
			TypeLanguageJavascript,
		},
	})
}

func AppendType(one *Type) {
	Types = append(Types, one)
}

func GetTypes() []*Type {
	return Types
}
