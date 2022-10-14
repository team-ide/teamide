package coder

import (
	"teamide/pkg/maker/coder/common"
	"teamide/pkg/maker/model"
)

type ApplicationGenerator struct {
	app              *model.Application
	generatorFactory common.GeneratorFactory
	appGenerator     common.AppGenerator
}

func NewApplicationGenerator(app *model.Application, generatorFactory common.GeneratorFactory) (applicationGenerator *ApplicationGenerator) {
	applicationGenerator = &ApplicationGenerator{
		app:              app,
		generatorFactory: generatorFactory,
		appGenerator:     generatorFactory.GetApp(),
	}
	return
}

func (this_ *ApplicationGenerator) GetService(model *model.ServiceModel) (codes []*common.Code, err error) {
	codeGenerators := this_.generatorFactory.GetService()
	for _, codeGenerator := range codeGenerators {
		codeGenerator.Model = model
		codeGenerator.App = this_.appGenerator
		codeGenerator.Code = &common.Code{}
		err = codeGenerator.Gen()
		if err != nil {
			return
		}
		codes = append(codes, codeGenerator.Code)
	}
	return
}

func (this_ *ApplicationGenerator) GetDao(model *model.DaoModel) (codes []*common.Code, err error) {
	codeGenerators := this_.generatorFactory.GetDao()
	for _, codeGenerator := range codeGenerators {
		codeGenerator.Model = model
		codeGenerator.App = this_.appGenerator
		codeGenerator.Code = &common.Code{}
		err = codeGenerator.Gen()
		if err != nil {
			return
		}
		codes = append(codes, codeGenerator.Code)
	}
	return
}

func (this_ *ApplicationGenerator) genSteps(code *common.Code, steps []interface{}) (err error) {
	codeGenerators := this_.generatorFactory.GetStep()
	for _, codeGenerator := range codeGenerators {
		for _, step := range steps {
			codeGenerator.Step = step
			codeGenerator.App = this_.appGenerator
			codeGenerator.Code = code
			err = codeGenerator.Gen()
			if err != nil {
				return
			}
		}
	}
	return
}
