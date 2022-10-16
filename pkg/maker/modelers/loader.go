package modelers

import (
	"go.uber.org/zap"
	"io/ioutil"
	"path/filepath"
	"strings"
	"teamide/pkg/util"
)

func newApplication() (app *Application) {
	app = &Application{}
	return
}

func Load(dir string) (app *Application) {
	dir_, err := filepath.Abs(dir)
	if err != nil {
		util.Logger.Error("filepath Abs error", zap.Any("dir", dir), zap.Error(err))
	} else {
		dir = dir_
	}
	app = newApplication()
	types := Types

	for _, modelType := range types {
		appendModelByType(dir, app, modelType)
	}
	return
}

func appendModelByType(dir string, app *Application, modelType *Type) {

	pathname := dir + "/" + modelType.Dir
	loadFiles(pathname, func(fileName string, fullName string) {
		appendModel(dir, app, modelType, fileName, fullName)
	})
	return
}

func appendModel(dir string, app *Application, modelType *Type, fileName string, fullName string) {
	if modelType.FileName != "" {
		if !strings.HasPrefix(fileName, modelType.FileName+".") {
			return
		}
	}
	path := strings.TrimLeft(fullName, dir)
	var err error
	defer func() {
		if err != nil {
			loadError := &LoadError{}
			loadError.Type = modelType
			loadError.Path = path
			loadError.Error = err.Error()
			app.LoadErrors = append(app.LoadErrors, loadError)
		}
	}()
	fullName_, err := filepath.Abs(fullName)
	if err != nil {
		util.Logger.Error("filepath Abs error", zap.Any("fullName", fullName), zap.Error(err))
	} else {
		fullName = fullName_
	}
	bs, err := ioutil.ReadFile(fullName)
	if err != nil {
		util.Logger.Error("appendModelByType ReadFile error", zap.Any("model", path), zap.Any("modelType", modelType), zap.Error(err))
		return
	}
	one, err := modelType.toModel(string(bs))
	if err != nil {
		util.Logger.Error("appendModelByType ToModel error", zap.Any("model", path), zap.Any("modelType", modelType), zap.Any("text", string(bs)), zap.Error(err))
		return
	}
	err = modelType.append(app, one)
	if err != nil {
		util.Logger.Error("appendModelByType Append error", zap.Any("model", path), zap.Any("modelType", modelType), zap.Any("model", one), zap.Error(err))
		return
	}
	return
}
func loadFiles(folder string, onLoad func(name string, pathname string)) {
	files, _ := ioutil.ReadDir(folder)
	for _, file := range files {
		if file.IsDir() {
			loadFiles(folder+"/"+file.Name(), onLoad)
		} else {
			if onLoad != nil {
				onLoad(file.Name(), folder+"/"+file.Name())
			}
		}
	}

}
