package model

type Type struct {
	Name    string `json:"name"`
	Comment string `json:"comment"`
	Dir     string `json:"dir"`

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
			return TextToConstants(text)
		},
		toText: func(model interface{}) (text string, err error) {
			return ConstantsToText(model.([]*ConstantModel))
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
			return TextToStruct(text)
		},
		toText: func(model interface{}) (text string, err error) {
			return StructToText(model)
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
			return TextToService(text)
		},
		toText: func(model interface{}) (text string, err error) {
			return ServiceToText(model)
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
			return TextToDao(text)
		},
		toText: func(model interface{}) (text string, err error) {
			return DaoToText(model)
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
			return TextToFunc(text)
		},
		toText: func(model interface{}) (text string, err error) {
			return FuncToText(model)
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
			return TextToErrors(text)
		},
		toText: func(model interface{}) (text string, err error) {
			return ErrorsToText(model.([]*ErrorModel))
		},
		append: func(app *Application, model interface{}) (err error) {
			list := model.([]*ErrorModel)
			for _, one := range list {
				err = app.AppendError(one)
			}
			return
		},
	}
)

func init() {
	Types = append(Types, TypeConstant)
	Types = append(Types, TypeStruct)
	Types = append(Types, TypeService)
	Types = append(Types, TypeDao)
	Types = append(Types, TypeError)
	Types = append(Types, TypeFunc)
}
