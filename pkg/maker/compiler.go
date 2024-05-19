package maker

func NewCompiler(app *Application) (runner *Compiler, err error) {
	runner = &Compiler{
		app: app,

		constantContext: make(map[string]interface{}),
	}

	err = runner.init()

	return
}

type Compiler struct {
	app *Application

	constantContext map[string]interface{}

	script *Script
}

func (this_ *Compiler) init() (err error) {
	return
}
