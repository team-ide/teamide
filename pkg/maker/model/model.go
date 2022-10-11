package model

type Type struct {
	Name    string                                                `json:"name"`
	Comment string                                                `json:"comment"`
	Dir     string                                                `json:"dir"`
	ToModel func(text string) (model interface{}, err error)      `json:"-"`
	ToText  func(model interface{}) (text string, err error)      `json:"-"`
	Append  func(app *Application, model interface{}) (err error) `json:"-"`
}

var (
	Types []*Type

	TypeStruct = &Type{
		Name:    "struct",
		Comment: "结构体",
		Dir:     "struct",
		ToModel: func(text string) (model interface{}, err error) {
			return TextToStruct(text)
		},
		ToText: func(model interface{}) (text string, err error) {
			return StructToText(model)
		},
		Append: func(app *Application, model interface{}) (err error) {
			err = app.AppendStruct(model.(*StructModel))
			return
		},
	}

	TypeService = &Type{
		Name:    "service",
		Comment: "服务",
		Dir:     "service",
		ToModel: func(text string) (model interface{}, err error) {
			return TextToService(text)
		},
		ToText: func(model interface{}) (text string, err error) {
			return ServiceToText(model)
		},
		Append: func(app *Application, model interface{}) (err error) {
			err = app.AppendService(model.(*ServiceModel))
			return
		},
	}
)

func init() {
	Types = append(Types, TypeStruct)
	Types = append(Types, TypeService)
}
