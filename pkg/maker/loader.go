package maker

import (
	"go.uber.org/zap"
	"io/ioutil"
	"path/filepath"
	"teamide/pkg/maker/model"
	"teamide/pkg/util"
)

func newApplication() (app *model.Application) {
	app = &model.Application{}
	return
}

func Load(dir string) (app *model.Application) {
	dir_, err := filepath.Abs(dir)
	if err != nil {
		util.Logger.Error("filepath Abs error", zap.Any("dir", dir), zap.Error(err))
	} else {
		dir = dir_
	}
	app = newApplication()
	types := model.Types

	for _, modelType := range types {
		appendModelByType(dir, app, modelType)
	}
	return
}

func appendModelByType(dir string, app *model.Application, modelType *model.Type) {

	pathname := dir + "/" + modelType.Dir
	LoadFiles(pathname, func(fullName string) {
		fullName_, err := filepath.Abs(fullName)
		if err != nil {
			util.Logger.Error("filepath Abs error", zap.Any("fullName", fullName), zap.Error(err))
		} else {
			fullName = fullName_
		}
		bs, err := ioutil.ReadFile(fullName)
		if err != nil {
			util.Logger.Error("appendModelByType ReadFile error", zap.Any("fullName", fullName), zap.Any("modelType", modelType), zap.Error(err))
			return
		}
		one, err := modelType.ToModel(string(bs))
		if err != nil {
			util.Logger.Error("appendModelByType ToModel error", zap.Any("fullName", fullName), zap.Any("modelType", modelType), zap.Any("text", string(bs)), zap.Error(err))
			return
		}
		err = modelType.Append(app, one)
		if err != nil {
			util.Logger.Error("appendModelByType Append error", zap.Any("fullName", fullName), zap.Any("modelType", modelType), zap.Any("model", one), zap.Error(err))
			return
		}
	})
	return
}

func LoadFiles(folder string, onLoad func(pathname string)) {
	files, _ := ioutil.ReadDir(folder)
	for _, file := range files {
		if file.IsDir() {
			LoadFiles(folder+"/"+file.Name(), onLoad)
		} else {
			if onLoad != nil {
				onLoad(folder + "/" + file.Name())
			}
		}
	}

}
