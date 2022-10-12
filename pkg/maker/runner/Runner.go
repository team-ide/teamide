package runner

import "teamide/pkg/maker/model"

func NewRunner(app *model.Application) (runner *Runner) {
	runner = &Runner{
		app: app,
	}
	return
}

type Runner struct {
	app *model.Application
}
