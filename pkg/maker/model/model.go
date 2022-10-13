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
			errorList := model.([]*ErrorModel)
			for _, one := range errorList {
				err = app.AppendError(one)
			}
			return
		},
	}
)

func init() {
	Types = append(Types, TypeStruct)
	Types = append(Types, TypeService)
	Types = append(Types, TypeDao)
	Types = append(Types, TypeError)
}
