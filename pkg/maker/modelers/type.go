package modelers

import (
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type Type struct {
	Name     string  `json:"name"`
	Comment  string  `json:"comment"`
	Dir      string  `json:"dir"`
	FileName string  `json:"fileName"`
	Children []*Type `json:"children"`

	toModel func(name, text string) (model interface{}, err error)
	toText  func(model interface{}) (text string, err error)
	append  func(app *Application, model interface{}) (err error)
}

var (
	Types []*Type
)

func init() {
	AppendType(&Type{
		Name:    "constant",
		Comment: "常量",
		Dir:     "constant",
		toModel: func(name, text string) (model interface{}, err error) {
			model = &ConstantModel{}
			err = yaml.Unmarshal([]byte(text), model)
			if err != nil {
				util.Logger.Error("yaml to constant error", zap.Any("yaml", text), zap.Error(err))
				return
			}
			model.(*ConstantModel).Name = name
			return
		},
		toText: func(model interface{}) (text string, err error) {
			bytes, err := yaml.Marshal(model)
			if err != nil {
				util.Logger.Error("constants to yaml error", zap.Any("constants", model), zap.Error(err))
				return
			}
			text = string(bytes)
			return
		},
		append: func(app *Application, model interface{}) (err error) {
			err = app.AppendConstant(model.(*ConstantModel))
			return
		},
	})

	AppendType(&Type{
		Name:    "error",
		Comment: "错误码",
		Dir:     "error",
		toModel: func(name, text string) (model interface{}, err error) {
			model = &ErrorModel{}
			err = yaml.Unmarshal([]byte(text), model)
			if err != nil {
				util.Logger.Error("yaml to error error", zap.Any("text", text), zap.Error(err))
				return
			}
			model.(*ErrorModel).Name = name
			return
		},
		toText: func(model interface{}) (text string, err error) {
			bytes, err := yaml.Marshal(model)
			if err != nil {
				util.Logger.Error("error to yaml error", zap.Any("model", model), zap.Error(err))
				return
			}
			text = string(bytes)
			return
		},
		append: func(app *Application, model interface{}) (err error) {
			err = app.AppendError(model.(*ErrorModel))
			return
		},
	})

	AppendType(&Type{
		Name:    "struct",
		Comment: "结构体",
		Dir:     "struct",
		toModel: func(name, text string) (model interface{}, err error) {
			model = &StructModel{}
			err = toModel(text, docTemplateStructName, model)
			if err != nil {
				util.Logger.Error("text to struct model error", zap.Any("text", text), zap.Error(err))
				return
			}
			model.(*StructModel).Name = name
			return
		},
		toText: func(model interface{}) (text string, err error) {
			text, err = toText(model, docTemplateStructName, &docOptions{
				outComment: true,
				omitEmpty:  false,
			})
			if err != nil {
				util.Logger.Error("struct model to text error", zap.Any("model", model), zap.Error(err))
				return
			}
			return
		},
		append: func(app *Application, model interface{}) (err error) {
			err = app.AppendStruct(model.(*StructModel))
			return
		},
	})

	AppendType(&Type{
		Name:    "service",
		Comment: "服务",
		Dir:     "service",
		toModel: func(name, text string) (model interface{}, err error) {
			model = &ServiceModel{}
			err = toModel(text, docTemplateServiceName, model)
			if err != nil {
				util.Logger.Error("text to service model error", zap.Any("text", text), zap.Error(err))
				return
			}
			model.(*ServiceModel).Name = name
			return
		},
		toText: func(model interface{}) (text string, err error) {
			text, err = toText(model, docTemplateServiceName, &docOptions{
				outComment: true,
				omitEmpty:  false,
			})
			if err != nil {
				util.Logger.Error("service model to text error", zap.Any("model", model), zap.Error(err))
				return
			}
			return
		},
		append: func(app *Application, model interface{}) (err error) {
			err = app.AppendService(model.(*ServiceModel))
			return
		},
	})

	AppendType(&Type{
		Name:    "dao",
		Comment: "数据层",
		Dir:     "dao",
		toModel: func(name, text string) (model interface{}, err error) {
			model = &DaoModel{}
			err = toModel(text, docTemplateDaoName, model)
			if err != nil {
				util.Logger.Error("text to dao model error", zap.Any("text", text), zap.Error(err))
				return
			}
			model.(*DaoModel).Name = name
			return
		},
		toText: func(model interface{}) (text string, err error) {
			text, err = toText(model, docTemplateDaoName, &docOptions{
				outComment: true,
				omitEmpty:  false,
			})
			if err != nil {
				util.Logger.Error("dao model to text error", zap.Any("model", model), zap.Error(err))
				return
			}
			return
		},
		append: func(app *Application, model interface{}) (err error) {
			err = app.AppendDao(model.(*DaoModel))
			return
		},
	})

	AppendType(&Type{
		Name:    "func",
		Comment: "函数",
		Dir:     "func",
		toModel: func(name, text string) (model interface{}, err error) {
			model = &FuncModel{}
			err = toModel(text, docTemplateFuncName, model)
			if err != nil {
				util.Logger.Error("text to func model error", zap.Any("text", text), zap.Error(err))
				return
			}
			model.(*FuncModel).Name = name
			return
		},
		toText: func(model interface{}) (text string, err error) {
			text, err = toText(model, docTemplateFuncName, &docOptions{
				outComment: true,
				omitEmpty:  false,
			})
			if err != nil {
				util.Logger.Error("func model to text error", zap.Any("model", model), zap.Error(err))
				return
			}
			return
		},
		append: func(app *Application, model interface{}) (err error) {
			err = app.AppendFunc(model.(*FuncModel))
			return
		},
	})

	AppendType(&Type{
		Name:    "config",
		Comment: "配置",
		Dir:     "config",
		Children: []*Type{

			{
				Name:    "configRedis",
				Comment: "Redis",
				Dir:     "redis",
				toModel: func(name, text string) (model interface{}, err error) {
					model = &ConfigRedisModel{}
					err = yaml.Unmarshal([]byte(text), model)
					if err != nil {
						util.Logger.Error("yaml to config redis error", zap.Any("text", text), zap.Error(err))
						return
					}
					model.(*ConfigRedisModel).Name = name
					return
				},
				toText: func(model interface{}) (text string, err error) {
					bytes, err := yaml.Marshal(model)
					if err != nil {
						util.Logger.Error("config redis to yaml error", zap.Any("model", model), zap.Error(err))
						return
					}
					text = string(bytes)
					return
				},
				append: func(app *Application, model interface{}) (err error) {
					err = app.AppendConfigRedis(model.(*ConfigRedisModel))
					return
				},
			},

			{
				Name:    "configDb",
				Comment: "Database",
				Dir:     "database",
				toModel: func(name, text string) (model interface{}, err error) {
					model = &ConfigDbModel{}
					err = yaml.Unmarshal([]byte(text), model)
					if err != nil {
						util.Logger.Error("yaml to config db error", zap.Any("text", text), zap.Error(err))
						return
					}
					model.(*ConfigDbModel).Name = name
					return
				},
				toText: func(model interface{}) (text string, err error) {
					bytes, err := yaml.Marshal(model)
					if err != nil {
						util.Logger.Error("config db to yaml error", zap.Any("model", model), zap.Error(err))
						return
					}
					text = string(bytes)
					return
				},
				append: func(app *Application, model interface{}) (err error) {
					err = app.AppendConfigDb(model.(*ConfigDbModel))
					return
				},
			},

			{
				Name:    "configZk",
				Comment: "Zookeeper",
				Dir:     "zookeeper",
				toModel: func(name, text string) (model interface{}, err error) {
					model = &ConfigZkModel{}
					err = yaml.Unmarshal([]byte(text), model)
					if err != nil {
						util.Logger.Error("yaml to config zk error", zap.Any("text", text), zap.Error(err))
						return
					}
					model.(*ConfigZkModel).Name = name
					return
				},
				toText: func(model interface{}) (text string, err error) {
					bytes, err := yaml.Marshal(model)
					if err != nil {
						util.Logger.Error("config zk to yaml error", zap.Any("model", model), zap.Error(err))
						return
					}
					text = string(bytes)
					return
				},
				append: func(app *Application, model interface{}) (err error) {
					err = app.AppendConfigZk(model.(*ConfigZkModel))
					return
				},
			},

			{
				Name:    "configKafka",
				Comment: "Kafka",
				Dir:     "kafka",
				toModel: func(name, text string) (model interface{}, err error) {
					model = &ConfigKafkaModel{}
					err = yaml.Unmarshal([]byte(text), model)
					if err != nil {
						util.Logger.Error("yaml to config kafka error", zap.Any("text", text), zap.Error(err))
						return
					}
					model.(*ConfigKafkaModel).Name = name
					return
				},
				toText: func(model interface{}) (text string, err error) {
					bytes, err := yaml.Marshal(model)
					if err != nil {
						util.Logger.Error("config kafka to yaml error", zap.Any("model", model), zap.Error(err))
						return
					}
					text = string(bytes)
					return
				},
				append: func(app *Application, model interface{}) (err error) {
					err = app.AppendConfigKafka(model.(*ConfigKafkaModel))
					return
				},
			},
			{
				Name:    "configMongodb",
				Comment: "Mongodb",
				Dir:     "mongodb",
				toModel: func(name, text string) (model interface{}, err error) {
					model = &ConfigMongodbModel{}
					err = yaml.Unmarshal([]byte(text), model)
					if err != nil {
						util.Logger.Error("yaml to config mongodb error", zap.Any("text", text), zap.Error(err))
						return
					}
					model.(*ConfigMongodbModel).Name = name
					return
				},
				toText: func(model interface{}) (text string, err error) {
					bytes, err := yaml.Marshal(model)
					if err != nil {
						util.Logger.Error("config mongodb to yaml error", zap.Any("model", model), zap.Error(err))
						return
					}
					text = string(bytes)
					return
				},
				append: func(app *Application, model interface{}) (err error) {
					err = app.AppendConfigMongodb(model.(*ConfigMongodbModel))
					return
				},
			},
		},
	})

	AppendType(&Type{
		Name:    "language",
		Comment: "导出语言",
		Dir:     "language",
		Children: []*Type{
			{
				Name:     "languageJavascript",
				Comment:  "JavaScript",
				Dir:      "",
				FileName: "javascript",
				toModel: func(name, text string) (model interface{}, err error) {
					model = &LanguageJavascriptModel{}
					err = yaml.Unmarshal([]byte(text), model)
					if err != nil {
						util.Logger.Error("yaml to language javascript error", zap.Any("text", text), zap.Error(err))
						return
					}
					return
				},
				toText: func(model interface{}) (text string, err error) {
					bytes, err := yaml.Marshal(model)
					if err != nil {
						util.Logger.Error("language javascript to yaml error", zap.Any("model", model), zap.Error(err))
						return
					}
					text = string(bytes)
					return
				},
				append: func(app *Application, model interface{}) (err error) {
					app.SetLanguageJavascript(model.(*LanguageJavascriptModel))
					return
				},
			},
		},
	})
}

func AppendType(one *Type) {
	Types = append(Types, one)
}

func GetTypes() []*Type {
	return Types
}
