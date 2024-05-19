package maker

import (
	"encoding/json"
	"github.com/dop251/goja"
	"github.com/team-ide/go-tool/util"
	"go.uber.org/zap"
	"os"
	"strings"
	"teamide/pkg/maker/modelers"
)

func newApplication() (app *Application) {
	app = &Application{
		elementCache:    make(map[string]*modelers.Element),
		modelTypeCaches: make(map[*modelers.Type]*util.Cache),
		modelTypeItems:  make(map[*modelers.Type][]modelers.ElementIFace),

		constantContext: make(map[string]interface{}),
		errorContext:    make(map[string]*Error),

		daoContext: make(map[string]interface{}),
		daoProgram: make(map[string]*goja.Program),

		serviceContext: make(map[string]interface{}),
		serviceProgram: make(map[string]*goja.Program),

		funcContext: make(map[string]interface{}),
		funcProgram: make(map[string]*goja.Program),
	}
	return
}

func Load(dir string) (app *Application) {
	app = newApplication()
	types := modelers.Types
	app.dir = util.FormatPath(dir)
	if app.dir == "" {
		app.dir = dir
	}
	if !strings.HasSuffix(app.dir, "/") {
		app.dir += "/"
	}

	for _, modelType := range types {
		app.loadByType(nil, modelType)
	}
	return
}

func (this_ *Application) loadByType(parent *modelers.Element, modelType *modelers.Type) {
	typeElement := this_.appendType(parent, modelType)
	if modelType.Children != nil {
		for _, one := range modelType.Children {
			this_.loadByType(typeElement, one)
		}
		return
	}
	modelTypePath := this_.getModeTypePath(modelType)
	if modelType.IsFile {
		if exist, _ := util.PathExists(modelTypePath + ".yml"); exist {
			this_.loadFile(typeElement, modelType, modelTypePath+".yml")
		}
	} else {
		this_.loadFiles(typeElement, modelType, modelTypePath)
	}
	return
}

func (this_ *Application) loadFiles(parent *modelers.Element, modelType *modelers.Type, folder string) {
	files, _ := os.ReadDir(folder)
	for _, file := range files {
		filePath := folder + file.Name()
		if file.IsDir() {
			packElement := this_.appendPack(parent, modelType, file.Name())
			this_.loadFiles(packElement, modelType, filePath+"/")
		} else {
			this_.loadFile(parent, modelType, filePath)
		}
	}
}

func (this_ *Application) loadFile(parent *modelers.Element, modelType *modelers.Type, filePath string) (model interface{}, element *modelers.Element) {
	if !(strings.HasSuffix(filePath, ".yml")) {
		return
	}

	path := strings.TrimPrefix(filePath, this_.dir)
	modelTypePath := this_.getModeTypePath(modelType)
	filename := filePath[strings.LastIndex(filePath, "/")+1:]
	filename = strings.TrimSuffix(filename, ".yml")

	var name = filename
	if !modelType.IsFile {
		name = strings.TrimPrefix(filePath, modelTypePath)
		name = strings.TrimSuffix(name, ".yml")
	}

	//fmt.Println("path:", path)
	//fmt.Println("name:", name)

	var err error
	defer func() {
		if err != nil {
			loadError := &LoadError{}
			loadError.Type = modelType
			loadError.Path = path
			loadError.Error = err.Error()
			this_.LoadErrors = append(this_.LoadErrors, loadError)
		}
	}()
	bs, err := os.ReadFile(filePath)
	if err != nil {
		util.Logger.Error("loadFile ReadFile error", zap.Any("model", path), zap.Any("modelType", modelType), zap.Error(err))
		return
	}
	if filename == packInfoFileName {
		parent.Pack = &modelers.Pack{}
		err = json.Unmarshal(bs, parent.Pack)
		if err != nil {
			util.Logger.Error("loadFile toPack error", zap.Any("path", path), zap.Error(err))
			return
		}
	} else {
		model, err = modelType.ToModel(name, string(bs))
		if err != nil {
			util.Logger.Error("loadFile ToModel error", zap.Any("model", path), zap.Any("modelType", modelType), zap.Any("text", string(bs)), zap.Error(err))
			return
		}
		element, err = this_.appendModel(parent, modelType, filename, name, model)
		if err != nil {
			util.Logger.Error("loadFile Append error", zap.Any("model", path), zap.Any("modelType", modelType), zap.Any("model", model), zap.Error(err))
			return
		}
	}
	return
}

var (
	packInfoFileName = "pack-info"
)
