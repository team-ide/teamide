package modelers

import (
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

type Type struct {
	Name     string `json:"name"`
	Comment  string `json:"comment"`
	Dir      string `json:"dir"`
	FileName string `json:"fileName"`

	toModel func(text string) (model interface{}, err error)
	toText  func(model interface{}) (text string, err error)
	append  func(app *Application, model interface{}) (err error)
}

var (
	Types []*Type

	TypeConstant = &Type{
		Name:    "constant",
		Comment: "常量",
		Dir:     "constant",
		toModel: func(text string) (model interface{}, err error) {
			var list []*ConstantModel
			err = yaml.Unmarshal([]byte(text), &list)
			if err != nil {
				util.Logger.Error("yaml to constants error", zap.Any("yaml", text), zap.Error(err))
				return
			}
			model = list
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
			list := model.([]*ConstantModel)
			for _, one := range list {
				err = app.AppendConstant(one)
			}
			return
		},
	}

	TypeStruct = &Type{
		Name:    "struct",
		Comment: "结构体",
		Dir:     "struct",
		toModel: func(text string) (model interface{}, err error) {
			model = &StructModel{}
			err = toModel(text, docTemplateStructName, model)
			if err != nil {
				util.Logger.Error("text to struct model error", zap.Any("text", text), zap.Error(err))
				return
			}
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
	}

	TypeService = &Type{
		Name:    "service",
		Comment: "服务",
		Dir:     "service",
		toModel: func(text string) (model interface{}, err error) {
			model = &ServiceModel{}
			err = toModel(text, docTemplateServiceName, model)
			if err != nil {
				util.Logger.Error("text to service model error", zap.Any("text", text), zap.Error(err))
				return
			}
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
	}

	TypeDao = &Type{
		Name:    "dao",
		Comment: "数据层",
		Dir:     "dao",
		toModel: func(text string) (model interface{}, err error) {
			model = &DaoModel{}
			err = toModel(text, docTemplateDaoName, model)
			if err != nil {
				util.Logger.Error("text to dao model error", zap.Any("text", text), zap.Error(err))
				return
			}
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
	}

	TypeFunc = &Type{
		Name:    "func",
		Comment: "函数",
		Dir:     "func",
		toModel: func(text string) (model interface{}, err error) {
			model = &FuncModel{}
			err = toModel(text, docTemplateFuncName, model)
			if err != nil {
				util.Logger.Error("text to func model error", zap.Any("text", text), zap.Error(err))
				return
			}
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
	}

	TypeError = &Type{
		Name:    "error",
		Comment: "错误码",
		Dir:     "error",
		toModel: func(text string) (model interface{}, err error) {
			var list []*ErrorModel
			err = yaml.Unmarshal([]byte(text), &list)
			if err != nil {
				util.Logger.Error("yaml to errors error", zap.Any("yaml", text), zap.Error(err))
				return
			}
			model = list
			return
		},
		toText: func(model interface{}) (text string, err error) {
			bytes, err := yaml.Marshal(model)
			if err != nil {
				util.Logger.Error("errors to yaml error", zap.Any("errors", model), zap.Error(err))
				return
			}
			text = string(bytes)
			return
		},
		append: func(app *Application, model interface{}) (err error) {
			list := model.([]*ErrorModel)
			for _, one := range list {
				err = app.AppendError(one)
			}
			return
		},
	}

	TypeLanguageJavascript = &Type{
		Name:     "languageJavascript",
		Comment:  "JavaScript语言",
		Dir:      "language",
		FileName: "javascript",
		toModel: func(text string) (model interface{}, err error) {
			model = &LanguageJavascriptModel{}
			err = yaml.Unmarshal([]byte(text), model)
			if err != nil {
				util.Logger.Error("yaml to language javascript error", zap.Any("yaml", text), zap.Error(err))
				return
			}
			return
		},
		toText: func(model interface{}) (text string, err error) {
			bytes, err := yaml.Marshal(model)
			if err != nil {
				util.Logger.Error("language javascript to yaml error", zap.Any("errors", model), zap.Error(err))
				return
			}
			text = string(bytes)
			return
		},
		append: func(app *Application, model interface{}) (err error) {
			app.SetLanguageJavascript(model.(*LanguageJavascriptModel))
			return
		},
	}

	TypeConfigRedis = &Type{
		Name:    "configRedis",
		Comment: "Redis配置",
		Dir:     "config/redis",
		toModel: func(text string) (model interface{}, err error) {
			model = &ConfigRedisModel{}
			err = yaml.Unmarshal([]byte(text), model)
			if err != nil {
				util.Logger.Error("yaml to config redis error", zap.Any("yaml", text), zap.Error(err))
				return
			}
			return
		},
		toText: func(model interface{}) (text string, err error) {
			bytes, err := yaml.Marshal(model)
			if err != nil {
				util.Logger.Error("config redis to yaml error", zap.Any("errors", model), zap.Error(err))
				return
			}
			text = string(bytes)
			return
		},
		append: func(app *Application, model interface{}) (err error) {
			err = app.AppendConfigRedis(model.(*ConfigRedisModel))
			return
		},
	}

	TypeConfigDb = &Type{
		Name:    "configDb",
		Comment: "Database配置",
		Dir:     "config/database",
		toModel: func(text string) (model interface{}, err error) {
			model = &ConfigDbModel{}
			err = yaml.Unmarshal([]byte(text), model)
			if err != nil {
				util.Logger.Error("yaml to config db error", zap.Any("yaml", text), zap.Error(err))
				return
			}
			return
		},
		toText: func(model interface{}) (text string, err error) {
			bytes, err := yaml.Marshal(model)
			if err != nil {
				util.Logger.Error("config db to yaml error", zap.Any("errors", model), zap.Error(err))
				return
			}
			text = string(bytes)
			return
		},
		append: func(app *Application, model interface{}) (err error) {
			err = app.AppendConfigDb(model.(*ConfigDbModel))
			return
		},
	}

	TypeConfigZk = &Type{
		Name:    "configZk",
		Comment: "Zookeeper配置",
		Dir:     "config/zookeeper",
		toModel: func(text string) (model interface{}, err error) {
			model = &ConfigZkModel{}
			err = yaml.Unmarshal([]byte(text), model)
			if err != nil {
				util.Logger.Error("yaml to config zk error", zap.Any("yaml", text), zap.Error(err))
				return
			}
			return
		},
		toText: func(model interface{}) (text string, err error) {
			bytes, err := yaml.Marshal(model)
			if err != nil {
				util.Logger.Error("config zk to yaml error", zap.Any("errors", model), zap.Error(err))
				return
			}
			text = string(bytes)
			return
		},
		append: func(app *Application, model interface{}) (err error) {
			err = app.AppendConfigZk(model.(*ConfigZkModel))
			return
		},
	}
)

func init() {
	Types = append(Types, TypeConstant)
	Types = append(Types, TypeStruct)
	Types = append(Types, TypeConfigRedis)
	Types = append(Types, TypeConfigDb)
	Types = append(Types, TypeConfigZk)
	Types = append(Types, TypeService)
	Types = append(Types, TypeDao)
	Types = append(Types, TypeError)
	Types = append(Types, TypeFunc)
	Types = append(Types, TypeLanguageJavascript)
}
