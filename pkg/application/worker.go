package application

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"teamide/pkg/application/base"
	"teamide/pkg/application/model"
)

type worker struct {
	Dir             string
	dirAbsolutePath string
}

func NewWorker(dir string) *worker {
	res := &worker{Dir: dir}
	res.init()
	return res
}

func (this_ *worker) init() {
	var abs string
	abs, _ = filepath.Abs(this_.Dir)
	this_.dirAbsolutePath = filepath.ToSlash(abs)
}

func (this_ *worker) GetModelPath(modelType *model.ModelType, name string) string {
	if name == "" {
		return this_.dirAbsolutePath + "/" + modelType.Dir + "/default.yaml"
	}
	return this_.dirAbsolutePath + "/" + modelType.Dir + "/" + name + ".yaml"
}

func (this_ *worker) Load() (context *model.ModelContext, err error) {

	var exists bool
	exists, err = base.PathExists(this_.Dir)
	if err != nil {
		return
	}
	if !exists {
		err = errors.New(fmt.Sprint("dir [" + this_.Dir + "] not exist"))
		return
	}

	fileContentMap := make(map[string]string)
	err = loadDirFiles(fileContentMap, this_.dirAbsolutePath)
	var namePaths []string
	for namePath := range fileContentMap {
		namePaths = append(namePaths, namePath)
	}

	sort.Strings(namePaths)

	context = &model.ModelContext{}
	for _, namePath := range namePaths {
		if strings.HasSuffix(namePath, ".md") {
			continue
		}
		text := fileContentMap[namePath]
		if strings.Index(namePath, "struct/") == 0 {
			var one *model.StructModel
			one, err = model.TextToStructModel(strings.TrimPrefix(namePath, "struct/"), text)
			if err != nil {
				return
			}
			context.AppendStruct(one)
		} else if strings.Index(namePath, "action/") == 0 {
			var one *model.ActionModel
			one, err = model.TextToActionModel(strings.TrimPrefix(namePath, "action/"), text)
			if err != nil {
				return
			}
			context.AppendAction(one)
		} else if strings.Index(namePath, "test/") == 0 {
			var one *model.TestModel
			one, err = model.TextToTestModel(strings.TrimPrefix(namePath, "test/"), text)
			if err != nil {
				return
			}
			context.AppendTest(one)
		} else if strings.Index(namePath, "teamide/server/web/") == 0 {
			var one *model.ServerWebModel
			one, err = model.TextToServerWebModel(strings.TrimPrefix(namePath, "teamide/server/web/"), text)
			if err != nil {
				return
			}
			context.AppendServerWeb(one)
		} else if strings.Index(namePath, "dictionary/") == 0 {
			var one *model.DictionaryModel
			one, err = model.TextToDictionaryModel(strings.TrimPrefix(namePath, "dictionary/"), text)
			if err != nil {
				return
			}
			context.AppendDictionary(one)
		} else if strings.Index(namePath, "constant/") == 0 {
			var one *model.ConstantModel
			one, err = model.TextToConstantModel(strings.TrimPrefix(namePath, "constant/"), text)
			if err != nil {
				return
			}
			context.AppendConstant(one)
		} else if strings.Index(namePath, "error/") == 0 {
			var one *model.ErrorModel
			one, err = model.TextToErrorModel(strings.TrimPrefix(namePath, "error/"), text)
			if err != nil {
				return
			}
			context.AppendError(one)
		} else if strings.Index(namePath, "datasource/database/") == 0 {
			var one *model.DatasourceDatabase
			one, err = model.TextToDatasourceDatabase(strings.TrimPrefix(namePath, "datasource/database/"), text)
			if err != nil {
				return
			}
			context.AppendDatasourceDatabase(one)
		} else if strings.Index(namePath, "datasource/redis/") == 0 {
			var one *model.DatasourceRedis
			one, err = model.TextToDatasourceRedis(strings.TrimPrefix(namePath, "datasource/redis/"), text)
			if err != nil {
				return
			}
			context.AppendDatasourceRedis(one)
		} else if strings.Index(namePath, "datasource/kafka/") == 0 {
			var one *model.DatasourceKafka
			one, err = model.TextToDatasourceKafka(strings.TrimPrefix(namePath, "datasource/kafka/"), text)
			if err != nil {
				return
			}
			context.AppendDatasourceKafka(one)
		} else if strings.Index(namePath, "datasource/zookeeper/") == 0 {
			var one *model.DatasourceZookeeper
			one, err = model.TextToDatasourceZookeeper(strings.TrimPrefix(namePath, "datasource/zookeeper/"), text)
			if err != nil {
				return
			}
			context.AppendDatasourceZookeeper(one)
		}
	}
	return
}

func (this_ *worker) Save(context *model.ModelContext) (err error) {
	for _, one := range context.Constants {
		err = this_.ModelSave(model.MODEL_TYPE_CONSTANT, one.Name, one)
		if err != nil {
			return
		}
	}
	for _, one := range context.Errors {
		err = this_.ModelSave(model.MODEL_TYPE_ERROR, one.Name, one)
		if err != nil {
			return
		}
	}
	for _, one := range context.Dictionaries {
		err = this_.ModelSave(model.MODEL_TYPE_DICTIONARY, one.Name, one)
		if err != nil {
			return
		}
	}
	for _, one := range context.Structs {
		err = this_.ModelSave(model.MODEL_TYPE_STRUCT, one.Name, one)
		if err != nil {
			return
		}
	}
	for _, one := range context.Actions {
		err = this_.ModelSave(model.MODEL_TYPE_ACTION, one.Name, one)
		if err != nil {
			return
		}
	}
	for _, one := range context.Tests {
		err = this_.ModelSave(model.MODEL_TYPE_TEST, one.Name, one)
		if err != nil {
			return
		}
	}
	for _, one := range context.ServerWebs {
		err = this_.ModelSave(model.MODEL_TYPE_SERVER_WEB, one.Name, one)
		if err != nil {
			return
		}
	}
	for _, one := range context.DatasourceDatabases {
		err = this_.ModelSave(model.MODEL_TYPE_DATASOURCE_DATABASE, one.Name, one)
		if err != nil {
			return
		}
	}
	for _, one := range context.DatasourceRedises {
		err = this_.ModelSave(model.MODEL_TYPE_DATASOURCE_REDIS, one.Name, one)
		if err != nil {
			return
		}
	}
	for _, one := range context.DatasourceKafkas {
		err = this_.ModelSave(model.MODEL_TYPE_DATASOURCE_KAFKA, one.Name, one)
		if err != nil {
			return
		}
	}
	for _, one := range context.DatasourceZookeepers {
		err = this_.ModelSave(model.MODEL_TYPE_DATASOURCE_ZOOKEEPER, one.Name, one)
		if err != nil {
			return
		}
	}
	return
}

func (this_ *worker) ModelInsert(modelType *model.ModelType, name string) (err error) {
	path := this_.GetModelPath(modelType, name)
	var exists bool
	exists, err = base.PathExists(path)
	if err != nil {
		return
	}
	if exists {
		err = errors.New(fmt.Sprint("应用模型[", modelType.Text, "][", name, "]已存在"))
		return
	} else {
		var file *os.File
		file, err = os.Create(path)
		if err != nil {
			return
		}
		defer file.Close()
	}
	return
}

func (this_ *worker) ModelDelete(modelType *model.ModelType, name string) (err error) {
	path := this_.GetModelPath(modelType, name)
	var exists bool
	exists, err = base.PathExists(path)
	if err != nil {
		return
	}
	if exists {
		err = os.Remove(path)
		return
	}
	return
}

func (this_ *worker) ModelRename(modelType *model.ModelType, name string, rename string) (err error) {
	path := this_.GetModelPath(modelType, name)
	renamePath := this_.GetModelPath(modelType, rename)
	var exists bool
	exists, err = base.PathExists(path)
	if err != nil {
		return
	}
	if !exists {
		err = errors.New(fmt.Sprint("应用模型[", modelType.Text, "][", name, "]不存在"))
		return
	}
	exists, err = base.PathExists(renamePath)
	if err != nil {
		return
	}
	if exists {
		err = errors.New(fmt.Sprint("应用模型[", modelType.Text, "][", name, "]已存在"))
		return
	}
	err = os.Rename(path, renamePath)
	return
}

func (this_ *worker) ModelSave(modelType *model.ModelType, name string, data interface{}) (err error) {

	text := ""
	if data != nil {
		dataMap, dataMapOk := data.(map[string]interface{})
		if dataMapOk {
			data, err = model.MapToModel(modelType, dataMap)
			if err != nil {
				return
			}
		}
		switch m := data.(type) {
		case *model.ConstantModel:
			text, err = model.ConstantModelToText(m)
		case *model.ErrorModel:
			text, err = model.ErrorModelToText(m)
		case *model.DictionaryModel:
			text, err = model.DictionaryModelToText(m)
		case *model.StructModel:
			text, err = model.StructModelToText(m)
		case *model.ActionModel:
			text, err = model.ActionModelToText(m)
		case *model.TestModel:
			text, err = model.TestModelToText(m)
		case *model.ServerWebModel:
			text, err = model.ServerWebModelToText(m)
		case *model.DatasourceDatabase:
			text, err = model.DatasourceDatabaseToText(m)
		case *model.DatasourceRedis:
			text, err = model.DatasourceRedisToText(m)
		case *model.DatasourceKafka:
			text, err = model.DatasourceKafkaToText(m)
		case *model.DatasourceZookeeper:
			text, err = model.DatasourceZookeeperToText(m)
		}
		if err != nil {
			return
		}
	}
	path := this_.GetModelPath(modelType, name)
	var exists bool
	exists, err = base.PathExists(path)
	if err != nil {
		return
	}
	var file *os.File
	if exists {
		file, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	} else {
		file, err = os.Create(path)
	}
	if err != nil {
		return
	}
	defer file.Close()
	file.WriteString(text)
	return
}
