package modelers

import (
	"encoding/json"
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
		appendModelByType(nil, app.Dir, app, modelType)
	}
	return
}

func appendModelByType(parent *Element, dir string, app *Application, modelType *Type) {
	baseDir := dir
	if modelType.Dir != "" {
		baseDir += modelType.Dir
	}
	if !strings.HasSuffix(baseDir, "/") {
		baseDir += "/"
	}
	typeElement := app.appendType(parent, modelType)
	if modelType.Children != nil {
		for _, one := range modelType.Children {
			appendModelByType(typeElement, baseDir, app, one)
		}
	} else {
		loadFiles(typeElement, app, modelType, baseDir, func(parent *Element, app *Application, fileName string, fullName string) {
			appendModel(parent, baseDir, app, modelType, fileName, fullName)
		})
	}
	return
}

var (
	packInfoFileName = "pack-info"
)

func appendModel(parent *Element, baseDir string, app *Application, modelType *Type, fileName string, fullName string) {
	if !(strings.HasSuffix(fileName, ".yml") || strings.HasSuffix(fileName, ".yaml")) {
		return
	}
	fileName = strings.TrimSuffix(fileName, ".yml")
	fileName = strings.TrimSuffix(fileName, ".yaml")
	if modelType.FileName != "" && fileName != packInfoFileName {
		if fileName != modelType.FileName {
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
	if fileName == packInfoFileName {
		parent.Pack = &Pack{}
		err = json.Unmarshal(bs, parent.Pack)
		if err != nil {
			util.Logger.Error("appendModelByType toPack error", zap.Any("path", path), zap.Error(err))
			return
		}
	} else {
		var one interface{}
		one, err = modelType.toModel(name, string(bs))
		if err != nil {
			util.Logger.Error("appendModelByType ToModel error", zap.Any("model", path), zap.Any("modelType", modelType), zap.Any("text", string(bs)), zap.Error(err))
			return
		}
		err = app.appendModel(parent, modelType, fileName, name, one)
		if err != nil {
			util.Logger.Error("appendModelByType Append error", zap.Any("model", path), zap.Any("modelType", modelType), zap.Any("model", one), zap.Error(err))
			return
		}
	}
	return
}
func loadFiles(parent *Element, app *Application, modelType *Type, folder string, onLoad func(parent *Element, app *Application, name string, pathname string)) {
	files, _ := os.ReadDir(folder)
	for _, file := range files {
		if file.IsDir() {
			packElement := app.appendPack(parent, modelType, file.Name())
			loadFiles(packElement, app, modelType, folder+file.Name()+"/", onLoad)
		} else {
			if onLoad != nil {
				onLoad(parent, app, file.Name(), folder+file.Name())
			}
		}
	}

}
