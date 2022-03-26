package application

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"teamide/pkg/application/base"
	"teamide/pkg/application/model"
)

type ContextLoader struct {
	Dir             string
	dirAbsolutePath string
}

func LoadContext(dir string) (context *model.ModelContext, err error) {
	loader := &ContextLoader{Dir: dir}
	return loader.Load()
}

func SaveContext(dir string, context *model.ModelContext) (err error) {
	loader := &ContextLoader{Dir: dir}
	return loader.Save(context)
}

func (this_ *ContextLoader) Load() (context *model.ModelContext, err error) {

	var exists bool
	exists, err = base.PathExists(this_.Dir)
	if err != nil {
		return
	}
	if !exists {
		err = errors.New(fmt.Sprint("dir [" + this_.Dir + "] not exist"))
		return
	}
	var abs string
	abs, err = filepath.Abs(this_.Dir)
	if err != nil {
		return
	}
	this_.dirAbsolutePath = filepath.ToSlash(abs)

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
				fmt.Println("to struct model error:", err)
				return
			}
			context.AppendStruct(one)
		} else if strings.Index(namePath, "action/") == 0 {
			var one *model.ActionModel
			one, err = model.TextToActionModel(strings.TrimPrefix(namePath, "action/"), text)
			if err != nil {
				fmt.Println("to action model error:", err)
				return
			}
			context.AppendAction(one)
		} else if strings.Index(namePath, "test/") == 0 {
			var one *model.TestModel
			one, err = model.TextToTestModel(strings.TrimPrefix(namePath, "test/"), text)
			if err != nil {
				fmt.Println("to test model error:", err)
				return
			}
			context.AppendTest(one)
		} else if strings.Index(namePath, "teamide/server/web/") == 0 {
			var one *model.ServerWebModel
			one, err = model.TextToServerWebModel(strings.TrimPrefix(namePath, "teamide/server/web/"), text)
			if err != nil {
				fmt.Println("to server web model error:", err)
				return
			}
			context.AppendServerWeb(one)
		} else if strings.Index(namePath, "dictionary/") == 0 {
			var one *model.DictionaryModel
			one, err = model.TextToDictionaryModel(strings.TrimPrefix(namePath, "dictionary/"), text)
			if err != nil {
				fmt.Println("to dictionary model error:", err)
				return
			}
			context.AppendDictionary(one)
		} else if strings.Index(namePath, "constant/") == 0 {
			var one *model.ConstantModel
			one, err = model.TextToConstantModel(strings.TrimPrefix(namePath, "constant/"), text)
			if err != nil {
				fmt.Println("to constant model error:", err)
				return
			}
			context.AppendConstant(one)
		} else if strings.Index(namePath, "error/") == 0 {
			var one *model.ErrorModel
			one, err = model.TextToErrorModel(strings.TrimPrefix(namePath, "error/"), text)
			if err != nil {
				fmt.Println("to error model error:", err)
				return
			}
			context.AppendError(one)
		} else if strings.Index(namePath, "datasource/database/") == 0 {
			var one *model.DatasourceDatabase
			one, err = model.TextToDatasourceDatabase(strings.TrimPrefix(namePath, "datasource/database/"), text)
			if err != nil {
				fmt.Println("to datasource database model error:", err)
				return
			}
			context.AppendDatasourceDatabase(one)
		} else if strings.Index(namePath, "datasource/redis/") == 0 {
			var one *model.DatasourceRedis
			one, err = model.TextToDatasourceRedis(strings.TrimPrefix(namePath, "datasource/redis/"), text)
			if err != nil {
				fmt.Println("to datasource redis model error:", err)
				return
			}
			context.AppendDatasourceRedis(one)
		} else if strings.Index(namePath, "datasource/kafka/") == 0 {
			var one *model.DatasourceKafka
			one, err = model.TextToDatasourceKafka(strings.TrimPrefix(namePath, "datasource/kafka/"), text)
			if err != nil {
				fmt.Println("to datasource kafka model error:", err)
				return
			}
			context.AppendDatasourceKafka(one)
		} else if strings.Index(namePath, "datasource/zookeeper/") == 0 {
			var one *model.DatasourceZookeeper
			one, err = model.TextToDatasourceZookeeper(strings.TrimPrefix(namePath, "datasource/zookeeper/"), text)
			if err != nil {
				fmt.Println("to datasource zookeeper model error:", err)
				return
			}
			context.AppendDatasourceZookeeper(one)
		}
	}
	return
}

func (this_ *ContextLoader) Save(context *model.ModelContext) (err error) {
	for _, one := range context.Constants {
		err = this_.SaveModel(model.MODEL_TYPE_CONSTANT, one.Name, one)
		if err != nil {
			return
		}
	}
	for _, one := range context.Errors {
		err = this_.SaveModel(model.MODEL_TYPE_ERROR, one.Name, one)
		if err != nil {
			return
		}
	}
	for _, one := range context.Dictionaries {
		err = this_.SaveModel(model.MODEL_TYPE_DICTIONARY, one.Name, one)
		if err != nil {
			return
		}
	}
	for _, one := range context.Structs {
		err = this_.SaveModel(model.MODEL_TYPE_STRUCT, one.Name, one)
		if err != nil {
			return
		}
	}
	for _, one := range context.Actions {
		err = this_.SaveModel(model.MODEL_TYPE_ACTION, one.Name, one)
		if err != nil {
			return
		}
	}
	for _, one := range context.Tests {
		err = this_.SaveModel(model.MODEL_TYPE_TEST, one.Name, one)
		if err != nil {
			return
		}
	}
	for _, one := range context.ServerWebs {
		err = this_.SaveModel(model.MODEL_TYPE_SERVER_WEB, one.Name, one)
		if err != nil {
			return
		}
	}
	for _, one := range context.DatasourceDatabases {
		err = this_.SaveModel(model.MODEL_TYPE_DATASOURCE_DATABASE, one.Name, one)
		if err != nil {
			return
		}
	}
	for _, one := range context.DatasourceRedises {
		err = this_.SaveModel(model.MODEL_TYPE_DATASOURCE_REDIS, one.Name, one)
		if err != nil {
			return
		}
	}
	for _, one := range context.DatasourceKafkas {
		err = this_.SaveModel(model.MODEL_TYPE_DATASOURCE_KAFKA, one.Name, one)
		if err != nil {
			return
		}
	}
	for _, one := range context.DatasourceZookeepers {
		err = this_.SaveModel(model.MODEL_TYPE_DATASOURCE_ZOOKEEPER, one.Name, one)
		if err != nil {
			return
		}
	}
	return
}

func (this_ *ContextLoader) SaveModel(modelType *model.ModelType, name string, data interface{}) (err error) {

	text := ""
	constantModel, constantModelOk := data.(*model.ConstantModel)
	if constantModelOk {
		text, err = model.ConstantModelToText(constantModel)
		if err != nil {
			return
		}
	}
	errorModel, errorModelOk := data.(*model.ErrorModel)
	if errorModelOk {
		text, err = model.ErrorModelToText(errorModel)
		if err != nil {
			return
		}
	}
	dictionaryModel, dictionaryModelOk := data.(*model.DictionaryModel)
	if dictionaryModelOk {
		text, err = model.DictionaryModelToText(dictionaryModel)
		if err != nil {
			return
		}
	}
	structModel, structModelOk := data.(*model.StructModel)
	if structModelOk {
		text, err = model.StructModelToText(structModel)
		if err != nil {
			return
		}
	}
	actionModel, actionModelOk := data.(*model.ActionModel)
	if actionModelOk {
		text, err = model.ActionModelToText(actionModel)
		if err != nil {
			return
		}
	}
	testModel, testModelOk := data.(*model.TestModel)
	if testModelOk {
		text, err = model.TestModelToText(testModel)
		if err != nil {
			return
		}
	}
	serverWebModel, serverWebModelOk := data.(*model.ServerWebModel)
	if serverWebModelOk {
		text, err = model.ServerWebModelToText(serverWebModel)
		if err != nil {
			return
		}
	}
	datasourceDatabase, datasourceDatabaseOk := data.(*model.DatasourceDatabase)
	if datasourceDatabaseOk {
		text, err = model.DatasourceDatabaseToText(datasourceDatabase)
		if err != nil {
			return
		}
	}
	datasourceRedis, datasourceRedisOk := data.(*model.DatasourceRedis)
	if datasourceRedisOk {
		text, err = model.DatasourceRedisToText(datasourceRedis)
		if err != nil {
			return
		}
	}
	datasourceKafka, datasourceKafkaOk := data.(*model.DatasourceKafka)
	if datasourceKafkaOk {
		text, err = model.DatasourceKafkaToText(datasourceKafka)
		if err != nil {
			return
		}
	}
	datasourceZookeeper, datasourceZookeeperOk := data.(*model.DatasourceZookeeper)
	if datasourceZookeeperOk {
		text, err = model.DatasourceZookeeperToText(datasourceZookeeper)
		if err != nil {
			return
		}
	}
	path := this_.Dir + "/" + modelType.Dir + "/" + name + ".yaml"
	if name == "" {
		path = this_.Dir + "/" + modelType.Dir + "/default.yaml"
	}
	var exists bool
	exists, err = base.PathExists(path)
	if err != nil {
		return
	}
	if exists {
		err = os.Remove(path)
		if err != nil {
			return
		}
	}
	var abs string
	abs, err = filepath.Abs(path)
	if err != nil {
		return
	}
	dirAbsolutePath := filepath.ToSlash(abs)
	parentPath := dirAbsolutePath[0:strings.LastIndex(dirAbsolutePath, "/")]
	exists, err = base.PathExists(parentPath)
	if err != nil {
		return
	}
	if !exists {
		err = os.MkdirAll(parentPath, 0777)
		if err != nil {
			return
		}
	}

	var file *os.File
	file, err = os.Create(path)
	if err != nil {
		return
	}
	defer file.Close()
	file.WriteString(text)
	return
}

func loadDirFiles(fileContentMap map[string]string, fileDir string) (err error) {
	//获取当前目录下的所有文件或目录信息
	err = filepath.Walk(fileDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {

		} else {
			var abs string
			abs, err = filepath.Abs(path)
			if err != nil {
				return err
			}
			fileAbsolutePath := filepath.ToSlash(abs)
			name := strings.TrimPrefix(fileAbsolutePath, fileDir)
			name = strings.TrimPrefix(name, "/")
			var f *os.File
			f, err = os.Open(path)
			if err != nil {
				return err
			}
			defer f.Close()
			r := bufio.NewReader(f)
			var line string
			var content string
			for {
				line, err = r.ReadString('\n')
				if err != nil && err != io.EOF {
					return err
				}
				content += line
				if err == io.EOF {
					err = nil
					break
				}
			}
			fileContentMap[name] = content
		}
		return nil
	})
	return
}
