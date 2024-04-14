package modelers

import (
	"fmt"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"os"
	"strings"
)

func newApplication() (app *Application) {
	app = &Application{}
	return
}

func Load(dir string) (app *Application) {
	app = newApplication()
	types := Types
	app.Dir = util.FormatPath(dir)
	if app.Dir == "" {
		app.Dir = dir
	}
	if !strings.HasSuffix(app.Dir, "/") {
		app.Dir += "/"
	}

	for _, modelType := range types {
		appendModelByType(app.Dir, app, modelType)
	}
	return
}

func appendModelByType(dir string, app *Application, modelType *Type) {
	baseDir := dir
	if modelType.Dir != "" {
		baseDir += modelType.Dir
	}
	if !strings.HasSuffix(baseDir, "/") {
		baseDir += "/"
	}
	if modelType.Children != nil {
		for _, one := range modelType.Children {
			appendModelByType(baseDir, app, one)
		}
	} else {
		loadFiles(baseDir, func(fileName string, fullName string) {
			appendModel(baseDir, app, modelType, fileName, fullName)
		})
	}
	return
}

func appendModel(baseDir string, app *Application, modelType *Type, fileName string, fullName string) {
	if !(strings.HasSuffix(fileName, ".yml") || strings.HasSuffix(fileName, ".yaml")) {
		return
	}
	if modelType.FileName != "" {
		if !strings.HasPrefix(fileName, modelType.FileName+".") {
			return
		}
	}
	if !strings.HasSuffix(baseDir, "/") {
		baseDir += "/"
	}
	path := strings.TrimPrefix(fullName, app.Dir)
	name := strings.TrimPrefix(fullName, baseDir)
	name = strings.TrimSuffix(name, ".yml")
	name = strings.TrimSuffix(name, ".yaml")
	fmt.Println("path:", path)
	fmt.Println("name:", name)
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
	bs, err := os.ReadFile(fullName)
	if err != nil {
		util.Logger.Error("appendModelByType ReadFile error", zap.Any("model", path), zap.Any("modelType", modelType), zap.Error(err))
		return
	}
	one, err := modelType.toModel(name, string(bs))
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
	files, _ := os.ReadDir(folder)
	for _, file := range files {
		if file.IsDir() {
			loadFiles(folder+file.Name()+"/", onLoad)
		} else {
			if onLoad != nil {
				onLoad(file.Name(), folder+file.Name())
			}
		}
	}

}
